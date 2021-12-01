package main

import (
	"scraping-news/scraper"
)

func main() {
	go scraper.StartScraping()
	scraper.WaitSignal()
}
