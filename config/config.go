package config

import (
	"encoding/json" // https://pkg.go.dev/encoding/json
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	keyPath     string = "./cfg/keys.json"
	configPath  string = "./cfg/config.json"
	loggingPath string = "./cfg/logging.json"
)

type Kakao_t struct {
	AppId                 string `json:"app_id"`
	Key                   string `json:"key"`
	Template              string `json:"template"`
	RedirectUrl           string `json:"redirect_url"`
	AuthCode              string `json:"authorization_code"`
	AccessToken           string `json:"access_token"`
	RefreshToken          string `json:"refresh_token"`
	ExpiresIn             int    `json:"expires_in"`
	RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"`
	ClientSecret          string `json:"client_secret"`
}

type Fcst_t struct {
	Encoding_key string `json:"encoding_key"`
	Decoding_key string `json:"decoding_key"`
}

type Naver_t struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type Newyo_t struct {
	Apikey string `json:"apikey"`
}

type CommKeys_t struct {
	Slack string  `json:"slack_key"` // 구조체 필드에 태그 지정
	Kakao Kakao_t `json:"kakao"`
	Fcst  Fcst_t  `json:"vilagefcst"`
	Naver Naver_t `json:"naver"`
	Newyo Newyo_t `json:"newyo"`
}

type News_t struct {
	Name    string `json:"name"`
	Flag    bool   `json:"send_flag"`
	SendCnt int
}

type CommCfg_t struct {
	SendHour   int      `json:"send_hour"`
	SendMin    int      `json:"send_min"`
	MaxSendCnt int      `json:"max_send_cnt"`
	Media      []News_t `json:"news"`
}

var Config CommCfg_t
var Keys CommKeys_t

func ChkSendCnt(m *News_t) {
	m.SendCnt += 1
	if m.SendCnt >= Config.MaxSendCnt {
		m.Flag = false
		log.Println("Maximum Send Count reached..")
	}
}

func ResetConfig(m *News_t) {
	m.Flag = true
	m.SendCnt = 0
}

// Load keys from json file
func LoadKeysConfig() CommKeys_t {
	var k CommKeys_t

	path, _ := filepath.Abs(keyPath)
	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&k)
	if err != nil {
		log.Println(err)
	}

	log.Error("< SCRAPER > Successful loading of Key Info ........")

	return k
}

// Load configuration from json file
func LoadCommConfig() CommCfg_t {
	c := CommCfg_t{Media: make([]News_t, 0, 3)}

	path, _ := filepath.Abs(configPath)
	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&c)
	if err != nil {
		log.Println(err)
	}

	log.Error("< SCRAPER > Configuration Informations ............")
	log.Errorf(" - Send Hour      = %d", c.SendHour)
	log.Errorf(" - Send Minute    = %d", c.SendMin)
	log.Errorf(" - Max Send Count = %d", c.MaxSendCnt)
	for i, m := range c.Media {
		log.Errorf(" - Media    < %d > = %s", i, m.Name)
		log.Errorf(" - Flag     < %d > = %t", i, m.Flag)
	}
	return c
}

/*
log.Trace("Something very low level.")
log.Debug("Useful debugging information.")
log.Info("Something noteworthy happened!")
log.Warn("You should probably take a look at this.")
log.Error("Something failed but I'm not quitting.")
// Calls os.Exit(1) after logging
log.Fatal("Bye.")
// Calls panic() after logging
log.Panic("I'm bailing.")
*/

// setup logger
func SetupLogger() {
	path, _ := filepath.Abs(loggingPath)
	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	var l *lumberjack.Logger

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&l)
	if err != nil {
		log.Println(err)
		return
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Error(err, "Err. Failed to get Hostname")
	}
	l.Filename = fmt.Sprintf(l.Filename, hostname)

	// Fork writing into two outputs
	multiWriter := io.MultiWriter(os.Stderr, l) // Stderr와 파일에 동시  출력

	logFormatter := new(log.TextFormatter)
	logFormatter.TimestampFormat = time.RFC1123Z // or RFC3339
	logFormatter.FullTimestamp = true

	log.SetFormatter(logFormatter)
	log.SetLevel(log.InfoLevel)
	log.SetOutput(multiWriter)
	log.SetReportCaller(true) // 해당 이벤트 발생 시 함수, 파일명 표기

	log.Error(" ")
	log.Error("===================================================")
	log.Error(" Scraping News with Go                 S T A R T   ")
	log.Error("===================================================")
	log.Error(" ")
	log.Error("< SCRAPER > Successful Logger setup ...............")
}

func RefreshKeyConfig(k CommKeys_t) {
	f_chg := false

	if k.Slack != Keys.Slack {
		Keys.Slack = k.Slack
		log.Info("Refresh Slack Key")
		f_chg = true
	}

	if Keys.Kakao.AppId != k.Kakao.AppId {
		Keys.Kakao.AppId = k.Kakao.AppId
		log.Info("Refresh Kakao AppId")
		f_chg = true
	}
	if Keys.Kakao.Key != k.Kakao.Key {
		Keys.Kakao.Key = k.Kakao.Key
		log.Info("Refresh Kakao Key")
		f_chg = true
	}
	if Keys.Kakao.RedirectUrl != k.Kakao.RedirectUrl {
		Keys.Kakao.RedirectUrl = k.Kakao.RedirectUrl
		log.Info("Refresh Kakao RedirectUrl")
		f_chg = true
	}
	if Keys.Kakao.AuthCode != k.Kakao.AuthCode {
		Keys.Kakao.AuthCode = k.Kakao.AuthCode
		log.Info("Refresh Kakao AuthCode")
		f_chg = true
	}
	if Keys.Kakao.AccessToken != k.Kakao.AccessToken {
		Keys.Kakao.AccessToken = k.Kakao.AccessToken
		log.Info("Refresh Kakao AccessToken")
		f_chg = true
	}
	if Keys.Kakao.RefreshToken != k.Kakao.RefreshToken {
		Keys.Kakao.RefreshToken = k.Kakao.RefreshToken
		log.Info("Refresh Kakao RefreshToken")
		f_chg = true
	}
	if Keys.Kakao.ExpiresIn != k.Kakao.ExpiresIn {
		Keys.Kakao.ExpiresIn = k.Kakao.ExpiresIn
		log.Info("Refresh Kakao ExpiresIn")
		f_chg = true
	}
	if Keys.Kakao.RefreshTokenExpiresIn != k.Kakao.RefreshTokenExpiresIn {
		Keys.Kakao.RefreshTokenExpiresIn = k.Kakao.RefreshTokenExpiresIn
		log.Info("Refresh Kakao RefreshTokenExpiresIn")
		f_chg = true
	}
	if Keys.Kakao.ClientSecret != k.Kakao.ClientSecret {
		Keys.Kakao.ClientSecret = k.Kakao.ClientSecret
		log.Info("Refresh Kakao ClientSecret")
		f_chg = true
	}

	if f_chg {
		log.Info("Update keys.json")
		enc, err := json.MarshalIndent(Keys, "", " ")
		if err != nil {
			log.Error(err)
		}

		path, _ := filepath.Abs(keyPath)
		f, err := os.Create(path)
		if err != nil {
			log.Error(err)
		}
		_, err = io.WriteString(f, string(enc))
		if err != nil {
			log.Error(err)
		}
	}
}

func init() {
	SetupLogger()
	Config = LoadCommConfig()
	Keys = LoadKeysConfig()
}
