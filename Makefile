.PHONY: build run clean install test lint fmt help release

BINARY_NAME=vex
VERSION=2.0.0
BUILD_DIR=dist
GO_FILES=$(shell find . -name '*.go' -type f)

# Build the application
build:
	@echo "Building $(BINARY_NAME) v$(VERSION)..."
	@go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BINARY_NAME) .
	@echo "Build complete!"

# Build optimized release binary
build-release:
	@echo "Building optimized release..."
	@go build -ldflags="-s -w -X main.version=$(VERSION)" -trimpath -o $(BINARY_NAME) .
	@echo "Release build complete!"

# Run with sample data
run: build
	@./$(BINARY_NAME) sample_data.csv

# Run tests
test:
	@echo "Running tests..."
	@go test -v -race -cover ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -race -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Code formatted!"

# Lint code
lint:
	@echo "Linting code..."
	@golangci-lint run --enable-all --disable exhaustruct,exhaustivestruct,gochecknoglobals,gochecknoinits
	@echo "Linting complete!"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies installed!"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@go clean
	@echo "Cleaned!"

# Install globally
install: build-release
	@echo "Installing $(BINARY_NAME)..."
	@go install
	@echo "Installed! Run '$(BINARY_NAME)' from anywhere"

# Create release builds for all platforms
release: clean
	@echo "Creating release builds v$(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w -X main.version=$(VERSION)" -trimpath -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	@GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w -X main.version=$(VERSION)" -trimpath -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	@GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X main.version=$(VERSION)" -trimpath -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	@GOOS=linux GOARCH=arm64 go build -ldflags="-s -w -X main.version=$(VERSION)" -trimpath -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 .
	@GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -X main.version=$(VERSION)" -trimpath -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	@cd $(BUILD_DIR) && sha256sum * > checksums.txt
	@echo "Release builds created in $(BUILD_DIR)/"

# Security audit
audit:
	@echo "Running security audit..."
	@go list -json -m all | nancy sleuth
	@echo "Security audit complete!"

# Benchmarks
bench:
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...

# Show help
help:
	@echo "Excel TUI v$(VERSION) - Makefile Help"
	@echo ""
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  build-release  - Build optimized release binary"
	@echo "  run            - Build and run with sample data"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  fmt            - Format code"
	@echo "  lint           - Lint code (requires golangci-lint)"
	@echo "  deps           - Install dependencies"
	@echo "  clean          - Remove build artifacts"
	@echo "  install        - Install globally"
	@echo "  release        - Create release builds for all platforms"
	@echo "  audit          - Run security audit (requires nancy)"
	@echo "  bench          - Run benchmarks"
	@echo "  help           - Show this help"
	@echo ""
	@echo "Examples:"
	@echo "  make build"
	@echo "  make test"
	@echo "  make release"

# Default target
.DEFAULT_GOAL := help
