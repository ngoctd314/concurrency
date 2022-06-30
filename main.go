package main

import (
	"fmt"
	"sync"
	"time"
)

func gen(a ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, v := range a {
			out <- v
		}
		close(out)
	}()

	return out
}

func sq(ch <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for v := range ch {
			out <- v * v
		}
		close(out)
	}()

	return out
}

func merge(done chan struct{}, ch ...<-chan int) <-chan int {
	out := make(chan int)

	wg := sync.WaitGroup{}
	output := func(ch <-chan int) {
		for v := range ch {
			out <- v
		}
		wg.Done()
	}

	wg.Add(len(ch))
	for _, v := range ch {
		go output(v)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				fmt.Println("RUN")
				time.Sleep(time.Second)
			}
		}
	}()

	return out
}

func main() {
	gen := gen(2, 3)
	sq1 := sq(gen)
	sq2 := sq(gen)

	done := make(chan struct{})
	merge := merge(done, sq1, sq2)

	fmt.Println(<-merge, <-merge)
	done <- struct{}{}
	time.Sleep(time.Minute)
}
