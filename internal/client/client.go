package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-kit/log"
	"github.com/prometheus/common/config"
)

const (
	defaultBaseURL    = "https://api.mila.ai"
	defaultTimeout    = 30 * time.Second
	defaultUserAgent  = "mila_exporter"
	apiPathLogin      = "/v1/auth/login"
	apiPathDevices    = "/v1/devices"
	apiPathDeviceData = "/v1/devices/%s/data"
)

// MilaClient represents the API client for Mila Air
type MilaClient struct {
	baseURL    string
	email      string
	password   config.Secret
	unitID     string
	httpClient *http.Client
	logger     log.Logger
	authToken  string
}

// MilaAuthResponse represents the authentication response
type MilaAuthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// MilaDevice represents a Mila device
type MilaDevice struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	SerialNumber string `json:"serial_number"`
	Firmware     string `json:"firmware_version"`
}

// MilaDeviceData represents the current data/state of a Mila device
type MilaDeviceData struct {
	Timestamp    time.Time          `json:"timestamp"`
	PM25         float64            `json:"pm25"`
	PM10         float64            `json:"pm10"`
	CO2          float64            `json:"co2"`
	Temperature  float64            `json:"temperature"`
	Humidity     float64            `json:"humidity"`
	FilterLife   float64            `json:"filter_life_percent"`
	FanSpeed     int                `json:"fan_speed"`
	Mode         string             `json:"mode"`
	PowerOn      bool               `json:"power_on"`
	ErrorCodes   []string           `json:"error_codes"`
	Warnings     []string           `json:"warnings"`
	Connectivity MilaConnectivity   `json:"connectivity"`
	Operations   []MilaOperation    `json:"operations"`
	UnitID       string             `json:"unit_id,omitempty"`
	UnitName     string             `json:"unit_name,omitempty"`
}

// MilaConnectivity represents device connectivity info
type MilaConnectivity struct {
	WiFiSignal int    `json:"wifi_signal_dbm"`
	SSID       string `json:"ssid"`
	IP         string `json:"ip_address"`
}

// MilaOperation represents a single operation record
type MilaOperation struct {
	Timestamp time.Time `json:"timestamp"`
	Event     string    `json:"event"`
	Value     string    `json:"value"`
}

// NewMilaClient creates a new Mila API client
func NewMilaClient(email string, password config.Secret, unitID string, logger log.Logger) (*MilaClient, error) {
	if email == "" {
		return nil, fmt.Errorf("email is required")
	}
	if password == "" {
		return nil, fmt.Errorf("password is required")
	}
	if unitID == "" {
		return nil, fmt.Errorf("unit_id is required")
	}

	client := &MilaClient{
		baseURL:    defaultBaseURL,
		email:      email,
		password:   password,
		unitID:     unitID,
		httpClient: &http.Client{Timeout: defaultTimeout},
		logger:     logger,
	}

	// Authenticate on startup
	if err := client.authenticate(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to authenticate: %w", err)
	}

	return client, nil
}

// authenticate logs into the Mila API and stores the auth token
func (c *MilaClient) authenticate(ctx context.Context) error {
	level.Debug(c.logger).Log("msg", "Authenticating with Mila API")

	authReq := map[string]string{
		"email":    c.email,
		"password": string(c.password),
	}

	body, err := json.Marshal(authReq)
	if err != nil {
		return fmt.Errorf("failed to marshal auth request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+apiPathLogin, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", defaultUserAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("authentication failed: %s", string(body))
	}

	var authResp MilaAuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return fmt.Errorf("failed to parse auth response: %w", err)
	}

	c.authToken = authResp.AccessToken
	level.Info(c.logger).Log("msg", "Authenticated successfully")

	return nil
}

// RefreshToken refreshes the authentication token if expired
func (c *MilaClient) refreshAuthToken(ctx context.Context) error {
	return c.authenticate(ctx)
}

// GetDeviceData fetches current data for the configured unit
func (c *MilaClient) GetDeviceData(ctx context.Context) (*MilaDeviceData, error) {
	level.Debug(c.logger).Log("msg", "Fetching device data", "unit_id", c.unitID)

	// Construct the API URL with unit ID
	url := fmt.Sprintf("%s%s", c.baseURL, fmt.Sprintf(apiPathDeviceData, c.unitID))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.authToken)
	req.Header.Set("User-Agent", defaultUserAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch device data: %w", err)
	}
	defer resp.Body.Close()

	// Handle 401 Unauthorized - token may have expired
	if resp.StatusCode == http.StatusUnauthorized {
		level.Info(c.logger).Log("msg", "Token expired, refreshing...")
		if err := c.refreshAuthToken(ctx); err != nil {
			return nil, fmt.Errorf("failed to refresh token: %w", err)
		}
		// Retry with new token
		req.Header.Set("Authorization", "Bearer "+c.authToken)
		resp, err = c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to retry fetch: %w", err)
		}
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch device data: %s", string(body))
	}

	var data MilaDeviceData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to parse device data: %w", err)
	}

	return &data, nil
}

// GetDevices fetches the list of devices associated with the account
func (c *MilaClient) GetDevices(ctx context.Context) ([]MilaDevice, error) {
	level.Debug(c.logger).Log("msg", "Fetching device list")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+apiPathDevices, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.authToken)
	req.Header.Set("User-Agent", defaultUserAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch devices: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch devices: %s", string(body))
	}

	var devices []MilaDevice
	if err := json.NewDecoder(resp.Body).Decode(&devices); err != nil {
		return nil, fmt.Errorf("failed to parse devices: %w", err)
	}

	return devices, nil
}
