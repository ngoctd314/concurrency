# Applying modern go concurrency patterns to data pipelines

## A simple pipeline

To kick things off, we will implement a simple pair of producer and consumer. The producer goes over a list of words and sends them to a channel. While the consumer is receiving values from that channel and printing them to the console.

```go
func producer(words ...string) <- chan string {
    out := make(chan string)

    go func(){
        for _, word := range words {
            out <- word
        }
        close(out)
    }()

    return out
}

func sink(ch <- chan string ) {
    for v := range ch {
        fmt.Println(v)
    }
}
```
## Graceful Shutdown With Context

Especially in Go web development it's common to thread a context value through all of your long running functions, so that you can cancel those functions gracefully and perform cleanup if necessary.

```go
func producer(ctx context.Context, words ...string) <-chan string {
	out := make(chan string)

	go func() {
		defer close(out)
		for _, word := range words {
			time.Sleep(time.Second * 1)
			select {
			case <-ctx.Done():
				return
			case out <- word:
			}
		}
	}()

	return out
}

func sink(ctx context.Context, ch <-chan string) {
	for {
		select {
		case <-ctx.Done():
			log.Print(ctx.Err().Error())
			return
		case val, ok := <-ch:
			log.Print("val: ", val)
			if ok {
				log.Printf("sink: %s", val)
			} else {
				log.Print("done")
				return
			}
		}
	}
}
```

## Adding Parallelism with Fan-out and Fan-in

Going straight from producer to consumer isn't really a pipeline, so let's add the second stage that transforms all strings to lower case. We can pretty much copy/past that producer stage and add strings.ToLower.

Imagine that this second stage took 10 seconds to run, but our producer can go through the entire list of strings in a fraction of a second. With only a single Go routine processing the items, we'd be waiting for 30s to process the three strings. On a machine with multiple CPU cores, we can do better!

Let's create a bunch of Go routines that all run the same pipeline step. That way we should be able to run the same function in parallel and bring total execution time down to around 10 seconds - the duration of a single transformation function.

The idea there is to run a loop that spawns as many Go routines as we have CPU cores avaiable. In each loop iteration, we create a Go routine that runs the same pipeline step function. That step function returns a channel, which we appen to a variable taht will contain all channels thus creates. Finally, we merge all of those channels together and pass the resulting single channel to sink.