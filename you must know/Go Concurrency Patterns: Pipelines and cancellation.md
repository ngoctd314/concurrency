# Go Concurrency Patterns: Pipelines and cancellation

## Introduction

Go's concurrency primitives make it easy to construct streaming data pipelines that make efficient use of I/O and multiple CPUs. 

## What is a pipeline?

There's no formal definition of a pipeline in Go; it's just one of many kinds of concurrent programs. Informally, a pipeline is a series of stages connected by channels, where each stage is a group of goroutines running the same function. In each stage, goroutines

- receive values from upstream via inbound channels
- perform some function on that data, usually producing new values
- send the value downstreawm via outbound channels

Each stage has any number of inbound and outbound channels, except the first and last stages, which have only outbound and inbound, respectively. The first stage is sometimes called the source or producer; the last stage, the sink or consumer.

## Squaring numbers

The first stage, gen, is a function that converts a list of integers to a channel that emits the integers in the list. The gen function starts a goroutine that sends the integers on the channel and closes the channel when all the values have been sent:

```go
func gen(nums ...int) <- chan int {
    out := make(chan int)

    go func(){
        defer close(out)

        for _, num := range nums {
            out <- num
        }
    }()

    return out
}
```
The second stage, sq, receives integers from a channel and returns a channel that emits the square of each received integer. After the inbound channel is closed and this stage has sent all the values downstream, it closes the outbound channel:

```go
func sq(in <- chan int) <- chan int {
    out := make(chan int)

    go func() {
        defer close(out)
        for v := range in {
            out <- in*in
        }
    }()

    return out
}
```

The main function sets up the pipeline and runs the final stage: it receives values from the second stage and prints each one, until the channel is closed:

```go
func main(){
    // set up the  pipeline
    c := gen(2,3)
    out := sq(c)

    fmt.Println(<- out, <- out)
}
```

## Fan-out, Fan-in

Multiple functions can read from the same channel until that channel is closed; this is called fan-out. This provides a way to distribute work amongst a group of workers to parallelize CPU use and I/O.

A function can read from multiple inputs and proceed until all race are closed by multiplexing the input channels onto a single channel that's closed when all inputs are closes. This is called fan-in.

We can change our pipeline to run two instances of sq, each reading from the same input channel. We introduce a new function, merge, to fan in the results:

```go
func main() {
    in := gen(2,3)

    // Distribute the sq work across two goroutines that both read from in.
    c1 := sq(in)
    c2 := sq(in)

    // consume the merged output from c1 and c2
    for n := range merge(c1, c2) {
        fmt.Println(n)
    }
}
```
The merge function converts a list of channels to a single channel by starting a goroutine for each inbound channel that copies the values to sole outbound channel.

Sends on closed channel panic, so it's important to ensure all sends are done before calling close. The sync.WaitGroup type provides a simple way to arrange this synchronization

```go
func merge (cs ...<- chan int) <- chan int {
    var wg sync.WaitGroup
    out := make(chan int)

    output := func(c chan int) {
        defer wg.Done()
        for v := range c {
            out <- v
        }
    }

    wg.Add(len(cs)) 
    for _, v := range cs {
        go output(v)
    }

    // start a goroutine to close out once all the ouput goroutines are done 
    // This must start after the wg.Add call
    go func(){
        wg.Wait()
        close(out)
    }()

    return out
}
```
## Stopping short

There is a pattern to our pipeline functions:

- Stage close their outbound channels when all the send operations are done
- Stage keep receiving values from inbound channels until those channels are close.

This pattern allows each receiving state to be written as a range loop and ensures that all goroutines exit once all values have been successfully send downstream.

But in real pipelines, stages don't always receive all the inbound values. Sometimes this is by design: the receiver may only need a subset of values to make progress. More often, a stage exists early because an inbound value represents an error in an earlier stage. In either case the receiver should not have to wait for the remaining values to arrive, and we want earlier stages to stop producing values that later stages don't need.

In our example pipeline, if a stage fails to consume all the inbound values, the goroutines attempting to send those values will block indefinitely:

```go
// Consume the first value from the output
out := merge(c1, c2)
fmt.Println(<-out)
return
// Since we didn't receive the second value from out
// one of the output goroutines is hung attempting to send it
```
This is a resource leak: goroutines consume memory and runtime resources, and heap references in goroutine stacks keep data from being garbage collected. Goroutines are not garbage collected; they must exit on their own.

We need to arrange for the upstream stages of our pipeline to exit even then the downstream stages fail to receive all the inbound values. One way to do this is to change the outbound channels to have a buffer. A buffer can hold a fixed numbers of values; send operations complete immediately if there's room in the buffer: 

```go
c := make(chan int, 2) // buffer size 2
c <- 1 // succeeds immediately
c <- 2 // succeeds immediately
c <- 3 // blocks until another goroutine does <- c and receive 1
```

When the number of values to be sent is known at channel creation time, a buffer can simplify the code.

```go
func gen(nums ...int) <-chan int {
    out := make(chan int, len(nums))
    for _, n := range nums {
        out <- n
    }
    close(out)

    return out
}
```

This fixes the blocked goroutine in this program, this is bad code.

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