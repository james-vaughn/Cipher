package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/james-vaughn/cipher/emailer"

	"github.com/james-vaughn/cipher/packetHandlers"
)

type AppConfiguration struct {
	Interface   string             `json:"interface"`
	DnsConfig   DnsConfiguration   `json:"dns"`
	EmailConfig EmailConfiguration `json:"email"`
}

type DnsConfiguration struct {
	CutoffMinutes          int `json:"cutoffMinutes"`
	TriggerThreshold       int `json:"triggerThreshold"`
	MinutesBetweenTriggers int `json:"minutesBetweenTriggers"`
}

type EmailConfiguration struct {
	ToAddress      string `json:"to"`
	FromAddress    string `json:"from"`
	Password       string `json:"password"`
	SmtpServerHost string `json:"smtpServerHostname"`
	SmtpServerPort int    `json:"smtpServerPort"`
}

var (
	Interface string
	DnsConfig packetHandlers.DnsPacketHandlerConfiguration
)

func init() {
	config := parseConfig()

	Interface = config.Interface
	emailer := mapEmailConfig(config.EmailConfig)
	DnsConfig = mapDnsConfig(config.DnsConfig, emailer)
}

func parseConfig() AppConfiguration {
	jsonFile, err := os.Open(CONFIG_LOCATION)

	if err != nil {
		log.Fatalf("Error opening configuration json: %v", err)
	}
	defer jsonFile.Close()

	var appConfig AppConfiguration
	configBytes, _ := ioutil.ReadAll(jsonFile)

	if err := json.Unmarshal(configBytes, &appConfig); err != nil {
		log.Fatalf("Error configuring application: %v", err)
	}

	return appConfig
}

func mapDnsConfig(config DnsConfiguration, emailer emailer.Emailer) packetHandlers.DnsPacketHandlerConfiguration {
	return packetHandlers.DnsPacketHandlerConfiguration{
		CutoffDuration:          -time.Duration(config.CutoffMinutes) * time.Minute,
		TriggerThreshold:        config.TriggerThreshold,
		DurationBetweenTriggers: time.Duration(config.MinutesBetweenTriggers) * time.Minute,
		Emailer:                 emailer,
	}
}

func mapEmailConfig(config EmailConfiguration) emailer.Emailer {
	return emailer.Emailer{
		ToAddress:      config.ToAddress,
		FromAddress:    config.FromAddress,
		Password:       config.Password,
		SmtpServerHost: config.SmtpServerHost,
		SmtpServerPort: config.SmtpServerPort,
	}
}
