package capture

import "wirecrab/internal/layers/dissect"

type CaptureSession struct {
	Packets []CapturedPacket
}

type CapturedPacket struct {
	Meta   PacketInfo            `json:"meta"`
	Parsed *dissect.ProtocolInfo `json:"parsed"`
}
