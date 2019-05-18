//https://www.devdungeon.com/content/packet-capture-injection-and-analysis-gopacket
package main

import (
	"flag"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"log"
	"time"
)

const (
	//Amount of packet collected; long enough for headers
	SNAPSHOT_LEN        = 1024
	//needed to snoop on the interface
	PROMISCUOUS_MODE 	= true
	//time period packets are collected before pushing the results to the Packets channel
	COLLECTION_INTERVAL = 10 * time.Second
)

var (
	INTERFACE_NAME string
)


func init() {
	parseFlags()
}

func main() {
	handle, err := pcap.OpenLive(INTERFACE_NAME, SNAPSHOT_LEN, true, COLLECTION_INTERVAL)
	if err != nil {
		log.Fatalf("Couldn't open device for capture: %v", err)
	}

	defer handle.Close()

	if err := capturePackets(handle); err != nil {
		log.Fatalf("Error encountered while capturing packets: %v", err)
	}
}

func parseFlags() {
	flag.StringVar(&INTERFACE_NAME, "i","", "Network interface name from e.g. ifconfig")
	flag.Parse()
}

func capturePackets(handle *pcap.Handle) error {
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		// Process packet here
		fmt.Println(packet)
	}

	return nil
}