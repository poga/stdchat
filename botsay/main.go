package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

var (
	appSecret   = ""
	botToken    = ""
	verify      = ""
	tlsCertFile = ""
	tlsKeyFile  = ""
	GraphAPI    = "https://graph.facebook.com"
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
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
		resp, err := sendSimpleMessage("1029299690525296", "hi")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(resp)
	}
}
