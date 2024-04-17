package main

import (
	"GoConcurrency/generators"
	"GoConcurrency/pipeline"
	"fmt"
	"math/rand"
	"runtime"
)

func main() {
	done := make(chan int)
	defer close(done)

	// GENERATE
	randomNumberFetcher := func() int { return rand.Intn(500000000) }
	streamGenerator := generators.RepeatFunc(done, randomNumberFetcher)

	// DIVIDE SEARCH FOR PRIMES
	cpuCount := runtime.NumCPU()
	primesStreamChannels := make([]<-chan int, cpuCount)
	for i := 0; i < cpuCount; i++ {
		primesStreamChannels[i] = pipeline.PrimeFinder(done, streamGenerator)
	}

	// UNIFY THE RESULTS
	fanResultStream := pipeline.FanIn(done, primesStreamChannels...)

	// TAKE ONLY WHAT IS NECESSARY (10 NUMBERS)
	for randNumber := range pipeline.TakeSome(done, fanResultStream, 10) {
		fmt.Println(randNumber)
	}
}
