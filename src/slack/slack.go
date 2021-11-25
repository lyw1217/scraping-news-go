package slack

import (
	"log"
	"scraping/cfg"

	"github.com/slack-go/slack"
)

// send message to slack
func SendMessageToSlack(media string, msg string) error {
	api := slack.New(cfg.Keys.Slack)
	attachment := slack.Attachment{
		Pretext: media,
		Text:    msg,
	}

	//channelID, timestamp, err := api.PostMessage(
	_, _, err := api.PostMessage(
		"python-trading-bot",
		//slack.MsgOptionText("Some text", false),
		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionAsUser(true), // Add this if you want that the bot would post message as a user, otherwise it will send response using the default slackbot
	)
	if err != nil {
		log.Println(err)
		return err
	}
	//fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)

	return err
}
