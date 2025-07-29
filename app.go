package main

import (
	"context"
	"fmt"
	"wirecrab/internal/capture"
	"wirecrab/internal/devices"
	"wirecrab/internal/tshark"
	"wirecrab/internal/types"
)

// App struct
type App struct {
	ctx            context.Context
	deviceService  *devices.DeviceService
	captureService *capture.CaptureService
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		deviceService:  devices.New(),
		captureService: capture.NewCaptureService(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// GetInterfaces returns a list of all interfaces
func (a *App) GetDevices() ([]devices.Device, error) {
	return a.deviceService.GetAvailableDevices()
}

// StartCapture starts sniffing packets
func (a *App) StartCapture(device string) {
	a.captureService.Start(device)
}

// GetCapturedPackets returns a list of captured packets
func (a *App) GetCapturedPackets(offset int, limit int) []types.CapturedPacket {
	if a.captureService == nil {
		return nil
	}
	return a.captureService.GetPackets(offset, limit)
}

func (a *App) GetPacketCount() int {
	if a.captureService == nil {
		return 0
	}
	return a.captureService.GetPacketCount()
}

// ClearCapturedPackets clears the list of captured packets
func (a *App) ClearCapturedPackets() {
	a.captureService.Clear()
}

func (a *App) GetPacketDetails(packetNumber int) (*tshark.PacketDetails, error) {
	if a.captureService == nil {
		return nil, fmt.Errorf("capture service not started")
	}
	return a.captureService.GetPacketDetails(packetNumber)
}
