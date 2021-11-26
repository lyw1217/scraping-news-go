package cfg

import (
	"encoding/json" // https://pkg.go.dev/encoding/json
	"log"
	"os"
	"path/filepath"
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

	path, _ := filepath.Abs("../.config_secret/keys.json")
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

	path, _ := filepath.Abs("../config/config.json")
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

func init() {
	Config = LoadCommConfig()
	Keys = LoadKeysConfig()
}
