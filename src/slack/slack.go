package slack

import (
	"fmt"
	"log"
	"scraping/cfg"

	"github.com/slack-go/slack"
)

// send file to slack
func SendFileToSlack(title string, url string) error {
	api := slack.New(cfg.Keys.Slack)
	params := slack.FileUploadParameters{
		Title: title,
		File:  url,
	}
	file, err := api.UploadFile(params)
	if err != nil {
		log.Println(err)
		return err
	}

	fmt.Printf("Name: %s, URL: %s\n", file.Name, file.URL)

	return err
}

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
