package datalink

import (
	"fmt"
	"wirecrab/internal/layers/dissect"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type ethernetDissector struct{}

var _ dissect.Dissector = ethernetDissector{}

func (ethernetDissector) Name() string {
	return "ethernet"
}

func (ethernetDissector) Match(payload []byte, _, _ string) bool {
	// Ethernet frame must be at least 14 bytes (MACs + EtherType)
	return len(payload) >= 14
}

func (ethernetDissector) Parse(payload []byte) (any, error) {
	packet := gopacket.NewPacket(payload, layers.LayerTypeEthernet, gopacket.NoCopy)
	layer := packet.Layer(layers.LayerTypeEthernet)
	if layer == nil {
		return nil, fmt.Errorf("ethernet: failed to parse")
	}

	eth, ok := layer.(*layers.Ethernet)
	if !ok {
		return nil, fmt.Errorf("ethernet: layer cast failed")
	}

	child := dissect.Detect(eth.Payload, "", "")

	return map[string]any{
		"src_mac":  eth.SrcMAC.String(),
		"dst_mac":  eth.DstMAC.String(),
		"eth_type": eth.EthernetType.String(),
		"payload":  eth.Payload,
		"child":    child,
	}, nil
}

func init() {
	dissect.Register(ethernetDissector{})
}
