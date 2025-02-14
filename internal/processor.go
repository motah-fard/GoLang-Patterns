package internal

import (
	"fmt"
	"sync"
)

// ProcessLogs processes logs using worker goroutines
func ProcessLogs(logCh <-chan string, processedCh chan<- *LogEntry, wg *sync.WaitGroup) {
	defer wg.Done() // Ensure Done is only called once per worker

	for logLine := range logCh {
		entry, err := ParseLog(logLine)
		if err != nil || entry == nil {
			continue
		}

		if entry.Level == "ERROR" {
			fmt.Println("Found error log:", entry)
		}

		processedCh <- entry // Send processed log
	}
}
