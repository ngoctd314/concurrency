# Pipelines

A pipeline is just another tool you can use to form an abstraction in your system. In particular, it is very powerful tool to use when your program needs to process streams or batches of data. A pipeline is nothing more than a series of things that take data in, perform an operation on it, and pass the data back out. We call each of these operations a stage of the pipeline.

By using a pipeline, you separate the concerns of each stage, which provides numerous benefits. You can modify stages independent of one another, you can mix and match how stages are combined independent of modifying the stages, you can process each stage concurrent to upstream or downstream stages, and you can fan-out or rate-limit portions of your pipeline.

As mentioned previously, a stage is just something that takes data in, performs a transformation on it, and sends the data back out.

```go
multily := func(values []int, multiplier int) []int {
    multipliedValues := make([]int, len(values))
    for i, v := range values {
        multipliedValues[i] = v * multiplier
    }

    return multipliedValues
}

add := func(values []int, additive int) []int {
    addedValues := make([]int, len(values))
    for i, v := range values {
        addedValues[i] = v + additive
    }
    return addedValues
}

fmt.Println(multily(add([]int{1, 1, 1}, 1), 2))
```

Add and multiple have the properties of a pipeline stage, we're able to combine them to form a pipeline. That's interesting; what are the properties of a pipeline stage?

- A stage consumes and returns the same type
- A stage must be reified by the language so that it may be passed around. Function in Go are reified and fit this purpose nicely.

Here, our add and multiply stages satisfy all the properties of a pipeline stage: they both consume a slice of int and return a slice of int, and because Go has reified functions, we can pass add and multiple around.

Notice how each stage is talking a slice of data and returing a slice of data? There stages are performing what we call batch processing. This just means that they operate on chunks of data all at once instead of one discrete value at a time. There is a type of pipeline stage that performs stream processing. This means that the stage receives and emits one element at a time.

There are pros and cons to batch processing versus stream processing.

```go
multiply := func(value, multiplier int) int {
    return value * multiplier
}

add := func(value, additive int) int {
    return value + additive
}

ints := []int{1, 2, 3, 4}
for _, v := range ints {
    fmt.Println(multiply(add(multiply(v,2), 1), 2))
}
```

## Best practices for constructing pipelines

Channels are uniquely suited to constructing pipelines in Go because fulfill all of our basic requirements. The can receive and emit values, they can safely be used to concurrently, they can be ranged over, and they are reified by the language.
