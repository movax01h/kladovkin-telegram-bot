package config

import (
	"log/slog"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_ValidInput(t *testing.T) {
	t.Setenv("ENVIRONMENT", "development")
	t.Setenv("LOG_LEVEL", "DEBUG")
	t.Setenv("LOG_FILE_PATH", "./logs/test.log")
	t.Setenv("DATABASE_PATH", "./data/telegram_bot.db")
	t.Setenv("TELEGRAM_BOT_TOKEN", "dummy_token")
	t.Setenv("TELEGRAM_ADMIN_ID", "123456")
	t.Setenv("NOTIFIER_INTERVAL", "10")
	t.Setenv("PARSER_URL", "https://kladovkin.ru/")
	t.Setenv("PARSER_INTERVAL", "10")

	cfg, err := NewConfig()
	require.NoError(t, err)
	assert.Equal(t, "development", cfg.Environment)
	assert.Equal(t, slog.LevelDebug, cfg.LoggerConfig.Level)
	assert.Equal(t, "./logs/app.log", cfg.LoggerConfig.FilePath)
	assert.Equal(t, "./data/test.db", cfg.DatabaseConfig.Path)
	assert.Equal(t, "dummy_token", cfg.TelegramConfig.BotToken)
	assert.Equal(t, int64(123456), cfg.TelegramConfig.AdminID)
	assert.Equal(t, int64(10), cfg.NotifierConfig.Interval)
	assert.Equal(t, "https://kladovkin.ru/", cfg.ParserConfig.URL)
	assert.Equal(t, int64(10), cfg.ParserConfig.Interval)
}

func TestValidateConfig(t *testing.T) {
	tempDir := t.TempDir() // Create a temporary directory for the test files

	tests := []struct {
		name        string
		cfg         *Config
		expectedErr bool
	}{
		{
			name: "valid configuration",
			cfg: &Config{
				Environment: "development",
				LoggerConfig: LoggerConfig{
					Level:    slog.LevelDebug,
					FilePath: filepath.Join(tempDir, "logfile.log"),
				},
				TelegramConfig: TelegramConfig{
					BotToken: "dummy_token",
					AdminID:  123456,
				},
			},
			expectedErr: false,
		},
		{
			name: "invalid environment",
			cfg: &Config{
				Environment: "invalid",
				LoggerConfig: LoggerConfig{
					Level:    slog.LevelDebug,
					FilePath: filepath.Join(tempDir, "logfile.log"),
				},
				TelegramConfig: TelegramConfig{
					BotToken: "dummy_token",
					AdminID:  123456,
				},
			},
			expectedErr: true,
		},
		{
			name: "invalid log level",
			cfg: &Config{
				Environment: "development",
				LoggerConfig: LoggerConfig{
					Level:    slog.Level(999),
					FilePath: filepath.Join(tempDir, "logfile.log")},
				TelegramConfig: TelegramConfig{
					BotToken: "dummy_token",
					AdminID:  123456,
				},
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.cfg)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIsValidEnvironment(t *testing.T) {
	tests := []struct {
		env      string
		expected bool
	}{
		{"development", true},
		{"production", true},
		{"invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.env, func(t *testing.T) {
			assert.Equal(t, tt.expected, isValidEnvironment(tt.env))
		})
	}
}

func TestIsValidLogLevel(t *testing.T) {
	tests := []struct {
		level    slog.Level
		expected bool
	}{
		{slog.LevelDebug, true},
		{slog.LevelInfo, true},
		{slog.LevelWarn, true},
		{slog.LevelError, true},
		{slog.Level(999), false},
	}

	for _, tt := range tests {
		t.Run(tt.level.String(), func(t *testing.T) {
			assert.Equal(t, tt.expected, isValidLogLevel(tt.level))
		})
	}
}

func TestValidateFilePath(t *testing.T) {
	// Using TempDir for testing file paths
	t.Run("valid path", func(t *testing.T) {
		tmpDir := t.TempDir()
		validPath := filepath.Join(tmpDir, "valid.log")

		assert.NoError(t, validateFilePath(validPath))
	})

	t.Run("invalid path", func(t *testing.T) {
		invalidPath := "/invalid/dir/logfile.log"
		assert.Error(t, validateFilePath(invalidPath))
	})
}

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected slog.Level
		hasError bool
	}{
		{LogLevelDebug, slog.LevelDebug, false},
		{LogLevelInfo, slog.LevelInfo, false},
		{LogLevelWarn, slog.LevelWarn, false},
		{LogLevelError, slog.LevelError, false},
		{"invalid", slog.LevelInfo, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			level, err := ParseLogLevel(tt.input)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, level)
			}
		})
	}
}
