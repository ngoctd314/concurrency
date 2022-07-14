# The first response wins

Sometimes, a piece of data can be received from several sources to avoid high latencies. For a lot of factors, the response durations of these sources may vary much. Even for a specified source, its response durations are also not constant. To make the response duration as short as possible, we can send a request to every source in a separated goroutine. Only the first response will be used, other slower ones will be discarded.

Note, if there are N sources, the capacity of the communication channel must be at least N-1, to avoid the goroutines corresponding the discarded responses being blocked for ever.

```go
func source(c chan<- int32) {
	ra, rb := rand.Int31(), rand.Intn(3)+1
	// Sleep 1s/2s/3s
	time.Sleep(time.Duration(rb) * time.Second)

	c <- ra
}

func main() {
	rand.Seed(time.Now().UnixNano())

	startTime := time.Now()
	// c must be a buffered channel
	c := make(chan int32, 5)
	for i := 0; i < cap(c); i++ {
		go source(c)
	}
	// Only the first response will be used
	rnd := <-c
	fmt.Println(time.Since(startTime))
	fmt.Println(rnd)
}
```