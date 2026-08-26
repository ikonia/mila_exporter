package collector

import (
	"context"
	"time"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"

	"github.com/mila-air-exporter/mila_exporter/internal/client"
)

const (
	namespace = "mila"
	subsystem = "air"
)

// MilaCollector implements the Prometheus Collector interface
type MilaCollector struct {
	client       client.MilaClient
	logger       log.Logger
	descChannels map[string]*prometheus.Desc

	// Metrics
	pm25Gauge           *prometheus.GaugeVec
	pm10Gauge           *prometheus.GaugeVec
	co2Gauge            *prometheus.GaugeVec
	temperatureGauge    *prometheus.GaugeVec
	humidityGauge       *prometheus.GaugeVec
	filterLifeGauge     *prometheus.GaugeVec
	fanSpeedGauge       *prometheus.GaugeVec
	powerStateGauge     *prometheus.GaugeVec
	wifiSignalGauge     *prometheus.GaugeVec
	errorCountGauge     *prometheus.GaugeVec
	warningCountGauge   *prometheus.GaugeVec
	upGauge             *prometheus.GaugeVec
	collectDurationGauge *prometheus.GaugeVec
}

// MilaClient defines the interface the collector needs from the API client
type MilaClient interface {
	GetDeviceData(ctx context.Context) (*client.MilaDeviceData, error)
	GetDevices(ctx context.Context) ([]client.MilaDevice, error)
}

// NewMilaCollector creates a new Prometheus collector for Mila Air data
func NewMilaCollector(client MilaClient, logger log.Logger) *MilaCollector {
	return &MilaCollector{
		client: client,
		logger: logger,
		descChannels: make(map[string]*prometheus.Desc),

		// Initialize metrics with appropriate help text and label names
		pm25Gauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "pm25_concentration",
				Help:      "PM2.5 particulate concentration in micrograms per cubic meter (μg/m³)",
			},
			[]string{"unit_id", "unit_name"},
		),
		pm10Gauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "pm10_concentration",
				Help:      "PM10 particulate concentration in micrograms per cubic meter (μg/m³)",
			},
			[]string{"unit_id", "unit_name"},
		),
		co2Gauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "co2_concentration",
				Help:      "Carbon dioxide concentration in parts per million (ppm)",
			},
			[]string{"unit_id", "unit_name"},
		),
		temperatureGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "temperature_celsius",
				Help:      "Temperature in degrees Celsius",
			},
			[]string{"unit_id", "unit_name"},
		),
		humidityGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "humidity_percent",
				Help:      "Relative humidity as a percentage (0-100)",
			},
			[]string{"unit_id", "unit_name"},
		),
		filterLifeGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "filter_life_remaining_percent",
				Help:      "Remaining filter life as a percentage (0-100)",
			},
			[]string{"unit_id", "unit_name"},
		),
		fanSpeedGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "fan_speed_rpm",
				Help:      "Current fan speed in RPM",
			},
			[]string{"unit_id", "unit_name"},
		),
		powerStateGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "power_state",
				Help:      "Power state: 1 = on, 0 = off",
			},
			[]string{"unit_id", "unit_name"},
		),
		wifiSignalGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "wifi_signal_dbm",
				Help:      "WiFi signal strength in dBm",
			},
			[]string{"unit_id", "unit_name"},
		),
		errorCountGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "error_count",
				Help:      "Number of active error conditions",
			},
			[]string{"unit_id", "unit_name"},
		),
		warningCountGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "warning_count",
				Help:      "Number of active warning conditions",
			},
			[]string{"unit_id", "unit_name"},
		),
		upGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "up",
				Help:      "Whether the Mila Air unit is reachable (1 = up, 0 = down)",
			},
			[]string{"unit_id"},
		),
		collectDurationGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "collector_duration_seconds",
				Help:      "Duration of the last metrics collection in seconds",
			},
			[]string{"unit_id"},
		),
	}
}

// Describe implements the Prometheus Collector interface
func (c *MilaCollector) Describe(ch chan<- *prometheus.Desc) {
	c.pm25Gauge.Describe(ch)
	c.pm10Gauge.Describe(ch)
	c.co2Gauge.Describe(ch)
	c.temperatureGauge.Describe(ch)
	c.humidityGauge.Describe(ch)
	c.filterLifeGauge.Describe(ch)
	c.fanSpeedGauge.Describe(ch)
	c.powerStateGauge.Describe(ch)
	c.wifiSignalGauge.Describe(ch)
	c.errorCountGauge.Describe(ch)
	c.warningCountGauge.Describe(ch)
	c.upGauge.Describe(ch)
	c.collectDurationGauge.Describe(ch)
}

// Collect implements the Prometheus Collector interface
func (c *MilaCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	startTime := time.Now()

	level.Debug(c.logger).Log("msg", "Starting metrics collection")

	// Fetch device data
	data, err := c.client.GetDeviceData(ctx)
	if err != nil {
		level.Error(c.logger).Log("msg", "Failed to fetch device data", "err", err)
		c.upGauge.WithLabelValues("unknown").Set(0)
		c.upGauge.Collect(ch)
		return
	}

	// Calculate collection duration
	duration := time.Since(startTime).Seconds()

	// Update up gauge
	c.upGauge.WithLabelValues(data.UnitID).Set(1)

	// Collect all metrics
	c.collectMetrics(data, ch)

	// Record collection duration
	c.collectDurationGauge.WithLabelValues(data.UnitID).Set(duration)

	level.Debug(c.logger).Log("msg", "Metrics collection complete", "duration_seconds", duration)
}

func (c *MilaCollector) collectMetrics(data *client.MilaDeviceData, ch chan<- prometheus.Metric) {
	unitID := data.UnitID
	unitName := data.UnitName

	// PM2.5 - typical air quality ranges: 0-12 μg/m³ good, 12-35.4 moderate, 35.4-55.4 unhealthy for sensitive
	c.pm25Gauge.WithLabelValues(unitID, unitName).Set(data.PM25)

	// PM10 - typical ranges: 0-54 μg/m³ good, 54-154 moderate
	c.pm10Gauge.WithLabelValues(unitID, unitName).Set(data.PM10)

	// CO2 - typical ranges: 350-450 ppm outdoor, 400-1000 ppm indoor acceptable, >1000 ppm unacceptable
	c.co2Gauge.WithLabelValues(unitID, unitName).Set(data.CO2)

	// Temperature in Celsius
	c.temperatureGauge.WithLabelValues(unitID, unitName).Set(data.Temperature)

	// Humidity percentage
	c.humidityGauge.WithLabelValues(unitID, unitName).Set(data.Humidity)

	// Filter life remaining
	c.filterLifeGauge.WithLabelValues(unitID, unitName).Set(data.FilterLife)

	// Fan speed in RPM
	c.fanSpeedGauge.WithLabelValues(unitID, unitName).Set(float64(data.FanSpeed))

	// Power state: 1 = on, 0 = off
	power := 0.0
	if data.PowerOn {
		power = 1.0
	}
	c.powerStateGauge.WithLabelValues(unitID, unitName).Set(power)

	// WiFi signal in dBm (typically -30 to -90 dBm)
	c.wifiSignalGauge.WithLabelValues(unitID, unitName).Set(float64(data.Connectivity.WiFiSignal))

	// Error and warning counts
	c.errorCountGauge.WithLabelValues(unitID, unitName).Set(float64(len(data.ErrorCodes)))
	c.warningCountGauge.WithLabelValues(unitID, unitName).Set(float64(len(data.Warnings)))

	// Expose mode as a label (for filtering/grouping)
	modeDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "mode"),
		"Current operating mode",
		[]string{"unit_id", "unit_name", "mode"},
		nil,
	)
	modeValue := 1.0
	ch <- prometheus.NewConstMetric(modeDesc, prometheus.GaugeValue, modeValue, unitID, unitName, data.Mode)

	// Expose fan speed level as a label
	fanLevelDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, subsystem, "fan_level"),
		"Fan level (1-4 approximate ranges)",
		[]string{"unit_id", "unit_name", "level"},
		nil,
	)
	fanLevel := "unknown"
	switch {
	case data.FanSpeed <= 200:
		fanLevel = "low"
	case data.FanSpeed <= 400:
		fanLevel = "medium"
	case data.FanSpeed <= 600:
		fanLevel = "high"
	default:
		fanLevel = "turbo"
	}
	ch <- prometheus.NewConstMetric(fanLevelDesc, prometheus.GaugeValue, 1.0, unitID, unitName, fanLevel)
}

// SetUnitID sets the unit ID for the collector (for dynamic unit switching)
func (c *MilaCollector) SetUnitID(unitID string) {
	// This can be used to update the collector's unit ID if needed
}

// SetUnitName sets the unit name for metric labeling
func (c *MilaCollector) SetUnitName(unitName string) {
	// This can be used to update the collector's unit name if needed
}
