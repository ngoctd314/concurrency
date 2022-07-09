package main

import (
	"fmt"
)

func sendCash(buyer chan int) {
	var i int
	for i = 0; i <= 3; i++ {
		buyer <- i
	}
	close(buyer)
	buyer <- i
}

func main() {
	money := make(chan int)

	go sendCash(money)

	for seller := range money {
		fmt.Println(seller)
	}
}
