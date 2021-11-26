package main

import (
	"log"
	"time"

	"scraping/cfg"
	"scraping/news"
	"scraping/slack"
)

func getMorningNews() {
	c := cfg.Config

	for {
		d_month := int(time.Now().Month())
		d_day := int(time.Now().Day())
		d_hour := int(time.Now().Hour())
		d_min := int(time.Now().Minute())

		for i, media := range c.Media {
			if c.SendHour == d_hour && c.SendMin <= d_min {
				if media.Flag {
					switch media.Name {
					// 한국 경제
					case "hankyung":
						StatusCode, contents := news.GetHankyungIssueToday(d_month, d_day)
						if StatusCode != 200 {
							log.Println("Err. news.GetHankyungIssueToday, StatusCode :", StatusCode)
							cfg.ChkSendCnt(&c.Media[i])
							continue
						}

						if err := slack.SendMessageToSlack("한국경제 Issue Today", contents); err != nil {
							log.Println("Err. slack.SendMessageToSlack")
							cfg.ChkSendCnt(&c.Media[i])
							continue
						}

						log.Println(contents)
						media.Flag = false

					// 매일 경제
					case "maekyung":
						log.Println("아직 미구현")

					default:
						log.Println("Err. Wrong Key")
					}
				}
			} else if c.SendHour != d_hour {
				cfg.ResetConfig(&c.Media[i])
			}
		}
		time.Sleep(10 * time.Second)
	}
}

func main() {
	done := make(chan bool)
	go getMorningNews()
	<-done
}
