package split

import (
	"log"
	"sync"
)

type SplitType interface {
	GetKey() string
}

// Split will split the incoming channel into n number of outgoing channels where n is the length of the keys slice.
// Split happens based on keys of the SplitType.
func Split[T SplitType](source <-chan T, keys []string) map[string]chan T {
	destinations := make(map[string]chan T)

	for _, key := range keys {
		destinations[key] = make(chan T)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for s := range source {
			if vwapCh, found := destinations[s.GetKey()]; found {
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
