package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/james-vaughn/cipher/packetHandlers"
)

var (
	Interface string
	DnsConfig packetHandlers.DnsPacketHandlerConfiguration
)

func init() {
	config := parseConfig()
	Interface = config.Interface
	DnsConfig = mapDnsConfig(config.DnsConfig)
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

func mapDnsConfig(config DnsConfiguration) packetHandlers.DnsPacketHandlerConfiguration {
	return packetHandlers.DnsPacketHandlerConfiguration{
		CutoffDuration:          -time.Duration(config.CutoffMinutes) * time.Minute,
		TriggerThreshold:        config.TriggerThreshold,
		DurationBetweenTriggers: time.Duration(config.MinutesBetweenTriggers) * time.Minute,
	}
}
