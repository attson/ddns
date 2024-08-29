package main

import (
	"encoding/json"
	"fmt"
	"github.com/attson/ddns/dns"
	"github.com/attson/ddns/dns/common"
	"github.com/attson/ddns/getip"
	"log"
	"os"
)

type IPProvider struct {
	Type  string                 `mapstructure:"type" json:"type"`
	Extra map[string]interface{} `mapstructure:"extra" json:"extra"`
}

type DNSProvider struct {
	Type  string                 `mapstructure:"type" json:"type"`
	Extra map[string]interface{} `mapstructure:"extra" json:"extra"`
}

type Config struct {
	IPProvider  IPProvider  `mapstructure:"ip_provider" json:"ip_provider"`
	DNSProvider DNSProvider `mapstructure:"dns_provider" json:"dns_provider"`
}

var c = &Config{}

func configuration() {
	file, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(file, c)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	configuration()

	getIp := getip.IpProvider(c.IPProvider.Type)
	ip, err := getIp(c.IPProvider.Extra)
	if err != nil {
		log.Fatal(err)
	}

	dnsProvider := dns.GetDNSProvider(c.DNSProvider.Type)
	r, err := dnsProvider.GetRecord(c.DNSProvider.Extra)
	if err != nil {
		if err == common.RecordNotFound {
			log.Println("record not found, adding...")
			err = dnsProvider.AddRecord(c.DNSProvider.Extra, ip)
			if err != nil {
				log.Fatal(err)
			}

			log.Println("record added")
			return
		}

		log.Fatal(err)
	}

	if r.Value() != ip {
		log.Println(fmt.Sprintf("record value %s is not equal to ip %s, updating...", r.Value(), ip))
		err = dnsProvider.UpdateRecord(c.DNSProvider.Extra, r, ip)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("record updated")
	} else {
		log.Println("record value is equal to ip, no need to update")
	}
}
