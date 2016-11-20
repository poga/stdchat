package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	messenger "github.com/maciekmm/messenger-platform-go-sdk"
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

	messenger := &messenger.Messenger{
		VerifyToken: verify,
		AppSecret:   appSecret,
		AccessToken: botToken,
	}
	messenger.MessageReceived = onMessageReceived

	http.HandleFunc("/webhook", messenger.Handler)
	log.Fatal(http.ListenAndServeTLS(":443", tlsCertFile, tlsKeyFile, nil))
}

func onMessageReceived(e messenger.Event, opts messenger.MessageOpts, msg messenger.ReceivedMessage) {
	out := Out{e, opts, msg}
	b, err := json.Marshal(out)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b[:]))
}

type Out struct {
	E    messenger.Event           `json:"event"`
	Opts messenger.MessageOpts     `json:"opts"`
	Msg  messenger.ReceivedMessage `json:"message"`
}
