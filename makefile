.PHONY: build test lint clean install

# Build configuration
BINARY_NAME=consensusforge
BUILD_DIR=build
GO_FILES=$(shell find . -name "*.go" -type f)

# Build the main binary
build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/consensusforge

# Run tests
test:
	go test -v -race ./...

# Run tests with coverage
test-coverage:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Lint code
lint:
	golangci-lint run

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Install dependencies
deps:
	go mod download
	go mod tidy

# Install the binary
install: build
	go install ./cmd/consensusforge

# Run integration tests
test-integration:
	go test -v -tags=integration ./...

# Generate protobuf files
proto:
	protoc --go_out=. --go-grpc_out=. api/proto/*.proto