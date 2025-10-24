# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOLINT=golangci-lint

# Binary name
BINARY_NAME=macinsight
BINARY_UNIX=$(BINARY_NAME)_unix

# Build flags
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev-$(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown')")
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

.PHONY: all build clean test deps fmt lint install

all: test build

# Build the binary
build:
	@mkdir -p bin
	@echo "Building version: $(VERSION)"
	$(GOBUILD) $(LDFLAGS) -o bin/$(BINARY_NAME) -v ./cmd/macinsight
# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -rf bin/

# Run tests
test:
	$(GOTEST) -v ./...

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Format code
fmt:
	$(GOFMT) -s -w .

# Run linter
lint:
	$(GOLINT) run

# Install the binary to GOPATH/bin
install:
	@echo "Installing version: $(VERSION)"
	$(GOCMD) install $(LDFLAGS) ./cmd/macinsight

# Install linter
install-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.54.2

# Run all checks
check: fmt lint test

# Generate JSON schema
schema:
	@echo "Generating JSON schema..."
	./bin/$(BINARY_NAME) schema --output schema/report.json

# Show help
help:
	@echo "Available targets:"
	@echo "  build        - Build the binary"
	@echo "  clean        - Clean build artifacts"
	@echo "  test         - Run tests"
	@echo "  deps         - Download dependencies"
	@echo "  fmt          - Format code"
	@echo "  lint         - Run linter"
	@echo "  install      - Install binary to GOPATH/bin"
	@echo "  install-lint - Install golangci-lint"
	@echo "  check        - Run fmt, lint, and test"
	@echo "  schema       - Generate JSON schema"
	@echo "  help         - Show this help"
	@echo ""
	@echo "Current version: $(VERSION)"
