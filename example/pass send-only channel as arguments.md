# Pass send-only channels as arguments

Same as the receive-only channel, the values of two arguments of the sumSquares function call are requested concurrently. Different to the last example, the longTimeRequest function takes a send-only channel as parameter instead of returning a receive-only channel result.

```go
func longTimeRequest(r chan<- int32) {
	// Simulate a workload
	time.Sleep(time.Second * 3)
	r <- rand.Int31n(100)
}

func sumSquare(a, b int32) int32 {
	return a*b + b*b
}

func main() {
	rand.Seed(time.Now().UnixNano())
	ra, rb := make(chan int32), make(chan int32)

	go longTimeRequest(ra)
	go longTimeRequest(rb)

	fmt.Println(sumSquare(<-ra, <-rb))
}
```

In fact, for the above specified example, we don't need two channels to transfer results. Using one channel is okay.

```go
results := make(chan int32, 2)

go longTimeRequest(results)
go longTimeRequest(results)

fmt.Println(sumSquares(<-results, <-results))
```