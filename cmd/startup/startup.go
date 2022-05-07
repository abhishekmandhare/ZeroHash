package startup

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func RunAppServer(ctx context.Context) func() error {
	return func() error { 

		return nil
	}
}

// RunSignalListener returns a function that starts a listener for system signals.
func RunSignalListener(ctx context.Context) func() error {
	return func() error {
		sigChan := make(chan os.Signal, 1)
		defer close(sigChan)

		signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

		select {
		case <-sigChan:
			return fmt.Errorf("Terminated by SIGTERM")
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
