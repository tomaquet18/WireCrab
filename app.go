package main

import (
	"context"
	"fmt"
	"wirecrab/internal/devices"
)

// App struct
type App struct {
	ctx           context.Context
	deviceService *devices.DeviceService
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		deviceService: devices.New(),
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
