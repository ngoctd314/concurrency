# Goroutines in Go

## Overview

Goroutines can be thought as a lightweight thread has a separate independent execution and which can execute concurrently with other goroutines. It is a function of method that executing concurrently with other goroutines. It is entirely managed by the Go runtime. Golang is  a concurrent language. Each goroutine is an independent execution. It is goroutine that helps achieve concurrency in golang.

## Main goroutine

The main function in the main package is the main goroutine. All goroutines are started from the main goroutine. These goroutines can then start multiple other goroutine and so on. 

The main goroutine represents the main program. Once it exits then it means that the program has exited.

Goroutines don't have parents or children. When you start a goroutine it just executes alongside all other running goroutines. Each goroutine exits only when its function returns. The only exeception to that is all goroutines exit when the main goroutine exits.

## Scheduling of the goroutines

Once the go program starts, go runtime will launch OS threads equivalent to the number of logical CPUs usable by the current process. There is one logical CPU per virtual core.

The go program will launch OS threads equal to the number of logical CPUs avaiable to it or the output of runtime.NumCPU(). These threads will be managed by the OS and scheduling of these threads on to CPU cores is the responsibility of OS only.

The go runtime has its own scheduler that will multiplex the goroutines on the OS level threads in the goruntime. So essentially each goroutine is running on an OS thread that is assigned to a logical CPU.

There are two queues involved for managing the goroutines and assigning it to the OS threads.

## Local run queue

Within go routine each of this OS thread will have one queue associated with it. It is called Local Run Queue. It contains all the goroutines that will be executed in the context of that thread. The go runtime will be doing the scheduling and context switching of the goroutines belonging to a particular LRQ to the corresponding OS level thread which owns this LRQ.

## Global Run Queue

It contains all the goroutines that haven't been moved to any LRQ of any OS thread. The Go scheduler will assign a goroutine from this queue to the LRQ of any OS thread.

## Golang scheduler is a cooperative scheduler

## Advantages of goroutines over threads