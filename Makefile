# Default version, can be overridden: make VERSION=1.2.0
VERSION ?= 0.1.0

# Go build variables; override as needed
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

# Main package (adjust if your main entry point differs)
MAIN_PKG = ./cmd/engine

# Binary name (you can customize if desired)
BINARY_NAME = argo

# Output directory for binaries
BIN_DIR = bin

# LDFLAGS to embed version info
LDFLAGS = -X main.version=$(VERSION)

.PHONY: all build clean install build-linux build-windows

all: build

build:
	@echo "Building $(BINARY_NAME) for $(GOOS)/$(GOARCH) version $(VERSION)"
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(VERSION)/$(BINARY_NAME) $(MAIN_PKG)

build-linux:
	@echo "Building $(BINARY_NAME) for linux/amd64 version $(VERSION)"
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(VERSION)/$(BINARY_NAME)-linux-amd64 $(MAIN_PKG)

build-windows:
	@echo "Building $(BINARY_NAME) for windows/amd64 version $(VERSION)"
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(VERSION)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PKG)

clean:
	@echo "Cleaning up..."
	rm -f $(BIN_DIR)/$(VERSION)/$(BINARY_NAME)*

install: build
	@echo "Installing $(BINARY_NAME) $(VERSION)"
	# Adjust this installation path as needed. Typically $GOPATH/bin or /usr/local/bin
	# If $GOBIN is set, `go install` uses it.
	cp $(BIN_DIR)/$(BINARY_NAME) $$GOPATH/bin/$(BINARY_NAME)
