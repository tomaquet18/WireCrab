package dissect

var dissectors []Dissector

// Register adds a new dissector to the global list
func Register(d Dissector) {
	dissectors = append(dissectors, d)
}
