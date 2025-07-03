package capture

type PacketInfo struct {
	Timestamp string
	SrcIP     string
	DstIP     string
	SrcPort   string
	DstPort   string
	Protocol  string
	Length    int
	RawData   []byte
}
