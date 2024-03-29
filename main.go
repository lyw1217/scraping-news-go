package main

import (
	"scraping-news/scraper"
	"scraping-news/util"
)

func main() {
	go scraper.StartScraping()
	go scraper.InitHandler()
	go util.KakaoCheckExpireToken()
	go scraper.GetZoneCode()
	scraper.WaitSignal()
}
