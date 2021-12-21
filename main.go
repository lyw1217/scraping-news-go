package main

import (
	"scraping-news/scraper"
)

func main() {
	go scraper.StartScraping()
	go scraper.InitHandler()
	scraper.WaitSignal()
}
