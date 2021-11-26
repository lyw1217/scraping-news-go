package news

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/djimenez/iconv-go" // https://pkg.go.dev/github.com/djimenez/iconv-go
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
)

type pageInfo struct {
	StatusCode int
	Contents   []string
}

func ConvEuckrToUtf8(input string) string {
	output, err := iconv.ConvertString(input, "euc-kr", "utf-8")
	if err != nil {
		log.Println(err)
	}
	return output
}

func TransEuckrToUtf8(input string) string {
	euckrDec := korean.EUCKR.NewDecoder()

	output, _, err := transform.String(euckrDec, input)
	if err != nil {
		log.Println(err)
	}
	return output
}

func GetHankyungIssueToday(d_month int, d_day int) (int, string) {
	list_url := "https://mobile.hankyung.com/apps/newsletter.view?topic=morning&gnb="

	p := &pageInfo{Contents: make([]string, 0, 10)}

	// Request
	resp, err := http.Get(list_url)
	if err != nil {
		log.Println(err, "Err, Failed to Get Request")
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		// HTML Read
		html, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Println(err, "Err. Failed to NewDocumentFromReader()")
		}

		// 파싱
		container := html.Find("td.stb-text-box")

		// text-box 순회하면서 문자열에 추가
		container.Each(func(i int, s *goquery.Selection) {
			content := ConvEuckrToUtf8(strings.ReplaceAll(s.Text(), "\u00a0", " "))

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
				return resp.StatusCode, strings.Join(p.Contents, "\r\n\n")
			} else {
				return resp.StatusCode, fmt.Sprintf("No article on %d-%d-%d", t_year, t_month, t_day)
			}
		}
	}

	return resp.StatusCode, "Err. Failed to get the Issue Today."
}
