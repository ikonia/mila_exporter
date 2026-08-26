# Makefile for Mila Air Exporter

.PHONY: all build test clean vet format

BINARY_NAME=mila_exporter
VERSION?=$(shell git describe --tags --always 2>/dev/null || echo "development")
ifdef BUILD_DATE
  BUILD_DATE?=$(BUILD_DATE)
else
  BUILD_DATE:=$(shell date -u +%Y-%m-%dT%H:%M:%SZ 2>/dev/null || date -u +%Y-%m-%dT%H:%M:%S 2>/dev/null || echo "unknown")
endif

all: build

build:
	go build -ldflags "-X main.version=$(VERSION) -X main.buildDate=$(BUILD_DATE)" \
		-o $(BINARY_NAME) ./cmd/mila_exporter

test:
	go test -v ./...

vet:
	go vet ./...

format:
	go fmt ./...

clean:
	rm -f $(BINARY_NAME)

docker:
	docker build -t mila-exporter:latest .
