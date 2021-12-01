package scraper

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"scraping-news/config"
	"scraping-news/util"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
)

type linkInfo struct {
	Title string
	Url   string
}

type pageInfo struct {
	StatusCode int
	Links      []linkInfo
	Contents   []string
}

func (p *pageInfo) AddLinks(t string, u string) []linkInfo {
	// https://stackoverflow.com/questions/18042439/go-append-to-slice-in-struct
	// https://stackoverflow.com/questions/34329441/golang-struct-array-values-not-appending-in-loop
	p.Links = append(p.Links, linkInfo{
		Title: t,
		Url:   u,
	})
	return p.Links
}

// url로 HTTP GET 요청하여 http.Response 객체 반환
func requestGetDocument(url string) (*http.Response, error) {
	// Request 객체 생성
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err, "Err, Failed to NewRequest()")
		return nil, err
	}

	req.Header.Add("User-Agent", "Mozilla/5.0") // 안티 크롤링 회피
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err, "Err, Failed to Get Request")
		return nil, err
	}

	return resp, err
}

// scrap maekyung M.S.G on the date as parameters
func GetMaekyungMSG(d_month int, d_day int) (int, string) {
	list_url := mkMSGUrl

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
			title := util.ConvEuckrToUtf8(s.Find("a").Text())
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

						content := util.ConvEuckrToUtf8(s.Text())
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

// scrap hankyung issue today on the date as a parameters
func GetHankyungIssueToday(d_month int, d_day int) (int, string) {
	list_url := hkIssueTodayUrl

	p := &pageInfo{Contents: make([]string, 0, 10)}

	// 1. Issue Today 조회
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
		container := html.Find("td.stb-text-box")

		// text-box 순회하면서 문자열에 추가
		container.Each(func(i int, s *goquery.Selection) {
			content := util.ConvEuckrToUtf8(strings.ReplaceAll(s.Text(), "\u00a0", " "))
			if content == "" {
				log.Warn("Err. Failed to convert content : ", s.Text())
				return
			}

			// https://stackoverflow.com/questions/65533097/replace-nbsp-or-0xao-with-space-in-a-string
			//content := TransEuckrToUtf8(strings.ReplaceAll(s.Text(), "\u00a0", " "))

			if !(strings.Contains(content, "카카오톡으로 공유하세요")) {
				p.Contents = append(p.Contents, content)
			}
		})

		t_date := strings.Split(p.Contents[0], ".")
		if len(t_date) >= 3 {
			t_year, _ := strconv.Atoi(strings.TrimSpace(t_date[0]))
			t_month, _ := strconv.Atoi(strings.TrimSpace(t_date[1]))
			t_day, _ := strconv.Atoi(strings.TrimSpace(t_date[2]))

			if d_month == t_month && d_day == t_day {
				return p.StatusCode, strings.Join(p.Contents, "\r\n\n")
			} else {
				return p.StatusCode, fmt.Sprintf("No article on %d-%d-%d", t_year, d_month, d_day)
			}
		}
	}

	log.Warn("Err. Failed to get the Issue Today.")
	return resp.StatusCode, err.Error()
}

var gSysClose bool

// start scraping
func StartScraping() {
	log.Error("< SCRAPER > Start Scraping Routine - Started ......")

	c := config.Config

	for !gSysClose {
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
						StatusCode, contents := GetHankyungIssueToday(d_month, d_day)
						if StatusCode != 200 {
							log.Warn("Err. news.GetHankyungIssueToday, StatusCode :", StatusCode)
							config.ChkSendCnt(&c.Media[i])
							continue
						}

						if err := util.SendMessageToSlack("한국경제 Issue Today", contents); err != nil {
							log.Warn("Err. slack.SendMessageToSlack")
							config.ChkSendCnt(&c.Media[i])
							continue
						}

						log.Info(contents)
						c.Media[i].Flag = false

					// 매일 경제
					case "maekyung":
						StatusCode, contents := GetMaekyungMSG(d_month, d_day)
						if StatusCode != 200 {
							log.Warn("Err. news.GetMaekyungMSG, StatusCode :", StatusCode)
							config.ChkSendCnt(&c.Media[i])
							continue
						}

						if err := util.SendMessageToSlack("매일경제 매.세.지", contents); err != nil {
							log.Warn("Err. slack.SendMessageToSlack")
							config.ChkSendCnt(&c.Media[i])
							continue
						}

						log.Info(contents)
						c.Media[i].Flag = false
					default:
						log.Warn("Err. Wrong Key")
					}
				}
			} else if c.SendHour != d_hour {
				config.ResetConfig(&c.Media[i])
			}
		}
		time.Sleep(time.Duration(1) * time.Second)
	}

	log.Error("< SCRAPER > Exit StartScraping Routine ...")
}
