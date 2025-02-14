package utils

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// SetupSignalHandler creates a context that is canceled when an interrupt signal is received
func SetupSignalHandler() context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalChan
		cancel() // Cancel the context when Ctrl+C is pressed
	}()

	return ctx
}
