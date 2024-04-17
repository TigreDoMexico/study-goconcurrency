package generators

// RepeatFunc get all data that came from fn and store in channel stream until done is complete
func RepeatFunc[T any, K any](done <-chan K, functionToRepeat func() T) <-chan T {
	stream := make(chan T)

	go func() {
		defer close(stream)

		// For Select Pattern
		for {
			select {
			case <-done:
				return
			case stream <- functionToRepeat():
			}
		}
	}()

	return stream
}
