package main

import (
	"scraping-news/scraper"
)

func main() {

	done := make(chan bool)
	go scraper.StartScraping()
	<-done
}
