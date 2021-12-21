package scraper

import (
	"log"
	"net/http"
)

const MkMSGUrl string = "https://www.mk.co.kr/premium/series/20007/"
const HkIssueTodayUrl string = "https://mobile.hankyung.com/apps/newsletter.view?topic=morning&gnb="
const HostName string = "http://localhost:9090/"

func InitHandler() {
	http.Handle("/hankyung", http.RedirectHandler(MkMSGUrl, http.StatusFound))
	http.Handle("/maekyung", http.RedirectHandler(HkIssueTodayUrl, http.StatusFound))
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
