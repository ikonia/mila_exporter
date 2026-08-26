# Mila Air Prometheus Exporter

[![License](https://img.shields.io/badge/license-Apache_2.0-blue.svg)](https://github.com/mila-air-exporter/mila_exporter/blob/main/LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.23+-orange.svg)](https://golang.org/)
[![Prometheus](https://img.shields.io/badge/Prometheus-2.x-green.svg)](https://prometheus.io/)

A Prometheus exporter for [Mila Air](https://mila.ai/) air purifier units, written in Go following [Prometheus exporter best practices](https://prometheus.io/docs/instrumenting/writing_exporters/).

This exporter connects to the Mila Cloud API to collect real-time air quality data and exposes it as Prometheus metrics.

## Features

- **PM2.5 & PM10** - Particulate matter concentration (μg/m³)
- **CO2** - Carbon dioxide levels (ppm)
- **Temperature** - Room temperature (°C)
- **Humidity** - Relative humidity (%)
- **Filter Life** - Remaining filter life (%)
- **Fan Speed** - Current fan speed (RPM)
- **Power State** - Device on/off status
- **WiFi Signal** - WiFi signal strength (dBm)
- **Error/Warning Count** - Device status alerts

## Building

```bash
cd mila_exporter
go build -o mila_exporter ./cmd/mila_exporter
```

## Running

### Command-line Flags

```bash
./mila_exporter \
  --mila.email="your-email@example.com" \
  --mila.password="your-password" \
  --mila.unit_id="your-unit-id" \
  --web.listen-address=":9100"
```

### Configuration File

Create a `config.yaml`:

```yaml
mila:
  email: "your-email@example.com"
  password: "your-password"
  unit_id: "your-unit-id"

server:
  listen_address: ":9100"
```

Then run:

```bash
./mila_exporter --config.file=config.yaml
```

### Docker

First, build the Docker image:

```bash
docker build -t mila-exporter:latest .
```

Then run:

```bash
docker run -d \
  -p 9100:9100 \
  -e MILA_EMAIL="your-email@example.com" \
  -e MILA_PASSWORD="your-password" \
  -e MILA_UNIT_ID="your-unit-id" \
  mila-exporter:latest
```

## Prometheus Configuration

Add this to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'mila-air'
    static_configs:
      - targets: ['localhost:9100']
    metrics_path: /metrics
```

## Example Metrics Output

When scraping the exporter, you'll see output like:

```
# HELP mila_air_collector_duration_seconds Duration of the last metrics collection in seconds
# TYPE mila_air_collector_duration_seconds gauge
mila_air_collector_duration_seconds{unit_id="unit-001"} 0.245
# HELP mila_air_co2_concentration Carbon dioxide concentration in parts per million (ppm)
# TYPE mila_air_co2_concentration gauge
mila_air_co2_concentration{unit_id="unit-001",unit_name="Living Room"} 450
# HELP mila_air_error_count Number of active error conditions
# TYPE mila_air_error_count gauge
mila_air_error_count{unit_id="unit-001",unit_name="Living Room"} 0
# HELP mila_air_filter_life_remaining_remaining_percent Remaining filter life as a percentage (0-100)
# TYPE mila_air_filter_life_remaining_percent gauge
mila_air_filter_life_remaining_percent{unit_id="unit-001",unit_name="Living Room"} 85
# HELP mila_air_fan_speed_rpm Current fan speed in RPM
# TYPE mila_air_fan_speed_rpm gauge
mila_air_fan_speed_rpm{unit_id="unit-001",unit_name="Living Room"} 350
# HELP mila_air_humidity_percent Relative humidity as a percentage (0-100)
# TYPE mila_air_humidity_percent gauge
mila_air_humidity_percent{unit_id="unit-001",unit_name="Living Room"} 45
# HELP mila_air_mode Current operating mode
# TYPE mila_air_mode gauge
mila_air_mode{unit_id="unit-001",unit_name="Living Room",mode="auto"} 1
# HELP mila_air_pm10_concentration PM10 particulate concentration in micrograms per cubic meter (g/m³)
# TYPE mila_air_pm10_concentration gauge
mila_air_pm10_concentration{unit_id="unit-001",unit_name="Living Room"} 8.1
# HELP mila_air_pm25_concentration PM2.5 particulate concentration in micrograms per cubic meter (g/m³)
# TYPE mila_air_pm25_concentration gauge
mila_air_pm25_concentration{unit_id="unit-001",unit_name="Living Room"} 5.2
# HELP mila_air_power_state Power state: 1 = on, 0 = off
# TYPE mila_air_power_state gauge
mila_air_power_state{unit_id="unit-001",unit_name="Living Room"} 1
# HELP mila_air_temperature_celsius Temperature in degrees Celsius
# TYPE mila_air_temperature_celsius gauge
mila_air_temperature_celsius{unit_id="unit-001",unit_name="Living Room"} 22.5
# HELP mila_air_up Whether the Mila Air unit is reachable (1 = up, 0 = down)
# TYPE mila_air_up gauge
mila_air_up{unit_id="unit-001"} 1
# HELP mila_air_warning_count Number of active warning conditions
# TYPE mila_air_warning_count gauge
mila_air_warning_count{unit_id="unit-001",unit_name="Living Room"} 0
# HELP mila_air_wifi_signal_dbm WiFi signal strength in dBm
# TYPE mila_air_wifi_signal_dbm gauge
mila_air_wifi_signal_dbm{unit_id="unit-001",unit_name="Living Room"} -65
```

## Metrics Reference

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `mila_air_pm25_concentration` | Gauge | unit_id, unit_name | PM2.5 concentration (μg/m³) |
| `mila_air_pm10_concentration` | Gauge | unit_id, unit_name | PM10 concentration (μg/m³) |
| `mila_air_co2_concentration` | Gauge | unit_id, unit_name | CO2 concentration (ppm) |
| `mila_air_temperature_celsius` | Gauge | unit_id, unit_name | Temperature (°C) |
| `mila_air_humidity_percent` | Gauge | unit_id, unit_name | Humidity (%) |
| `mila_air_filter_life_remaining_percent` | Gauge | unit_id, unit_name | Filter life remaining (%) |
| `mila_air_fan_speed_rpm` | Gauge | unit_id, unit_name | Fan speed (RPM) |
| `mila_air_power_state` | Gauge | unit_id, unit_name | Power state (1=on, 0=off) |
| `mila_air_wifi_signal_dbm` | Gauge | unit_id, unit_name | WiFi signal (dBm) |
| `mila_air_error_count` | Gauge | unit_id, unit_name | Active error count |
| `mila_air_warning_count` | Gauge | unit_id, unit_name | Active warning count |
| `mila_air_mode` | Gauge | unit_id, unit_name, mode | Current operating mode |
| `mila_air_fan_level` | Gauge | unit_id, unit_name, level | Fan level (low/medium/high/turbo) |
| `mila_air_up` | Gauge | unit_id | Up status (1=up, 0=down) |
| `mila_air_collector_duration_seconds` | Gauge | unit_id | Collection duration (seconds) |

## API

The Mila API client is implemented in `internal/client/client.go`. It handles:

- Authentication with email/password
- Automatic token refresh
- Device data fetching
- Error handling with retries

## Project Structure

```
mila_exporter/
├── cmd/
│   └── mila_exporter/
│       └── main.go           # Main entry point
├── internal/
│   ├── client/               # Mila API client
│   │   ├── client.go
│   │   └── client_test_helpers.go
│   └── collector/            # Prometheus collector
│       └── collector.go
├── pkg/
│   └── flags/                # Command-line flags
│       └── flags.go
├── config.example.yaml       # Configuration example
├── web-config.example.yaml   # Web server config example
├── go.mod                    # Go module definition
└── README.md
```

## Best Practices Followed

This exporter follows [Prometheus exporter best practices](https://prometheus.io/docs/instrumenting/writing_exporters/):

- [x] Uses Go with proper error handling
- [x] Implements the Prometheus Collector interface
- [x] Provides clear metric names and help text
- [x] Uses appropriate metric types (Gauge, Counter, etc.)
- [x] Includes unit IDs as labels for multi-unit support
- [x] Exposes an `/up` metric for availability monitoring
- [x] Includes collection duration metrics
- [x] Uses proper HTTP status codes and error responses
- [x] Implements TLS support via exporter-toolkit
- [x] Provides comprehensive documentation

## License

Apache 2.0 - see [LICENSE](LICENSE) for details.

## Acknowledgments

- Built following the patterns from [node_exporter](https://github.com/prometheus/node_exporter) and [mysqld_exporter](https://github.com/prometheus/mysqld_exporter)
- Uses the [Prometheus Go client library](https://github.com/prometheus/client_golang)
- Uses [exporter-toolkit](https://github.com/prometheus/exporter-toolkit) for web server and TLS
