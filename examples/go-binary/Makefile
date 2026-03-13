# Provenix Go Binary Example Makefile

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT  := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

BINARY_NAME := app
BUILD_DIR   := build
LDFLAGS     := -ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildDate=$(BUILD_DATE)"

.PHONY: all build build-multi test-provenix clean run version help

all: clean build

## build: Build the binary for the current platform
build:
	@echo "Building $(BINARY_NAME) version $(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "Binary built at $(BUILD_DIR)/$(BINARY_NAME)"

## build-multi: Build for all supported platforms (requires Unix shell: macOS/Linux/WSL/Git Bash)
build-multi:
	@echo "Building for multiple architectures..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux   GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	GOOS=linux   GOARCH=arm64 CGO_ENABLED=0 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 .
	GOOS=darwin  GOARCH=arm64 CGO_ENABLED=0 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	@echo "Multi-arch builds complete:"
	@ls -lh $(BUILD_DIR)/ 2>/dev/null || dir $(BUILD_DIR)

## test-provenix: Run the full Provenix attestation workflow locally
test-provenix: build
	@echo "Testing Provenix attestation workflow..."
	@echo "Step 1: Generate attestation"
	provenix attest $(BUILD_DIR)/$(BINARY_NAME) --output attestation.json --skip-transparency
	@echo "Step 2: Generate VEX"
	provenix vex generate --output vex.json
	@echo "Step 3: Policy check"
	provenix policy check
	@echo "Provenix workflow test complete!"

## clean: Remove build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f attestation.json vex.json
	@echo "Clean complete"

## run: Build and run the binary
run: build
	@$(BUILD_DIR)/$(BINARY_NAME)

## version: Show version information from the built binary
version: build
	@$(BUILD_DIR)/$(BINARY_NAME) version

## help: Show this help message
help:
	@grep -E '^## ' Makefile | sed 's/## /  /'
