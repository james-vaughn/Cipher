package packetHandlers

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/gopacket/layers"
)

type DnsPacketHandlerConfiguration struct {
	CutoffDuration          time.Duration `json:"cutoffDuration"`
	TriggerThreshold        int           `json:"triggerThreshold"`
	DurationBetweenTriggers time.Duration `json:"durationBetweenTriggers"`
}

type dnsInfo struct {
	timestamp time.Time
	questions []layers.DNSQuestion
	answers   []layers.DNSResourceRecord
}

var (
	dnsPacketInfo []dnsInfo
	lastTrigger   time.Time
)

func HandleDnsPacket(dnsPacket layers.DNS, config DnsPacketHandlerConfiguration) {
	removeOldEntries(config.CutoffDuration)
	addNewEntry(dnsPacket)

	nextTriggerTime := lastTrigger.Add(config.DurationBetweenTriggers)
	if time.Now().After(nextTriggerTime) {
		triggerIfThresholdIsMet(config.TriggerThreshold)
	}
}

func removeOldEntries(cutoffDuration time.Duration) {
	if len(dnsPacketInfo) == 0 {
		return
	}

	//sort packet info by timestamp
	//not needed because always sorted by time due to appending to the back?
	//sort.Slice(dnsPacketInfo, func(i, j int) bool {
	//	return dnsPacketInfo[i].timestamp.Before(
	//		dnsPacketInfo[j].timestamp)
	//})

	//remove entries older than the cutoff time
	cutoffTime := time.Now().Add(cutoffDuration)

	cutoffIndex := 0
	for i := 0; i < len(dnsPacketInfo); i++ {
		info := dnsPacketInfo[i]

		if info.timestamp.Before(cutoffTime) {
			cutoffIndex++
		} else {
			break
		}
	}

	dnsPacketInfo = dnsPacketInfo[cutoffIndex:]
}

func addNewEntry(dnsPacket layers.DNS) {
	dnsInfo := dnsInfo{
		timestamp: time.Now(),
		questions: make([]layers.DNSQuestion, len(dnsPacket.Questions)),
		answers:   make([]layers.DNSResourceRecord, len(dnsPacket.Answers)),
	}

	//need to make a copy to avoid copying the reference
	copy(dnsInfo.questions, dnsPacket.Questions)
	copy(dnsInfo.answers, dnsPacket.Answers)

	dnsPacketInfo = append(dnsPacketInfo, dnsInfo)
}

func triggerIfThresholdIsMet(threshold int) {
	if len(dnsPacketInfo) < threshold {
		return
	}

	log.Println("Trigger Hit")
	for _, info := range dnsPacketInfo {
		log.Println(info)
	}

	log.Println("----------------------")
}

//TODO make better
func (d dnsInfo) String() string {
	var dnsQuestions []string

	for _, question := range d.questions {
		timestamp := d.timestamp
		dnsQuestions = append(dnsQuestions, fmt.Sprintf("%v: %s", timestamp, question))
	}

	return strings.Join(dnsQuestions, ", ")
}
