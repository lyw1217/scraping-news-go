package main

import (
	"fmt"
	"log"
	"time"

	"scraping/news"
	"scraping/slack"
)

func main() {

	d_month := int(time.Now().Month())
	d_day := time.Now().Day()

	StatusCode, contents := news.GetHankyungIssueToday(d_month, d_day)
	if StatusCode != 200 {
		log.Println(StatusCode)
		return
	}

	err := slack.SendMessageToSlack("hankyung", contents)
	fmt.Println(err)
}
