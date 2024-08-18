# Variables
SERVICE_DIR := .
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/)
GOTESTSUM_PATH := $(shell go env GOPATH)/bin/gotestsum

# Default Goal
.DEFAULT_GOAL := help

# Help
.PHONY: help
help:  ## 💬 Display this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Dependency Management
.PHONY: install
install: ## 📦 Install dependencies
	@echo "Installing dependencies..."
	@go mod download -x
	@go install gotest.tools/gotestsum@latest

# Linting and Formatting
.PHONY: lint
lint: install ## 📜 Lint & format code, will try to fix errors and modify code
	@echo "Running linter..."
	@golangci-lint run ./...

.PHONY: lint-fix
lint-fix: install ## 📜 Lint & format code, automatically fix issues
	@echo "Running linter with auto-fix..."
	@golangci-lint run ./... --fix

.PHONY: format
format: ## 🧽 Format Go code
	@echo "Formatting code..."
	@go fmt ./...

# Running the Application
.PHONY: run
run: ## 🚀 Run the application
	@echo "Running the application..."
	@set -a && source .env && go run $(SERVICE_DIR)/cmd/main.go

# Testing
.PHONY: test
test: install ## 🧪 Run tests with coverage
	@echo "Running tests..."
	@go test ./... -cover -v

.PHONY: test-pretty
test-pretty: install ## 🧪 Run tests with enhanced output
	@echo "Running tests with enhanced output..."
	@$(GOTESTSUM_PATH) --format=short-verbose -- ./...

.PHONY: test-coverage
test-coverage: install ## 🧪 Run tests and generate a coverage report
	@echo "Running tests and generating coverage report..."
	@go test ./... -coverprofile=coverage.out
	@go tool cover -func=coverage.out

# Clean
.PHONY: clean
clean: ## 🧹 Clean up the workspace
	@echo "Cleaning up..."
	@go clean
	@rm -rf coverage.out
