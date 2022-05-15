package pipeline

import (
	"log"

	"github.com/abhishekmandhare/zeroHash/internal/app/models"
	"github.com/abhishekmandhare/zeroHash/internal/app/stream"
	"github.com/abhishekmandhare/zeroHash/internal/app/vwap"
	"github.com/abhishekmandhare/zeroHash/internal/arch/merge"
	"github.com/abhishekmandhare/zeroHash/internal/arch/split"
)

// Connect splits incoming sourceTrade channel into n channels where n is number of products.
// It then feeds each channel to a seperate vwap calculator routine.
// It also merges each of the vwap calculator output channels into a single channel which is fed to stream.
func Connect(sourceTradeChan <-chan models.Trade, products []string, vwapCalculationWindowSize int) {
	prodVwapMap := map[string]*vwap.Vwap{}
	splitTradeChans := split.Split(sourceTradeChan, products)
	streamerChans := make([]chan stream.StreamData, 0)

	for product, tradeChan := range splitTradeChans {
		streamerChan := make(chan stream.StreamData)
		prodVwapMap[product] = vwap.NewVwap(vwapCalculationWindowSize, tradeChan, streamerChan)
		prodVwapMap[product].RunCalculator()
		streamerChans = append(streamerChans, streamerChan)
	}

	mergedStreamerChan := merge.Merge(streamerChans...)
	writeStreamer := stream.NewStreamer()

	writeStreamer.StreamData(mergedStreamerChan)
	log.Println("Closing connect")
}
