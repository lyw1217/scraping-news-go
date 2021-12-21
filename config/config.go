package config

import (
	"encoding/json" // https://pkg.go.dev/encoding/json
	"io"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	keyPath     string = "./config/keys.json"
	configPath  string = "./config/config.json"
	loggingPath string = "./config/logging.json"
)

type Kakao_t struct {
	AppId       string `json:"app_id"`
	Key         string `json:"key"`
	Template    string `json:"template"`
	RedirectUrl string `json:"redirect_url"`
	Token       string `json:"token"`
}

type CommKeys struct {
	Slack string  `json:"slack_key"` // 구조체 필드에 태그 지정
	Kakao Kakao_t `json:"kakao"`
}

type News struct {
	Name    string `json:"name"`
	Flag    bool   `json:"send_flag"`
	SendCnt int
}

type CommCfg struct {
	SendHour   int    `json:"send_hour"`
	SendMin    int    `json:"send_min"`
	MaxSendCnt int    `json:"max_send_cnt"`
	Media      []News `json:"news"`
}

var Config CommCfg
var Keys CommKeys

func ChkSendCnt(m *News) {
	m.SendCnt += 1
	if m.SendCnt >= Config.MaxSendCnt {
		m.Flag = false
		log.Println("Maximum Send Count reached..")
	}
}

func ResetConfig(m *News) {
	m.Flag = true
	m.SendCnt = 0
}

// Load keys from json file
func LoadKeysConfig() CommKeys {
	var k CommKeys

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
func LoadCommConfig() CommCfg {
	c := CommCfg{Media: make([]News, 0, 2)}

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

func init() {
	SetupLogger()
	Config = LoadCommConfig()
	Keys = LoadKeysConfig()
}
