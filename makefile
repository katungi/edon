# Build variables
BINARY_NAME=halo
VERSION=$(shell git describe --tags --always --dirty)
COMMIT=$(shell git rev-parse --short HEAD)
BUILD_DATE=$(shell date -u '+%Y-%m-%d %H:%M:%S')
LD_FLAGS='-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(BUILD_DATE)'

# Go related variables
GOFILES=$(shell find . -type f -name '*.go')

# Make targets
.PHONY: all build build-runtime build-web build-all clean run test help

all: clean build

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p bin/
	@go build -ldflags $(LD_FLAGS) -o bin/$(BINARY_NAME) ./cmd/edon

build-runtime:
	@echo "Building runtime..."
	@mkdir -p bin/
	@go build -o bin/$(BINARY_NAME)-runtime ./cmd/runtime

build-web:
	@echo "Building web server..."
	@mkdir -p bin/
	@go build -o bin/$(BINARY_NAME)-web ./cmd/web

build-all: build build-runtime build-web

clean:
	@echo "Cleaning..."
	@go clean

run: build
	@./bin/$(BINARY_NAME)

test:
	@echo "Running tests..."
	@go test -v ./tests/...

dev: build
	@./bin/$(BINARY_NAME)

install: build
	@echo "Installing $(BINARY_NAME)..."
	@if [ -z "$(GOPATH)" ]; then \
		echo "Error: GOPATH is not set"; \
		exit 1; \
	fi
	@echo "Using GOPATH: $(GOPATH)"
	@if [ ! -d "$(GOPATH)/bin" ]; then \
		echo "Creating directory: $(GOPATH)/bin"; \
		mkdir -p "$(GOPATH)/bin" || { echo "Error: Failed to create $(GOPATH)/bin"; exit 1; }; \
	fi
	@echo "Moving binary from bin/$(BINARY_NAME) to $(GOPATH)/bin/$(BINARY_NAME)"
	@mv bin/$(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME) || { echo "Error: Failed to move binary"; exit 1; }
	@chmod +x $(GOPATH)/bin/$(BINARY_NAME) || { echo "Error: Failed to set executable permissions"; exit 1; }
	@echo "Successfully installed $(BINARY_NAME) to $(GOPATH)/bin"
	@SHELL_CONFIG=""; \
	if [ -f "$$HOME/.zshrc" ]; then \
		SHELL_CONFIG="$$HOME/.zshrc"; \
		echo "Found shell config: $$HOME/.zshrc"; \
	elif [ -f "$$HOME/.bashrc" ]; then \
		SHELL_CONFIG="$$HOME/.bashrc"; \
		echo "Found shell config: $$HOME/.bashrc"; \
	fi; \
	if [ -z "$$SHELL_CONFIG" ]; then \
		echo "Warning: No shell config file found (.zshrc or .bashrc)"; \
	else \
		if ! grep -q "export PATH=\"\$$GOPATH/bin:\$$PATH\"" "$$SHELL_CONFIG"; then \
			echo "Adding GOPATH/bin to PATH in $$SHELL_CONFIG"; \
			echo '\nexport PATH="$$GOPATH/bin:$$PATH"' >> "$$SHELL_CONFIG" || { echo "Error: Failed to update $$SHELL_CONFIG"; exit 1; }; \
			echo "Added GOPATH/bin to PATH in $$SHELL_CONFIG"; \
			echo "Please run 'source $$SHELL_CONFIG' to update your current shell"; \
		else \
			echo "GOPATH/bin already in PATH (found in $$SHELL_CONFIG)"; \
		fi \
	fi

lint:
	@echo "Linting..."
	@golangci-lint run

help:
	@echo "Available targets:"
	@echo "  build         - Build the main CLI binary"
	@echo "  build-runtime - Build the standalone runtime binary"
	@echo "  build-web     - Build the web server binary"
	@echo "  build-all     - Build all binaries"
	@echo "  clean         - Clean build artifacts"
	@echo "  run           - Run the binary"
	@echo "  test          - Run tests"
	@echo "  dev           - Build and run for development"
	@echo "  install       - Install binary to GOPATH"
	@echo "  lint          - Run linter"
