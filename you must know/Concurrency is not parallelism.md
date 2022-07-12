# Concurrency is not parallelism

[Via](https://go.dev/blog/waza-talk)

When people hear the work concurrency they often think of parallelism, a related but quite distinct concept, in programming, concurrency is the composition of independently executing processes, while parallelism is the simultaneous execution of (possibly related) computations. Concurrency is about dealing with lots of things at once. Paralelism is about doing lots of thing at once.

## Go supports concurrency

Go provides:

- concurrent execution (goroutines)
- synchronization and messaging (channels)
- multi-way concurrent control(select)

## Concurrency

Programming as the composition of independently execution processes.

## Paralleslism

Programming as the simultaneous execution of (possible related) computations

## Concurrency vs parallelism

Concurency is about dealing with lots of things at once.
Parallelism is about doing lots of things at once.
Not the same, but related. 
Concurrency is about structure, parallelism is about execution.
Concurrency provides a way to structure a solution to solve a problem that may (but not necessarily) be parallelizable.

## An analogy

Concurrent: Mouse, keyboard, display, and disk drivers
Parallel: Vector dot product

## Concurrency plus communication

Concurrency is a way to structure a program by breaking it into pieces that can be executed independently.
Communication is the means to coordinate the independent executions.

## Goroutines

A goroutine is a function running independently in the same address space as other goroutines.
Like launching a function with shell's & notation.

## Goroutines are not threads

(They're a bit like threads, but they're much cheaper).
Goroutines are multiplexed onto OS threads as required.
When a goroutine blocks, that thread blocks but no other goroutine blocks.

## Channels

Channels are typed values that allow goroutines to synchronize and exchange information.

```go
timerChan := make(chan time.Time)
go func() {
    time.Sleep(deltaT)
    timerChan <- time.Now() // send time on timerChan
}()
// Do something else; when ready, receive
// Receive will block until timerChan delivers
// Value sent is other goroutine's completion time
completeAt := <- timerChan
```

## Go really supports concurrency

Really.
It's routine to create thousands of goroutines in one program
Stack start small, but grow and shrink as required.
Goroutines aren't free, but they've very cheap.

## Launching daemons

Use a closure to wrap a background operation.
This copies items from the input channel to the output channel:

```go
go func() {
    for val := range input {
        output <- val
    }
}()
```