# Makefile for vic-go

.PHONY: build clean install test run help

# Variables
BINARY_NAME=vic
GO=go
GOFLAGS=-ldflags="-s -w"
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Default target
help:
	@echo "vic-go Makefile"
	@echo "================"
	@echo ""
	@echo "Targets:"
	@echo "  build       Build the vic binary (local)"
	@echo "  build-all  Build for all platforms"
	@echo "  clean      Clean build artifacts"
	@echo "  install    Install to \$$GOPATH/bin"
	@echo "  test       Run tests"
	@echo "  run        Build and run locally"
	@echo "  help       Show this help"

# Build for local platform
build:
	@echo "Building vic $(VERSION)..."
	$(GO) build $(GOFLAGS) -o $(BINARY_NAME) .
	@echo "Built: ./$(BINARY_NAME)"

# Build for all platforms
build-all:
	@echo "Building for all platforms..."
	@mkdir -p dist
	GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -o dist/vic-linux-amd64 .
	GOOS=linux GOARCH=arm64 $(GO) build $(GOFLAGS) -o dist/vic-linux-arm64 .
	GOOS=darwin GOARCH=amd64 $(GO) build $(GOFLAGS) -o dist/vic-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 $(GO) build $(GOFLAGS) -o dist/vic-darwin-arm64 .
	GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) -o dist/vic-windows-amd64.exe .
	@echo "Built: ./dist/"

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -rf dist

# Install to GOPATH/bin
install: build
	cp $(BINARY_NAME) $$GOPATH/bin/
	@echo "Installed to: \$$GOPATH/bin/$(BINARY_NAME)"

# Run tests
test:
	$(GO) test -v ./...

# Run locally
run: build
	./$(BINARY_NAME) $(ARGS)

# Cross-compile for Windows
build-win:
	GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) -o vic.exe .
	@echo "Built: ./vic.exe"
