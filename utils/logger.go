package utils

import (
	"log"
	"os"
)

// Logger is a global logger instance
var Logger *log.Logger

func InitLogger() {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	Logger = log.New(logFile, "LOG_PIPELINE: ", log.Ldate|log.Ltime|log.Lshortfile)
}
