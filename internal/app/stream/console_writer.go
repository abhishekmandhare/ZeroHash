package stream

import (
	"fmt"
	"log"
)

type ConsoleStreamer struct {
}

func (c ConsoleStreamer) StreamData(chanWriterData chan StreamData) {
	for data := range chanWriterData {
		fmt.Printf("Currency: %v, VWAP: %v\n", data.Currency, data.VWAP)
	}
	log.Println("Closing Writer")
}
