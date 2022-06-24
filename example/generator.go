package example

import (
	"fmt"
	"math/rand"
	"time"
)

// boring is a function that returns a channel to communicate with it
//
// generator technique:
// will return reference to channel (receive only channel)
// and it process will exec in a goroutine
func boring(msg string) <-chan string {
	outbound := make(chan string)
	// we launch goroutine inside a functino
	// that sends the data to channel
	go func() {
		// infinite seender
		for i := 0; i < 10; i++ {
			outbound <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}

		// the sender should close the channel
		close(outbound)
	}()

	return outbound
}

func foo() <-chan string {
	ch := make(chan string)

	go func() {
		for i := 0; ; i++ {
			ch <- fmt.Sprintf("%s %d", "Counter at : ", i)
		}
	}()

	// return immediately
	return ch

}

func execBoring() {
	joe := boring("Joe")
	ahn := boring("Ahn")

	for i := range joe {
		fmt.Println(i)
	}
	for j := range ahn {
		fmt.Println(j)
	}

	fmt.Println("You are both boring. I'm leaving")
}

func execFoo() {
	ch := foo()
	for i := 0; i < 5; i++ {
		fmt.Printf("%q\n", <-ch)
	}
	fmt.Println("Done with Counter")
}

func updatePosition(name string) <-chan string {
	positionCh := make(chan string)
	go func() {
		for i := 0; ; i++ {
			positionCh <- fmt.Sprintf("%s %d", name, i)
		}
	}()

	return positionCh
}

func execUpdatePosition() {
	ch1 := updatePosition("Legolas: ")
	ch2 := updatePosition("Gandalf: ")

	for i := 0; i < 5; i++ {
		fmt.Println(<-ch1) // blocking <-ch1
		fmt.Println(<-ch2) // blocking <-ch2
	}
}

func fibonaciGenerator(n int) <-chan int {
	outbound := make(chan int, n)

	f1, f2 := 0, 1
	go func() {
		defer close(outbound)

		for i := 0; i < n; i++ {
			outbound <- f1
			f1, f2 = f2, f1+f2
		}
	}()

	return outbound
}

func execFibonaciGenerator() {
	fibo := fibonaciGenerator(5)
	for f := range fibo {
		fmt.Println(f)
	}
}
