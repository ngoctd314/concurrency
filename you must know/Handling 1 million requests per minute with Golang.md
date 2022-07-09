# Handling 1 milion requests per minute with Golang

[via](http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/)

## Trying again

We needed to find a different way. Since the beginning we started discussing how we needed to keep the lifetime of the request handler very short and spawn processing in the background.

So the second iteration was to create a buffered channel where we could queue up some jobs and upload them to S3, and since we could control the maximum number of items in our queue and we had plenty of RAM to queue up jobs in memory, we though it would be okay to just buffer jobs in the channel queue.

```go
var Queue chan Payload

func init() {
    Queue = make(chan Payload, MAX_QUEUE)
}

func payloadHandler(w http.ResponseWriter, r *http.Request) {
    // Go through each payload and queue items
    for _, payload := range content.Payloads {
        Queue <- payload
    }
}
```

And then to actually dequeue jobs and process them, we were using something similar to this:

This approach didn't buy us anything, we have traded flawed concurrency with a buffered queue that was simply postponing the problem. Our sychronous processor was only uploading one payload at a time to S3, and since the rate of incoming requests were much larger than the ability of the single processor to update to S3, our buffered channel was quickly reaching its limit and blocking the requestr handler ability to queue more times.

We were simply avoiding the problem and started a count-down to the death of our system eventually. Our latency rates kept increasing in a constant rate minutes after we deployed this flawed version.

```go
func StartProcessor() {
    for {
        select {
            case job:= <- Queue:
                job.payload.UploadToS3()
        }
    }
}
```

## The Better Solution

We have decided to utilize a common pattern when using Go channels, in order to create a 2-tier channel system, one for queuing jobs and another to control how many workers operate on the JobQueue concurrently.

The idea was to parallelize the updates to S3 to somewhat sustainable rate, one that would not cripple the machine nor start generating connections error from S3. So we have opted for creating a Job/Worker pattern. Think about this as the Golang way of implementing a Worker Thread-Pool utilizing channels instead.

```go
var (
    MaxWorker = os.Getenv("MAX_WORKERS")
    MaxQueue = os.Getenv("MAX_QUEUE")
)

// Job represents the job to be run
type Job struct {
    Payload Payload
}

// A buffered channel that we can send work requests on
var JobQueue chan Job

// Worker represents the worker that executes the job
type Worker struct {
    WorkerPool chan chan Job
    JobChannel chan Job
    quit chan bool
}

func NewWorker(workerPool chan chan Job) Worker {
    return Worker{
        WorkerPool: workerPool,
        JobChannel: make(chan Job),
        quit: make(chan bool)
    }
}

// start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w Worker) Start() {
    go func() {
        for {
            // register the current worker into the worker queue
            w.WorkerPool <- w.JobChannel
            select {
            case job := <- w.JobChannel:
                // we have received a work request 
                if err := job.Payload.UplaodToS3(); err != nil {
                    log.Errorf("Error uploading to S3: %s")
                }
            }
            case <- w.quit:
                // we have received a signal to stop
                return

        }
    }()
}

func(w Worker) Stop() {
    go func() {
        w.quit <- true
    }()
}
```

We have modified our Web request handler to create an instance of Job struct with payload and send into the JobQueue channel for the workers to pickup.

```go
func payloadHandler(w http.ResponseWriter, r *http.Request) {
    // Go through each payload and queue items individually to be posted to S3
    for _, payload := range content.Payloads {
        // let's create a job with the payload
        work := Job{Payload: payload}
        // Push the work onto the queue
        JobQueue <- work
    }
}
```

During our web server initialization we create a Dispatcher and call Run() to create the pool of workers and start listening for jobs that would appear in the JobQueue.

```go
dispatcher := NewDispatcher(MaxWorker)
dispatcher.Run()
```

Dispatcher implementation

```go
type Dispatcher struct {
    // A pool of workers channels that are registered with the dispatcher
    WorkerPool chan chan Job
}

func NewDispatcher(maxWorkers int) *Dispatcher {
    // create a pool
    pool := make(chan chan Job, maxWorkers)
    return &Dispatcher{WorkerPool: pool}
}

func (d *Dispatcher) Run() {
    // starting n number of workers
    for i := 0 ; i < d.MaxWorkers; i++ {
        worker := NewWorker(d.pool)
        worker.Start()
    }

    go d.dispatch()
}

func (d *Dispatcher) dispatch() {
    for {
        select {
        case job := <- JobQueue:
            go func(job Job) {
                // try to obtain a worker job channel that is avaiable 
                // this will block until a worker is idle
                jobChannel := <- d.WorkerPool
                // dispatch the job to the worker job channel
                jobChannel <- job
            }(job)
        }
    }
}
```
