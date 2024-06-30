package logger

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestMustGetFile(t *testing.T) {
	t.Run("should return a valid file", func(t *testing.T) {
		filePath := "test.log"
		file := MustGetFile(filePath)
		defer func() {
			file.Close()
			os.Remove(filePath)
		}()

		assert.NotNil(t, file)
		info, err := file.Stat()
		assert.NoError(t, err)
		assert.False(t, info.IsDir())
	})

	t.Run("should panic if failed to open the file", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic, but did not")
			}
		}()
		MustGetFile("/invalid/path")
	})
}

func TestMustSetup(t *testing.T) {
	t.Run("set up the logger with different environments", func(t *testing.T) {
		var tests = []struct {
			name     string
			env      string
			expected slog.Level
			panic    bool
		}{
			{"for development environment should be debug level", "development", slog.LevelDebug, false},
			{"for production environment should be info level", "production", slog.LevelInfo, false},
			{"for unknown environment should panic", "unknown", slog.LevelInfo, true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				filePath := "test.log"
				file := MustGetFile(filePath)
				defer func() {
					file.Close()
					os.Remove(filePath)
				}()

				if tt.panic {
					defer func() {
						if r := recover(); r == nil {
							t.Error("Expected panic, but did not")
						}
					}()
				}

				MustSetup(file, tt.env)
				assert.True(t, slog.Default().Enabled(context.Background(), tt.expected))
			})
		}
	})

	t.Run("should log stack trace for error", func(t *testing.T) {
		filePath := "test.log"
		file := MustGetFile(filePath)
		defer func() {
			file.Close()
			os.Remove(filePath)
		}()

		MustSetup(file, "development")

		err := errors.New("test error")
		slog.Error("test error message", slog.Any("error", err))

		// Close the file to ensure all content is written
		file.Close()

		// Reopen the file in read-only mode
		file, err = os.Open(filePath)
		assert.NoError(t, err)
		defer file.Close()

		// Read the file content
		content := make([]byte, 1024)
		n, err := file.Read(content)
		assert.NoError(t, err)

		// Check if the log file contains the error message and stack trace
		output := string(content[:n])
		assert.Contains(t, output, "test error message")
		assert.Contains(t, output, "msg")
		assert.Contains(t, output, "test error")
		assert.Contains(t, output, "trace")
	})
}

func TestGetLogLevel(t *testing.T) {
	assert.Equal(t, slog.LevelDebug, getLogLevel("development"))
	assert.Equal(t, slog.LevelInfo, getLogLevel("production"))
	assert.Panics(t, func() { getLogLevel("unknown") })
}

func TestReplaceAttr(t *testing.T) {
	attr := slog.Attr{
		Key:   "error",
		Value: slog.AnyValue(errors.New("test error")),
	}

	replaced := replaceAttr(nil, attr)
	assert.Equal(t, "error", replaced.Key)
	assert.Equal(t, slog.KindGroup, replaced.Value.Kind())
}

func TestFormatError(t *testing.T) {
	err := errors.New("test error")
	value := formatError(err)

	assert.Equal(t, slog.KindGroup, value.Kind())
	assert.Contains(t, value.String(), "msg")
	assert.Contains(t, value.String(), "test error")
}

func TestGetStackTracer(t *testing.T) {
	t.Run("should return stack tracer", func(t *testing.T) {
		err := errors.New("test error")
		st, ok := getStackTracer(err)

		assert.NotNil(t, st)
		assert.True(t, ok)
	})
	t.Run("should return stack tracer from wrapped error", func(t *testing.T) {
		err := errors.Wrap(errors.New("test error"), "wrapped error")
		st, ok := getStackTracer(err)

		assert.NotNil(t, st)
		assert.True(t, ok)
	})
}

func TestFormatStackTrace(t *testing.T) {
	stack := errors.New("test error").(stackTracer).StackTrace()
	formatted := formatStackTrace(stack)

	assert.NotEmpty(t, formatted)
	assert.Contains(t, formatted[0], "logger_test.go")
}
