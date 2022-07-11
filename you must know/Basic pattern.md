# Basic concurrency pattern

## Generator pattern

Let's study our first pattern of concurrency, which is about generators that return channels as returning argument.

Generators return the next value in a sequence each time they are called. This means that each value is avaiable as an output before the generator computes the next value.

```go
// we use a generator as a a function which returns a channel
// Both the goroutine and the main routine can execute concurrently as we
// print the value of the counter as soon as we receive it while the goroutine
// simultaneously computes the next value.
func foo() <- chan string {
    ch := make(chan string) 

    go func() {
        for i := 0 ; ; i++ {
            fmt.Println("send")
            ch <- fmt.Sprintf("%s %d", "Counter at: ", i)
        }
    }()

    return ch // return before computed
}

func main() {
    ch := foo() // foo() returns a channel

    for i := 0 ; i < 5 ; ;i++ {
        fmt.Printf("%q\n", <- ch)
    }

    fmt.Println("Done with Counter")
}
```
I hope you are with me so far. Let's extend this pattern to a concept where we think of the channel as a handle of a service.

```go
func updatePosition(name string) <-chan string {
	ch := make(chan string)

	go func() {
		for i := 0; ; i++ {
			ch <- fmt.Sprintf("%s %d", name, i)
		}
	}()

	return ch
}

func main() {
	ch1 := updatePosition("Legolas: ")
	ch2 := updatePosition("Gandalf: ")

	for i := 0; i < 5; i++ {
		fmt.Println(<-ch1)
		fmt.Println(<-ch2)
	}

	fmt.Println("Done with getting updates on positions.", runtime.NumGoroutine())
}
```

In the code above, we are getting updates on the position of Legolas and Gandalf using the updatePosition function. We again launch the goroutine from inside the function, which is the key thing to notice in both the examples we have seen so far. Thus by returning a channel, the function enables us to communicate with the service it provides which, in this case, is giving position updates.

However, we still have a slight problem, which is when the following statements are blocking each other:

```go
fmt.Println(<- ch1)
fmt.Println(<- ch2)
```

## Fan-in, Fan-out

In this lesson, we will get familiar with Fan-in, Fan-out techniques which are used to multiplex and demultiplex data in Go.

Fan-in refers to a technique in which you join data from multiple inputs into a single entity. On the other hand, Fan-out means to divide the data from a single source into multiple smaller chunks.

The code below is from the previous lesson where two receiving operations were blocking each other

```go
fmt.Println(<-positionChannel1)
fmt.Println(<-positionChannel2)
```

But what if you want to get position updates as soon as they are updated? This is where the fan-in technique comes into play. By using this technique, we'll combine the inputs from both channels and send them through a single channel. Look at the code below to see how it's done.

```go
func updatePosition(name string, t time.Duration) <-chan string {
	ch := make(chan string)

	service := func() {
		for i := 0; ; i++ {
			ch <- fmt.Sprintf("%s %d", name, i)
			time.Sleep(t * time.Second)
		}
	}

	go service()

	return ch
}

func fanIn(list ...<-chan string) <-chan string {
	ch := make(chan string)
	wg := sync.WaitGroup{}

	service := func(k <-chan string) {
		defer wg.Done()
		for {
			ch <- <-k
		}
	}

	wg.Add(len(list))
	for _, v := range list {
        // data will be passed to ch as soon as it is received by element in list
        // because the goroutines are running concurrently
        // element in list commnunicating with ch now instead of directly communicating with
        // the main routine as before. You will realize from the output that the position
        // updates are no longer sequential. Thus by using this technique, we can solve the
        // blocking issue that we were previously facing.
		go service(v)
	}

    closer := func() {
		wg.Wait()
		close(ch)
    }
    go closer()

	return ch

}

func main() {
	ch1 := updatePosition("Legolas: ", 1)
	ch2 := updatePosition("Gandalf: ", 2)

	f := fanIn(ch1, ch2)

	for k := range f {
		fmt.Println(k)
	}
}
```

Let's jump to the Fan-out technique now. The code below generates an array of random numbers and prints all the values after doubling them.

```go
func main() {
	var myNumbers [10]int
	for i := 0; i < 10; i++ {
		rand.Seed(time.Now().UnixNano())
		myNumbers[i] = rand.Intn(50)
	}

	myChannelOut := channelGenerator(myNumbers)

	ch1 := fanout(myChannelOut, "channel 1")
	ch2 := fanout(myChannelOut, "channel 2")

	ch := fanin(ch1, ch2)
	for k := range ch {
		fmt.Println(k)
	}
}

func channelGenerator(numbers [10]int) <-chan string {
	channel := make(chan string)

	service := func() {
		for _, i := range numbers {
			channel <- strconv.Itoa(i)
		}
		close(channel)
	}

	go service()

	return channel
}

func fanout(ch <-chan string, msg string) <-chan string {
	chout := make(chan string)

	service := func() {
		for v := range ch {
			num, err := strconv.Atoi(v)
			if err != nil {
			}
			chout <- fmt.Sprintf("%s %d*2 = %d", msg, num, num*2)
		}

		close(chout)
	}

	go service()

	return chout
}

func fanin(ins ...<-chan string) <-chan string {
	out := make(chan string)
	wg := sync.WaitGroup{}

	service := func(ch <-chan string) {
		defer wg.Done()
		for k := range ch {
			out <- k
		}
	}

	wg.Add(len(ins))
	for _, in := range ins {
		go service(in)
	}

	closer := func() {
		wg.Wait()
		close(out)
	}
	go closer()

	return out
}
```

The Fan-in, Fan-out techniques can be pretty useful when we have to divide work among jobs and then combine the results from those jobs

## Sequencing

This lesson will teach you about a pattern used for sequencing in a program by sending channel over a channel

Remembder the code from the last lesson where we unblock the two receive operations. But operations is not sequencing.

What if we don't want to block the code and introduce sequence to our program instead of randomness? Let's see how we approach this problem

Imagine a cooking competition. You are participating in it with your partner.

The rules of the same are:

1. There are three rounds in the competition
2. In each round, both partners will have to come up with their own dishes.
3. A player cannot move on to the next round until their partner is done with their dish.
4. The judge will decide the entry to the next round after tasting food from both the team members.

Hope the rules are clear. Now in order to achieve this scenario, we'll send a channel over a channel!

Surprised? Let's see how it's done:

```go
```

## Range and Close

This lesson will teach you a way to close a channel

If you remember the lesson on channels, you are already familiar with this pattern

In Go, we have a range function which lets us iterate over elements in different data structs. Using this function, we can range over the items we receive on a channel until it is closed. Also, note that only the sender, not the receiver, should close the channel when it feels that it has no more values to send.

```go
type money struct {
	amount int
	year   int
}

func sendMoney(parent chan money) {
	for i := 0; i <= 18; i++ {
		parent <- money{5000, i}
	}
	close(parent)
}

func main() {
	money := make(chan money)

	go sendMoney(money)

	for kidMoney := range money {
		fmt.Printf("Money received by kid in year %d : %d\n", kidMoney.year, kidMoney.amount)
	}
}
```

**Note: If you are sending on a closed channel, it will cause a panic**

## Quit Channel

In this lesson, we'll study a pattern related to quiting from a select statement

## Timeout using select statement

This lesson will introduce you a pattern which uses the time.After functino in a select statement. What if we want to break out of channel communications after a certain period of time?

This will be done using time.After function which is imported from the time package. It returns a channel that bloacks the code for the specified duration. After that duration, the channel delivers the current time but only once.