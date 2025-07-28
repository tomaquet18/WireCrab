package storage

import (
	"sync"
	"wirecrab/internal/types"
)

type PacketStore interface {
	Push(packet types.CapturedPacket)
	GetRange(offset, limit int) []types.CapturedPacket
	Count() int
	Clear()
}

type MemoryPacketStore struct {
	packets []types.CapturedPacket
	mutex   sync.RWMutex
	maxSize int
}

func NewMemoryPacketStore(maxSize int) *MemoryPacketStore {
	return &MemoryPacketStore{
		packets: make([]types.CapturedPacket, 0, 1000),
		maxSize: maxSize,
	}
}

func (s *MemoryPacketStore) Push(packet types.CapturedPacket) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if len(s.packets) >= s.maxSize {
		// Remove oldest packet when reaching max size
		s.packets = s.packets[1:]
	}
	s.packets = append(s.packets, packet)
}

func (s *MemoryPacketStore) GetRange(offset, limit int) []types.CapturedPacket {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if offset >= len(s.packets) {
		return []types.CapturedPacket{}
	}

	end := offset + limit
	if end > len(s.packets) {
		end = len(s.packets)
	}

	result := make([]types.CapturedPacket, end-offset)
	copy(result, s.packets[offset:end])
	return result
}

func (s *MemoryPacketStore) Count() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.packets)
}

func (s *MemoryPacketStore) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.packets = s.packets[:0]
}
