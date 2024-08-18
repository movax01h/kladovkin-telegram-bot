package tools

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// OpenLogFile opens the log file and returns an io.Writer for logging.
func OpenLogFile(logFilePath string) (io.Writer, func() error, error) {
	// If no log file path is provided, return stdout
	if logFilePath == "" {
		return os.Stdout, func() error { return nil }, nil
	}

	// Ensure the directory exists
	dir := filepath.Dir(logFilePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, nil, fmt.Errorf("log directory does not exist: %s", dir)
	}

	// Open the log file
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open log file: %w", err)
	}

	multiWriter := io.MultiWriter(os.Stdout, file)

	// Return the MultiWriter and a cleanup function to close the file
	return multiWriter, file.Close, nil
}
