package main

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
)

func producer(ctx context.Context, words []string) (<-chan string, error) {
	outbound := make(chan string)
	go func() {
		defer close(outbound)

		for _, s := range words {
			select {
			case <-ctx.Done():
				fmt.Println("DONE")
				return
			case outbound <- s:
				// default:
				// 	outbound <- s
			}
		}
	}()

	return outbound, nil
}

func transformToLower(ctx context.Context, values <-chan string) (<-chan string, error) {
	outbound := make(chan string)

	go func() {
		defer close(outbound)

		for s := range values {
			select {
			case <-ctx.Done():
				return
			default:
				outbound <- strings.ToLower(s)
			}
		}
	}()

	return outbound, nil
}

func transformToTitle(ctx context.Context, values <-chan string) (<-chan string, error) {
	outbound := make(chan string)

	go func() {
		for s := range values {
			select {
			case <-ctx.Done():
				return
			default:
				outbound <- strings.ToTitle(s)
			}
		}
	}()

	return outbound, nil
}

func mergeStringChans(ctx context.Context, cs ...<-chan string) <-chan string {
	var wg sync.WaitGroup
	outbound := make(chan string)

	output := func(c <-chan string) {
		defer wg.Done()
		for n := range c {
			select {
			case <-ctx.Done():
				return
			default:
				outbound <- n
			}
		}
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(outbound)
	}()

	return outbound
}

func sink(ctx context.Context, values <-chan string) {
	for {
		select {
		case <-ctx.Done():
			return
		case val, ok := <-values:
			if ok {
				fmt.Printf("sink: %s", val)
			} else {
				fmt.Println("Done")
				return
			}
		}
	}
}

func mainn() {
	// source := []string{"FOO", "BAR", "BAX"}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		time.Sleep(time.Second * 2)
		cancel()
	}()
	select {
	case <-ctx.Done():
		fmt.Println("RUNN")
	}

	time.Sleep(time.Second * 10)

	// producer, err := producer(ctx, source)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// time.Sleep(time.Second * 10)
	// for v := range producer {
	// 	fmt.Println(v)
	// }

	// stage1Channels := []<-chan string{}
	// for i := 0; i < runtime.NumCPU(); i++ {
	// 	lowerCaseChannel, err := transformToLower(ctx, producer)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	stage1Channels = append(stage1Channels, lowerCaseChannel)
	// }

}
