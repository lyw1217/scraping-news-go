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

//var kakaoLoginUrl string = "https://kauth.kakao.com/oauth/authorize?response_type=code&client_id=" + config.Keys.Kakao.Key + "&redirect_uri=" + config.Keys.Kakao.RedirectUrl

const kakaoSendToMeUrl string = "https://kapi.kakao.com/v2/api/talk/memo/default/send"

type Link_t struct {
	WebUrl       string `json:"web_url"`
	MobileWebUrl string `json:"mobile_web_url"`
}

type TemplateObject struct {
	ObjectType  string `json:"object_type"`
	Text        string `json:"text"`
	Link        Link_t `json:"link"`
	ButtonTitle string `json:"button_title"`
}

func KakaoSendToMe(news string, msg string, lnk string) error {
	log.Println("send msg to me by kakaotalk")
	k := config.Keys.Kakao

	var l = Link_t{
		WebUrl:       lnk,
		MobileWebUrl: lnk,
	}
	var tobj = TemplateObject{
		ObjectType:  "text",
		Text:        msg,
		Link:        l,
		ButtonTitle: news,
	}

	jsonBytes, _ := json.Marshal(tobj)
	buff := bytes.NewBuffer(jsonBytes)

	data := url.Values{}
	data.Set("template_object", buff.String())

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

	rspBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		str := string(rspBody)
		log.Fatal(str)
	}
	println(rspBody)

	return err
}
