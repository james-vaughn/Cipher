//Infomation on gopacket https://www.devdungeon.com/content/packet-capture-injection-and-analysis-gopacket
//Example code https://gist.github.com/dvas0004/3d280f95a83a3084b314ea208f19f9e9
package main

import (
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/james-vaughn/cipher/packetHandlers"
)

type AppConfiguration struct {
	Interface string           `json:"interface"`
	DnsConfig DnsConfiguration `json:"dns"`
}

type DnsConfiguration struct {
	CutoffMinutes          int `json:"cutoffMinutes"`
	TriggerThreshold       int `json:"triggerThreshold"`
	MinutesBetweenTriggers int `json:"minutesBetweenTriggers"`
}

const (
	//Amount of packet collected; long enough for headers
	SNAPSHOT_LEN    = 4096
	CONFIG_LOCATION = "appsettings.config"
)

func main() {
	log.Printf("Opening %s for capturing...\n", Interface)
	handle, err := pcap.OpenLive(Interface, SNAPSHOT_LEN, true, pcap.BlockForever)
	if err != nil {
		log.Fatalf("Couldn't open device for capture: %v", err)
	}

	defer handle.Close()

	capturePackets(handle)
}

func capturePackets(handle *pcap.Handle) {
	//Need all of the layers up to DNS to parse out just DNS
	var eth layers.Ethernet
	var ip4 layers.IPv4
	var ip6 layers.IPv6
	var tcp layers.TCP
	var udp layers.UDP
	var dns layers.DNS
	var payload gopacket.Payload

	parser := gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, &eth, &ip4, &ip6, &tcp, &udp, &dns, &payload)
	decodedLayers := make([]gopacket.LayerType, 0, 10)

	for {
		data, _, err := handle.ReadPacketData()

		if err != nil {
			log.Println("Error reading packet data: %w", err)
			continue
		}

		parser.DecodeLayers(data, &decodedLayers)

		for _, layerType := range decodedLayers {
			switch layerType {
			case layers.LayerTypeDNS:
				packetHandlers.HandleDnsPacket(dns, DnsConfig)
			}
		}
	}
}
