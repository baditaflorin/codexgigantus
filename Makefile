.PHONY: help build build-cli build-web test clean run-cli run-web docker-build docker-up docker-down install deps lint fmt vet

# Variables
APP_NAME=codexgigantus
CLI_BINARY=$(APP_NAME)-cli
WEB_BINARY=$(APP_NAME)-web
DOCKER_IMAGE=$(APP_NAME)
GO=go
GOFLAGS=-v

# Default target
help: ## Show this help message
	@echo "CodexGigantus Makefile"
	@echo "====================="
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Build targets
build: build-cli build-web ## Build both CLI and Web binaries

build-cli: ## Build CLI binary
	@echo "Building CLI binary..."
	$(GO) build $(GOFLAGS) -o $(CLI_BINARY) ./cmd/cli

build-web: ## Build Web GUI binary
	@echo "Building Web GUI binary..."
	$(GO) build $(GOFLAGS) -o $(WEB_BINARY) ./cmd/web

# Run targets
run-cli: build-cli ## Run CLI application
	@echo "Running CLI..."
	./$(CLI_BINARY) --help

run-web: build-web ## Run Web GUI application
	@echo "Running Web GUI..."
	./$(WEB_BINARY)

# Test targets
test: ## Run all tests
	@echo "Running tests..."
	$(GO) test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-unit: ## Run unit tests only
	@echo "Running unit tests..."
	$(GO) test -v -short ./...

test-race: ## Run tests with race detector
	@echo "Running tests with race detector..."
	$(GO) test -v -race ./...

# Code quality targets
lint: ## Run linter
	@echo "Running golangci-lint..."
	@command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; exit 1; }
	golangci-lint run ./...

fmt: ## Format code
	@echo "Formatting code..."
	$(GO) fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	$(GO) vet ./...

check: fmt vet lint test ## Run all checks (fmt, vet, lint, test)

# Dependency management
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	$(GO) mod download

deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	$(GO) get -u ./...
	$(GO) mod tidy

deps-tidy: ## Tidy dependencies
	@echo "Tidying dependencies..."
	$(GO) mod tidy

# Docker targets
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):latest .

docker-up: ## Start services with docker-compose
	@echo "Starting services..."
	docker-compose up -d

docker-down: ## Stop services with docker-compose
	@echo "Stopping services..."
	docker-compose down

docker-logs: ## View docker-compose logs
	docker-compose logs -f

docker-clean: ## Remove Docker image and containers
	@echo "Cleaning Docker resources..."
	docker-compose down -v
	docker rmi $(DOCKER_IMAGE):latest || true

# Installation targets
install: build ## Install binaries to GOPATH/bin
	@echo "Installing binaries..."
	$(GO) install ./cmd/cli
	$(GO) install ./cmd/web

install-cli: build-cli ## Install CLI binary to GOPATH/bin
	@echo "Installing CLI binary..."
	$(GO) install ./cmd/cli

install-web: build-web ## Install Web binary to GOPATH/bin
	@echo "Installing Web binary..."
	$(GO) install ./cmd/web

# Clean targets
clean: ## Remove build artifacts
	@echo "Cleaning build artifacts..."
	rm -f $(CLI_BINARY) $(WEB_BINARY)
	rm -f coverage.out coverage.html
	rm -rf dist/

clean-all: clean docker-clean ## Remove all build artifacts and Docker resources
	@echo "Cleaning all..."
	$(GO) clean -cache -testcache -modcache

# Example configurations
example-config: ## Create example configuration files
	@echo "Creating example configurations..."
	@mkdir -p configs
	@echo '{\n  "source_type": "filesystem",\n  "directories": ["."],\n  "recursive": true,\n  "output_file": "output.txt"\n}' > configs/filesystem.json
	@echo '{\n  "source_type": "database",\n  "db_type": "postgres",\n  "db_host": "localhost",\n  "db_port": 5432,\n  "db_name": "codex",\n  "db_user": "postgres",\n  "db_table_name": "code_files",\n  "db_column_path": "file_path",\n  "db_column_content": "content",\n  "output_file": "output.txt"\n}' > configs/database.json
	@echo "Example configs created in configs/"

# Development targets
dev-cli: ## Run CLI in development mode with race detector
	$(GO) run -race ./cmd/cli

dev-web: ## Run Web GUI in development mode with race detector
	$(GO) run -race ./cmd/web

watch: ## Watch for changes and rebuild (requires entr)
	@command -v entr >/dev/null 2>&1 || { echo "entr not installed. Install it with your package manager."; exit 1; }
	find . -name '*.go' | entr -r make build

# Documentation
docs: ## Generate documentation
	@echo "Generating documentation..."
	$(GO) doc -all ./...

docs-server: ## Start documentation server
	@echo "Starting documentation server on http://localhost:6060"
	godoc -http=:6060

# Release targets
release: clean check build ## Prepare for release
	@echo "Release build complete!"
	@echo "CLI binary: $(CLI_BINARY)"
	@echo "Web binary: $(WEB_BINARY)"

# Version info
version: ## Show version information
	@echo "Go version: $(shell $(GO) version)"
	@echo "App version: $(shell git describe --tags --always --dirty 2>/dev/null || echo 'dev')"

# Quick shortcuts
.PHONY: b r t c
b: build ## Shortcut for build
r: run-cli ## Shortcut for run-cli
t: test ## Shortcut for test
c: clean ## Shortcut for clean
