package internal

import (
	"fmt"
	"sync"
)

func ProcessLogs(logCh <-chan string, processedCh chan<- *LogEntry, wg *sync.WaitGroup) {
	defer wg.Done()

	for logLine := range logCh {
		entry, err := ParseLog(logLine)
		if err != nil || entry == nil {
			continue
		}

		if entry.Level == "ERROR" { // Example: Filter ERROR logs
			fmt.Println("Found error log:", entry)
		}

		processedCh <- entry // Send structured logs to next stage
	}
}
