# Fan-In Fan-Out

Fan-In Fan-Out techniques which are used to multiplex and demultiplex data in Go. Fan-In refers to a techniques in which you join data from multiple inputs into a single entity. On other hand, Fan-Out means to divide the data from a single source into multiple smaller chunks.

Sometimes, stages in your pipeline can be particularly computationally expensive. When this happens, upstream stages in your pipeline can become blocked while waiting for your expensive stages to complete. Not only that, but the pipeline itself can take a long time to execute as a whole. How can we address this?

Fan-out is a term to describe the process of starting multiple goroutines to handle input from the pipeline, and fan-in is a term to describe the process of combining multiple results into one channel.

So what makes a stage of a pipeline suited for utilizing this pattern? You might consider fanning out one of your stages if both of the following apply:

- It doesn't rely on values that the stage had calculated before.
- It takes to long time to run.

The property of order-independence is important because you have no guarantee in what order concurrent copies of your stage will run, nor in what order they will return.

This is a relatively simple example, so we only have two stages: random number generation and prime serving.

The process of fanning out a stage in a pipeline is extraordinarily easy. All we have to do is start multiple versions of that stage

```go
numFinder := runtime.NumCPU()
finders := make([]<- chan int, numFinders)

for i := 0 ; i < numFinders; ++ {
    finders[i] = primeFinder(done, randIntStream)
}

ch := make([]<-chan int, 0)
for i := 0; i < runtime.NumCPU(); i++ {
    ch = append(ch, fn(i))
}
for _, v := range ch {
    for k := range v {
        fmt.Println(k)
    }
}
```

We still have a problem though: new that we have four goroutines, we also have four channels, but our range over primes is only expecting one channel. This brings us to the fan-in portion of the pattern.

As we discussed earlier, fanning in means multiplexing or joining together multiple streams of data into a single stream.

```go
fanIn := func(done <- chan any, channels ...<- chan any) <-chan any{
    var (
        wg = sync.WaitGroup{}
        multiplexedStream = make(chan any)
    )

    multiplex := func(c <- chan any) {
        defer wg.Done()
        for i := range c {
            select {
                case <- done:
                    return
                case multiplexedStream <- i:
            }
        }
    }

    for _, c := range channels {
        wg.Add(1)
        go multiplex(c)
    }

    go func(){
        wg.Wait()
        close(multiplexedStream)
    }()

    return multiplexedStream
}
```