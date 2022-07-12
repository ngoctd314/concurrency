package main

type work struct {
	x, y, z int
}

func worker(in <-chan *work, out chan<- *work) {
	for w := range in {
		w.z = w.x * w.y
	}
}

func main() {
}
