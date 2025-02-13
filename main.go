package main

import (
	"fmt"
	"log_pipeline/internal"
	"sync"
)

func main() {
	logCh := make(chan string, 100)
	processedCh := make(chan *internal.LogEntry, 100)

	var wg sync.WaitGroup

	// Start log reader
	go func() {
		err := internal.ReadLogs("logs.txt", logCh)
		if err != nil {
			fmt.Println("Error reading logs:", err)
		}
		close(logCh) // Close after reading
	}()

	// Start worker pool
	numWorkers := 5
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go internal.ProcessLogs(logCh, processedCh, &wg)
	}

	// Collect processed logs
	go func() {
		wg.Wait()
		close(processedCh) // Close after processing
	}()

	// Write logs to JSON
	var processedLogs []*internal.LogEntry
	for entry := range processedCh {
		processedLogs = append(processedLogs, entry)
	}

	err := internal.WriteLogsToJSON(processedLogs, "output.json")
	if err != nil {
		fmt.Println("Error writing to JSON:", err)
	}

	fmt.Println("Log processing complete.")
}
