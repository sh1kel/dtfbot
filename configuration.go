package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Mail struct {
	Host     string `json:"Host"`
	ImapPort string `json:"Port"`
	User     string `json:"Username"`
	Password string `json:"Password"`
}

type Db struct {
	Host     string `json:"Host"`
	Port     string `json:"Port"`
	User     string `json:"Username"`
	Password string `json:"Password"`
}

type Configuration struct {
	MailServer Mail `json:"MailAccount"`
	Postgres   Db   `json:"DbAccount"`
}

func loadConfig() Configuration {
	var configData Configuration
	config, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Printf("Can't read config file: %v\n", err)
		os.Exit(1)
	}
	json.Unmarshal(config, &configData)
	fmt.Printf("%s\n", configData)
	return configData
}
