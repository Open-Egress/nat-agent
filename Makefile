# Project Variables
BINARY_NAME=open-egress-agent
GO_FILES=$(shell find . -name "*.go")
VERSION?=0.1.0
BUILD_DIR=dist

# Build Flags
# -s -w: Strips debug information to reduce binary size
# -extldflags "-static": Ensures a completely static binary (no CGO dependencies)
LDFLAGS=-ldflags="-s -w -X main.Version=$(VERSION)"

.PHONY: all clean build-all build-amd64 build-arm64 help install-deps test test-coverage

all: help

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## clean: Remove build artifacts
clean:
	@echo "Cleaning build directory..."
	@rm -rf $(BUILD_DIR)

## build-amd64: Build for standard 64-bit Linux (Intel/AMD)
build-amd64:
	@echo "Building for Linux AMD64..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/agent

## build-arm64: Build for ARM 64-bit Linux (AWS Graviton / Oracle Ampere)
build-arm64:
	@echo "Building for Linux ARM64..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/agent

## build-all: Build for all supported architectures
build-all: build-amd64 build-arm64
	@echo "Build complete. Check the $(BUILD_DIR) directory."

## install-deps: Install Go dependencies
install-deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download

## test: Run unit tests with race detection
test:
	@echo "Running unit tests..."
	@go test -v -race ./...

## test-coverage: Run unit tests and generate coverage report
test-coverage:
	@echo "Running tests with coverage..."
	@mkdir -p $(BUILD_DIR)
	@go test -v -race -coverprofile=$(BUILD_DIR)/coverage.out ./...
	@go tool cover -html=$(BUILD_DIR)/coverage.out -o $(BUILD_DIR)/coverage.html
	@echo "Coverage report generated at $(BUILD_DIR)/coverage.html"
