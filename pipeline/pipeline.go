package pipeline

import (
	"sync"
)

func PrimeFinder(done <-chan int, intStream <-chan int) <-chan int {
	primeStream := make(chan int)

	go func() {
		defer close(primeStream)
		for {
			select {
			case <-done:
				return
			case number := <-intStream:
				if IsPrime(number) {
					primeStream <- number
				}
			}
		}
	}()

	return primeStream
}

func FanIn[T any](done <-chan int, channels ...<-chan T) <-chan T {
	var wg = sync.WaitGroup{}
	fannedInStream := make(chan T)

	transfer := func(channel <-chan T) {
		defer wg.Done()
		for values := range channel {
			select {
			case <-done:
				return
			case fannedInStream <- values:
			}
		}
	}

	for _, channel := range channels {
		wg.Add(1)
		go transfer(channel)
	}

	go func() {
		wg.Wait()
		close(fannedInStream)
	}()

	return fannedInStream
}

func TakeSome[T any, K any](done <-chan K, stream <-chan T, number int) <-chan T {
	itemsTaken := make(chan T)

	go func() {
		defer close(itemsTaken)
		for i := 0; i < number; i++ {
			select {
			case <-done:
				return
			case itemsTaken <- <-stream:
			}
		}
	}()

	return itemsTaken
}

func IsPrime(number int) bool {
	for i := number - 1; i > 1; i-- {
		if number%i == 0 {
			return false
		}
	}
	return true
}
