package vwap

import (
	"fmt"

	"github.com/abhishekmandhare/zeroHash/internal/models"
	"github.com/abhishekmandhare/zeroHash/internal/queue"
)

type Vwap struct {
	vwapWindow      *queue.Queue[models.Trade]
	chanTrade       chan models.Trade
	vwapSum         float64
	vwapQuantitySum float64
}

const WindowSize int = 200

func NewVwap() *Vwap {
	chanTrade := make(chan models.Trade)
	return &Vwap{
		vwapWindow: queue.NewQueue[models.Trade](),
		chanTrade:  chanTrade,
	}
}

func (v *Vwap) SendTradeToChannel(trade models.Trade) {
	v.chanTrade <- trade
}

func (v *Vwap) CloseChannel() {
	close(v.chanTrade)
}

func (v *Vwap) RunCalculator() error {

	for t := range v.chanTrade {
		var removeElem models.Trade

		if v.vwapWindow.Len() < WindowSize {
			v.vwapWindow.Push(t)
		} else {
			var err error
			removeElem, err = v.vwapWindow.Pop()
			if err != nil {
				return err
			}
			v.vwapWindow.Push(t)

		}
		vwap := v.Calculate(&t, &removeElem)
		if vwap != 0 {
			fmt.Printf("Currency %v, Vwap : %v\n", t.Currency, vwap)
		}

	}
	return nil
}

func (v *Vwap) Calculate(addElement *models.Trade, removeElement *models.Trade) float64 {
	v.vwapSum += addElement.Price * addElement.Quantity
	v.vwapQuantitySum += addElement.Quantity

	if removeElement != nil {
		v.vwapSum -= removeElement.Price * removeElement.Quantity
		v.vwapQuantitySum -= removeElement.Quantity
	}
	var vwap float64
	if v.vwapWindow.Len() == WindowSize {
		vwap = v.vwapSum / v.vwapQuantitySum
	}
	return vwap
}
