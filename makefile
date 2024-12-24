# Build variables
BINARY_NAME=halo
VERSION=$(shell git describe --tags --always --dirty)
COMMIT=$(shell git rev-parse --short HEAD)
BUILD_DATE=$(shell date -u '+%Y-%m-%d %H:%M:%S')
LD_FLAGS=-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(BUILD_DATE)

# Go related variables
GOFILES=$(shell find . -type f -name '*.go')

# Make targets
.PHONY: all build clean run test help

all: clean build

build:
	@echo "Building $(BINARY_NAME)..."
	@go build -ldflags "$(LD_FLAGS)" -o bin/$(BINARY_NAME) cmd/edon/main.go

clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@go clean

run: build
	@./bin/$(BINARY_NAME)

test:
	@echo "Running tests..."
	@go test -v ./...

dev: build
	@./bin/$(BINARY_NAME)

install: build
	@echo "Installing $(BINARY_NAME)..."
	@mv bin/$(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)

lint:
	@echo "Linting..."
	@golangci-lint run

help:
	@echo "Available targets:"
	@echo "  build    - Build the binary"
	@echo "  clean    - Clean build artifacts"
	@echo "  run      - Run the binary"
	@echo "  test     - Run tests"
	@echo "  dev      - Build and run for development"
	@echo "  install  - Install binary to GOPATH"
	@echo "  lint     - Run linter"
