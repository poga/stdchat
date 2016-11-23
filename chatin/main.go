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
	appSecret   string
	botToken    string
	verify      string
	tlsCertFile string
	tlsKeyFile  string
	bot         *messenger.Messenger

	appDescription string
	appPrivacy     string
)

func main() {
	viper.SetConfigName("chat")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Read config error: %s\n", err)
	}
	viper.AutomaticEnv()
	viper.SetDefault("APP_DESC", "JUST A BOT")
	viper.SetDefault("APP_PRIVACY", "ALL YOUR DATA ARE BELONG TO US")

	appSecret = viper.GetString("APP_SECRET")
	botToken = viper.GetString("TOKEN")
	verify = viper.GetString("VERIFY")
	tlsCertFile = viper.GetString("TLS_CERT")
	tlsKeyFile = viper.GetString("TLS_KEY")
	appDescription = viper.GetString("APP_DESC")
	appPrivacy = viper.GetString("APP_PRIVACY")

	bot = &messenger.Messenger{
		VerifyToken: verify,
		AppSecret:   appSecret,
		AccessToken: botToken,
	}
	bot.MessageReceived = onMessageReceived

	http.HandleFunc("/webhook", bot.Handler)
	http.HandleFunc("/privacy", privacy)
	http.HandleFunc("/whatisthis", index)
	log.Fatal(http.ListenAndServeTLS(":443", tlsCertFile, tlsKeyFile, nil))
}

func onMessageReceived(e messenger.Event, opts messenger.MessageOpts, msg messenger.ReceivedMessage) {
	profile, err := bot.GetProfile(opts.Sender.ID)
	if err != nil {
		log.Fatal(err)
	}
	out := Out{e, opts, msg, *profile}
	b, err := json.Marshal(out)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b[:]))
}

type Out struct {
	E       messenger.Event           `json:"event"`
	Opts    messenger.MessageOpts     `json:"opts"`
	Msg     messenger.ReceivedMessage `json:"message"`
	Profile messenger.Profile         `json:"profile"`
}

func privacy(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ALL YOUR DATA ARE BELONG TO US"))
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("JUST A BOT"))
}
