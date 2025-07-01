package dissect

// Detect runs all registered dissectors and returns the first match
func Detect(payload []byte, srcPort, dstPort string) *ProtocolInfo {
	for _, d := range dissectors {
		if d.Match(payload, srcPort, dstPort) {
			parsed, err := d.Parse(payload)
			if err != nil {
				continue
			}
			proto := &ProtocolInfo{
				Name:   d.Name(),
				Detail: parsed,
			}

			if nested, ok := parsed.(map[string]any); ok {
				if payload, ok := nested["payload"].([]byte); ok {
					proto.child = Detect(payload, srcPort, dstPort)
				}
				if child, ok := nested["child"].(*ProtocolInfo); ok && proto.child == nil {
					proto.child = child
				}
			}
			return proto
		}
	}
	return nil
}
