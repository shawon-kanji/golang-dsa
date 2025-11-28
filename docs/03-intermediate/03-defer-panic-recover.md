# Defer, Panic, and Recover

## Defer

### Defer Basics
```go
// Defer delays execution until surrounding function returns
func example() {
    defer fmt.Println("World")
    fmt.Println("Hello")
}
// Output: Hello, World
```

### Multiple Defers (LIFO - Last In, First Out)
```go
func example() {
    defer fmt.Println("1")
    defer fmt.Println("2")
    defer fmt.Println("3")
}
// Output: 3, 2, 1
```

### Defer Use Cases

#### Resource Cleanup
```go
func processFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()  // Ensures file is closed

    // Process file
    return nil
}
```

#### Mutex Unlock
```go
var mu sync.Mutex

func criticalSection() {
    mu.Lock()
    defer mu.Unlock()  // Ensures unlock even if panic

    // Critical code
}
```

#### Database Transactions
```go
func transfer(from, to Account, amount int) error {
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()  // Rolls back if not committed

    if err := debit(tx, from, amount); err != nil {
        return err
    }
    if err := credit(tx, to, amount); err != nil {
        return err
    }

    return tx.Commit()  // Rollback is no-op after commit
}
```

### Defer with Arguments
```go
func example() {
    x := 10
    defer fmt.Println(x)  // Captures value 10
    x = 20
    fmt.Println(x)        // Prints 20
}
// Output: 20, 10
```

### Defer with Functions
```go
func example() {
    x := 10
    defer func() {
        fmt.Println(x)  // Closure captures variable, not value
    }()
    x = 20
}
// Output: 20
```

### Defer in Loops
```go
// Bad: Defers accumulate
func processFiles(filenames []string) {
    for _, filename := range filenames {
        file, _ := os.Open(filename)
        defer file.Close()  // All defers execute at function end!
    }
}

// Good: Use function to scope defer
func processFiles(filenames []string) {
    for _, filename := range filenames {
        processFile(filename)
    }
}

func processFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()  // Executes when this function returns

    // Process file
    return nil
}
```

## Panic

### Panic Basics
```go
// Panic stops normal execution
func example() {
    fmt.Println("Before panic")
    panic("something went wrong")
    fmt.Println("After panic")  // Never executes
}
```

### When to Panic
```go
// Unrecoverable errors
func divide(a, b int) int {
    if b == 0 {
        panic("division by zero")
    }
    return a / b
}

// Programming errors (bugs)
func process(data []int) {
    if len(data) == 0 {
        panic("data cannot be empty")  // Should never happen
    }
    // Process data
}

// Initialization failures
func init() {
    if !canInitialize() {
        panic("initialization failed")
    }
}
```

### Panic with defer
```go
func example() {
    defer fmt.Println("1")
    defer fmt.Println("2")
    panic("error")
    defer fmt.Println("3")  // Never registered
}
// Output: 2, 1, panic: error
```

## Recover

### Recover Basics
```go
// Recover stops panic and returns panic value
func example() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered from:", r)
        }
    }()

    panic("something went wrong")
    fmt.Println("After panic")  // Never executes
}
// Output: Recovered from: something went wrong
```

### Recover in Goroutines
```go
func safeLaunch(f func()) {
    go func() {
        defer func() {
            if r := recover(); r != nil {
                log.Printf("Recovered in goroutine: %v", r)
            }
        }()
        f()
    }()
}

// Usage
safeLaunch(func() {
    panic("error in goroutine")
})
```

### Recover and Return Error
```go
func safeExecute(f func()) (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic occurred: %v", r)
        }
    }()

    f()
    return nil
}

// Usage
err := safeExecute(func() {
    panic("something went wrong")
})
if err != nil {
    fmt.Println("Error:", err)
}
```

### Re-panic
```go
func example() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovering:", r)
            // Clean up
            // Re-panic if needed
            panic(r)
        }
    }()

    panic("error")
}
```

## Common Patterns

### Safe Goroutine Wrapper
```go
func safeGo(fn func()) {
    go func() {
        defer func() {
            if r := recover(); r != nil {
                log.Printf("Goroutine panic: %v\n%s", r, debug.Stack())
            }
        }()
        fn()
    }()
}
```

### HTTP Handler Panic Recovery
```go
func recoverMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if r := recover(); r != nil {
                log.Printf("Panic: %v\n%s", r, debug.Stack())
                http.Error(w, "Internal Server Error", 500)
            }
        }()
        next.ServeHTTP(w, r)
    })
}
```

### Defer with Named Returns
```go
func example() (result string, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic: %v", r)
        }
    }()

    // Code that might panic
    panic("error")

    return "success", nil
}
```

### Defer with Cleanup
```go
func processWithTimeout(data []byte) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()  // Always cancel context

    // Process with context
    return process(ctx, data)
}
```

## Best Practices

### Defer Order Matters
```go
// Good: Close in reverse order of opening
func example() {
    f1, _ := os.Open("file1.txt")
    defer f1.Close()

    f2, _ := os.Open("file2.txt")
    defer f2.Close()

    // f2 closes first, then f1
}
```

### Error Handling in Defer
```go
func example() (err error) {
    file, err := os.Open("file.txt")
    if err != nil {
        return err
    }
    defer func() {
        closeErr := file.Close()
        if err == nil {
            err = closeErr
        }
    }()

    // Process file
    return nil
}
```

### Don't Panic in Libraries
```go
// Bad: Library code panics
func LibraryFunction(data []byte) {
    if len(data) == 0 {
        panic("empty data")
    }
}

// Good: Return error
func LibraryFunction(data []byte) error {
    if len(data) == 0 {
        return errors.New("empty data")
    }
    return nil
}
```

## Gotchas

❌ **Defer evaluation timing**
```go
func example() {
    x := 1
    defer fmt.Println(x)  // Prints 1, not 2
    x = 2
}
```

❌ **Recover only works in defer**
```go
func example() {
    recover()  // Does nothing!
    panic("error")
}

func example() {
    defer recover()  // Still does nothing!
    panic("error")
}

func example() {
    defer func() {
        recover()  // Works!
    }()
    panic("error")
}
```

❌ **Defer in loop**
```go
// Memory leak: defers accumulate
for _, file := range files {
    f, _ := os.Open(file)
    defer f.Close()  // Defers pile up!
}

// Fix: Use separate function
for _, file := range files {
    processFile(file)
}
```

## When to Use Each

### Use Defer When:
- Cleaning up resources (files, locks, connections)
- Logging function entry/exit
- Recovering from panics
- Ensuring code runs regardless of return path

### Use Panic When:
- Unrecoverable errors (can't continue)
- Programming errors (bugs)
- Initialization failures
- Fatal configuration errors

### Use Recover When:
- Top-level error handling (servers, goroutines)
- Converting panics to errors at boundaries
- Logging panics before crashing
- Testing code that might panic

## Real-World Examples

### Database Connection Pool
```go
func executeQuery(query string) error {
    conn, err := pool.Get()
    if err != nil {
        return err
    }
    defer pool.Put(conn)

    return conn.Execute(query)
}
```

### Timing Function Execution
```go
func timeTrack(name string) func() {
    start := time.Now()
    return func() {
        fmt.Printf("%s took %v\n", name, time.Since(start))
    }
}

func example() {
    defer timeTrack("example")()
    // Function logic
}
```

### Trace Function Calls
```go
func trace(name string) func() {
    fmt.Printf("Entering %s\n", name)
    return func() {
        fmt.Printf("Exiting %s\n", name)
    }
}

func example() {
    defer trace("example")()
    // Function logic
}
```
