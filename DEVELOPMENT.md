# Mila Air Exporter - Development Guide

## Project Overview

This is a Prometheus exporter for Mila Air purifier units. It connects to the Mila Cloud API and exposes air quality metrics for monitoring with Prometheus.

## Quick Start

### Prerequisites

- Go 1.23 or later
- A Mila Air account with email/password
- The unit ID of your Mila Air purifier

### Getting Your Unit ID

1. Log into the Mila app or website
2. The unit ID can be found in the device settings or by using the developer tools in your browser to inspect API calls
3. Alternatively, run the exporter without `unit_id` to list available devices (future enhancement)

### Build and Run

```bash
# Clone the repository
git clone https://github.com/mila-air-exporter/mila_exporter.git
cd mila_exporter/mila_exporter

# Build
go build -o mila_exporter ./cmd/mila_exporter

# Run
./mila_exporter \
  --mila.email="your-email@example.com" \
  --mila.password="your-password" \
  --mila.unit_id="your-unit-id" \
  --web.listen-address=":9100"
```

### Verify

Open http://localhost:9100/metrics in your browser to see the exported metrics.

## Adding New Metrics

1. Add a new field to `MilaDeviceData` struct in `internal/client/client.go`
2. Add a new metric in `internal/collector/collector.go`
3. Update the `collectMetrics` function to expose the new metric
4. Add documentation in `README.md`

## Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Test a specific package
go test ./internal/client/
```

## Debugging

Enable debug logging:

```bash
./mila_exporter --log.level=debug
```

## Project Structure

```
internal/
├── client/          # Mila API client
│   ├── client.go    # Main client implementation
│   └── test_helpers.go  # Test utilities
└── collector/       # Prometheus collector
    └── collector.go # Metrics collection logic

pkg/
└── flags/           # Command-line flag definitions
    └── flags.go
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests
5. Submit a pull request

## Troubleshooting

### Authentication Failed

- Verify your email and password are correct
- Check that your account has access to the Mila API

### Unit ID Not Found

- Verify the unit ID is correct
- Check that the unit is associated with your Mila account

### Connection Errors

- Check your internet connection
- Verify the Mila API is accessible (https://api.mila.ai)

## License

Apache 2.0
