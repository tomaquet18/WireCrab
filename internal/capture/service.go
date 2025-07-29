package capture

import (
	"fmt"
	"strconv"
	"sync"
	"wirecrab/internal/tshark"
	"wirecrab/internal/types"
)

type CaptureService struct {
	tshark    *tshark.TsharkLive
	stopChan  chan struct{}
	waitGroup sync.WaitGroup
}

func NewCaptureService() *CaptureService {
	return &CaptureService{}
}

func (s *CaptureService) Start(device string) error {
	if s.tshark != nil {
		s.Stop() // Stop previous capture
	}

	tshark, err := tshark.StartTsharkLive(device)
	if err != nil {
		return err
	}
	s.tshark = tshark

	return nil
}

func (s *CaptureService) Stop() {
	if s.tshark != nil {
		_ = s.tshark.Close()
	}
	s.tshark = nil
}

func (s *CaptureService) GetPackets(offset, limit int) []types.CapturedPacket {
	if s.tshark == nil {
		return nil
	}

	protos, err := s.tshark.GetPacketList(offset, limit)
	if err != nil {
		return nil
	}

	packets := make([]types.CapturedPacket, 0, len(protos))
	for i := range protos {
		packets = append(packets, types.CapturedPacket{
			Meta:   extractMetaFromParsed(&protos[i]),
			Parsed: &protos[i],
		})
	}
	return packets
}

func (s *CaptureService) GetPacketCount() int {
	if s.tshark == nil {
		return 0
	}

	count, err := s.tshark.GetPacketCount()
	if err != nil {
		return 0
	}
	return count
}

func (s *CaptureService) GetPacketDetails(packetNumber int) (*tshark.PacketDetails, error) {
	if s.tshark == nil {
		return nil, fmt.Errorf("tshark not running")
	}
	return s.tshark.GetPacketDetails(packetNumber)
}

// ---------------------- helpers ----------------------

func extractMetaFromParsed(parsed *tshark.ProtocolInfo) types.PacketMeta {
	meta := types.PacketMeta{}

	// Helper function to get values from detail map
	if detailMap, ok := parsed.Detail.(map[string]any); ok {
		get := func(key string) string {
			if v, ok := detailMap[key]; ok {
				if m, ok := v.(map[string]any); ok {
					if val, ok := m["value"].(string); ok {
						return val
					}
				}
			}
			return ""
		}

		meta.Timestamp = get("timestamp")
		meta.SrcIP = get("ip.src")
		meta.DstIP = get("ip.dst")
		meta.Protocol = parsed.Name
		meta.Length = atoi(get("frame.len"))
		meta.Info = get("_ws.col.info")
	}

	return meta
}

func atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
