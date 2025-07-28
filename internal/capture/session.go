package capture

import (
	"wirecrab/internal/tshark"
)

type CaptureSession struct {
	Packets []CapturedPacket
}

type CapturedPacket struct {
	Meta   PacketMeta           `json:"meta"`
	Parsed *tshark.ProtocolInfo `json:"parsed"`
}
