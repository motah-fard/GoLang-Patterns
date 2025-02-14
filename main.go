package main

import (
	"fmt"
	"log_pipeline/internal"
	"log_pipeline/utils"
	"sync"

	"github.com/spf13/viper"
)

// loadConfig initializes Viper to read settings from config.yaml
func loadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".") // Look in the project root directory

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file:", err)
	}
}

func main() {
	// Initialize logger
	utils.InitLogger()
	utils.Logger.Println("Starting Log Pipeline...")

	// Load configuration
	loadConfig()

	// Read values from config.yaml
	logFile := viper.GetString("log_file")
	outputFile := viper.GetString("output_file")
	numWorkers := viper.GetInt("num_workers")

	if numWorkers <= 0 {
		numWorkers = 5 // Default value if not set in config.yaml
	}

	utils.Logger.Printf("Loaded Config: log_file=%s, output_file=%s, num_workers=%d\n", logFile, outputFile, numWorkers)

	// Setup context for graceful shutdown
	ctx := utils.SetupSignalHandler()
	defer utils.Logger.Println("Shutting down gracefully...")

	logCh := make(chan string, 100)
	processedCh := make(chan *internal.LogEntry, 100)

	var wg sync.WaitGroup

	// Start log reader
	go func() {
		err := internal.ReadLogs(logFile, logCh)
		if err != nil {
			utils.Logger.Printf("Error reading logs: %v\n", err)
		}
		close(logCh) // Close after reading
	}()

	// Start worker pool
	// Add workers to WaitGroup before starting
	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go internal.ProcessLogs(logCh, processedCh, &wg) // No defer wg.Done() here, it's inside ProcessLogs
	}

	// Collect processed logs
	go func() {
		wg.Wait()
		close(processedCh) // Close after processing
	}()

	// Write logs to JSON with context handling
	var processedLogs []*internal.LogEntry
	done := make(chan struct{})

	go func() {
		for entry := range processedCh {
			processedLogs = append(processedLogs, entry)
		}
		done <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		utils.Logger.Println("Received shutdown signal, stopping log processing...")
	case <-done:
		utils.Logger.Println("All logs processed successfully.")
	}

	err := internal.WriteLogsToJSON(processedLogs, outputFile)
	if err != nil {
		utils.Logger.Printf("Error writing to JSON: %v\n", err)
	}

	utils.Logger.Println("Log processing complete.")
}
