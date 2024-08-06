package logger

import (
	"fmt"
	"io"
	"runtime"
	"strings"

	"log/slog"

	"github.com/pkg/errors"
)

// stackTracer is an interface that may be implemented by types to provide stack traces.
type stackTracer interface {
	StackTrace() errors.StackTrace
}

// Setup sets up the logger with the given environment and optional output file.
// Loggers are set up with a JSON handler that writes to the given file and
// standard output. The logger will replace any error values with a formatted
// error message and stack trace.
func Setup(l slog.Level, w io.Writer) {
	logger := slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
		ReplaceAttr: replaceAttr,
		Level:       l,
	}))

	slog.SetDefault(logger)
	slog.Info("Logger is set up with level", "level", l)
}

// replaceAttr is a handler function that replaces the value of an attribute.
func replaceAttr(_ []string, attr slog.Attr) slog.Attr {
	if attr.Value.Kind() == slog.KindAny {
		if err, ok := attr.Value.Any().(error); ok {
			attr.Value = formatError(err)
		}
	}
	return attr
}

// formatError formats an error value into a slog.Value.
func formatError(err error) slog.Value {
	attrs := []slog.Attr{slog.String("msg", err.Error())}

	if st, ok := getStackTracer(err); ok {
		attrs = append(attrs, slog.Any("trace", formatStackTrace(st.StackTrace())))
	}

	return slog.GroupValue(attrs...)
}

// getStackTracer returns the stackTracer interface from the error.
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
