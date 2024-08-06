package logger

import (
	"bytes"
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

// TestSetup tests the Setup function.
// It should set up the logger with the given log level and output file.
// The logger should replace any error values with a formatted error message and stack trace.
// The logger should write log messages to the output file if specified.
func TestSetup(t *testing.T) {
	testCases := []struct {
		name     string
		level    slog.Level
		expected slog.Level
	}{
		{"for debug level should be debug level", slog.LevelDebug, slog.LevelDebug},
		{"for info level should be info level", slog.LevelInfo, slog.LevelInfo},
		{"for warn level should be warn level", slog.LevelWarn, slog.LevelWarn},
		{"for error level should be error level", slog.LevelError, slog.LevelError},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			Setup(tc.level, os.Stdout)

			// Assert
			assert.True(t, slog.Default().Enabled(context.Background(), tc.expected))
		})
	}

	t.Run("write log message to the buffer", func(t *testing.T) {
		// Arrange
		var buf bytes.Buffer
		Setup(slog.LevelDebug, &buf)

		// Act
		slog.Info("test message")

		// Assert
		assert.Contains(t, buf.String(), "msg")
		assert.Contains(t, buf.String(), "test message")
	})

	t.Run("should replace error with formatted error message and stack trace", func(t *testing.T) {
		// Arrange
		var buf bytes.Buffer
		Setup(slog.LevelDebug, &buf)

		// Act
		err := errors.New("test error")
		slog.Error("test error message", slog.Any("error", err))

		// Assert
		assert.Contains(t, buf.String(), "msg")
		assert.Contains(t, buf.String(), "test error")
		assert.Contains(t, buf.String(), "trace")
	})
}

// TestReplaceAttr tests the replaceAttr function.
// It should replace the value of an attribute with a formatted error message and stack trace.
// The attribute value should be a group containing the error message and stack trace.
func TestReplaceAttr(t *testing.T) {
	// Act
	attr := replaceAttr(nil, slog.Attr{
		Key:   "error",
		Value: slog.AnyValue(errors.New("test error")),
	})

	// Assert
	assert.Equal(t, "error", attr.Key)
	assert.Equal(t, slog.KindGroup, attr.Value.Kind())
	assert.Contains(t, attr.Value.String(), "msg")
	assert.Contains(t, attr.Value.String(), "test error")
	assert.Contains(t, attr.Value.String(), "trace")
}

// TestFormatError tests the formatError function.
// It should return a group value containing the error message and stack trace.
func TestFormatError(t *testing.T) {
	// Act
	err := errors.New("test error")
	value := formatError(err)

	// Assert
	assert.Equal(t, slog.KindGroup, value.Kind())
	assert.Contains(t, value.String(), "msg")
	assert.Contains(t, value.String(), "test error")
}

// TestGetStackTracer tests the getStackTracer function.
// It should return the stackTracer interface from the error.
// If the error is wrapped, it should return the stackTracer interface from the wrapped error.
func TestGetStackTracer(t *testing.T) {
	t.Run("should return stack tracer", func(t *testing.T) {
		// Act
		err := errors.New("test error")
		st, ok := getStackTracer(err)

		// Assert
		assert.NotNil(t, st)
		assert.True(t, ok)
	})
	t.Run("should return stack tracer from wrapped error", func(t *testing.T) {
		// Act
		err := errors.Wrap(errors.New("test error"), "wrapped error")
		st, ok := getStackTracer(err)

		// Assert
		assert.NotNil(t, st)
		assert.True(t, ok)
	})
}

// TestFormatStackTrace tests the formatStackTrace function.
// It should return a slice of strings containing the stack trace.
// The stack trace should contain the file name and line number.
func TestFormatStackTrace(t *testing.T) {
	// Act
	stack := errors.New("test error").(stackTracer).StackTrace()
	formatted := formatStackTrace(stack)

	// Assert
	assert.NotEmpty(t, formatted)
	assert.Contains(t, formatted[0], "logger_test.go")
}
