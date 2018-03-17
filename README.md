# go-rtm2http-slackbot
This bot can receive direst RTM messages from slack, send it to http server and send reply to to slack user.
It can translate slack user messages to 1C (example http://infostart.ru/public/581824/).

Bot sends GET-queries to server url with params: **user**, **message** and **email**.


Binaries are here https://github.com/nikepan/go-rtm2http-slackbot/releases

Set settings in config.json (you can use config.sample.json as example):

`"SlackToken": "xoxb-xxxxxxxxxx-xxxxxxxxxxxxxxxxxxxxxxxx"`

`"HttpPath": "http://127.0.0.1/slackbot/"`

You can take token at https://my.slack.com/services/new/bot

Bot can use basic http auth of http-server

`"BasicUser": "slackbot" // "" if no auth`

`"BasicPassword": "slackbot"`

Uses slack access package "github.com/nlopes/slack"
