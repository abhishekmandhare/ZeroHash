package main

import (
	"context"
	"log"
	"os"

	"github.com/abhishekmandhare/zeroHash/cmd/startup"
	"github.com/abhishekmandhare/zeroHash/internal/config"
	"golang.org/x/sync/errgroup"
)

func main() {

	log.Println("Zero-Hash VWAP calculator")

	_, err := config.Load()
	if err != nil {
		log.Fatalf("Unable to load config file : %v", err)
		os.Exit(1)
	}


	ctx := context.Background()
	errGrp, gCtx := errgroup.WithContext(ctx)
	errGrp.Go(startup.RunAppServer(gCtx))
	errGrp.Go(startup.RunSignalListener(gCtx))
	if err := errGrp.Wait(); err != nil {
		log.Fatalf("Terminating with error: %s\n", err)
	}
}
