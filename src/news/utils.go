package news

import (
	"log"
	"net/http"

	"github.com/djimenez/iconv-go" // https://pkg.go.dev/github.com/djimenez/iconv-go
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
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

// euc-kr 문자열을 utf-8 문자열로 변환 (iconv-go)
func ConvEuckrToUtf8(input string) string {
	output, err := iconv.ConvertString(input, "euc-kr", "utf-8")
	if err != nil {
		log.Println(err)
		return ""
	}
	return output
}

// euc-kr 문자열을 utf-8 문자열로 변환 (korean, transform)
func TransEuckrToUtf8(input string) string {
	euckrDec := korean.EUCKR.NewDecoder()

	output, _, err := transform.String(euckrDec, input)
	if err != nil {
		log.Println(err)
		return ""
	}
	return output
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
