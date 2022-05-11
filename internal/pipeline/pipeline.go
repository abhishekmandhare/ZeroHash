package pipeline

import (
	"context"

	"github.com/abhishekmandhare/zeroHash/internal/models"
	"github.com/abhishekmandhare/zeroHash/internal/vwap"
)

type Pipeline struct {
	vwapChannels map[string]*vwap.Vwap
	ctx          context.Context
}

func NewPipeline(ctx context.Context, products []string) *Pipeline {
	outChannels := make(map[string]*vwap.Vwap)

	for _, product := range products {
		outChannels[product] = vwap.NewVwap()
		go outChannels[product].RunCalculator()
	}
	return &Pipeline{ctx: ctx, vwapChannels: outChannels}
}

func (p *Pipeline) Close() {
	for _, vwapch := range p.vwapChannels {
		vwapch.CloseChannel()
	}
}

func (p *Pipeline) SendTrade(t models.Trade) {

	if vwapCh, found := p.vwapChannels[t.Currency]; found {

		vwapCh.SendTradeToChannel(t)

	}
}
