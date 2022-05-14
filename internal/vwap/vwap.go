package vwap

import (
	"log"

	"github.com/abhishekmandhare/zeroHash/internal/models"
	"github.com/abhishekmandhare/zeroHash/internal/queue"
	"github.com/abhishekmandhare/zeroHash/internal/writer"
)

type Vwap struct {
	vwapWindow      *queue.Queue[models.Trade]
	chanTradeIn     <-chan models.Trade
	chanTradeOut    chan<- writer.WriterData
	vwapSum         float64
	vwapQuantitySum float64
	windowSize      int
	}

func NewVwap(vwapWindowSize int, chanTradeIn <-chan models.Trade, chanTradeOut chan<- writer.WriterData) *Vwap {

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
			v.chanTradeOut <- writer.WriterData{Currency: newTrade.Currency, VWAP: vwap}
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

	var vwap float64
	vwap = v.vwapSum / v.vwapQuantitySum

	return vwap
}
