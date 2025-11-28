# Sync Package

## sync.Mutex

### Basic Mutex
```go
var (
    mu      sync.Mutex
    counter int
)

func increment() {
    mu.Lock()
    counter++
    mu.Unlock()
}

// With defer
func increment() {
    mu.Lock()
    defer mu.Unlock()
    counter++
}
```

## sync.RWMutex

### Read-Write Lock
```go
var (
    mu    sync.RWMutex
    data  map[string]string
)

// Read lock (multiple readers allowed)
func read(key string) string {
    mu.RLock()
    defer mu.RUnlock()
    return data[key]
}

// Write lock (exclusive)
func write(key, value string) {
    mu.Lock()
    defer mu.Unlock()
    data[key] = value
}
```

## sync.WaitGroup

### Basic WaitGroup
```go
var wg sync.WaitGroup

for i := 0; i < 5; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        fmt.Printf("Goroutine %d\n", id)
    }(i)
}

wg.Wait()  // Wait for all goroutines
```

## sync.Once

### Execute Once
```go
var (
    once     sync.Once
    instance *Singleton
)

func GetInstance() *Singleton {
    once.Do(func() {
        instance = &Singleton{}
        // Expensive initialization
    })
    return instance
}
```

## sync.Map

### Concurrent Map
```go
var m sync.Map

// Store
m.Store("key", "value")

// Load
val, ok := m.Load("key")

// LoadOrStore
actual, loaded := m.LoadOrStore("key", "value")

// Delete
m.Delete("key")

// Range
m.Range(func(key, value interface{}) bool {
    fmt.Printf("%v: %v\n", key, value)
    return true  // Continue iteration
})
```

## sync.Pool

### Object Pooling
```go
var pool = sync.Pool{
    New: func() interface{} {
        return &Buffer{}
    },
}

// Get from pool
buf := pool.Get().(*Buffer)
defer pool.Put(buf)  // Return to pool

// Use buf
buf.Reset()
buf.Write(data)
```

## sync.Cond

### Condition Variables
```go
var (
    cond  = sync.NewCond(&sync.Mutex{})
    ready = false
)

// Waiter
func waiter() {
    cond.L.Lock()
    for !ready {
        cond.Wait()
    }
    // Proceed when ready
    cond.L.Unlock()
}

// Signaler
func signaler() {
    cond.L.Lock()
    ready = true
    cond.Signal()  // or cond.Broadcast()
    cond.L.Unlock()
}
```

## Best Practices

1. **Use defer for unlock** - ensures unlock even with panic
2. **Keep critical sections small** - minimize time holding lock
3. **Prefer channels** for communication
4. **Use RWMutex** when reads dominate
5. **Avoid nested locks** - can cause deadlock
6. **Don't copy mutexes** - pass by pointer
7. **Use sync.Once** for lazy initialization
8. **sync.Pool for frequently allocated objects**
