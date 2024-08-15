package config

import (
	"log/slog"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig_ValidInput(t *testing.T) {
	t.Setenv("ENVIRONMENT", "development")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("TELEGRAM_BOT_TOKEN", "dummy_token")
	t.Setenv("ERROR_NOTIFICATION_CHAT_ID", "123456")

	cfg, err := NewConfig()
	require.NoError(t, err)
	assert.Equal(t, "development", cfg.Environment)
	assert.Equal(t, slog.LevelDebug, cfg.LoggerConfig.Level)
	assert.Equal(t, "dummy_token", cfg.TelegramConfig.Token)
}

func TestParseConfig_ValidInput(t *testing.T) {
	t.Setenv("ENVIRONMENT", "production")
	t.Setenv("LOG_LEVEL", "info")
	t.Setenv("TELEGRAM_BOT_TOKEN", "dummy_token")
	t.Setenv("ERROR_NOTIFICATION_CHAT_ID", "123456")

	cfg, err := parseConfig()
	require.NoError(t, err)
	assert.Equal(t, "production", cfg.Environment)
	assert.Equal(t, slog.LevelInfo, cfg.LoggerConfig.Level)
	assert.Equal(t, "dummy_token", cfg.TelegramConfig.Token)
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
					Token:                   "dummy_token",
					ErrorNotificationChatID: 123456,
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
					Token:                   "dummy_token",
					ErrorNotificationChatID: 123456,
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
					Token:                   "dummy_token",
					ErrorNotificationChatID: 123456,
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
