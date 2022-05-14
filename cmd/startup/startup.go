package startup

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/abhishekmandhare/zeroHash/internal/client"
	"github.com/abhishekmandhare/zeroHash/internal/config"
	"github.com/abhishekmandhare/zeroHash/internal/models"
	"github.com/abhishekmandhare/zeroHash/internal/pipeline"
)

func RunAppClient(ctx context.Context, config *config.AppConfiguration, tradeChannel chan<- models.Trade) func() error {
	return func() error {
		cl := client.NewClient(ctx, config.Spec.Products, config.Spec.Websocket)
		if err := cl.Subscribe(); err != nil {
			return err
		}

		defer close(tradeChannel)
		for {
			select {
			case <-ctx.Done():
				log.Println("App Client terminated by upstream")
				cl.Close()

				return nil
			default:
				trade, err := cl.Read()
				if err != nil {
					return err
				}
				tradeChannel <- *trade
			}
		}
	}
}

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
