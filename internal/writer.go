package internal

import (
	"encoding/json"
	"os"
)

func WriteLogsToJSON(logs []*LogEntry, outputFile string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(logs)
}
