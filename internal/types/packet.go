package types

import "wirecrab/internal/tshark"

type PacketMeta struct {
	Timestamp string `json:"Timestamp"`
	SrcIP     string `json:"SrcIP,omitempty"`
	DstIP     string `json:"DstIP,omitempty"`
	SrcPort   string `json:"SrcPort,omitempty"`
	DstPort   string `json:"DstPort,omitempty"`
	Protocol  string `json:"Protocol,omitempty"`
	Length    int    `json:"Length,omitempty"`
	Info      string `json:"Info,omitempty"`
}

type CapturedPacket struct {
	Meta   PacketMeta           `json:"meta"`
	Parsed *tshark.ProtocolInfo `json:"parsed"`
}
