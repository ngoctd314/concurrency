package example

// Exec ...
func Exec() {
	// execBoring()
	// execFoo()
	// execUpdatePosition()
	// execFibonaciGenerator()
	orDone := func(done, c <-chan any) <-chan any {
		valStream := make(chan any)
		go func() {
			defer close(valStream)
			for {
				select {
				case <-done:
					return
				case v, ok := <-c:
					if !ok {
						return
					}
					select {
					case valStream <- v:
					case <-done:
					}
				case <-done:
				}
			}
		}()

		return valStream
	}

	tee := func(done <-chan interface{}, in <-chan interface{}) (<-chan any, <-chan any) {
		out1 := make(chan any)
		out2 := make(chan any)
		go func() {
			defer close(out1)
			defer close(out2)
			for val := range orDone(done, in) {
				// we will want to use local versions of out1 and out2, so we shadow these variables
				var out1, out2 = out1, out2
				// we're going to use one select statement so that writes to out1 and out2
				// don't block each other. To ensure both are written to, we'll perform two iterations
				// of the select statement: one for each outbound channel
				for i := 0; i < 2; i++ {
					select {
					case <-done:
					case out1 <- val:
						// once we've written to a channel, we set its shadowed copy to nil
						// so that further writes will block and the other channel may continue
						out1 = nil
					case out2 <- val:
						// once we've written to a channel, we set its shadowned copy to nil
						// so the further writes will block and the other channel may continue
						out2 = nil
					}
				}
			}
		}()
		return out1, out2
	}
	_ = tee

	bridge := func(done <-chan any, chanStream <-chan <-chan any) <-chan any {
		// this channel will return all values from bridge
		valStream := make(chan any)
		go func() {
			defer close(valStream)
			// This loop is responsible for pulling channels off of chanStream and providing
			// them to a nested loop for use
			for {
				var stream <-chan any
				select {
				case mayBeStream, ok := <-chanStream:
					if ok == false {
						return
					}
					stream = mayBeStream
				case <-done:
					return
				}
				// This loop is responsible for reading values off the channel it has been given and
				// repeating those values onto valStream. When the stream we're currently looping
				// over is closed, we break out of the loop performing the reads from this channel,
				// and continue with the next iteration of the loop, selecting channels to read from.
				// This provides us with an unbroken stream of values.
				for val := range stream {
					select {
					case valStream <- val:
					case <-done:
					}
				}

			}
		}()
		return valStream
	}
	_ = bridge
}
