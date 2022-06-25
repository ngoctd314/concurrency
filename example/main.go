package example

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// Exec ...
func Exec() {
	// execBoring()
	// execFoo()
	// execUpdatePosition()
	// execFibonaciGenerator()
	fn := func(t int) <-chan int {
		outbound := make(chan int)
		go func() {
			time.Sleep(time.Second * time.Duration(t))
			outbound <- t
			close(outbound)
		}()

		return outbound
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

	fanIn := func(done <-chan any, channels ...<-chan any) <-chan any {
		var wg sync.WaitGroup
		multiplexedStream := make(chan any)

		multiplex := func(c <-chan any) {
			defer wg.Done()
			for i := range c {
				select {
				case <-done:
					return
				case multiplexedStream <- i:
				}
			}
		}

		wg.Add(len(channels))
		for _, c := range channels {
			go multiplex(c)
		}

		go func() {
			wg.Wait()
			close(multiplexedStream)
		}()

		return multiplexedStream
	}

	_ = fanIn
}
