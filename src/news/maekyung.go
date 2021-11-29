package news

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
)

func GetMaekyungMSG(d_month int, d_day int) (int, string) {
	list_url := "https://www.mk.co.kr/premium/series/20007/"

	p := &pageInfo{
		Links:    make([]linkInfo, 0, 10),
		Contents: make([]string, 0, 10),
	}

	// 1. 매세지 첫 page 목록 조회
	resp, err := requestGetDocument(list_url)
	if err != nil {
		log.Warn(err, "Err, Failed to Get Request")
		return resp.StatusCode, err.Error()
	}
	defer resp.Body.Close()

	if p.StatusCode = resp.StatusCode; p.StatusCode == 200 {
		// HTML Read
		html, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Warn(err, "Err. Failed to NewDocumentFromReader()")
			return p.StatusCode, err.Error()
		}

		// 파싱
		container := html.Find("dt.tit") // returns a new Selection ojbect

		// 하이퍼링크 순회 및 저장
		container.Each(func(i int, s *goquery.Selection) {
			title := ConvEuckrToUtf8(s.Find("a").Text())
			href, ok := s.Find("a").Attr("href")
			if !ok {
				log.Info(ok, "Err. No Exist href in", p.Links[i].Title)
				return
			}

			p.AddLinks(title, href)
		})

		// 날짜에 맞는 article 확인
		for _, lnk := range p.Links {
			if strings.Contains(lnk.Title, strconv.Itoa(d_month)) &&
				strings.Contains(lnk.Title, strconv.Itoa(d_day)) {

				resp_link, err := requestGetDocument(lnk.Url)
				if err != nil {
					log.Warn(err, "Err, Failed to Get Request")
					break
				}
				defer resp_link.Body.Close()

				if p.StatusCode = resp_link.StatusCode; p.StatusCode == 200 {
					html, err := goquery.NewDocumentFromReader(resp_link.Body)
					if err != nil {
						log.Warn(err, "Err. Failed to NewDocumentFromReader()")
						break
					}

					container := html.Find("#content > div.content_left > div.view_txt")

					// view_txt 순회하면서 문자열에 추가
					container.Each(func(i int, s *goquery.Selection) {

						content := ConvEuckrToUtf8(s.Text())
						if content == "" {
							log.Warn("Err. Failed to convert content : ", s.Text())
							return
						}

						p.Contents = append(p.Contents, content)
					})

					return p.StatusCode, strings.Join(p.Contents, "")

				} else {
					log.Warn("Err. Failed to get M.S.G Article.")
					return p.StatusCode, "Err. Failed to get M.S.G Article."
				}
			}
		}

		return p.StatusCode, fmt.Sprintf("No article on %d-%d", d_month, d_day)
	}

	log.Warn("Err. Failed to get the M.S.G list.")
	return resp.StatusCode, err.Error()
}
