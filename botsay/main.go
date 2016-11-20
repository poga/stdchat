package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	messenger "github.com/maciekmm/messenger-platform-go-sdk"
	"github.com/spf13/viper"
)

var (
	botToken = ""
)

func main() {
	viper.SetConfigName("bot")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Read config error: %s\n", err)
	}
	viper.AutomaticEnv()

	botToken = viper.GetString("TOKEN")

	messenger := &messenger.Messenger{
		AccessToken: botToken,
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		txt, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		log.Println("read", txt)
		data, err := parseInput(txt)
		if err != nil {
			log.Fatal("parse error", err)
		}
		fmt.Printf("%v\n", data)
		resp, err := messenger.SendSimpleMessage(data["recipient"], data["message"])
		if err != nil {
			log.Fatal("send msg error", err)
		}
		fmt.Printf("%v\n", resp)
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
