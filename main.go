package main

import (
	"github.com/nlopes/slack"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const SLACK_TOKEN = "xoxb-xxxxxxxxxx-xxxxxxxxxxxxxxxxxxxxxxxx"
const HTTP_PATH = "http://127.0.0.1/slackbot/"
const BASIC_USER = "slackbot" // "" if no auth
const BASIC_PASSWORD = "slackbot"

func main() {

	log.Println("Slackbot started")

	api := slack.New(SLACK_TOKEN)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

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
				req, err := http.NewRequest("GET", HTTP_PATH+"?user="+ev.Msg.User+"&message="+ev.Msg.Text+"&email="+userInfo.Profile.Email, nil)

				if BASIC_USER != "" {
					req.SetBasicAuth(BASIC_USER, BASIC_PASSWORD)
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
