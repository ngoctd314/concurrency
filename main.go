package main

import (
	"fmt"
	"time"
)

func main() {
	dynamite := make(chan string)

	go func() {
		time.Sleep(time.Second * 3)
		dynamite <- "Dynamite Diffused!"
	}()

	for {
		fmt.Println("RUN")
		select {
		case s := <-dynamite:
			fmt.Println(s)
			return
		case <-time.After(time.Second * 2):
			fmt.Println("Time expired")
			return
		}
	}
}
