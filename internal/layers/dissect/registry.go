package dissect

// layer 2
var datalinkDissectors []Dissector

// layer 3
var networkDissectors []Dissector

// layer 4
var transportDissectors []Dissector

// layer 7
var applicationDissectors []Dissector

// -------------------
// Registry functions
// -------------------

func RegisterDataLink(d Dissector) {
	datalinkDissectors = append(datalinkDissectors, d)
}

func RegisterNetwork(d Dissector) {
	networkDissectors = append(networkDissectors, d)
}

func RegisterTransport(d Dissector) {
	transportDissectors = append(transportDissectors, d)
}

func RegisterApplication(d Dissector) {
	applicationDissectors = append(applicationDissectors, d)
}
