package vwap

import (
	"log"

	"github.com/abhishekmandhare/zeroHash/internal/app/models"
	"github.com/abhishekmandhare/zeroHash/internal/app/stream"
	"github.com/abhishekmandhare/zeroHash/internal/arch/queue"
)

type Vwap struct {
	vwapWindow      *queue.Queue[models.Trade]
	chanTradeIn     <-chan models.Trade
	chanTradeOut    chan<- stream.StreamData
	vwapSum         float64
	vwapQuantitySum float64
	windowSize      int
}

func NewVwap(vwapWindowSize int, chanTradeIn <-chan models.Trade, chanTradeOut chan<- stream.StreamData) *Vwap {

	return &Vwap{
		vwapWindow:   queue.NewQueue[models.Trade](),
		chanTradeIn:  chanTradeIn,
		chanTradeOut: chanTradeOut,
		windowSize:   vwapWindowSize,
	}
}

func (v *Vwap) RunCalculator() {
	go func() {
		defer close(v.chanTradeOut)
		for newTrade := range v.chanTradeIn {
			vwap := v.calculate(newTrade)
			v.chanTradeOut <- stream.StreamData{Currency: newTrade.Currency, VWAP: vwap}
		}
	}()
}

func (v *Vwap) calculate(newTrade models.Trade) float64 {
	if v.vwapWindow.Len() >= v.windowSize {
		removedTrade, err := v.vwapWindow.Pop()
		if err != nil {
			log.Fatalf("Error running VWAP calculator: %v", err)
			return 0
		}

		v.vwapSum -= removedTrade.Price * removedTrade.Quantity
		v.vwapQuantitySum -= removedTrade.Quantity
	}
	v.vwapWindow.Push(newTrade)
	v.vwapSum += newTrade.Price * newTrade.Quantity
	v.vwapQuantitySum += newTrade.Quantity

	return v.vwapSum / v.vwapQuantitySum
}
