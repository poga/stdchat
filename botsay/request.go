package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

func doRequest(method string, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	query := req.URL.Query()
	query.Set("access_token", botToken)
	req.URL.RawQuery = query.Encode()
	return http.DefaultClient.Do(req)
}

func sendMessage(mq MessageQuery) (*MessageResponse, error) {
	byt, err := json.Marshal(mq)
	if err != nil {
		return nil, err
	}
	resp, err := doRequest("POST", GraphAPI+"/v2.6/me/messages", bytes.NewReader(byt))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	read, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		er := new(rawError)
		json.Unmarshal(read, er)
		return nil, errors.New("Error occured: " + er.Error.Message)
	}
	response := &MessageResponse{}
	err = json.Unmarshal(read, response)
	return response, err
}

//SendSimpleMessage :
func sendSimpleMessage(recipient string, message string) (*MessageResponse, error) {
	return sendMessage(MessageQuery{
		Recipient: Recipient{
			ID: recipient,
		},
		Message: SendMessage{
			Text: message,
		},
	})
}

type MessageResponse struct {
	RecipientID string `json:"recipient_id"`
	MessageID   string `json:"message_id"`
}

type rawMessage struct {
	Recipient
	MessageQuery
}
