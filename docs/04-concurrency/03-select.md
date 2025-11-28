# Select Statement

## Select Basics

### Basic Select
```go
select {
case msg1 := <-ch1:
    fmt.Println("Received from ch1:", msg1)
case msg2 := <-ch2:
    fmt.Println("Received from ch2:", msg2)
case ch3 <- 23:
    fmt.Println("Sent to ch3")
}
// Blocks until one case can proceed
// Random selection if multiple cases ready
```

### Select with Default
```go
select {
case msg := <-ch:
    fmt.Println("Received:", msg)
default:
    fmt.Println("No message available")
}
// Non-blocking - executes default if no case ready
```

## Common Patterns

### Timeout Pattern
```go
select {
case result := <-ch:
    fmt.Println("Got result:", result)
case <-time.After(5 * time.Second):
    fmt.Println("Timeout after 5 seconds")
}
```

### Cancellation Pattern
```go
func worker(done <-chan struct{}, work <-chan int) {
    for {
        select {
        case <-done:
            fmt.Println("Cancelled")
            return
        case w := <-work:
            process(w)
        }
    }
}
```

### Multiple Channels
```go
for {
    select {
    case msg := <-input1:
        handle1(msg)
    case msg := <-input2:
        handle2(msg)
    case msg := <-input3:
        handle3(msg)
    case <-quit:
        return
    }
}
```

### Non-Blocking Send/Receive
```go
// Non-blocking receive
select {
case msg := <-ch:
    fmt.Println(msg)
default:
    fmt.Println("Nothing to receive")
}

// Non-blocking send
select {
case ch <- value:
    fmt.Println("Sent")
default:
    fmt.Println("Can't send")
}
```

### Fairness with Select
```go
// Select chooses randomly if multiple cases ready
// To ensure fairness, use buffered channels or separate goroutines
```

## Advanced Select Patterns

### Priority Select
```go
// Give priority to done channel
select {
case <-done:
    return
default:
    select {
    case msg := <-messages:
        process(msg)
    case <-done:
        return
    }
}
```

### Rate Limiting
```go
rate := time.Tick(100 * time.Millisecond)
for req := range requests {
    <-rate  // Wait for rate limiter
    go handle(req)
}
```

### Heartbeat
```go
func heartbeat(interval time.Duration) <-chan struct{} {
    ch := make(chan struct{})
    go func() {
        ticker := time.NewTicker(interval)
        defer ticker.Stop()
        for {
            select {
            case <-ticker.C:
                ch <- struct{}{}
            }
        }
    }()
    return ch
}

// Usage
hb := heartbeat(1 * time.Second)
for {
    select {
    case <-hb:
        fmt.Println("Heartbeat")
    case work := <-workCh:
        process(work)
    }
}
```

### Fanout
```go
func fanout(input <-chan int, outputs []chan<- int) {
    for val := range input {
        for _, out := range outputs {
            select {
            case out <- val:
            case <-time.After(time.Second):
                // Skip slow consumer
            }
        }
    }
}
```

## Best Practices

1. **Random selection** - select randomly if multiple cases ready
2. **Use default carefully** - makes operations non-blocking
3. **Avoid busy waiting** - don't use select in tight loop with default
4. **Timeout critical operations** - use time.After
5. **Prioritize cancellation** - check done/context first
6. **Empty select blocks forever** - `select {}` for waiting
7. **Combine with context** - for better cancellation
