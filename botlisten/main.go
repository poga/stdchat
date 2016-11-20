package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

var (
	appSecret   = ""
	botToken    = ""
	verify      = ""
	tlsCertFile = ""
	tlsKeyFile  = ""
)

func main() {
	viper.SetConfigName("bot")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Read config error: %s\n", err)
	}
	viper.AutomaticEnv()

	appSecret = viper.GetString("APP_SECRET")
	botToken = viper.GetString("TOKEN")
	verify = viper.GetString("VERIFY")
	tlsCertFile = viper.GetString("TLS_CERT")
	tlsKeyFile = viper.GetString("TLS_KEY")

	http.HandleFunc("/webhook", handler)
	log.Fatal(http.ListenAndServeTLS(":443", tlsCertFile, tlsKeyFile, nil))
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

// secret check
func checkIntegrity(appSecret string, bytes []byte, expectedSignature string) bool {
	mac := hmac.New(sha1.New, []byte(appSecret))
	mac.Write(bytes)
	if fmt.Sprintf("%x", mac.Sum(nil)) != expectedSignature {
		return false
	}
	return true
}
