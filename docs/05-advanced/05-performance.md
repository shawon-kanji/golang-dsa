# Performance Optimization

## Profiling

### CPU Profiling
```go
import (
    "runtime/pprof"
    "os"
)

func main() {
    f, _ := os.Create("cpu.prof")
    defer f.Close()

    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()

    // Your code
}

// Analyze: go tool pprof cpu.prof
```

### Memory Profiling
```go
import "runtime/pprof"

func writeMemProfile() {
    f, _ := os.Create("mem.prof")
    defer f.Close()

    runtime.GC()  // Get up-to-date statistics
    pprof.WriteHeapProfile(f)
}
```

### HTTP Profiling
```go
import _ "net/http/pprof"
import "net/http"

func main() {
    go func() {
        http.ListenAndServe("localhost:6060", nil)
    }()

    // Visit http://localhost:6060/debug/pprof/
}
```

## Benchmarking

### Writing Benchmarks
```go
func BenchmarkConcatenate(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = "hello" + "world"
    }
}

func BenchmarkStringsBuilder(b *testing.B) {
    for i := 0; i < b.N; i++ {
        var sb strings.Builder
        sb.WriteString("hello")
        sb.WriteString("world")
        _ = sb.String()
    }
}
```

### Benchmark Comparison
```bash
# Run benchmarks
go test -bench=. -benchmem

# Compare benchmarks
go test -bench=. -benchmem > old.txt
# Make changes
go test -bench=. -benchmem > new.txt
benchcmp old.txt new.txt
```

## Memory Optimization

### Pre-allocate Slices
```go
// Bad: Multiple allocations
var items []int
for i := 0; i < 1000; i++ {
    items = append(items, i)
}

// Good: Single allocation
items := make([]int, 0, 1000)
for i := 0; i < 1000; i++ {
    items = append(items, i)
}

// Best: No append needed
items := make([]int, 1000)
for i := 0; i < 1000; i++ {
    items[i] = i
}
```

### String Builder
```go
// Bad: Creates many strings
var s string
for i := 0; i < 1000; i++ {
    s += "a"
}

// Good: Efficient concatenation
var sb strings.Builder
sb.Grow(1000)  // Pre-allocate
for i := 0; i < 1000; i++ {
    sb.WriteString("a")
}
s := sb.String()
```

### Avoid Pointer Indirection
```go
// Consider value vs pointer based on size

// Small struct - value might be faster
type Small struct {
    a, b int
}

// Large struct - pointer likely faster
type Large struct {
    data [1000]int
}
```

### sync.Pool for Reuse
```go
var pool = sync.Pool{
    New: func() interface{} {
        return &bytes.Buffer{}
    },
}

func processData(data []byte) {
    buf := pool.Get().(*bytes.Buffer)
    defer pool.Put(buf)

    buf.Reset()
    buf.Write(data)
    // Process buffer
}
```

## CPU Optimization

### Avoid Allocations in Hot Paths
```go
// Bad: Allocates in loop
for i := 0; i < 1000000; i++ {
    s := fmt.Sprintf("value: %d", i)
    process(s)
}

// Good: Reuse buffer
var buf bytes.Buffer
for i := 0; i < 1000000; i++ {
    buf.Reset()
    buf.WriteString("value: ")
    buf.WriteString(strconv.Itoa(i))
    process(buf.String())
}
```

### Use Buffered I/O
```go
// Bad: Unbuffered
file, _ := os.Open("large.txt")
defer file.Close()
// Direct read/write

// Good: Buffered
file, _ := os.Open("large.txt")
defer file.Close()
reader := bufio.NewReader(file)
// Use reader
```

### Concurrent Processing
```go
func processParallel(items []Item) {
    numWorkers := runtime.NumCPU()
    jobs := make(chan Item, len(items))
    results := make(chan Result, len(items))

    // Start workers
    for w := 0; w < numWorkers; w++ {
        go worker(jobs, results)
    }

    // Send jobs
    for _, item := range items {
        jobs <- item
    }
    close(jobs)

    // Collect results
    for i := 0; i < len(items); i++ {
        <-results
    }
}
```

## Compiler Optimizations

### Inline Functions
```go
// Small functions may be inlined
//go:inline
func add(a, b int) int {
    return a + b
}
```

### Bounds Check Elimination
```go
// Compiler can eliminate bounds checks
func sum(s []int) int {
    total := 0
    for i := 0; i < len(s); i++ {
        total += s[i]  // Bounds check eliminated
    }
    return total
}
```

### Escape Analysis
```go
// Allocated on stack (fast)
func localVar() int {
    x := 10
    return x
}

// Escapes to heap (slower)
func heapVar() *int {
    x := 10
    return &x  // Escapes
}
```

## Best Practices

1. **Profile before optimizing** - measure don't guess
2. **Optimize hot paths** - focus on bottlenecks
3. **Pre-allocate when size known** - avoid resizing
4. **Use string.Builder** for concatenation
5. **Reuse with sync.Pool** - for frequently allocated objects
6. **Buffered I/O** - for file operations
7. **Limit allocations** - in tight loops
8. **Use goroutines wisely** - they have overhead
9. **Cache expensive operations** - when possible
10. **Benchmark changes** - verify improvements

## Common Pitfalls

❌ **Premature optimization**
```go
// Don't optimize before profiling
```

❌ **Over-engineering**
```go
// Simple code is often fast enough
```

❌ **Ignoring readability**
```go
// Maintainable code > slightly faster code
```

## Tools

```bash
# Profiling
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof

# Memory profiling
go test -memprofile=mem.prof -bench=.
go tool pprof mem.prof

# Race detection
go test -race

# Escape analysis
go build -gcflags="-m" file.go

# Assembly output
go tool compile -S file.go
```
