# Channels inside channels pattern in Golang

Sometimes I struggles especially with request-response cases. It appeared that I should sill use mutexes, but I didn't want to. This post is about how to not do that and put channels inside channels instead.

## The problem

It sound odd first, so let me show a real-life example. Imagine you have a piece of data, changing from time to time. You might assign an owner method running in a go-routine to conduce the updates. Say your process is a network service and the subject data is served via its API when a client request comes. That means there will be data reads from the same data structure, which may changes inside its owner go-routine meanwhile.

```go
type Data struct {
    secretOfTheSecond int
}

// Run writes the data periodically to the struct. This meant to run in 
// it's own goroutine
func (d *Data) Run() {
    seed := rand.NewSource(time.Now().UnixNano())
    gen := rand.New(seed)
    ticker := time.NewTicker(1 * time.Second)
}
```

```go
type Data struct {
    secretOfTheSecond int
    ReadRequest chan struct{}
    ReadResponse chan int
}
```
We are getting closer. If any go-routine wants the data, it sends anything into the ReadRequest channel and starts listening on the ReadResponse channel. Run() waits for anything pushed into th ReadRequest channel and pushed the current data to the ReadResponse channel.

But we now have a new problem...

A channel can be listened to freely by multiple go-routines at the same time. If more than one receiver is interested in the update coming out from the response channel, only one of them will be able to read an update. Of course, if there are multiple receivers, they all should send data requests, but which response is for whom.

## Invert channel in channel

Channels are first-class citizens,a kind of basic data type in Go. As part of the declaration, they get another arbitrary data type assigned, which they can transport later on. This being said, a channel can hold any data type, which again can be a channel!

Ow new data structure can look like this:

```go
type Data struct {
    secretOfTheSecond int
    readRequest chan chan int
}
```

Now we can re-implement the Get() method for the struct to create an exclusive channel (of type int) in the reader's context and send it into the readRequest channel. The data-owner go-routine can dispatch this event and can extract t;he channel created by the data requester to send the response.

```go
type Data int {
    secreateOfTheSecond int
    readRequest chan chan int
}

// Run writes the data periodically to the struct. This meat to run in it's own goroutine
func (d *Data) Run() {
    seed := rand.NewResource(time.Now().UnixNano())
    gen := rand.New(seed)
    ticker := time.NewTicker(1 * time.Second)
    for {
        select {
        case <- ticker.C:
            d.secretOfTheSecond = gen.Int()
        case responseChan := <- d.readRequest
            responseChan <- d.secretOfTheSecond
        }
    }
}
// Get does an ad-hoc read of the data. This is the user interface
func (d *Dat) Get() int {
    responseChan := make(chan int)

    // send addr to channel
    d.readRequest <- responseChan
    response : <- responseChan

    return response
}
```