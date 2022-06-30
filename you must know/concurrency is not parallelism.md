# Concurrency is not parallelism

[Via](https://go.dev/blog/waza-talk)

When people hear the work concurrency they often think of parallelism, a related but quite distinct concept, in programming, concurrency is the composition of independently executing processes, while parallelism is the simultaneous execution of (possibly related) computations. Concurrency is about dealing with lots of things at once. Paralelism is about doing lots of thing at once. Not the same, but related. Concurrency is about structure, parallelism is about execution. Concurrency provides a way to structure a solution to solve problem that may (but not necessarily) be parallelizable.

## Concurrency plus communication

Concurrency is a way to structure a program by breaking it into pieces that can be executed independently.

Communication is the means to coordinate the independent executions

## Lession

There are many ways to break the processing down

That's concurrent design

Once we have breakdown, parallelization can fall out and correctness is easy.

## Explicit cancellation

When main decides to exit without receiving all the values from out, it must tell the goroutines in the upstream stages to abandon the values they're trying to send. It does so by sending values on a channel called done.

It sends two values since there are potentially two blocked senders

```go
func main() {
    in := gen(2, 3)
    // Distribute the sq work arcoss two goroutines that both read from in
    c1 := sq(in)
    c2 := sq(in)

    // consume the first value 
    done := make(chan struct{}, 2)
    out := merge(done, c1, c2)
    fmt.Println(<-out)

    // Tell the remaining senders we're leaving
    done <- struct{}{}
    done <- struct{}{}
}

func merge(done <- chan struct{}, cs ...<- chan int) <- chan int {
    var wg sync.WaitGroup
    out := make(chan int)

    output := func(ch chan int) {
        defer wg.Done()
        for v := range ch {
            select {
                case out <- v:
                case <- done:
                    return
            }
        }
    }

    return out
}
```

This approach has a problem: each downstream receiver needs to know the number of petentially blocked upstream senders and arrange to signal those senders on early return. Keeping track of these counts is tedious and error-prone.

We need a way to tell an unknown and unbounded number of goroutines to stop sending their values downstream. In Go, we can do this by closing a channel, because a receive operation on a closed channel can always proceed immediately, yieldying the element type's zero value.

This means that main can unblock all the senders simply by closing the done channel. This close is effectively a broadcast signal to the senders. We extend each of our pipeline functions to accept done as a parameter and arrange for the close to happen via a defer statement, so that all return paths from main will signal the pipeline stages to exit.

```go
func main() {
    // set up a done channel that's shared by the whole pipeline,
    // and close that channel when this pipeline exits, as a signal for all
    // the goroutines we started to exit
    //
    // 1-to-N notification
    done := make(chan struct{})
    defer close(done)

    in := gen(done, 2, 3)

    // Distribute the sq work across two goroutines that both read from in
    // fan-out
    c1 := sq(done, in)
    c2 := sq(done, in)

    // Consume the first value from output
    out := merge(done, c1, c2)
    fmt.Println(<-out) // 4 or 9

    // done will be closed by the deferred call
}
```

Each of out pipeline stages is not free to return as soon as done is closed. The output routine in merge can return without draning its inbound channel, since it knowns the upstream sender, sq will stop attempting to send when done is closed. output ensures wg.Done() is called on all return paths via a defer statement:

```go
func merge(done <- chan struct{}, cs ...<- chan int) <- chan int {
    var wg sync.WaitGroup
    out := make(chan int)

    output := func (c <- chan int)  {
        defer wg.Done()
        for v := range c {
            select {
                case out <- v:
                case <-done:
                    return
            }
        }
    }
}

func sq(done <- chan struct{}, in <- chan int) <- chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for v := range in {
            select {
                case out <- v*v:
                case <- done:
                    return
            }
        }
    }()

    return out
}
```
**Here are the guidelines for pipeline construction**
- stages close their outbound channels when all the send operations are done
- stages keep receiving values from inbound channels until those channel are closed or the senders are unblocked.

## Digesting a tree