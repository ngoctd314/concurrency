package example

import (
	"fmt"
	"runtime"
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
}
