package startup

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/abhishekmandhare/zeroHash/internal/config"
	"github.com/abhishekmandhare/zeroHash/internal/models"
	"github.com/abhishekmandhare/zeroHash/internal/vwap"
	"github.com/gorilla/websocket"
)

type MatchMsg struct {
	Type          string `json:"type"`
	TradeID       int    `json:"trade_id"`
	MarkerOrderId string `json:"marker_order_id"`
	TakerOrderId  string `json:"taker_order_id"`
	Side          string `json:"side"`
	Size          string `json:"size"`
	Price         string `json:"price"`
	ProductID     string `json:"product_id"`
	Sequence      int    `json:"sequence"`
	Time          string `json:"time"`
}

type CoinbaseSubscribeMsg struct {
	Type       string   `json:"type"`
	ProductIDs []string `json:"product_ids"`
	Channels   []string `json:"channels"`
}

func RunAppServer(ctx context.Context, config *config.AppConfiguration) func() error {
	return func() error {

		c, _, err := websocket.DefaultDialer.Dial(config.Spec.Websocket, nil)
		if err != nil {
			log.Fatalf("dial: %v", err)
		}
		defer c.Close()

		// sub
		subMessage := CoinbaseSubscribeMsg{
			Type:       "subscribe",
			Channels:   []string{"matches"},
			ProductIDs: []string{"ETH-USD"},
		}

		subEvent, err := json.Marshal(subMessage)
		if err != nil {
			log.Fatalf("Unable to marshal : %v", err)
			return err
		}

		err = c.WriteMessage(websocket.TextMessage, subEvent)
		if err != nil {
			log.Fatalf("write err : %v", err)
			return err
		}

		done := make(chan interface{})

		tradeChan := make(chan models.Trade)
		vwapCalc := vwap.NewVwap(tradeChan)
		go vwapCalc.RunCalculator()

		go func() {
			defer close(done)
			for {
				m := &MatchMsg{}
				err := c.ReadJSON(m)
				if err != nil {
					log.Printf("read err: %v", err)
					return
				}

				price, err := strconv.ParseFloat(m.Price, 64)
				if err != nil {
					continue
				}

				quantity, err := strconv.ParseFloat(m.Size, 64)
				if err != nil {
					continue
				}

				trade := &models.Trade{
					Price:    price,
					Quantity: quantity,
				}

				tradeChan <- *trade

				log.Printf("received: %v", m)
			}
		}()

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return nil
			case <-ctx.Done():
				log.Println("Terminated by upstream")
				return nil

			}
		}
	}
}

// RunSignalListener returns a function that starts a listener for system signals.
func RunSignalListener(ctx context.Context) func() error {
	return func() error {
		sigChan := make(chan os.Signal, 1)
		defer close(sigChan)

		signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

		select {
		case <-sigChan:
			return fmt.Errorf("Terminated by SIGTERM")
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
