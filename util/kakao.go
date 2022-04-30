package util

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"scraping-news/config"

	log "github.com/sirupsen/logrus"
)

const kakaoSendToMeUrl string = "https://kapi.kakao.com/v2/api/talk/memo/default/send"
const kakaoReqCodeUrl string = "https://kauth.kakao.com/oauth/authorize"
const kakaoReqTokenUrl string = "https://kauth.kakao.com/oauth/token"
const kakaoInfoTokenUrl string = "https://kauth.kakao.com/v1/user/access_token_info"

type Code_t struct {
	ClientId     string `json:"client_id"`
	RedirectUri  string `json:"redirect_uri"`
	ResponseType string `json:"response_type"`
}

type Token_t struct {
	AccessToken           string `json:"access_token"`
	RefreshToken          string `json:"refresh_token"`
	ExpiresIn             string `json:"expires_in"`
	RefreshTokenExpiresIn string `json:"refresh_token_expires_in"`
	Scope                 string `json:"scope"`
}

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

/*
인가 코드 받기
https://developers.kakao.com/docs/latest/ko/kakaologin/rest-api#request-code
*/
func requestKakaoCode() error {
	log.Info("Request KakaoTalk Api authorize code .")
	k := config.Keys.Kakao

	var c = Code_t{
		ClientId:     k.Key,
		RedirectUri:  k.RedirectUrl,
		ResponseType: "code",
	}

	params := url.Values{}
	params.Add("client_id", c.ClientId)
	params.Add("redirect_uri", c.RedirectUri)
	params.Add("response_type", c.ResponseType)

	// https://golang.cafe/blog/how-to-make-http-url-form-encoded-request-golang.html
	req, err := http.NewRequest("GET", kakaoReqCodeUrl, strings.NewReader(params.Encode())) // URL-encoded payload
	if err != nil {
		log.Error(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	rsp, err := client.Do(req)
	if err != nil {
		log.Error(err)
	}
	defer rsp.Body.Close()

	rspBody, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		str := string(rspBody)
		log.Error(str)
	}

	str := string(rspBody)
	log.Error(str)

	return err
}

/*
토큰 받기
https://developers.kakao.com/docs/latest/ko/kakaologin/rest-api#request-token
*/
func requestKakaoToken() error {
	log.Info("Request KakaoTalk Api Token.")
	k := config.Keys.Kakao

	params := url.Values{}
	params.Add("grant_type", "authorization_code")
	params.Add("client_id", k.Key)
	params.Add("redirect_uri", k.RedirectUrl)
	params.Add("code", k.AuthCode)
	params.Add("client_secret", k.ClientSecret)

	// https://golang.cafe/blog/how-to-make-http-url-form-encoded-request-golang.html
	req, err := http.NewRequest("POST", kakaoReqCodeUrl, strings.NewReader(params.Encode())) // URL-encoded payload
	if err != nil {
		log.Error(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

	client := &http.Client{}
	rsp, err := client.Do(req)
	if err != nil {
		log.Error(err)
	}
	defer rsp.Body.Close()

	rspBody, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		str := string(rspBody)
		log.Error(str)
	}

	str := string(rspBody)
	log.Error(str)

	return err
}

/*
토큰 정보 보기
https://developers.kakao.com/docs/latest/ko/kakaologin/rest-api#get-token-info
*/
func getKakaoTokenInfo() error {
	log.Info("Get KakaoTalk Api Token Info.")
	k := config.Keys.Kakao

	// https://golang.cafe/blog/how-to-make-http-url-form-encoded-request-golang.html
	req, err := http.NewRequest("GET", kakaoReqCodeUrl, nil) // URL-encoded payload
	if err != nil {
		log.Error(err)
	}
	req.Header.Add("Authorization", "Bearer "+k.AccessToken)

	client := &http.Client{}
	rsp, err := client.Do(req)
	if err != nil {
		log.Error(err)
	}
	defer rsp.Body.Close()

	rspBody, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		str := string(rspBody)
		log.Error(str)
	}

	str := string(rspBody)
	log.Error(str)

	return err
}

/*
토큰 갱신하기
koe010 error : https://devtalk.kakao.com/t/react-invalid-client-koe010/114139
*/
func refreshKakaoToken() error {
	log.Info("Refresh KakaoTalk Api Access Token.")

	k := config.Keys

	params := url.Values{}
	params.Add("grant_type", "refresh_token")
	params.Add("client_id", k.Kakao.Key)
	params.Add("refresh_token", k.Kakao.RefreshToken)
	params.Add("client_secret", k.Kakao.ClientSecret)

	// https://golang.cafe/blog/how-to-make-http-url-form-encoded-request-golang.html
	req, err := http.NewRequest("POST", kakaoReqTokenUrl, strings.NewReader(params.Encode())) // URL-encoded payload
	if err != nil {
		log.Error(err)
	}
	req.Header.Add("Content-type", "application/x-www-form-urlencoded;charset=utf-8")

	client := &http.Client{}
	rsp, err := client.Do(req)
	if err != nil {
		log.Error(err)
	}
	defer rsp.Body.Close()

	rspBody, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		str := string(rspBody)
		log.Error(str)
	}

	err = json.Unmarshal(rspBody, &k.Kakao)
	if err != nil {
		log.Error(err)
	}

	config.RefreshKeyConfig(k)

	return err
}

func KakaoSendToMe(news string, msg string, lnk string) error {
	log.Info("Send me a message through KakaoTalk.")
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
	//fmt.Println(buff.String())
	data := url.Values{}
	data.Set("template_object", buff.String())

	// https://golang.cafe/blog/how-to-make-http-url-form-encoded-request-golang.html
	req, err := http.NewRequest("POST", kakaoSendToMeUrl, strings.NewReader(data.Encode())) // URL-encoded payload
	if err != nil {
		log.Error(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Bearer "+k.AccessToken)

	client := &http.Client{}
	rsp, err := client.Do(req)
	if err != nil {
		log.Error(err)
	}
	defer rsp.Body.Close()

	rspBody, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		str := string(rspBody)
		log.Error(str)
	}

	return err
}

var SysClose bool

func KakaoCheckExpireToken() {
	log.Error("< SCRAPER > Start Check Expire Token Routine - Started ......")

	err := refreshKakaoToken()
	if err != nil {
		log.Error(err)
	}

	exp := config.Keys.Kakao.ExpiresIn
	refresh_exp := config.Keys.Kakao.RefreshTokenExpiresIn

	for !SysClose {
		exp--
		refresh_exp--

		if exp < 3600 {
			err := refreshKakaoToken()
			if err != nil {
				log.Error(err)
			}
			exp = config.Keys.Kakao.ExpiresIn

		}

		if refresh_exp < 3600 {
			err := refreshKakaoToken()
			if err != nil {
				log.Error(err)
			}
			refresh_exp = config.Keys.Kakao.RefreshTokenExpiresIn
		}

		time.Sleep(time.Duration(1) * time.Second)
	}

	log.Error("< SCRAPER > Exit KakaoCheckExpireToken Routine ...")
}
