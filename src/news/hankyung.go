package news

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
)

func GetHankyungIssueToday(d_month int, d_day int) (int, string) {
	list_url := "https://mobile.hankyung.com/apps/newsletter.view?topic=morning&gnb="

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
			content := ConvEuckrToUtf8(strings.ReplaceAll(s.Text(), "\u00a0", " "))
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
