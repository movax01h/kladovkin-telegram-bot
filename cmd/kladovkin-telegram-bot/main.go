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
	logFile := logger.MustGetFile(cfg.LogFilePath)
	defer logFile.Close()
	logger.MustSetup(logFile, cfg.Env)

	// Log a configuration
	slog.Info("configuration is loaded", "config", cfg)

	// Test error stack trace
	if err := test(); err != nil {
		slog.Error("test error", "error", err)
	}

	// Initialize storage

	// Start site parser

	// Start telegram bot
}

func test() error {
	return errors.New("test error")
}
