package startup

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/abhishekmandhare/zeroHash/internal/app/models"
	"github.com/abhishekmandhare/zeroHash/internal/app/pipeline"
	"github.com/abhishekmandhare/zeroHash/internal/client"
	"github.com/abhishekmandhare/zeroHash/internal/config"
)

// RunAppClient runs the websocket client and terminates on context closure or error.
func RunAppClient(ctx context.Context, config *config.AppConfiguration, tradeChannel chan<- models.Trade) func() error {
	return func() error {
		client := client.NewClient(ctx, config.Spec.Products, config.Spec.Websocket)
		if err := client.Subscribe(); err != nil {
			return err
		}

		defer close(tradeChannel)
		for {
			select {
			case <-ctx.Done():
				log.Println("App Client terminated by upstream")
				client.Close()

				return nil
			default:
				trade, err := client.Read()
				if err != nil {
					return err
				}
				if trade != nil {
					tradeChannel <- *trade
				}
			}
		}
	}
}

// RunPipeline runs the pipeline which connects channels into calculator and streams and terminates on context closure or error. 
func RunPipeline(ctx context.Context, config *config.AppConfiguration, tradeChannel <-chan models.Trade) func() error {
	return func() error {
		donePipe := make(chan interface{})
		go func() {
			defer func() { donePipe <- 1 }()
			pipeline.Connect(tradeChannel, config.Spec.Products, config.Spec.VwapWindowSize)
		}()
		for {
			select {
			case <-ctx.Done():
				log.Println("Pipeline terminated by upstream")
				return nil
			case <-donePipe:
				log.Println("Pipeline closed.")
				return nil
			}
		}
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
