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
	Name string `json:"name"`
	Flag bool   `json:"send_flag"`
}

type CommCfg struct {
	SendHour int    `json:"send_hour"`
	Media    []News `json:"news"`
}

var Config CommCfg
var Keys CommKeys

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
	var c CommCfg

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
