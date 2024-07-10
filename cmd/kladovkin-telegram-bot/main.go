package main

import (
	"log/slog"

	"github.com/pkg/errors"

	"github.com/movax01h/kladovkin-telegram-bot/internal/config"
	"github.com/movax01h/kladovkin-telegram-bot/internal/logger"
)

func main() {
	// Set up configuration
	cfg := config.MustLoad()

	// Set up the logger
	logFile := logger.OpenLogFile(cfg.LogFilePath, logger.OSFileOpener{})
	defer logFile.Close()
	logger.MustSetup(logFile, cfg.Env)

	// Initialize storage

	// Start site parser

	// Start telegram bot
}

func test() error {
	return errors.New("test error")
}
