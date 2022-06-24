# Generator technique

Generators return the next value in a sequence each time they are called. This is means that each value is avaiable as an output before the generator computes the next value.

```go
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
```
Let's extend this pattern to a concept where we think of the channel as a handle on a service.
