# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /build

# Install git for go mod download
RUN apk add --no-cache git

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" \
    -o mila_exporter ./cmd/mila_exporter

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary
COPY --from=builder /build/mila_exporter .

# Create non-root user
RUN adduser -D -g '' exporter

USER exporter

# Expose the metrics port
EXPOSE 9100

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -q --spider http://localhost:9100/metrics || exit 1

ENTRYPOINT ["/app/mila_exporter"]

# Default command can be overridden
CMD ["--help"]
