package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	appSecret = os.Getenv("SECRET")
	botToken  = os.Getenv("TOKEN")
	verify    = os.Getenv("VERIFY")
)

func main() {
	http.HandleFunc("/webhook", handler)
	log.Fatal(http.ListenAndServeTLS(":443", "cert.pem", "key.pem", nil))
}

func handler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		query := req.URL.Query()
		verifyToken := query.Get("hub.verify_token")
		if verifyToken != verify {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		io.WriteString(w, req.URL.Query().Get("hub.challenge"))
	} else if req.Method == "POST" {
		read, err := ioutil.ReadAll(req.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if appSecret != "" {
			if req.Header.Get("x-hub-signature") == "" || !checkIntegrity(appSecret, read, req.Header.Get("x-hub-signature")[5:]) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		fmt.Println(string(read[:]))
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// event struct

type upstreamEvent struct {
	Object  string          `json:"object"`
	Entries []*MessageEvent `json:"entry"`
}

type Event struct {
	ID   json.Number `json:"id"`
	Time int64       `json:"time"`
}

type MessageOpts struct {
	Sender struct {
		ID string `json:"id"`
	} `json:"sender"`
	Recipient struct {
		ID string `json:"id"`
	} `json:"recipient"`
	Timestamp int64 `json:"timestamp"`
}

type MessageEvent struct {
	Event
	Messaging []struct {
		MessageOpts
		Message  *ReceivedMessage `json:"message,omitempty"`
		Delivery *Delivery        `json:"delivery,omitempty"`
		Postback *Postback        `json:"postback,omitempty"`
		Optin    *Optin           `json:"optin,empty"`
	} `json:"messaging"`
}

type ReceivedMessage struct {
	ID          string        `json:"mid"`
	Text        string        `json:"text,omitempty"`
	Attachments []*Attachment `json:"attachments,omitempty"`
	Seq         int           `json:"seq"`
}

type Delivery struct {
	MessageIDS []string `json:"mids"`
	Watermark  int64    `json:"watermark"`
	Seq        int      `json:"seq"`
}

type Postback struct {
	Payload string `json:"payload"`
}

type Optin struct {
	Ref string `json:"ref"`
}

// attachment

type AttachmentType string

const (
	AttachmentTypeTemplate AttachmentType = "template"
	AttachmentTypeImage    AttachmentType = "image"
	AttachmentTypeVideo    AttachmentType = "video"
	AttachmentTypeAudio    AttachmentType = "audio"
	AttachmentTypeLocation AttachmentType = "location"
)

type Attachment struct {
	Type    AttachmentType `json:"type"`
	Payload interface{}    `json:"payload,omitempty"`
}

type Resource struct {
	URL string `json:"url"`
}

type Location struct {
	Coordinates Coordinates `json:"coordinates"`
}

type Coordinates struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"long"`
}

// secret check
func checkIntegrity(appSecret string, bytes []byte, expectedSignature string) bool {
	mac := hmac.New(sha1.New, []byte(appSecret))
	mac.Write(bytes)
	if fmt.Sprintf("%x", mac.Sum(nil)) != expectedSignature {
		return false
	}
	return true
}
