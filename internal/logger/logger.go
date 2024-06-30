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

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func MustGetFile(path string) *os.File {
	// Open the log file
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Panicf("failed to open log file: %v", err)
	}
	return file
}

func MustSetup(file *os.File, env string) {
	level := getLogLevel(env)
	mw := io.MultiWriter(os.Stdout, file)

	logger := slog.New(slog.NewJSONHandler(mw, &slog.HandlerOptions{
		ReplaceAttr: replaceAttr,
		Level:       level,
	}))

	slog.SetDefault(logger)
	slog.Info("logger is set up")
}

func getLogLevel(env string) slog.Level {
	switch env {
	case "development":
		return slog.LevelDebug
	case "production":
		return slog.LevelInfo
	default:
		log.Panicf("unknown environment: \"%s\". Use \"development\" or \"production\"", env)
		return slog.LevelInfo // Unreachable, but keeps the compiler happy.
	}
}

func replaceAttr(_ []string, attr slog.Attr) slog.Attr {
	if attr.Value.Kind() == slog.KindAny {
		if err, ok := attr.Value.Any().(error); ok {
			attr.Value = formatError(err)
		}
	}
	return attr
}

func formatError(err error) slog.Value {
	var attrs []slog.Attr
	attrs = append(attrs, slog.String("msg", err.Error()))

	if st, ok := getStackTracer(err); ok {
		attrs = append(attrs, slog.Any("trace", formatStackTrace(st.StackTrace())))
	}

	return slog.GroupValue(attrs...)
}

func getStackTracer(err error) (stackTracer, bool) {
	for err != nil {
		if st, ok := err.(stackTracer); ok {
			return st, true
		}
		err = errors.Unwrap(err)
	}
	return nil, false
}

func formatStackTrace(frames errors.StackTrace) []string {
	lines := make([]string, len(frames))
	var skipped int
	skipping := true
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
		} else {
			skipping = false
		}

		file, line := fn.FileLine(pc)
		lines[i] = fmt.Sprintf("%s %s:%d", name, file, line)
	}
	return lines[:len(lines)-skipped]
}
