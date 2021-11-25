package main

import (
	"fmt"

	"scraping/cfg"
	"scraping/slack"
)

func main() {
	fmt.Println(cfg.Config)
	fmt.Println(cfg.Keys)
	err := slack.SendMessageToSlack("hankyung", "test with go")
	fmt.Println(err)
}
