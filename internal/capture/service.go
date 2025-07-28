package capture

import (
	"fmt"
	"strconv"
	"sync"
	"time"
	"wirecrab/internal/tshark"
)

var defaultTimeout = 2 * time.Second

type CaptureService struct {
	session   *CaptureSession
	tshark    *tshark.TsharkLive
	stopChan  chan struct{}
	waitGroup sync.WaitGroup
}

func NewCaptureService() *CaptureService {
	return &CaptureService{}
}

func (s *CaptureService) Start(device string) error {
	if s.session != nil {
		s.Stop() // Stop previous capture
	}

	s.session = &CaptureSession{
		Packets: make([]CapturedPacket, 0, 500), // Preallocate space to reduce allocations
	}
	s.stopChan = make(chan struct{})

	tshark, err := tshark.StartTsharkLive(device)
	if err != nil {
		return err
	}
	s.tshark = tshark

	s.waitGroup.Add(1)
	go func() {
		defer s.waitGroup.Done()
		for {
			select {
			case <-s.stopChan:
				return
			default:
				proto, err := tshark.Next()
				if err != nil {
					fmt.Println(err)
					continue
				}

				meta := extractMetaFromParsed(proto)

				s.session.Packets = append(s.session.Packets, CapturedPacket{
					Meta:   meta,
					Parsed: proto,
				})
			}
		}
	}()

	return nil
}

func (s *CaptureService) Stop() {
	if s.tshark != nil {
		_ = s.tshark.Close()
	}
	if s.stopChan != nil {
		close(s.stopChan)
	}
	s.waitGroup.Wait()
	s.session = nil
	s.tshark = nil
	s.stopChan = nil
}

func (s *CaptureService) GetPackets() []CapturedPacket {
	if s.session == nil {
		return nil
	}
	out := make([]CapturedPacket, len(s.session.Packets))
	copy(out, s.session.Packets)
	return out
}

func (s *CaptureService) Clear() {
	s.session = nil
}

// ---------------------- helpers ----------------------

func extractMetaFromParsed(parsed *tshark.ProtocolInfo) PacketMeta {
	meta := PacketMeta{
		Timestamp: time.Now().Format(time.RFC3339),
	}

	var lastProto string
	var walk func(p *tshark.ProtocolInfo)
	walk = func(p *tshark.ProtocolInfo) {
		if p == nil {
			return
		}

		lastProto = p.Name

		if detailMap, ok := p.Detail.(map[string]any); ok {
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

			switch p.Name {
			case "ip":
				meta.SrcIP = get("ip.src")
				meta.DstIP = get("ip.dst")
				if meta.Length == 0 {
					meta.Length = atoi(get("ip.len"))
				}
			case "tcp", "udp":
				meta.SrcPort = get(p.Name + ".srcport")
				meta.DstPort = get(p.Name + ".dstport")
			case "frame":
				if meta.Length == 0 {
					meta.Length = atoi(get("frame.len"))
				}
			}
		}

		walk(p.Child)
	}

	walk(parsed)

	meta.Protocol = lastProto
	return meta
}

func atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
