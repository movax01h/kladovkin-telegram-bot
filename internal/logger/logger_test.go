package logger

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

// MockFileOpener is a mock implementation of FileOpener for testing.
type MockFileOpener struct {
	File io.WriteCloser
	Err  error
}

// OpenFile returns the mock file and error.
func (m MockFileOpener) OpenFile(name string, flag int, perm os.FileMode) (io.WriteCloser, error) {
	return m.File, m.Err
}

// MockWriteCloser is a mock implementation of io.WriteCloser for testing.
type MockWriteCloser struct {
	bytes.Buffer
}

// Close is a no-op for MockWriteCloser.
func (m *MockWriteCloser) Close() error {
	return nil
}

func TestOpenLogFile(t *testing.T) {
	t.Run("should return a valid file", func(t *testing.T) {
		mockFile := &MockWriteCloser{}
		opener := MockFileOpener{File: mockFile, Err: nil}
		file := OpenLogFile("test.log", opener)

		assert.NotNil(t, file)
	})

	t.Run("should panic if failed to open the file", func(t *testing.T) {
		opener := MockFileOpener{File: nil, Err: errors.New("failed to open file")}
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic, but did not")
			}
		}()
		OpenLogFile("/invalid/path", opener)
	})
}

func TestMustSetup(t *testing.T) {
	t.Run("set up the logger with different environments", func(t *testing.T) {
		tests := []struct {
			name        string
			env         string
			expected    slog.Level
			shouldPanic bool
		}{
			{"development environment should be debug level", "development", slog.LevelDebug, false},
			{"production environment should be info level", "production", slog.LevelInfo, false},
			{"unknown environment should panic", "unknown", slog.LevelInfo, true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mockFile := &MockWriteCloser{}
				opener := MockFileOpener{File: mockFile, Err: nil}
				file := OpenLogFile("test.log", opener)

				if tt.shouldPanic {
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
		var buffer bytes.Buffer
		mockFile := &MockWriteCloser{}
		opener := MockFileOpener{File: mockFile, Err: nil}
		file := OpenLogFile("test.log", opener)
		mw := io.MultiWriter(&buffer, file)

		MustSetup(mw, "development")

		err := errors.New("test error")
		slog.Error("test error message", slog.Any("error", err))

		output := buffer.String()
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
	err := errors.New("test error")
	stack, ok := getStackTracer(err)
	assert.True(t, ok)
	formatted := formatStackTrace(stack.StackTrace())

	assert.NotEmpty(t, formatted)
	assert.Contains(t, formatted[0], "logger_test.go")
}
