package logger

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

// FileOpener defines an interface for opening files.
type FileOpener interface {
	OpenFile(name string, flag int, perm os.FileMode) (io.WriteCloser, error)
}

// OSFileOpener is a concrete implementation of FileOpener using the os package.
type OSFileOpener struct{}

// OpenFile opens a file using os.OpenFile.
func (o OSFileOpener) OpenFile(name string, flag int, perm os.FileMode) (io.WriteCloser, error) {
	return os.OpenFile(name, flag, perm)
}

// stackTracer is an interface that wraps the StackTrace method.
type stackTracer interface {
	StackTrace() errors.StackTrace
}

// OpenLogFile opens the log file specified by path, creating it if necessary.
func OpenLogFile(path string, opener FileOpener) io.WriteCloser {
	file, err := opener.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Panicf("failed to open log file: %v", err)
	}
	return file
}

// MustSetup initializes the logger with the specified log writer and environment.
func MustSetup(writer io.Writer, env string) {
	level := getLogLevel(env)
	mw := io.MultiWriter(os.Stdout, writer)

	logger := slog.New(slog.NewJSONHandler(mw, &slog.HandlerOptions{
		ReplaceAttr: replaceAttr,
		Level:       level,
	}))

	slog.SetDefault(logger)
	slog.Info("logger is set up")
}

// getLogLevel returns the appropriate log level based on the environment.
func getLogLevel(env string) slog.Level {
	switch env {
	case "development":
		return slog.LevelDebug
	case "production":
		return slog.LevelInfo
	default:
		log.Panicf("unknown environment: %q. Use \"development\" or \"production\"", env)
		return slog.LevelInfo // Unreachable, but keeps the compiler happy.
	}
}

// replaceAttr replaces the attribute with a formatted error if applicable.
func replaceAttr(_ []string, attr slog.Attr) slog.Attr {
	if attr.Value.Kind() == slog.KindAny {
		if err, ok := attr.Value.Any().(error); ok {
			attr.Value = formatError(err)
		}
	}
	return attr
}

// formatError formats an error into a slog.Value.
func formatError(err error) slog.Value {
	attrs := []slog.Attr{slog.String("msg", err.Error())}

	if st, ok := getStackTracer(err); ok {
		attrs = append(attrs, slog.Any("trace", formatStackTrace(st.StackTrace())))
	}

	return slog.GroupValue(attrs...)
}

// getStackTracer returns the stack tracer if the error implements stackTracer.
func getStackTracer(err error) (stackTracer, bool) {
	for err != nil {
		if st, ok := err.(stackTracer); ok {
			return st, true
		}
		err = errors.Unwrap(err)
	}
	return nil, false
}

// formatStackTrace formats a stack trace into a slice of strings.
func formatStackTrace(frames errors.StackTrace) []string {
	var (
		lines    = make([]string, len(frames))
		skipped  int
		skipping = true
	)

	for i := len(frames) - 1; i >= 0; i-- {
		pc := uintptr(frames[i]) - 1
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			lines[i] = "unknown"
			skipping = false
			continue
		}

		name := fn.Name()
		if skipping && strings.HasPrefix(name, "runtime.") {
			skipped++
			continue
		}
		skipping = false

		file, line := fn.FileLine(pc)
		lines[i] = fmt.Sprintf("%s %s:%d", name, file, line)
	}
	return lines[:len(lines)-skipped]
}
