package capture

import (
	"time"
	"wirecrab/internal/layers/dissect"
)

var defaultTimeout = 2 * time.Second

type CaptureService struct {
	session *CaptureSession
}

func NewCaptureService() *CaptureService {
	return &CaptureService{}
}

func (s *CaptureService) Start(device string) {
	cfg := SnifferConfig{
		Device:      device,
		SnapshotLen: 65536,
		Promiscuous: true,
		Timeout:     defaultTimeout,
		PacketLimit: 0,
	}

	packetChan, _ := StartSniffing(cfg)

	s.session = &CaptureSession{}

	go func() {
		for pkt := range packetChan {
			parsed, _ := dissect.ParseWithTshark(pkt.RawData)

			s.session.Packets = append(s.session.Packets, CapturedPacket{
				Meta:   pkt,
				Parsed: parsed,
			})
		}
	}()
}

func (s *CaptureService) GetPackets() []CapturedPacket {
	if s.session == nil {
		return nil
	}
	return s.session.Packets
}

func (s *CaptureService) Clear() {
	s.session = nil
}
