# Concurrency is not parallelism

[Via](https://go.dev/blog/waza-talk)

When people hear the work concurrency they often think of parallelism, a related but quite distinct concept, in programming, concurrency is the composition of independently executing processes, while parallelism is the simultaneous execution of (possibly related) computations. Concurrency is about dealing with lots of things at once. Paralelism is about doing lots of thing at once.

## Concurrency

Programming as the composition of independently execution processes.

## Paralleslism

Programming as the simultaneous execution of (possible related) computations

## Concurrency vs parallelism

Concurency is about dealing with lots of things at once. Parallelism is about doing lots of things at once. Not the same, but related. Concurrency is about structure, parallelism is about execution. Concurrency provides a way to structure a solution to solve a problem that may (but not necessarily) be parallelizable.
