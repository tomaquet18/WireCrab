package dissect

// EnrichWithDissectors traverses the ProtocolInfo tree and enriches unknown or empty nodes
func EnrichWithDissectors(proto *ProtocolInfo, payload []byte, srcPort, dstPort string) {
	if proto == nil {
		return
	}

	if proto.Detail == nil || isUnknown(proto.Name) {
		if d := GetDissectorForLayer(proto.Name, payload, srcPort, dstPort); d != nil {
			parsed, err := d.Parse(payload)
			if err == nil {
				proto.Name = d.Name()
				proto.Detail = parsed
			}
		}
	}

	// Always recurse into child
	if proto.Child != nil {
		EnrichWithDissectors(proto.Child, payload, srcPort, dstPort)
	}
}

func isUnknown(name string) bool {
	return name == "data" || name == "unknown" || name == "raw"
}
