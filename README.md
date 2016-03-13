# go-rtm2http-slackbot
This bot can receive direst RTM messages from slack, send it to http server and send reply to to slack user

Bot can use basic http auth of http-server

Bot sends GET-queries to server url with params: user, message and email.

Uses slack acces package "github.com/nlopes/slack"
