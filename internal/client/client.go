package client

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/gorilla/websocket"

	"github.com/abhishekmandhare/zeroHash/internal/app/models"
	"github.com/abhishekmandhare/zeroHash/internal/client/model"
)

// Client struct is used to read "Matches" data stream from the given websocket in JSON format.
type Client struct {
	ProductIds []string
	ctx        context.Context
	Websocket  string
	connection *websocket.Conn
}

// NewClient creates and returns a Client struct.
func NewClient(ctx context.Context, productIds []string, websocket string) *Client {
	return &Client{
		ProductIds: productIds,
		ctx:        ctx,
		Websocket:  websocket,
	}
}

// Subscribe sends a "matches" subscription message to websocket in order to receive match data.
func (client *Client) Subscribe() error {

	connection, _, err := websocket.DefaultDialer.Dial(client.Websocket, nil)
	if err != nil {
		log.Fatalf("dial: %v", err)
	}

	client.connection = connection

	subMessage := model.CoinbaseSubscribeMsg{
		Type:       "subscribe",
		Channels:   []string{"matches"},
		ProductIDs: client.ProductIds,
	}

	err = client.connection.WriteJSON(subMessage)
	if err != nil {
		log.Fatalf("write err : %v", err)
		return err
	}

	return nil
}

// Close closes the websocket stream.
func (client *Client) Close() {
	client.connection.Close()
}

// Read reads the "match" messages from websocket. In case of errors, the error is returned.
// All other messages are ignored.
func (client *Client) Read() (*models.Trade, error) {

	m := &model.MatchMsg{}
	err := client.connection.ReadJSON(m)
	if err != nil {
		return nil, err
	}

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
		return nil, nil
	}
}
