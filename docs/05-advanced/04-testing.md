# Testing

## Test Basics

### Simple Test
```go
// in math.go
func Add(a, b int) int {
    return a + b
}

// in math_test.go
import "testing"

func TestAdd(t *testing.T) {
    result := Add(2, 3)
    expected := 5
    if result != expected {
        t.Errorf("Add(2, 3) = %d; want %d", result, expected)
    }
}
```

### Table-Driven Tests
```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive", 2, 3, 5},
        {"negative", -1, -2, -3},
        {"zero", 0, 0, 0},
        {"mixed", -5, 10, 5},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Add(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("got %d, want %d", result, tt.expected)
            }
        })
    }
}
```

## Testing Functions

### t.Error vs t.Fatal
```go
func TestSomething(t *testing.T) {
    t.Error("continues test")   // Log error, continue
    t.Fatal("stops test")        // Log error, stop test
    t.Log("info message")        // Log message
    t.Skip("skip this test")     // Skip test
}
```

### Helper Functions
```go
func assertEqual(t *testing.T, got, want interface{}) {
    t.Helper()  // Mark as helper
    if got != want {
        t.Errorf("got %v, want %v", got, want)
    }
}

func TestAdd(t *testing.T) {
    assertEqual(t, Add(2, 3), 5)
}
```

## Benchmarking

### Basic Benchmark
```go
func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Add(2, 3)
    }
}

// Run: go test -bench=.
// Output: BenchmarkAdd-8   1000000000   0.25 ns/op
```

### Benchmark with Setup
```go
func BenchmarkFibonacci(b *testing.B) {
    data := setupData()  // Setup before timer

    b.ResetTimer()  // Reset timer after setup

    for i := 0; i < b.N; i++ {
        fibonacci(data)
    }
}
```

### Benchmark Reporting
```go
func BenchmarkProcess(b *testing.B) {
    b.ReportAllocs()  // Report allocations

    for i := 0; i < b.N; i++ {
        process(data)
    }

    b.SetBytes(1024)  // Set processed bytes
}
```

## Test Coverage

```bash
# Run tests with coverage
go test -cover

# Generate coverage profile
go test -coverprofile=coverage.out

# View coverage in browser
go tool cover -html=coverage.out
```

## Mocking

### Interface Mocking
```go
type Database interface {
    Query(string) ([]byte, error)
}

type MockDB struct {
    QueryFunc func(string) ([]byte, error)
}

func (m *MockDB) Query(q string) ([]byte, error) {
    return m.QueryFunc(q)
}

// Test
func TestService(t *testing.T) {
    mock := &MockDB{
        QueryFunc: func(q string) ([]byte, error) {
            return []byte("mocked"), nil
        },
    }

    service := NewService(mock)
    result := service.DoWork()
    // Assert result
}
```

## Testing HTTP

### HTTP Handler Testing
```go
func TestHandler(t *testing.T) {
    req := httptest.NewRequest("GET", "/api/users", nil)
    w := httptest.NewRecorder()

    handler(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("got status %d, want %d", w.Code, http.StatusOK)
    }

    expected := `{"status":"ok"}`
    if w.Body.String() != expected {
        t.Errorf("got %s, want %s", w.Body.String(), expected)
    }
}
```

### HTTP Server Testing
```go
func TestServer(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(handler))
    defer server.Close()

    resp, err := http.Get(server.URL)
    if err != nil {
        t.Fatal(err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        t.Errorf("got status %d", resp.StatusCode)
    }
}
```

## Test Organization

### Setup and Teardown
```go
func TestMain(m *testing.M) {
    // Setup
    setup()

    // Run tests
    code := m.Run()

    // Teardown
    teardown()

    os.Exit(code)
}
```

### Subtests
```go
func TestFeature(t *testing.T) {
    t.Run("subtest1", func(t *testing.T) {
        // Test case 1
    })

    t.Run("subtest2", func(t *testing.T) {
        // Test case 2
    })
}
```

## Best Practices

1. **Test file naming** - `*_test.go`
2. **Test function naming** - `TestXxx`
3. **Use table-driven tests** - for multiple cases
4. **Test edge cases** - empty, nil, zero values
5. **Use t.Helper()** for test helpers
6. **Parallel tests** - use `t.Parallel()`
7. **Clear test names** - descriptive subtests
8. **Mock external dependencies**
9. **Test public API** - not implementation
10. **Keep tests simple** - one assertion per test

## Common Commands

```bash
# Run all tests
go test ./...

# Run specific test
go test -run TestAdd

# Verbose output
go test -v

# Run benchmarks
go test -bench=.

# Coverage
go test -cover

# Race detection
go test -race

# Short mode (skip long tests)
go test -short
```
