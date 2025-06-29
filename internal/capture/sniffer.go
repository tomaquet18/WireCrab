package capture

import (
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

type SnifferConfig struct {
	Device      string // if empty, run on all
	SnapshotLen int32
	Promiscuous bool
	Timeout     time.Duration
	PacketLimit int // for testing, 0 = infinite
}

// StartSniffing begins capturing packets on the given interface
func StartSniffing(cfg SnifferConfig) (<-chan PacketInfo, <-chan error) {
	out := make(chan PacketInfo)
	errs := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errs)

		var err error
		if cfg.Device == "" {
			err = sniffAllInterfacesStream(cfg, out)
		} else {
			err = sniffInterfaceStream(cfg, out)
		}

		if err != nil {
			errs <- err
		}
	}()

	return out, errs
}

func sniffInterfaceStream(cfg SnifferConfig, out chan<- PacketInfo) error {
	handle, err := pcap.OpenLive(cfg.Device, cfg.SnapshotLen, cfg.Promiscuous, cfg.Timeout)
	if err != nil {
		return fmt.Errorf("pcap open error: %w", err)
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	count := 0

	for packet := range packetSource.Packets() {
		info := PacketInfo{
			Timestamp: packet.Metadata().Timestamp,
			Length:    packet.Metadata().Length,
			RawData:   packet.Data(),
		}

		if netLayer := packet.NetworkLayer(); netLayer != nil {
			src, dst := netLayer.NetworkFlow().Endpoints()
			info.SrcIP = src.String()
			info.DstIP = dst.String()
			info.Protocol = netLayer.LayerType().String()
		}

		if transLayer := packet.TransportLayer(); transLayer != nil {
			srcPort, dstPort := transLayer.TransportFlow().Endpoints()
			info.SrcPort = srcPort.String()
			info.DstPort = dstPort.String()
			info.Protocol = transLayer.LayerType().String()
		}

		out <- info

		count++
		if cfg.PacketLimit > 0 && count >= cfg.PacketLimit {
			break
		}
	}
	return nil
}

func sniffAllInterfacesStream(cfg SnifferConfig, out chan<- PacketInfo) error {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		return fmt.Errorf("error listing interfaces: %w", err)
	}
	if len(devices) == 0 {
		return fmt.Errorf("no interfaces found")
	}

	for _, dev := range devices {
		go func(deviceName string) {
			err := sniffInterfaceStream(SnifferConfig{
				Device:      dev.Name,
				SnapshotLen: cfg.SnapshotLen,
				Promiscuous: cfg.Promiscuous,
				Timeout:     cfg.Timeout,
				PacketLimit: cfg.PacketLimit,
			}, out)
			if err != nil {
				log.Printf("Sniff error on %s: %v", dev.Name, err)
			}
		}(dev.Name)
	}

	return nil
}
