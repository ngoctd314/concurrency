package example

// data pass through outbound as soon as it is received by ch elements
// input channel will communicate with fanIn's channel instead of directly
// communicating with the main routine as before
func fanIn(ch ...chan string) <-chan string {
	outbound := make(chan string)

	// handle fanIn
	handle := func(ch <-chan string) {
		for i := range ch {
			outbound <- i
		}
	}

	for _, v := range ch {
		go handle(v)
	}

	return outbound
}
