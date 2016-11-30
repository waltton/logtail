package main

import (
	"github.com/go-chat-bot/bot/slack"

	_ "github.com/waltton/logtail/slackbot/plugin" // Initialize logtail plugin
)

func main() {
	const slackToken string = "xoxb-110355919763-q0Tk4US3UgCkKs1Pm4WXMt4m"

	slack.Run(slackToken)
}
