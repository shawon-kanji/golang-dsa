# Concurrency Patterns

## Worker Pool
```go
func workerPool(jobs <-chan Job, results chan<- Result, numWorkers int) {
    var wg sync.WaitGroup
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for job := range jobs {
                results <- process(job)
            }
        }()
    }
    wg.Wait()
    close(results)
}
```

## Pipeline
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
```

## Fan-Out, Fan-In
```go
func fanOut(input <-chan int, n int) []<-chan int {
    chs := make([]<-chan int, n)
    for i := 0; i < n; i++ {
        ch := make(chan int)
        chs[i] = ch
        go func(out chan<- int) {
            for val := range input {
                out <- process(val)
            }
            close(out)
        }(ch)
    }
    return chs
}

func fanIn(chs ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup
    for _, ch := range chs {
        wg.Add(1)
        go func(c <-chan int) {
            for val := range c {
                out <- val
            }
            wg.Done()
        }(ch)
    }
    go func() {
        wg.Wait()
        close(out)
    }()
    return out
}
```

## Rate Limiting
```go
func rateLimiter(rate time.Duration, burst int) <-chan time.Time {
    ticker := time.NewTicker(rate)
    tokens := make(chan time.Time, burst)

    go func() {
        defer ticker.Stop()
        for t := range ticker.C {
            select {
            case tokens <- t:
            default:
            }
        }
    }()

    return tokens
}
```

## Timeout and Cancellation
```go
func doWorkWithTimeout(ctx context.Context, timeout time.Duration) error {
    ctx, cancel := context.WithTimeout(ctx, timeout)
    defer cancel()

    result := make(chan error, 1)
    go func() {
        result <- doWork()
    }()

    select {
    case err := <-result:
        return err
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

## Error Handling
```go
type Result struct {
    Value interface{}
    Err   error
}

func processAsync(items []Item) <-chan Result {
    results := make(chan Result)
    go func() {
        defer close(results)
        for _, item := range items {
            val, err := process(item)
            results <- Result{Value: val, Err: err}
        }
    }()
    return results
}
```

## Graceful Shutdown
```go
func server(ctx context.Context) error {
    srv := &http.Server{Addr: ":8080"}

    go func() {
        <-ctx.Done()
        shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        srv.Shutdown(shutdownCtx)
    }()

    return srv.ListenAndServe()
}
```

## Semaphore Pattern
```go
type Semaphore chan struct{}

func NewSemaphore(n int) Semaphore {
    return make(Semaphore, n)
}

func (s Semaphore) Acquire() {
    s <- struct{}{}
}

func (s Semaphore) Release() {
    <-s
}
```

## Circuit Breaker
```go
type CircuitBreaker struct {
    maxFailures int
    resetTime   time.Duration
    failures    int
    lastFail    time.Time
    mu          sync.Mutex
}

func (cb *CircuitBreaker) Call(fn func() error) error {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    if cb.failures >= cb.maxFailures {
        if time.Since(cb.lastFail) < cb.resetTime {
            return errors.New("circuit breaker open")
        }
        cb.failures = 0
    }

    err := fn()
    if err != nil {
        cb.failures++
        cb.lastFail = time.Now()
    } else {
        cb.failures = 0
    }

    return err
}
```

## Pub/Sub Pattern
```go
type PubSub struct {
    mu   sync.RWMutex
    subs map[string][]chan<- interface{}
}

func (ps *PubSub) Subscribe(topic string) <-chan interface{} {
    ps.mu.Lock()
    defer ps.mu.Unlock()

    ch := make(chan interface{}, 1)
    ps.subs[topic] = append(ps.subs[topic], ch)
    return ch
}

func (ps *PubSub) Publish(topic string, msg interface{}) {
    ps.mu.RLock()
    defer ps.mu.RUnlock()

    for _, ch := range ps.subs[topic] {
        go func(c chan<- interface{}) {
            c <- msg
        }(ch)
    }
}
```
