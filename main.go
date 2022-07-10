package main

func sendCash(buyer chan int) {
	var i int
	for i = 0; i <= 3; i++ {
		buyer <- i
	}
	close(buyer)
	buyer <- i
}
