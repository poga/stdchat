package main

import (
	"bufio"
	"encoding/json"
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
		txt := scanner.Text()
		fmt.Println(txt)
		data, err := parseInput(txt)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v\n", data)
		resp, err := sendSimpleMessage(data["recipient"], data["message"])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(resp)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func parseInput(input string) (map[string]string, error) {
	var data map[string]string
	err := json.Unmarshal([]byte(input), &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
