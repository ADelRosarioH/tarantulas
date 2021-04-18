package utils

import (
	"os"

	"github.com/slack-go/slack"
)

func Notify(text string) error {
	SLACK_API_TOKEN := os.Getenv("SLACK_API_TOKEN")
	SLACK_DEFAULT_CHANNEL := os.Getenv("SLACK_DEFAULT_CHANNEL")

	api := slack.New(SLACK_API_TOKEN)

	_, _, _, err := api.SendMessage(SLACK_DEFAULT_CHANNEL, slack.MsgOptionText(text, false))
	if err != nil {
		return err
	}

	return nil
}
