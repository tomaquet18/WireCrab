package dissect

// Dissector defines the interface that each protocol dissector must implement
type Dissector interface {
	Name() string                                       // Unique name of the protocol (e.g. "http", "tcp")
	Match(Payload []byte, srcPort, dstPort string) bool // Lightweight match function to identify protocol
	Parse(payload []byte) (any, error)                  // Main parser for protocol data
}

// ProtocolInfo holds the output of a parsed protocol and potentially nested child protocols
type ProtocolInfo struct {
	Name   string      // Protocol name (e.g, "http")
	Detail interface{} // Map with extracted fields
	child  *ProtocolInfo
}
