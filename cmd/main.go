package main

import (
	"context"
	"log"
	"os"

	"github.com/abhishekmandhare/zeroHash/cmd/startup"
	"github.com/abhishekmandhare/zeroHash/internal/app/models"
	"github.com/abhishekmandhare/zeroHash/internal/config"
	"golang.org/x/sync/errgroup"
)

func main() {

	log.Println("Zero-Hash VWAP calculator")

	config, err := config.Load()
	if err != nil {
		log.Fatalf("Unable to load config file : %v", err)
		os.Exit(1)
	}

	ctx := context.Background()
	errGrp, gCtx := errgroup.WithContext(ctx)
	tradeChannel := make(chan models.Trade, len(config.Spec.Products))
	errGrp.Go(startup.RunAppClient(gCtx, config, tradeChannel))
	errGrp.Go(startup.RunPipeline(ctx, config, tradeChannel))
	errGrp.Go(startup.RunSignalListener(gCtx))

	if err := errGrp.Wait(); err != nil {
		log.Fatalf("Terminating with error: %s\n", err)
	}
}
