.PHONY: build run test test-coverage clean docker-build docker-run help test-endpoints

# Build the application
build:
	go build -o bin/listenarr ./cmd/listenarr

# Run the application
run:
	go run ./cmd/listenarr/main.go

# Run tests
test:
	go test ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests with verbose output
test-verbose:
	go test -v ./...

# Test API endpoints (requires server to be running)
test-endpoints:
	@if [ -f scripts/test-endpoints.sh ]; then \
		bash scripts/test-endpoints.sh; \
	else \
		go run scripts/test-endpoints.go; \
	fi

# Run all checks (format, vet, test)
check: fmt vet test

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./...

# Lint code (requires golangci-lint)
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean

# Build Docker image
docker-build:
	docker build -t listenarr:latest .

# Run Docker container
docker-run:
	docker-compose up -d

# Stop Docker container
docker-stop:
	docker-compose down

# View logs
docker-logs:
	docker-compose logs -f

# Install dependencies
deps:
	go mod download
	go mod tidy

# Run build-all script
build-all:
	bash .ai/build/build-all.sh

help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  run            - Run the application"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  test-verbose   - Run tests with verbose output"
	@echo "  test-endpoints - Test API endpoints (server must be running)"
	@echo "  check          - Run format, vet, and tests"
	@echo "  fmt            - Format code"
	@echo "  vet            - Run go vet"
	@echo "  lint           - Run linter (requires golangci-lint)"
	@echo "  clean          - Clean build artifacts"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  docker-stop    - Stop Docker container"
	@echo "  docker-logs    - View Docker logs"
	@echo "  deps           - Install dependencies"
	@echo "  build-all      - Run comprehensive build and test script"
