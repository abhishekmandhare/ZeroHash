package writer

import (
	"fmt"
	"log"
)

type ConsoleWriter struct {
}

func (c ConsoleWriter) WriteData(chanWriterData chan WriterData) {
	for data := range chanWriterData {
		fmt.Printf("Currency: %v, VWAP: %v\n", data.Currency, data.VWAP)
	}
	log.Println("Closing Writer")
}
