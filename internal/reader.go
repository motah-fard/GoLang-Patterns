package internal

import (
	"bufio"
	"os"
)

func ReadLogs(filePath string, logCh chan<- string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		logCh <- scanner.Text() // Send log line to channel
	}

	return scanner.Err()
}

// Pattern Used: Producer-Consumer
// Concurrency: Runs as a goroutine, streaming logs to the next stage.
