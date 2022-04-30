package scraper

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

const MkMSGUrl string = "https://www.mk.co.kr/premium/series/20007/"
const HkIssueTodayUrl string = "https://mobile.hankyung.com/apps/newsletter.view?topic=morning&gnb="
const HostName string = "http://mumeog.site:9090/"

func InitHandler() {
	http.Handle("/maekyung", http.RedirectHandler(MkMSGUrl, http.StatusFound))
	http.Handle("/hankyung", http.RedirectHandler(HkIssueTodayUrl, http.StatusFound))
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Error("ListenAndServe: ", err)
	}
}
