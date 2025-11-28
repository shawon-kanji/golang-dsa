# Type Assertions and Type Switches

## Type Assertions

### Basic Type Assertion
```go
var i interface{} = "hello"

// Type assertion
s := i.(string)
fmt.Println(s)  // "hello"

// Panic if wrong type
// n := i.(int)  // Panic: interface conversion

// Safe type assertion
s, ok := i.(string)
if ok {
    fmt.Println(s)
} else {
    fmt.Println("not a string")
}
```

### Checking Interface Implementation
```go
type Reader interface {
    Read([]byte) (int, error)
}

var r Reader = &bytes.Buffer{}

// Check if implements another interface
if w, ok := r.(io.Writer); ok {
    w.Write([]byte("hello"))
}
```

### Asserting to Concrete Type
```go
type Dog struct {
    Name string
}

func (d Dog) Speak() string {
    return "Woof!"
}

var animal interface{} = Dog{Name: "Buddy"}

// Assert to concrete type
if dog, ok := animal.(Dog); ok {
    fmt.Println(dog.Name)  // "Buddy"
}
```

## Type Switches

### Basic Type Switch
```go
func describe(i interface{}) {
    switch v := i.(type) {
    case int:
        fmt.Printf("Integer: %d\n", v)
    case string:
        fmt.Printf("String: %s\n", v)
    case bool:
        fmt.Printf("Boolean: %t\n", v)
    default:
        fmt.Printf("Unknown type: %T\n", v)
    }
}

describe(42)      // Integer: 42
describe("hello") // String: hello
describe(true)    // Boolean: true
describe(3.14)    // Unknown type: float64
```

### Multiple Types in One Case
```go
func checkType(i interface{}) {
    switch i.(type) {
    case int, int8, int16, int32, int64:
        fmt.Println("Some kind of integer")
    case float32, float64:
        fmt.Println("Some kind of float")
    case string:
        fmt.Println("String")
    default:
        fmt.Println("Unknown")
    }
}
```

### Type Switch with Value
```go
func process(i interface{}) {
    switch v := i.(type) {
    case int:
        fmt.Println("Int:", v*2)
    case string:
        fmt.Println("String:", strings.ToUpper(v))
    case []int:
        fmt.Println("Slice sum:", sum(v))
    case nil:
        fmt.Println("Nil value")
    default:
        fmt.Printf("Unknown type: %T\n", v)
    }
}
```

### Type Switch Without Assignment
```go
func checkType(i interface{}) {
    switch i.(type) {
    case string:
        fmt.Println("It's a string")
    case int:
        fmt.Println("It's an int")
    default:
        fmt.Println("Unknown type")
    }
}
```

## Common Patterns

### JSON Unmarshal Dynamic Types
```go
func handleJSON(data []byte) error {
    var v interface{}
    if err := json.Unmarshal(data, &v); err != nil {
        return err
    }

    switch val := v.(type) {
    case map[string]interface{}:
        // Handle object
        for k, v := range val {
            fmt.Printf("%s: %v\n", k, v)
        }
    case []interface{}:
        // Handle array
        for i, item := range val {
            fmt.Printf("[%d]: %v\n", i, item)
        }
    case string:
        fmt.Println("String:", val)
    case float64:
        // JSON numbers are float64
        fmt.Println("Number:", val)
    case bool:
        fmt.Println("Boolean:", val)
    case nil:
        fmt.Println("Null")
    }

    return nil
}
```

### Type-safe Wrapper
```go
type Value struct {
    data interface{}
}

func (v Value) AsString() (string, error) {
    if s, ok := v.data.(string); ok {
        return s, nil
    }
    return "", fmt.Errorf("not a string")
}

func (v Value) AsInt() (int, error) {
    if i, ok := v.data.(int); ok {
        return i, nil
    }
    return 0, fmt.Errorf("not an int")
}
```

### Interface Method Checking
```go
func writeData(w io.Writer, data []byte) error {
    // Check if supports WriterTo for efficiency
    if wt, ok := w.(io.WriterTo); ok {
        _, err := wt.WriteTo(w)
        return err
    }

    // Fallback to standard Write
    _, err := w.Write(data)
    return err
}
```

### Optional Interface Methods
```go
type Closer interface {
    Close() error
}

func cleanup(v interface{}) {
    if closer, ok := v.(Closer); ok {
        closer.Close()
    }
}
```

## Error Type Assertions

```go
// Check specific error type
if pathErr, ok := err.(*os.PathError); ok {
    fmt.Println("Path:", pathErr.Path)
    fmt.Println("Op:", pathErr.Op)
    fmt.Println("Err:", pathErr.Err)
}

// Using errors.As (preferred in Go 1.13+)
var pathErr *os.PathError
if errors.As(err, &pathErr) {
    fmt.Println("Path error:", pathErr.Path)
}
```

## Reflection Alternative

```go
// Type assertion is faster than reflection
// Prefer type assertion when possible

// Type assertion - fast
if s, ok := v.(string); ok {
    fmt.Println(s)
}

// Reflection - slower
import "reflect"
if reflect.TypeOf(v).Kind() == reflect.String {
    fmt.Println(v)
}
```

## Best Practices

### Prefer Type Switches for Multiple Types
```go
// Good
func process(i interface{}) {
    switch v := i.(type) {
    case int:
        // handle int
    case string:
        // handle string
    case bool:
        // handle bool
    }
}

// Less ideal
func process(i interface{}) {
    if v, ok := i.(int); ok {
        // handle int
    } else if v, ok := i.(string); ok {
        // handle string
    } else if v, ok := i.(bool); ok {
        // handle bool
    }
}
```

### Always Check Type Assertions
```go
// Bad - can panic
s := i.(string)

// Good - safe
if s, ok := i.(string); ok {
    // use s
}
```

### Use Specific Interfaces
```go
// Bad - too generic
func process(data interface{}) {
    // Need many type assertions
}

// Good - specific interface
func process(r io.Reader) {
    // No type assertions needed
}
```

## Gotchas

❌ **Type assertion on nil interface**
```go
var i interface{}  // nil
// s, ok := i.(string)  // ok is false, s is ""
```

❌ **Forgetting to check ok**
```go
s := i.(string)  // Panics if i is not string

// Always check
if s, ok := i.(string); ok {
    // Safe to use s
}
```

❌ **Type switch with fallthrough**
```go
switch v := i.(type) {
case int:
    fmt.Println("Int")
    // fallthrough not allowed in type switch!
case string:
    fmt.Println("String")
}
```

## Performance Considerations

```go
// Type assertions are fast
// Type switches are optimized by compiler
// Both are much faster than reflection

// Benchmark: type assertion vs reflection
// Type assertion: ~1-2 ns/op
// Reflection: ~100+ ns/op
```

## When to Use

✅ **Use type assertions when:**
- Working with interface{} or empty interfaces
- Need to extract concrete type
- Checking optional interface methods
- Parsing dynamic data (JSON, etc.)

❌ **Avoid when:**
- Can use specific types instead
- Design can use interfaces properly
- Becomes too complex with many assertions
