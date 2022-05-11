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
		cl.Subscribe()

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
		pipe := pipeline.NewPipeline(ctx, config.Spec.Products)

		for {
			select {
			case <-ctx.Done():
				log.Println("Pipeline terminated by upstream")
				pipe.Close()

				return nil
			case t, ok := <-tradeChannel:
				if !ok {
					log.Println("Closing Pipeline, trade channel is closed")
					pipe.Close()
					return nil
				}
				pipe.SendTrade(t)
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
