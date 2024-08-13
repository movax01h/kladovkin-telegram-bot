# Don't change
SERIVCE_DIR := .

.PHONY: test
.DEFAULT_GOAL := help

help:  ## 💬 This help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

install: ## 📦 Install dependencies
	@go mod download -x

lint: ## 📜 Lint & format, will try to fix errors and modify code
	@golangci-lint run

test: install  ## 🧪 Run tests
	@go test ./... -cover
