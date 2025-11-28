# Goroutines

## Goroutine Basics

### Starting a Goroutine
```go
// Sequential execution
func main() {
    doSomething()
    doSomethingElse()
}

// Concurrent execution
func main() {
    go doSomething()      // Runs concurrently
    doSomethingElse()
}
```

### Anonymous Function Goroutine
```go
go func() {
    fmt.Println("Running in goroutine")
}()

// With parameters
go func(msg string) {
    fmt.Println(msg)
}("Hello from goroutine")
```

### Waiting for Goroutines

#### Using sync.WaitGroup
```go
var wg sync.WaitGroup

for i := 0; i < 5; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        fmt.Printf("Goroutine %d\n", id)
    }(i)
}

wg.Wait()  // Wait for all goroutines to finish
```

#### Using Channels
```go
done := make(chan bool)

go func() {
    // Do work
    done <- true
}()

<-done  // Wait for completion
```

## Goroutine Patterns

### Worker Pool
```go
func worker(id int, jobs <-chan int, results chan<- int) {
    for j := range jobs {
        fmt.Printf("Worker %d processing job %d\n", id, j)
        time.Sleep(time.Second)
        results <- j * 2
    }
}

func main() {
    jobs := make(chan int, 100)
    results := make(chan int, 100)

    // Start workers
    for w := 1; w <= 3; w++ {
        go worker(w, jobs, results)
    }

    // Send jobs
    for j := 1; j <= 9; j++ {
        jobs <- j
    }
    close(jobs)

    // Collect results
    for a := 1; a <= 9; a++ {
        <-results
    }
}
```

### Fan-Out, Fan-In
```go
// Fan-out: Multiple goroutines read from same channel
func fanOut(input <-chan int, n int) []<-chan int {
    channels := make([]<-chan int, n)
    for i := 0; i < n; i++ {
        ch := make(chan int)
        channels[i] = ch
        go func() {
            for val := range input {
                ch <- process(val)
            }
            close(ch)
        }()
    }
    return channels
}

// Fan-in: Multiple channels merge into one
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

### Pipeline
```go
func generator(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()
    return out
}

func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * n
        }
        close(out)
    }()
    return out
}

// Usage
nums := generator(1, 2, 3, 4)
squared := square(nums)
for result := range squared {
    fmt.Println(result)
}
```

## Goroutine Scheduling

### GOMAXPROCS
```go
import "runtime"

// Get number of CPUs
numCPU := runtime.NumCPU()

// Set max number of CPUs to use
runtime.GOMAXPROCS(numCPU)

// Default is number of CPUs available
```

### Yielding Execution
```go
// Yield processor to other goroutines
runtime.Gosched()

// Example use case
for i := 0; i < 10; i++ {
    go func(id int) {
        for j := 0; j < 5; j++ {
            fmt.Printf("Goroutine %d: %d\n", id, j)
            runtime.Gosched()  // Allow other goroutines to run
        }
    }(i)
}
```

## Goroutine Lifecycle

### Starting
```go
go func() {
    // Goroutine starts immediately
}()
```

### Completion
```go
// Goroutine completes when:
// 1. Function returns
// 2. main() exits (kills all goroutines)
// 3. Panic (can be recovered)
```

### Graceful Shutdown
```go
func worker(ctx context.Context, id int) {
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("Worker %d shutting down\n", id)
            return
        default:
            // Do work
            time.Sleep(time.Second)
        }
    }
}

func main() {
    ctx, cancel := context.WithCancel(context.Background())

    for i := 0; i < 5; i++ {
        go worker(ctx, i)
    }

    time.Sleep(5 * time.Second)
    cancel()  // Signal shutdown
    time.Sleep(time.Second)  // Wait for cleanup
}
```

## Common Patterns

### Rate Limiting
```go
func rateLimited() {
    rate := time.Second / 10  // 10 requests per second
    throttle := time.Tick(rate)

    for req := range requests {
        <-throttle  // Wait for rate limit
        go handleRequest(req)
    }
}
```

### Timeout
```go
func doWorkWithTimeout(timeout time.Duration) error {
    done := make(chan bool)

    go func() {
        doWork()
        done <- true
    }()

    select {
    case <-done:
        return nil
    case <-time.After(timeout):
        return errors.New("timeout")
    }
}
```

### Periodic Task
```go
func periodicTask(interval time.Duration, task func()) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            task()
        }
    }
}
```

## Goroutine Best Practices

### Always Handle Panics
```go
func safeGoroutine(fn func()) {
    go func() {
        defer func() {
            if r := recover(); r != nil {
                log.Printf("Goroutine panic: %v", r)
            }
        }()
        fn()
    }()
}
```

### Avoid Goroutine Leaks
```go
// Bad: Goroutine never exits
func leak() {
    ch := make(chan int)
    go func() {
        val := <-ch  // Blocks forever
        fmt.Println(val)
    }()
    // ch never receives value, goroutine leaks
}

// Good: Use context or timeout
func noLeak(ctx context.Context) {
    ch := make(chan int)
    go func() {
        select {
        case val := <-ch:
            fmt.Println(val)
        case <-ctx.Done():
            return
        }
    }()
}
```

### Don't Create Too Many Goroutines
```go
// Bad: May create millions of goroutines
for _, item := range millionItems {
    go process(item)
}

// Good: Use worker pool
func processWithPool(items []Item, numWorkers int) {
    jobs := make(chan Item, len(items))
    var wg sync.WaitGroup

    // Start workers
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for item := range jobs {
                process(item)
            }
        }()
    }

    // Send jobs
    for _, item := range items {
        jobs <- item
    }
    close(jobs)
    wg.Wait()
}
```

## Gotchas

❌ **Loop variable capture**
```go
// Wrong: All goroutines print same value
for i := 0; i < 10; i++ {
    go func() {
        fmt.Println(i)  // Captures reference, not value
    }()
}

// Right: Pass as parameter
for i := 0; i < 10; i++ {
    go func(id int) {
        fmt.Println(id)
    }(i)
}

// Or shadow variable
for i := 0; i < 10; i++ {
    i := i  // Create new variable
    go func() {
        fmt.Println(i)
    }()
}
```

❌ **main exits before goroutines**
```go
func main() {
    go fmt.Println("Hello")
    // main exits, goroutine may not run
}

// Fix: Wait for goroutines
func main() {
    done := make(chan bool)
    go func() {
        fmt.Println("Hello")
        done <- true
    }()
    <-done
}
```

❌ **Goroutine leaks**
```go
// Goroutine waits forever
func doWork() {
    ch := make(chan int)
    go func() {
        <-ch  // Blocks forever if nothing sends
    }()
}
```

## Debugging Goroutines

### Get Goroutine Count
```go
import "runtime"

count := runtime.NumGoroutine()
fmt.Printf("Number of goroutines: %d\n", count)
```

### Stack Traces
```go
import "runtime"

// Print all goroutine stacks
buf := make([]byte, 1<<20)
stackSize := runtime.Stack(buf, true)
fmt.Printf("%s\n", buf[:stackSize])
```

### pprof for Goroutine Profiling
```go
import _ "net/http/pprof"
import "net/http"

go func() {
    http.ListenAndServe("localhost:6060", nil)
}()

// Visit http://localhost:6060/debug/pprof/goroutine
```

## Performance Tips

1. **Goroutines are cheap** but not free (~2KB stack)
2. **Use worker pools** for high-volume tasks
3. **Buffered channels** can reduce blocking
4. **Avoid excessive context switching**
5. **Profile before optimizing**
6. **Consider goroutine overhead** for tiny tasks
7. **Clean up goroutines** when done
8. **Use sync.Pool** for frequently allocated objects
9. **Batch operations** when possible
10. **Monitor goroutine count** in production
