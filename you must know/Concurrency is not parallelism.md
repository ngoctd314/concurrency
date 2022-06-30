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
