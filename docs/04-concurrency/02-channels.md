# Channels

## Channel Basics

### Creating Channels
```go
// Unbuffered channel
ch := make(chan int)

// Buffered channel
ch := make(chan int, 100)

// Channel types
var ch chan int          // Bidirectional
var sendCh chan<- int    // Send-only
var recvCh <-chan int    // Receive-only
```

### Sending and Receiving
```go
ch := make(chan int)

// Send
ch <- 42

// Receive
value := <-ch

// Receive and discard
<-ch

// Check if channel is closed
value, ok := <-ch
if !ok {
    // Channel is closed
}
```

### Closing Channels
```go
ch := make(chan int)

// Close channel (only sender should close)
close(ch)

// Send to closed channel panics
// ch <- 1  // Panic!

// Receive from closed channel returns zero value
val := <-ch  // Returns 0 immediately

// Check if closed
val, ok := <-ch
if !ok {
    fmt.Println("Channel closed")
}
```

## Buffered vs Unbuffered

### Unbuffered Channels (Synchronous)
```go
ch := make(chan int)  // No buffer

// Blocks until receiver is ready
go func() {
    ch <- 42  // Blocks here
}()

val := <-ch  // Receives 42
```

### Buffered Channels (Asynchronous)
```go
ch := make(chan int, 2)  // Buffer size 2

ch <- 1  // Doesn't block
ch <- 2  // Doesn't block
// ch <- 3  // Would block (buffer full)

val1 := <-ch  // 1
val2 := <-ch  // 2
```

## Channel Patterns

### Generator Pattern
```go
func generateNumbers(max int) <-chan int {
    ch := make(chan int)
    go func() {
        defer close(ch)
        for i := 0; i < max; i++ {
            ch <- i
        }
    }()
    return ch
}

// Usage
for num := range generateNumbers(10) {
    fmt.Println(num)
}
```

### Fan-Out (Multiple consumers)
```go
func fanOut(input <-chan int, workers int) []<-chan int {
    channels := make([]<-chan int, workers)
    for i := 0; i < workers; i++ {
        ch := make(chan int)
        channels[i] = ch
        go func(out chan<- int) {
            for val := range input {
                out <- process(val)
            }
            close(out)
        }(ch)
    }
    return channels
}
```

### Fan-In (Merge channels)
```go
func fanIn(channels ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup

    for _, ch := range channels {
        wg.Add(1)
        go func(c <-chan int) {
            defer wg.Done()
            for val := range c {
                out <- val
            }
        }(ch)
    }

    go func() {
        wg.Wait()
        close(out)
    }()

    return out
}
```

### Pipeline Pattern
```go
func stage1(nums <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range nums {
            out <- n * 2
        }
    }()
    return out
}

func stage2(nums <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range nums {
            out <- n + 1
        }
    }()
    return out
}

// Usage
input := generateNumbers(10)
doubled := stage1(input)
incremented := stage2(doubled)
for result := range incremented {
    fmt.Println(result)
}
```

### Done Channel Pattern
```go
func worker(done <-chan struct{}) {
    for {
        select {
        case <-done:
            fmt.Println("Worker stopping")
            return
        default:
            // Do work
            time.Sleep(time.Second)
        }
    }
}

func main() {
    done := make(chan struct{})
    go worker(done)

    time.Sleep(5 * time.Second)
    close(done)  // Signal stop
    time.Sleep(time.Second)
}
```

### Request-Response Pattern
```go
type Request struct {
    Data     string
    Response chan<- string
}

func handler(requests <-chan Request) {
    for req := range requests {
        // Process request
        result := process(req.Data)
        req.Response <- result
    }
}

// Usage
requests := make(chan Request)
go handler(requests)

response := make(chan string)
requests <- Request{
    Data:     "hello",
    Response: response,
}
result := <-response
```

## Channel Directions

### Send-Only Channels
```go
func sender(ch chan<- int) {
    ch <- 42
    // val := <-ch  // Error: can't receive from send-only channel
}
```

### Receive-Only Channels
```go
func receiver(ch <-chan int) {
    val := <-ch
    // ch <- 42  // Error: can't send to receive-only channel
}
```

### Bidirectional to Unidirectional
```go
func process(ch chan int) {
    sender(ch)    // Converts to chan<- int
    receiver(ch)  // Converts to <-chan int
}
```

## Select Statement

### Basic Select
```go
select {
case msg := <-ch1:
    fmt.Println("Received from ch1:", msg)
case msg := <-ch2:
    fmt.Println("Received from ch2:", msg)
case ch3 <- 42:
    fmt.Println("Sent to ch3")
}
```

### Select with Default
```go
select {
case msg := <-ch:
    fmt.Println("Received:", msg)
default:
    fmt.Println("No message")
}
```

### Select with Timeout
```go
select {
case result := <-ch:
    fmt.Println("Got result:", result)
case <-time.After(5 * time.Second):
    fmt.Println("Timeout")
}
```

### Select in Loop
```go
for {
    select {
    case msg := <-ch:
        process(msg)
    case <-done:
        return
    }
}
```

## Nil Channels

```go
var ch chan int  // nil channel

// Send to nil channel blocks forever
// ch <- 1  // Blocks forever

// Receive from nil channel blocks forever
// <-ch  // Blocks forever

// Useful for disabling cases in select
select {
case msg := <-ch:  // Never selected if ch is nil
    process(msg)
}
```

## Channel Axioms

```go
// 1. Send to nil channel blocks forever
var ch chan int
// ch <- 1  // Blocks forever

// 2. Receive from nil channel blocks forever
// <-ch  // Blocks forever

// 3. Send to closed channel panics
ch = make(chan int)
close(ch)
// ch <- 1  // Panic!

// 4. Receive from closed channel returns zero value
val := <-ch  // Returns 0

// 5. Close nil channel panics
var ch2 chan int
// close(ch2)  // Panic!

// 6. Close closed channel panics
// close(ch)  // Panic!
```

## Common Patterns

### Semaphore (Limit concurrency)
```go
sem := make(chan struct{}, maxConcurrent)

for _, task := range tasks {
    sem <- struct{}{}  // Acquire
    go func(t Task) {
        defer func() { <-sem }()  // Release
        process(t)
    }(task)
}

// Wait for all
for i := 0; i < maxConcurrent; i++ {
    sem <- struct{}{}
}
```

### Worker Queue
```go
jobs := make(chan Job, 100)
results := make(chan Result, 100)

// Start workers
for w := 0; w < numWorkers; w++ {
    go func() {
        for job := range jobs {
            results <- process(job)
        }
    }()
}

// Send jobs
go func() {
    for _, job := range allJobs {
        jobs <- job
    }
    close(jobs)
}()

// Collect results
for i := 0; i < len(allJobs); i++ {
    result := <-results
    handleResult(result)
}
```

### Cancellation
```go
func doWork(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return
        default:
            // Do work
        }
    }
}
```

## Best Practices

1. **Only sender should close** channels
2. **Close channels to signal completion**
3. **Use buffered channels** to avoid blocking
4. **Don't send/receive on same side** of bidirectional channel
5. **Check `ok` when receiving** from closed channels
6. **Use `select` with `default`** for non-blocking ops
7. **Avoid closing channels** when not necessary
8. **Use context for cancellation**
9. **Document channel ownership**
10. **Avoid channel leaks**

## Gotchas

❌ **Closing channel multiple times**
```go
ch := make(chan int)
close(ch)
// close(ch)  // Panic!
```

❌ **Sending to closed channel**
```go
ch := make(chan int)
close(ch)
// ch <- 1  // Panic!
```

❌ **Not ranging over closed channel**
```go
ch := make(chan int)
go func() {
    for i := 0; i < 5; i++ {
        ch <- i
    }
    close(ch)  // Important!
}()

for val := range ch {
    fmt.Println(val)
}
```

❌ **Deadlock**
```go
ch := make(chan int)
ch <- 1  // Deadlock! No receiver
<-ch
```

## Debugging

### Detect Goroutine Leaks
```go
// Monitor goroutine count
before := runtime.NumGoroutine()
// Run code
after := runtime.NumGoroutine()
if after > before {
    log.Printf("Potential goroutine leak: %d -> %d", before, after)
}
```

### Channel Buffer Status
```go
ch := make(chan int, 10)
ch <- 1
ch <- 2

fmt.Println("Length:", len(ch))  // 2
fmt.Println("Capacity:", cap(ch))  // 10
```
