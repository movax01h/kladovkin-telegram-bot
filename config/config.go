package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Environment      string `env:"ENVIRONMENT,required"`
	LogLevel         slog.Level
	LogFilePath      string `env:"LOG_FILE_PATH"` // Optional, if not provided, logs will be written to stdout, otherwise to both stdout and the file
	TelegramBotToken string `env:"TELEGRAM_BOT_TOKEN,required"`
}

const (
	// Environment constants
	EnvDevelopment = "development"
	EnvProduction  = "production"

	// LogLevel strings
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
)

// NewConfig creates a new configuration instance by parsing environment variables
// and validating the parsed data. It returns the configuration instance or an error
// if the environment variables are not properly set.
func NewConfig() (*Config, error) {
	cfg, err := parseConfig()
	if err != nil {
		return nil, err
	}

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// parseConfig handles environment variable parsing and initial setup.
func parseConfig() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse the environment: %w", err)
	}

	level, err := ParseLogLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		return nil, err
	}
	cfg.LogLevel = level

	return &cfg, nil
}

// validateConfig validates the configuration values.
func validateConfig(cfg *Config) error {
	if !isValidEnvironment(cfg.Environment) {
		return fmt.Errorf("invalid environment: %s", cfg.Environment)
	}

	if !isValidLogLevel(cfg.LogLevel) {
		return fmt.Errorf("invalid log level: %s", cfg.LogLevel)
	}

	if cfg.LogFilePath != "" {
		if err := validateFilePath(cfg.LogFilePath); err != nil {
			return err
		}
	}

	return nil
}

// isValidEnvironment checks if the environment is valid.
func isValidEnvironment(env string) bool {
	return env == EnvDevelopment || env == EnvProduction
}

// isValidLogLevel checks if the log level is valid.
func isValidLogLevel(level slog.Level) bool {
	return level == slog.LevelDebug || level == slog.LevelInfo ||
		level == slog.LevelWarn || level == slog.LevelError
}

// validateFilePath validates the log file path.
func validateFilePath(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path of the log file: %w", err)
	}

	dir := filepath.Dir(absPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("log directory does not exist: %s", dir)
	}

	file, err := os.OpenFile(absPath, os.O_WRONLY|os.O_CREATE, 0o0600)
	if err != nil {
		return fmt.Errorf("log file cannot be opened or created: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			slog.Error("failed to close log file", "error", err)
		}
	}()

	return nil
}

// ParseLogLevel converts a string log level to slog.Level.
func ParseLogLevel(levelStr string) (slog.Level, error) {
	logLevels := map[string]slog.Level{
		LogLevelDebug: slog.LevelDebug,
		LogLevelInfo:  slog.LevelInfo,
		LogLevelWarn:  slog.LevelWarn,
		LogLevelError: slog.LevelError,
	}

	level, exists := logLevels[levelStr]
	if !exists {
		return slog.LevelInfo, fmt.Errorf("invalid log level: %s", levelStr)
	}

	return level, nil
}
