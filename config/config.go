package config

import (
	"encoding/json" // https://pkg.go.dev/encoding/json
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
	"time"
)

const (
	keyPath     string = "./config/keys.json"
	configPath  string = "./config/config.json"
	loggingPath string = "./config/logging.json"
)

type CommKeys struct {
	Slack string `json:"SLACK_KEY"` // 구조체 필드에 태그 지정
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
		log.Println("Err. Max Send Count")
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
}

func init() {
	Config = LoadCommConfig()
	Keys = LoadKeysConfig()
	SetupLogger()
}
