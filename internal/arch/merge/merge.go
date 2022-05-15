package merge

import (
	"log"
	"sync"
)

func Merge[T any](sources ...chan T) chan T {
	destination := make(chan T, len(sources))
	wg := sync.WaitGroup{}

	wg.Add(len(sources))

	for _, source := range sources {
		go func(chIn <-chan T) {
			defer wg.Done()

			for n := range chIn {
				destination <- n
			}
		}(source)
	}

	go func() {
		wg.Wait()
		log.Println("Merge: Closing destination channels.")
		close(destination)
	}()

	return destination
}
