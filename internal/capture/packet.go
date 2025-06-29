package capture

import "time"

type PacketInfo struct {
	Timestamp time.Time
	SrcIP     string
	DstIP     string
	SrcPort   string
	DstPort   string
	Protocol  string
	Length    int
	RawData   []byte
}
