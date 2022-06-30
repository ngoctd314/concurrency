# The bridge channel

In some circumstances, you may find yourself wanting to consume values from a sequence of channels:

```go
<- chan <- chan any
```
This is sightly then a slice of channels into a single channel. A sequence of channels suggest an ordered write from different sources.

As a consumer, the code may not care about the fact that its values come from a sequence of channels. In that case, dealing with a channel of channels can be cumbersome. If we instead define a function that can destructure the channel of channels into a simple channel - a technique called bridging the channels.

```go
bridge := func(done <- chan any) {}
```