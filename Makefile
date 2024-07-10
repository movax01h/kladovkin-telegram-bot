# Variables
APP_NAME = kladovkin-telegram-bot
CMD_PATH = ./cmd/$(APP_NAME)/main.go
DOCKERFILE_PATH = ./build/package/DockerFile

# Go commands
GO_CMD = go
GOFMT_CMD = gofmt
DOCKER_CMD = docker

# Targets
.PHONY: all run test docker fmt clean install help
.DEFAULT_GOAL := help

install:  ## 📦 Install dependencies
	@echo "Installing dependencies..."
	@$(GO_CMD) mod download

run:  ## 🏃 Run the application
	@echo "Running the application..."
	@$(GO_CMD) run $(CMD_PATH)

run-docker:  ## 🏃 Build and run docker container
	@echo "Building and running docker container..."
	@$(DOCKER_CMD) build -t $(APP_NAME) -f $(DOCKERFILE_PATH) .
	@$(DOCKER_CMD) run $(APP_NAME)

test:  ## 🧪 Run unit tests
	@echo "Running unit tests..."
	@$(GO_CMD) test ./...

fmt:  ## 📝 Run go formatter
	@echo "Running go formatter..."
	@$(GOFMT_CMD) -w .


clean:  ## 🧹 Clean up
	@echo "Cleaning up..."
	@$(GO_CMD) clean

help:  ## 💬 This help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
