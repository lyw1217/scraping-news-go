package news

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

type pageInfo struct {
	StatusCode int
	Contents   string
}

func GetHankyungIssueToday(d_month int, d_day int) (int, string) {
	list_url := "https://mobile.hankyung.com/apps/newsletter.view?topic=morning&gnb="

	var p pageInfo

	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: mobile.hankyung.com
		colly.AllowedDomains("mobile.hankyung.com"),
		colly.MaxDepth(1),
	)

	c.OnHTML("table.stb-container", func(e *colly.HTMLElement) {
		fmt.Println("HERE")
		e.ForEach("td.stb-text-box", func(_ int, el *colly.HTMLElement) {
			p.Contents += el.Text
			fmt.Println(el.Text)
		})
		fmt.Println(p.Contents)
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// extract status code
	c.OnResponse(func(r *colly.Response) {
		log.Println("response received", r.StatusCode)
		p.StatusCode = r.StatusCode
	})
	c.OnError(func(r *colly.Response, err error) {
		log.Println("error:", r.StatusCode, err)
		p.StatusCode = r.StatusCode
	})

	// Start scraping on list_url
	c.Visit(list_url)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", " ")
	enc.Encode(p.Contents)
	return p.StatusCode, p.Contents
}
