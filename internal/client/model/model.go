package model

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
