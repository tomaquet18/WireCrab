package capture

import (
	"fmt"
	"strconv"
	"sync"
	"time"
	"wirecrab/internal/storage"
	"wirecrab/internal/tshark"
	"wirecrab/internal/types"
)

var defaultTimeout = 2 * time.Second

type CaptureService struct {
	store     storage.PacketStore
	tshark    *tshark.TsharkLive
	stopChan  chan struct{}
	waitGroup sync.WaitGroup
}

func NewCaptureService() *CaptureService {
	return &CaptureService{
		store: storage.NewMemoryPacketStore(1000000), // Store up to 1 million packets
	}
}

func (s *CaptureService) Start(device string) error {
	if s.tshark != nil {
		s.Stop() // Stop previous capture
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
				s.store.Push(types.CapturedPacket{
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
	s.tshark = nil
	s.stopChan = nil
}

func (s *CaptureService) GetPackets(offset, limit int) []types.CapturedPacket {
	return s.store.GetRange(offset, limit)
}

func (s *CaptureService) GetPacketCount() int {
	return s.store.Count()
}

func (s *CaptureService) Clear() {
	s.store.Clear()
}

// ---------------------- helpers ----------------------

func extractMetaFromParsed(parsed *tshark.ProtocolInfo) types.PacketMeta {
	meta := types.PacketMeta{
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
