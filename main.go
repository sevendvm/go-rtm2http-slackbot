package main

import (
	"encoding/json"
	"github.com/nlopes/slack"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Config struct {
	SlackToken    string
	HttpPath      string
	BasicUser     string
	BasicPassword string
}

func readJSON(fn string, v interface{}) {
	file, _ := os.Open(fn)
	defer file.Close()
	decoder := json.NewDecoder(file)
	err := decoder.Decode(v)
	if err != nil {
		log.Println("error:", err)
	}
}

var config Config

func main() {

	config = Config{}
	readJSON("config.json", &config)

	api := slack.New(config.SlackToken)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	log.Println("Slackbot started")

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.HelloEvent:
				// Ignore hello
			case *slack.ConnectedEvent:
				log.Println("Connection counter:", ev.ConnectionCount)
			case *slack.MessageEvent:
				if ev.Msg.User == "" {
					continue
				}

				// only direct channels reply
				if !strings.HasPrefix(ev.Msg.Channel, "D") {
					continue
				}
				_, _, userChannel, err := rtm.OpenIMChannel(ev.Msg.User)
				if err != nil {
					log.Printf("Get channel Error: %v\n", err)
					continue
				}
				if userChannel != ev.Msg.Channel {
					continue
				}

				userInfo, err := rtm.GetUserInfo(ev.Msg.User)
				if err != nil {
					log.Printf("User info Error: %v\n", err)
					continue
				}

				log.Printf("User: %v; Message: %v\n", userInfo.Profile.Email, ev.Msg.Text) // ev.Msg.User, ev.Msg.Channel

				client := &http.Client{}
				req, err := http.NewRequest("GET", config.HttpPath+"?user="+ev.Msg.User+"&message="+ev.Msg.Text+"&email="+userInfo.Profile.Email, nil)

				if config.BasicUser != "" {
					req.SetBasicAuth(config.BasicUser, config.BasicPassword)
				}

				resp, err := client.Do(req)
				if err != nil {
					log.Printf("HTTP request Error: %v\n", err)
					continue
				}
				defer resp.Body.Close()

				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Printf("HTTP read Error %v\n", err)
				} else {
					log.Printf("HTTP response: %+v\n", string(body))
					params := slack.PostMessageParameters{}
					rtm.PostMessage(ev.Msg.Channel, string(body), params)
				}

			case *slack.InvalidAuthEvent:
				log.Printf("Invalid credentials")
				break Loop

			default:
				// Ignore other events..
			}
		}
	}
}
