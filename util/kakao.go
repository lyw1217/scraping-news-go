package util

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"scraping-news/config"

	log "github.com/sirupsen/logrus"
)

const kakaoSendToMeUrl string = "https://kapi.kakao.com/v2/api/talk/memo/default/send"

type Link_t struct {
	WebUrl       string `json:"web_url"`
	MobileWebUrl string `json:"mobile_web_url"`
}

type TemplateObject_t struct {
	ObjectType  string `json:"object_type"`
	Text        string `json:"text"`
	Link        Link_t `json:"link"`
	ButtonTitle string `json:"button_title"`
}

func KakaoSendToMe(news string, msg string, lnk string) error {
	log.Println("Send me a message through KakaoTalk.")
	k := config.Keys.Kakao

	var l = Link_t{
		WebUrl:       lnk,
		MobileWebUrl: lnk,
	}
	var tobj = TemplateObject_t{
		ObjectType:  "text",
		Text:        msg,
		Link:        l,
		ButtonTitle: news,
	}

	jsonBytes, _ := json.Marshal(tobj)
	buff := bytes.NewBuffer(jsonBytes)
	println(buff.String())
	data := url.Values{}
	data.Set("template_object", buff.String())

	// https://golang.cafe/blog/how-to-make-http-url-form-encoded-request-golang.html
	req, err := http.NewRequest("POST", kakaoSendToMeUrl, strings.NewReader(data.Encode())) // URL-encoded payload
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Bearer "+k.Token)

	client := &http.Client{}
	rsp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer rsp.Body.Close()

	rspBody, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		str := string(rspBody)
		log.Fatal(str)
	}

	return err
}
