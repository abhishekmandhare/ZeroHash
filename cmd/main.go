package main

import (
	"context"
	"log"

	"github.com/abhishekmandhare/zeroHash/cmd/startup"
	"golang.org/x/sync/errgroup"
)

func main() {

	log.Println("Zero-Hash VWAP calculator")
	ctx := context.Background()
	errGrp, gCtx := errgroup.WithContext(ctx)
	errGrp.Go(startup.RunAppServer(gCtx))
	errGrp.Go(startup.RunSignalListener(gCtx))
	if err := errGrp.Wait(); err != nil {
		log.Fatalf("Terminating with error: %s\n", err)
	}
}
