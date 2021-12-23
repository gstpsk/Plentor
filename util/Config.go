package util

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type config struct {
	BCRYPT_COST int
	MONGODB_URI string
}

var Config config

func LoadConfig() {
	content, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Failed to read config.json: %s", err)
	}
	json.Unmarshal(content, &Config)
}
