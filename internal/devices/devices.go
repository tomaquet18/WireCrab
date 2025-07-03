package devices

import (
	"strings"

	"github.com/google/gopacket/pcap"
)

// DeviceService provides access to available network devices
type DeviceService struct{}

// New returns a new instance of DeviceService
func New() *DeviceService {
	return &DeviceService{}
}

// GetAvailableDevices returns a list of devices with their type identified
func (s *DeviceService) GetAvailableDevices() ([]Device, error) {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		return nil, err
	}

	var result []Device

	for _, dev := range devices {
		ifaceType := detectDeviceType(dev)
		result = append(result, Device{
			Name:        dev.Name,
			Description: dev.Description,
			Type:        ifaceType,
		})
	}

	return result, nil
}

// detectDeviceType attempts to classify the device type based on its name or description
func detectDeviceType(dev pcap.Interface) string {
	name := strings.ToLower(dev.Name)
	desc := strings.ToLower(dev.Description)

	switch {
	case strings.Contains(name, "wlan") || strings.Contains(name, "wi-fi") || strings.Contains(desc, "wi-fi") || strings.Contains(desc, "wireless"):
		return "wifi"
	case strings.Contains(name, "eth") || strings.Contains(desc, "ethernet"):
		return "ethernet"
	case strings.Contains(name, "bluetooth") || strings.Contains(desc, "bluetooth"):
		return "bluetooth"
	case strings.Contains(name, "loopback") || strings.Contains(desc, "loopback"):
		return "loopback"
	case strings.Contains(name, "vmnet") || strings.Contains(name, "docker") || strings.Contains(desc, "virtual"):
		return "virtual"
	default:
		return "unknown"
	}
}
