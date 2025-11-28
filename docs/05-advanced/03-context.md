# Context

## Context Basics

### Creating Contexts
```go
import "context"

// Background context (root)
ctx := context.Background()

// TODO context (placeholder)
ctx := context.TODO()
```

### Context with Cancel
```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()  // Always call cancel

go func() {
    <-ctx.Done()
    fmt.Println("Cancelled:", ctx.Err())
}()

// Cancel the context
cancel()
```

### Context with Timeout
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

select {
case <-time.After(3 * time.Second):
    fmt.Println("Work completed")
case <-ctx.Done():
    fmt.Println("Timeout:", ctx.Err())
}
```

### Context with Deadline
```go
deadline := time.Now().Add(10 * time.Second)
ctx, cancel := context.WithDeadline(context.Background(), deadline)
defer cancel()

select {
case <-doWork():
    fmt.Println("Done")
case <-ctx.Done():
    fmt.Println("Deadline exceeded:", ctx.Err())
}
```

## Context with Values

### Storing Values
```go
type key string

const userKey key = "user"

ctx := context.WithValue(context.Background(), userKey, "Alice")

// Retrieve value
if user, ok := ctx.Value(userKey).(string); ok {
    fmt.Println("User:", user)
}
```

### Type-Safe Context Values
```go
type contextKey int

const (
    userIDKey contextKey = iota
    requestIDKey
)

func WithUserID(ctx context.Context, userID int) context.Context {
    return context.WithValue(ctx, userIDKey, userID)
}

func GetUserID(ctx context.Context) (int, bool) {
    userID, ok := ctx.Value(userIDKey).(int)
    return userID, ok
}
```

## Common Patterns

### HTTP Request Context
```go
func handler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    // Add values
    ctx = context.WithValue(ctx, "requestID", generateID())

    // Pass to functions
    result, err := fetchData(ctx)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    w.Write(result)
}
```

### Database Operations
```go
func queryWithContext(ctx context.Context, query string) ([]Row, error) {
    rows, err := db.QueryContext(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var results []Row
    for rows.Next() {
        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        default:
            var row Row
            if err := rows.Scan(&row); err != nil {
                return nil, err
            }
            results = append(results, row)
        }
    }

    return results, nil
}
```

### Worker with Context
```go
func worker(ctx context.Context) {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            fmt.Println("Worker stopping:", ctx.Err())
            return
        case <-ticker.C:
            fmt.Println("Working...")
        }
    }
}

// Usage
ctx, cancel := context.WithCancel(context.Background())
go worker(ctx)

time.Sleep(5 * time.Second)
cancel()  // Stop worker
```

### Fan-Out with Context
```go
func fanOut(ctx context.Context, input <-chan int, n int) []<-chan int {
    channels := make([]<-chan int, n)
    for i := 0; i < n; i++ {
        ch := make(chan int)
        channels[i] = ch
        go func(out chan<- int) {
            defer close(out)
            for {
                select {
                case <-ctx.Done():
                    return
                case val, ok := <-input:
                    if !ok {
                        return
                    }
                    out <- process(val)
                }
            }
        }(ch)
    }
    return channels
}
```

## Best Practices

1. **Pass context as first parameter** - `func DoSomething(ctx context.Context, ...)`
2. **Don't store contexts in structs** - pass explicitly
3. **Always call cancel()** - use defer
4. **Use context.Background() at top level**
5. **Use context.TODO() as placeholder**
6. **Context values for request-scoped** data only
7. **Don't use context for optional parameters**
8. **Check ctx.Done() in loops**
9. **Propagate context** through call chain
10. **Document context usage** in functions

## Common Errors

❌ **Storing context in struct**
```go
// Bad
type Service struct {
    ctx context.Context
}

// Good
func (s *Service) DoWork(ctx context.Context) {}
```

❌ **Not canceling context**
```go
// Bad
ctx, _ := context.WithCancel(parent)

// Good
ctx, cancel := context.WithCancel(parent)
defer cancel()
```

❌ **Using context for optional parameters**
```go
// Bad
ctx := context.WithValue(ctx, "timeout", 5)

// Good - use function parameter
func DoWork(ctx context.Context, timeout time.Duration) {}
```

## Context Patterns

### Chaining Contexts
```go
ctx := context.Background()
ctx, cancel1 := context.WithTimeout(ctx, 10*time.Second)
defer cancel1()

ctx, cancel2 := context.WithCancel(ctx)
defer cancel2()
```

### Context with Multiple Cancellations
```go
func workWithMultipleCancels(ctx context.Context) error {
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    go func() {
        select {
        case <-time.After(5 * time.Second):
            cancel()  // Cancel after timeout
        case <-externalEvent:
            cancel()  // Cancel on external event
        }
    }()

    return doWork(ctx)
}
```

### Context Cleanup
```go
func processWithCleanup(ctx context.Context) error {
    resource, err := allocate()
    if err != nil {
        return err
    }

    done := make(chan struct{})
    go func() {
        select {
        case <-ctx.Done():
            resource.Close()
        case <-done:
        }
    }()

    defer close(done)
    return process(resource)
}
```
