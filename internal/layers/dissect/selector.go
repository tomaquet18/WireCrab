package dissect

// GetDissectorForLayer tries to find a matching dissector for the given layer type and payload
func GetDissectorForLayer(layer string, payload []byte, srcPort, dstPort string) Dissector {
	var list []Dissector
	switch layer {
	case "datalink":
		list = datalinkDissectors
	case "network":
		list = networkDissectors
	default:
		return nil
	}

	for _, d := range list {
		if d.Match(payload, srcPort, dstPort) {
			return d
		}
	}
	return nil
}
