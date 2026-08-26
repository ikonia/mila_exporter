package client

import (
	"context"
	"fmt"
	"time"
)

// MockMilaClient is a mock implementation for testing
type MockMilaClient struct {
	DeviceData *MilaDeviceData
	Devices    []MilaDevice
	Err        error
}

// NewMockMilaClient creates a mock client for testing
func NewMockMilaClient(data *MilaDeviceData, devices []MilaDevice, err error) *MockMilaClient {
	return &MockMilaClient{
		DeviceData: data,
		Devices:    devices,
		Err:        err,
	}
}

// GetDeviceData returns mock data
func (m *MockMilaClient) GetDeviceData(ctx context.Context) (*MilaDeviceData, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.DeviceData, nil
}

// GetDevices returns mock device list
func (m *MockMilaClient) GetDevices(ctx context.Context) ([]MilaDevice, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.Devices, nil
}

// Example mock data for testing
func ExampleMockData() *MilaDeviceData {
	return &MilaDeviceData{
		Timestamp:  time.Now(),
		PM25:       5.2,
		PM10:       8.1,
		CO2:        450.0,
		Temperature: 22.5,
		Humidity:    45.0,
		FilterLife:  85.0,
		FanSpeed:    2,
		Mode:        "auto",
		PowerOn:     true,
		ErrorCodes:  []string{},
		Warnings:    []string{},
		Connectivity: MilaConnectivity{
			WiFiSignal: -65,
			SSID:       "home-network",
			IP:         "192.168.1.100",
		},
	}
}

// ExampleDevices returns sample devices
func ExampleDevices() []MilaDevice {
	return []MilaDevice{
		{
			ID:           "unit-001",
			Name:         "Living Room",
			Type:         "air_purifier",
			SerialNumber: "MA12345678",
			Firmware:     "1.2.3",
		},
	}
}

// ExampleAuthResponse returns a sample auth response
func ExampleAuthResponse() *MilaAuthResponse {
	return &MilaAuthResponse{
		AccessToken: "test-token-12345",
		ExpiresIn:   3600,
		TokenType:   "Bearer",
	}
}

// ExampleAuthError returns an error for testing
func ExampleAuthError() error {
	return fmt.Errorf("authentication failed: invalid credentials")
}
