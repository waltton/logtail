package main

import (
	"github.com/go-chat-bot/bot/slack"

	_ "github.com/waltton/logtail/slackbot/plugin" // Initialize logtail plugin
)

func main() {
	const slackToken string = ""

	slack.Run(slackToken)
}
