package merge

import (
	"log"
	"sync"

	"github.com/abhishekmandhare/zeroHash/internal/writer"
)

func Merge(sources ...chan writer.WriterData) chan writer.WriterData {
	destination := make(chan writer.WriterData)
	wg := sync.WaitGroup{}

	wg.Add(len(sources))

	for _, source := range sources {
		go func(chIn <-chan writer.WriterData) {
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
