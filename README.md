# go-rtm2http-slackbot
This bot can receive direst RTM messages from slack, send it to http server and send reply to to slack user

Bot sends GET-queries to server url with params: **user**, **message** and **email**.

`const SLACK_TOKEN = "xoxb-xxxxxxxxxx-xxxxxxxxxxxxxxxxxxxxxxxx"`

`const HTTP_PATH = "http://127.0.0.1/slackbot/"`

Bot can use basic http auth of http-server

`const BASIC_USER = "slackbot" // "" if no auth`

`const BASIC_PASSWORD = "slackbot"`

Uses slack access package "github.com/nlopes/slack"
