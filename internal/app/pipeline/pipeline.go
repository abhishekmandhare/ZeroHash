package pipeline

import (
	"log"
	"sync"

	"github.com/abhishekmandhare/zeroHash/internal/app/models"
	"github.com/abhishekmandhare/zeroHash/internal/app/stream"
	"github.com/abhishekmandhare/zeroHash/internal/app/vwap"
	"github.com/abhishekmandhare/zeroHash/internal/arch/merge"
	"github.com/abhishekmandhare/zeroHash/internal/arch/split"
)

func Connect(sourceTrade <-chan models.Trade, products []string, vwapWindowSize int) {
	vwapChannels := make(map[string]*vwap.Vwap)
	outChannels := split.Split(sourceTrade, products)
	chanWriters := make([]chan stream.StreamData, 0)

	for product, outChannel := range outChannels {
		chanWriter := make(chan stream.StreamData)
		vwapChannels[product] = vwap.NewVwap(vwapWindowSize, outChannel, chanWriter)
		vwapChannels[product].RunCalculator()
		chanWriters = append(chanWriters, chanWriter)
	}

	mergedChanWriter := merge.Merge(chanWriters...)
	writeStreamer := stream.NewStreamer()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		writeStreamer.StreamData(mergedChanWriter)
	}()

	wg.Wait()
	log.Println("Closing connect")
}
