package devices

// Interface represents a network interface with a type classification
type Device struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"` // e.g., "wifi", "ethernet", "bluetooth", "loopback", "virtual", "unknown"
}
