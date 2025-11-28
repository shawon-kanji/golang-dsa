# Error Handling

## Error Basics

### The error Interface
```go
type error interface {
    Error() string
}
```

### Creating Errors
```go
import "errors"

// Simple error
err := errors.New("something went wrong")

// Formatted error
err := fmt.Errorf("failed to process %s", filename)

// With value wrapping (Go 1.13+)
err := fmt.Errorf("process failed: %w", originalError)
```

### Returning Errors
```go
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

// Usage
result, err := divide(10, 0)
if err != nil {
    fmt.Println("Error:", err)
    return
}
fmt.Println("Result:", result)
```

## Error Checking

### Standard Pattern
```go
file, err := os.Open("file.txt")
if err != nil {
    return err  // Or handle it
}
defer file.Close()

// Use file
```

### Multiple Errors
```go
file, err := os.Open("file.txt")
if err != nil {
    return fmt.Errorf("failed to open file: %w", err)
}
defer file.Close()

data, err := io.ReadAll(file)
if err != nil {
    return fmt.Errorf("failed to read file: %w", err)
}

// Process data
```

### Early Return Pattern
```go
func processData(filename string) error {
    data, err := loadData(filename)
    if err != nil {
        return err
    }

    if err := validate(data); err != nil {
        return err
    }

    if err := transform(data); err != nil {
        return err
    }

    if err := save(data); err != nil {
        return err
    }

    return nil
}
```

## Custom Errors

### Simple Custom Error
```go
type MyError struct {
    Code    int
    Message string
}

func (e MyError) Error() string {
    return fmt.Sprintf("error %d: %s", e.Code, e.Message)
}

func doSomething() error {
    return MyError{Code: 404, Message: "not found"}
}
```

### Error with Context
```go
type ValidationError struct {
    Field string
    Value interface{}
    Msg   string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation failed for %s (%v): %s",
        e.Field, e.Value, e.Msg)
}

func validateAge(age int) error {
    if age < 0 {
        return ValidationError{
            Field: "age",
            Value: age,
            Msg:   "must be non-negative",
        }
    }
    return nil
}
```

### Sentinel Errors
```go
var (
    ErrNotFound    = errors.New("not found")
    ErrUnauthorized = errors.New("unauthorized")
    ErrInvalidInput = errors.New("invalid input")
)

func findUser(id int) (*User, error) {
    // ... search logic ...
    if !found {
        return nil, ErrNotFound
    }
    return user, nil
}

// Check specific error
user, err := findUser(123)
if err == ErrNotFound {
    // Handle not found
} else if err != nil {
    // Handle other errors
}
```

## Error Wrapping (Go 1.13+)

### Wrapping Errors
```go
func readConfig(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return fmt.Errorf("failed to open config: %w", err)
    }
    defer file.Close()

    // ... read and parse ...

    if err := parse(data); err != nil {
        return fmt.Errorf("failed to parse config: %w", err)
    }

    return nil
}
```

### Unwrapping Errors
```go
// errors.Is - check error in chain
if errors.Is(err, os.ErrNotExist) {
    fmt.Println("File not found")
}

// errors.As - get specific error type
var pathErr *os.PathError
if errors.As(err, &pathErr) {
    fmt.Println("Failed path:", pathErr.Path)
}

// errors.Unwrap - get wrapped error
originalErr := errors.Unwrap(err)
```

### Custom Unwrap
```go
type QueryError struct {
    Query string
    Err   error
}

func (e *QueryError) Error() string {
    return fmt.Sprintf("query error: %v", e.Err)
}

func (e *QueryError) Unwrap() error {
    return e.Err
}

// Now works with errors.Is and errors.As
```

## Multiple Errors

### Collecting Errors
```go
func validateAll(items []Item) []error {
    var errors []error
    for i, item := range items {
        if err := validate(item); err != nil {
            errors = append(errors, fmt.Errorf("item %d: %w", i, err))
        }
    }
    return errors
}

// Check if any errors
errors := validateAll(items)
if len(errors) > 0 {
    for _, err := range errors {
        fmt.Println(err)
    }
    return
}
```

### MultiError Type
```go
type MultiError []error

func (m MultiError) Error() string {
    if len(m) == 0 {
        return "no errors"
    }
    if len(m) == 1 {
        return m[0].Error()
    }
    return fmt.Sprintf("%d errors occurred: %v", len(m), m)
}

func (m MultiError) OrNil() error {
    if len(m) == 0 {
        return nil
    }
    return m
}
```

## Error Handling Patterns

### Retry Pattern
```go
func retryOperation(maxRetries int, operation func() error) error {
    var err error
    for i := 0; i < maxRetries; i++ {
        err = operation()
        if err == nil {
            return nil
        }
        time.Sleep(time.Second * time.Duration(i+1))
    }
    return fmt.Errorf("failed after %d retries: %w", maxRetries, err)
}
```

### Fallback Pattern
```go
func fetchData() ([]byte, error) {
    data, err := fetchFromPrimary()
    if err != nil {
        log.Printf("Primary failed: %v, trying fallback", err)
        data, err = fetchFromFallback()
        if err != nil {
            return nil, fmt.Errorf("both primary and fallback failed: %w", err)
        }
    }
    return data, nil
}
```

### Cleanup with Errors
```go
func processFile(filename string) (err error) {
    file, err := os.Open(filename)
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

### Error Aggregation
```go
func processAll(items []string) error {
    var firstErr error
    for _, item := range items {
        if err := process(item); err != nil {
            if firstErr == nil {
                firstErr = err
            }
            log.Printf("Failed to process %s: %v", item, err)
        }
    }
    return firstErr
}
```

## Error Context

### Adding Context
```go
func loadUser(id int) (*User, error) {
    user, err := db.FindByID(id)
    if err != nil {
        return nil, fmt.Errorf("load user %d: %w", id, err)
    }
    return user, nil
}
```

### Structured Context
```go
type ErrorContext struct {
    Operation string
    UserID    int
    Time      time.Time
    Err       error
}

func (e ErrorContext) Error() string {
    return fmt.Sprintf("[%s] operation=%s user=%d: %v",
        e.Time.Format(time.RFC3339), e.Operation, e.UserID, e.Err)
}
```

## Error Types

### Temporary Errors
```go
type temporary interface {
    Temporary() bool
}

func isTemporary(err error) bool {
    te, ok := err.(temporary)
    return ok && te.Temporary()
}

// Usage
if isTemporary(err) {
    // Retry
} else {
    // Give up
}
```

### Timeout Errors
```go
type timeout interface {
    Timeout() bool
}

func isTimeout(err error) bool {
    te, ok := err.(timeout)
    return ok && te.Timeout()
}
```

## Best Practices

### Error Messages
```go
// Good: Lowercase, no punctuation
err := errors.New("connection failed")

// Bad: Uppercase, with punctuation
err := errors.New("Connection failed.")

// Good: Add context
err := fmt.Errorf("connection to %s failed: %w", host, err)
```

### Error Checking
```go
// Good: Check immediately
data, err := readFile()
if err != nil {
    return err
}

// Bad: Delay checking
data, err := readFile()
process(data)  // Might be nil!
if err != nil {
    return err
}
```

### Error Return Position
```go
// Good: Error is last return value
func doSomething() (result string, err error)

// Bad: Error not last
func doSomething() (err error, result string)
```

### Named Error Variables
```go
// Good: Start with Err
var (
    ErrNotFound   = errors.New("not found")
    ErrDuplicate  = errors.New("duplicate")
)

// Bad: Inconsistent naming
var (
    NotFoundError = errors.New("not found")
    ErrDuplicate  = errors.New("duplicate")
)
```

## Common Patterns

### Guard Clauses
```go
func process(data []byte) error {
    if len(data) == 0 {
        return errors.New("empty data")
    }

    if !isValid(data) {
        return errors.New("invalid data")
    }

    // Main logic
    return nil
}
```

### Error Factory
```go
func newDatabaseError(operation string, err error) error {
    return fmt.Errorf("database %s failed: %w", operation, err)
}

// Usage
if err := db.Insert(record); err != nil {
    return newDatabaseError("insert", err)
}
```

### Error Logging
```go
func handleRequest(r *http.Request) error {
    err := processRequest(r)
    if err != nil {
        log.Printf("Request processing failed: %v", err)
        return fmt.Errorf("request failed: %w", err)
    }
    return nil
}
```

## Gotchas

❌ **Ignoring errors**
```go
// Bad
file, _ := os.Open("file.txt")

// Good
file, err := os.Open("file.txt")
if err != nil {
    return err
}
```

❌ **Shadowing err variable**
```go
// Bad
var result string
if data, err := fetch(); err == nil {
    result = data  // err shadowed, not available outside
}
// Can't check err here!

// Good
data, err := fetch()
if err == nil {
    result = data
}
```

❌ **Comparing wrapped errors with ==**
```go
err := fmt.Errorf("wrapped: %w", ErrNotFound)
// err == ErrNotFound  // false!

// Use errors.Is instead
errors.Is(err, ErrNotFound)  // true
```

## Testing Errors

```go
func TestDivideByZero(t *testing.T) {
    _, err := divide(10, 0)
    if err == nil {
        t.Fatal("expected error, got nil")
    }

    expected := "division by zero"
    if err.Error() != expected {
        t.Errorf("expected %q, got %q", expected, err.Error())
    }
}

// Testing sentinel errors
func TestNotFound(t *testing.T) {
    _, err := findUser(999)
    if !errors.Is(err, ErrNotFound) {
        t.Errorf("expected ErrNotFound, got %v", err)
    }
}

// Testing error types
func TestCustomError(t *testing.T) {
    err := doSomething()
    var myErr *MyError
    if !errors.As(err, &myErr) {
        t.Fatal("expected MyError type")
    }
    if myErr.Code != 404 {
        t.Errorf("expected code 404, got %d", myErr.Code)
    }
}
```

## When Not to Use Errors

```go
// Don't use errors for control flow
// Bad
func find(items []int, target int) (int, error) {
    for i, v := range items {
        if v == target {
            return i, nil
        }
    }
    return -1, errors.New("not found")
}

// Good - use bool or special value
func find(items []int, target int) (int, bool) {
    for i, v := range items {
        if v == target {
            return i, true
        }
    }
    return -1, false
}
```
