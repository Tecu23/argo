# Build variables
VERSION ?= 0.7.0
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
CGO_ENABLED ?= 0

# Project structure
MAIN_PKG := ./cmd/argo
BINARY_NAME := argo
BIN_DIR := bin
DIST_DIR := dist

# Installation paths
ifeq ($(GOOS),windows)
    INSTALL_PATH := $(GOPATH)/bin/$(BINARY_NAME).exe
else
    INSTALL_PATH := $(GOPATH)/bin/$(BINARY_NAME)
endif

# Build tags & flags
BUILD_TAGS := 
LDFLAGS := -w -s \
    -X 'main.version=$(VERSION)' \
    -X 'main.buildDate=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")' \
    -X 'main.gitCommit=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")'

# Test variables
TEST_TIMEOUT := 10m
COVERAGE_DIR := coverage
COVERAGE_FILE := $(COVERAGE_DIR)/coverage.out
COVERAGE_HTML := $(COVERAGE_DIR)/coverage.html

# Docker variables
DOCKER_IMAGE := argo-chess
DOCKER_TAG := $(VERSION)

.PHONY: all build build-all clean install test coverage lint docker-build help

# Default target
all: clean build test

# Help target
help:
	@echo "Available targets:"
	@echo "  build       - Build for current OS/ARCH"
	@echo "  build-all   - Build for all supported platforms"
	@echo "  clean       - Remove built binaries and artifacts"
	@echo "  install     - Install binary to GOPATH/bin"
	@echo "  test        - Run tests"
	@echo "  coverage    - Generate test coverage report"
	@echo "  lint        - Run linters"
	@echo "  docker-build- Build Docker image"
	@echo "  help        - Show this help message"

# Build targets
build:
	@echo "Building $(BINARY_NAME) $(VERSION) for $(GOOS)/$(GOARCH)"
	@mkdir -p $(BIN_DIR)/$(VERSION)
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
		-tags "$(BUILD_TAGS)" \
		-ldflags "$(LDFLAGS)" \
		-o $(BIN_DIR)/$(VERSION)/$(BINARY_NAME)_$(GOOS)_$(GOARCH)$(if $(filter windows,$(GOOS)),.exe,) \
		$(MAIN_PKG)

build-all: build-linux build-windows build-darwin

build-linux:
	@echo "Building for Linux..."
	@GOOS=linux GOARCH=amd64 $(MAKE) build
	@GOOS=linux GOARCH=arm64 $(MAKE) build

build-windows:
	@echo "Building for Windows..."
	@GOOS=windows GOARCH=amd64 $(MAKE) build

build-darwin:
	@echo "Building for macOS..."
	@GOOS=darwin GOARCH=amd64 $(MAKE) build
	@GOOS=darwin GOARCH=arm64 $(MAKE) build

# Release target
release: build-all
	@echo "Creating release archives..."
	@mkdir -p $(DIST_DIR)
	@cd $(BIN_DIR)/$(VERSION) && find . -type f -name "$(BINARY_NAME)*" -exec zip -9 ../../$(DIST_DIR)/{}.zip {} \;

# Test targets
test:
	@echo "Running tests..."
	@go test -race -timeout $(TEST_TIMEOUT) ./...

coverage:
	@echo "Generating coverage report..."
	@mkdir -p $(COVERAGE_DIR)
	@go test -race -coverprofile=$(COVERAGE_FILE) -covermode=atomic ./...
	@go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Coverage report generated at $(COVERAGE_HTML)"

# Lint target
lint:
	@echo "Running linters..."
	@if command -v golangci-lint >/dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed"; \
		exit 1; \
	fi

# Docker target
docker-build:
	@echo "Building Docker image $(DOCKER_IMAGE):$(DOCKER_TAG)..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ") \
		--build-arg GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown") \
		.

# Clean target
clean:
	@echo "Cleaning..."
	@rm -rf $(BIN_DIR)
	@rm -rf $(DIST_DIR)
	@rm -rf $(COVERAGE_DIR)

# Install target
install: build
	@echo "Installing to $(INSTALL_PATH)"
	@mkdir -p $(dir $(INSTALL_PATH))
	@cp $(BIN_DIR)/$(VERSION)/$(BINARY_NAME)$(if $(filter windows,$(GOOS)),.exe,) $(INSTALL_PATH)

# Quick test target for development
quick-test:
	@go test -short ./...
