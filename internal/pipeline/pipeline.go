package pipeline

import (
	"log"
	"sync"

	"github.com/abhishekmandhare/zeroHash/internal/merge"
	"github.com/abhishekmandhare/zeroHash/internal/models"
	"github.com/abhishekmandhare/zeroHash/internal/split"
	"github.com/abhishekmandhare/zeroHash/internal/vwap"
	"github.com/abhishekmandhare/zeroHash/internal/writer"
)

func Connect(sourceTrade <-chan models.Trade, products []string, vwapWindowSize int) {
	vwapChannels := make(map[string]*vwap.Vwap)
	outChannels := split.Split(sourceTrade, products)
	chanWriters := make([]chan writer.WriterData, 0)

	for product, outChannel := range outChannels {
		chanWriter := make(chan writer.WriterData)
		vwapChannels[product] = vwap.NewVwap(vwapWindowSize, outChannel, chanWriter)
		vwapChannels[product].RunCalculator()
		chanWriters = append(chanWriters, chanWriter)
	}

	mergedChanWriter := merge.Merge(chanWriters...)
	writeStreamer := writer.NewWriter()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		writeStreamer.WriteData(mergedChanWriter)
	}()

	wg.Wait()
	log.Println("Closing connect")
}
