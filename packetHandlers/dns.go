package packetHandlers

import (
	"github.com/google/gopacket/layers"
	"sort"
	"time"
)

type dnsInfo struct {
	timestamp time.Time
	questions []layers.DNSQuestion
	answers   []layers.DNSResourceRecord
}

const (
	CUTOFF_DURATION = -1 * time.Hour
)

var (
	dnsPacketInfo []dnsInfo
)

func HandleDnsPacket(dnsPacket layers.DNS) {
	removeOldEntries()
	addNewEntry(dnsPacket)
}

func removeOldEntries() {
	if len(dnsPacketInfo) == 0 {
		return
	}

	//sort packet info by timestamp
	sort.Slice(dnsPacketInfo, func(i, j int) bool {
		return dnsPacketInfo[i].timestamp.Before(
			dnsPacketInfo[j].timestamp)
	})

	//remove entries older than the cutoff time
	hourAgo := time.Now().Add(CUTOFF_DURATION)

	cutoffIndex := 0
	for i := 0; i < len(dnsPacketInfo); i++ {
		info := dnsPacketInfo[i]

		if info.timestamp.Before(hourAgo) {
			cutoffIndex++;
		} else {
			break
		}
	}

	dnsPacketInfo = dnsPacketInfo[cutoffIndex:]
}

func addNewEntry(dnsPacket layers.DNS) {
	dnsInfo := dnsInfo{
		time.Now(),
		dnsPacket.Questions,
		dnsPacket.Answers,
	}
	dnsPacketInfo = append(dnsPacketInfo, dnsInfo)
}