package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/abhishekmandhare/zeroHash/internal/client/model"
	"github.com/abhishekmandhare/zeroHash/internal/models"
	"github.com/gorilla/websocket"
)

type Client struct {
	ProductIds []string
	ctx        context.Context
	Websocket  string
	connection *websocket.Conn
}

func NewClient(ctx context.Context, productIds []string, websocket string) *Client {
	return &Client{
		ProductIds: productIds,
		ctx:        ctx,
		Websocket:  websocket,
	}
}

func (client *Client) Subscribe() error {

	c, _, err := websocket.DefaultDialer.Dial(client.Websocket, nil)
	if err != nil {
		log.Fatalf("dial: %v", err)
	}

	client.connection = c

	// sub
	subMessage := model.CoinbaseSubscribeMsg{
		Type:       "subscribe",
		Channels:   []string{"matches"},
		ProductIDs: client.ProductIds,
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

	return nil
}

func (client *Client) Close() {
	client.connection.Close()
}

func (client *Client) Read() (*models.Trade, error) {

	m := &model.MatchMsg{}
	err := client.connection.ReadJSON(m)
	if err != nil {
		return nil, err
	}
	//log.Printf("received: %v", m)
	switch m.Type {
	case "error":
		return nil, fmt.Errorf("Received error from websocket")
	case "subscriptions":
		log.Printf("Received subscriptions message")
		return &models.Trade{}, nil
	case "match":
		price, err := strconv.ParseFloat(m.Price, 64)
		if err != nil {
			return nil, err
		}

		quantity, err := strconv.ParseFloat(m.Size, 64)
		if err != nil {
			return nil, err
		}

		trade := &models.Trade{
			Price:    price,
			Quantity: quantity,
			Currency: m.ProductID,
		}

		return trade, nil
	default:
		log.Printf("Received unhandled message: %v", m.Type)
		return &models.Trade{}, nil
	}

}