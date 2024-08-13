package tools

import (
	"fmt"
	"io"
	"os"
)

// OpenLogFile opens the log file and returns an io.Writer for logging.
// If a log file path is provided, it opens the file and returns a MultiWriter
// that writes to both the file and stdout. If no file path is provided, it returns stdout.
func OpenLogFile(logFilePath string) (io.Writer, func() error, error) {
	if logFilePath == "" {
		return os.Stdout, func() error { return nil }, nil
	}

	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open log file: %w", err)
	}

	multiWriter := io.MultiWriter(os.Stdout, file)

	// Return the MultiWriter and a cleanup function to close the file
	return multiWriter, file.Close, nil
}
