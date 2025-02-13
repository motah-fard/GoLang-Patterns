package internal

import (
	"errors"
	"strings"
)

// LogEntry represents structured log data
type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

// ParseLog converts raw log lines into structured data
func ParseLog(logLine string) (*LogEntry, error) {
	parts := strings.SplitN(logLine, " ", 4) // Correctly splitting into 4 parts: date, time, level, message

	if len(parts) < 4 {
		return nil, errors.New("invalid log format")
	}

	// Concatenate date and time to get full timestamp
	timestamp := parts[0] + " " + parts[1] // "2025-02-11 15:30:45"

	// Log level is now correctly assigned
	level := parts[2] // "INFO", "ERROR", etc.

	// Remaining message
	message := parts[3] // "User logged in successfully"

	return &LogEntry{
		Timestamp: timestamp,
		Level:     level,
		Message:   message,
	}, nil
}
