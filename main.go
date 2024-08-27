package main

import (
	"encoding/json"
	"github.com/attson/ddns/getip"
	"log"
	"os"
)

type IpProvider struct {
	Type  string                 `mapstructure:"type" json:"type"`
	Extra map[string]interface{} `mapstructure:"extra" json:"extra"`
}

type Config struct {
	IpProvider IpProvider `mapstructure:"ip_provider" json:"ip_provider"`
}

var c = &Config{}

func main() {
	file, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(file, c)
	if err != nil {
		log.Fatal(err)
	}

	getIp := getip.IpProvider(c.IpProvider.Type)
	ip, err := getIp(c.IpProvider.Extra)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(ip)
}
