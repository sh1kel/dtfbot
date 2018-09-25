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
	DbName   string `json:"DbName"`
	User     string `json:"Username"`
	Password string `json:"Password"`
}

type Configuration struct {
	MailServer Mail `json:"MailAccount"`
	Database   Db   `json:"DbAccount"`
}

type User struct {
	id        int    `db:"id"`
	FullName  string `db:"Name"`
	Email     string `db:"Email"`
	Password  string `db:"Password"`
	Cookie    string `db:"Cookie"`
	Confirmed bool   `db:"Confirmed"`
}

func loadConfig() Configuration {
	var configData Configuration
	config, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Printf("Can't read config file: %v\n", err)
		os.Exit(1)
	}
	json.Unmarshal(config, &configData)
	return configData
}
