package split

import (
	"log"
	"sync"

	"github.com/abhishekmandhare/zeroHash/internal/app/models"
)

func Split(source <-chan models.Trade, currencies []string) map[string]chan models.Trade {
	destinations := make(map[string]chan models.Trade)

	for _, currency := range currencies {
		destinations[currency] = make(chan models.Trade)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for s := range source {
			if vwapCh, found := destinations[s.Currency]; found {
				vwapCh <- s
			}
		}
	}()

	go func() {
		wg.Wait()
		log.Println("Split: Closing destination channels.")
		for _, destinationChannel := range destinations {
			close(destinationChannel)
		}
	}()

	return destinations
}
