# Variables and Types

## Basic Types

### Numeric Types
```go
// Integers
int8, int16, int32, int64  // Signed
uint8, uint16, uint32, uint64  // Unsigned
int, uint  // Platform dependent (32 or 64 bit)
byte  // alias for uint8
rune  // alias for int32 (Unicode code point)

// Floating point
float32, float64

// Complex numbers
complex64, complex128
```

### String and Boolean
```go
string  // UTF-8 encoded
bool    // true or false
```

## Variable Declaration

### Method 1: var keyword
```go
var name string
var age int
var isActive bool

// With initialization
var name string = "John"
var age int = 30

// Type inference
var name = "John"  // Type inferred as string
var age = 30       // Type inferred as int

// Multiple declarations
var (
    name   string = "John"
    age    int    = 30
    active bool   = true
)
```

### Method 2: Short declaration `:=` (Inside functions only)
```go
name := "John"
age := 30
isActive := true

// Multiple assignment
x, y := 10, 20
name, age := "John", 30
```

### Method 3: Multiple variables
```go
var a, b, c int = 1, 2, 3
x, y, z := 10, 20, 30
```

## Zero Values
Variables without explicit initialization get zero values:
```go
var i int       // 0
var f float64   // 0.0
var b bool      // false
var s string    // "" (empty string)
var p *int      // nil
```

## Constants
```go
const Pi = 3.14159
const (
    StatusOK    = 200
    StatusError = 500
)

// Typed constants
const Pi float64 = 3.14159

// iota - auto-incrementing constant
const (
    Sunday    = iota  // 0
    Monday           // 1
    Tuesday          // 2
    Wednesday        // 3
)

const (
    _  = iota        // Skip 0
    KB = 1 << (10 * iota)  // 1024
    MB                      // 1048576
    GB                      // 1073741824
)
```

## Type Conversion
Go requires explicit type conversion:
```go
var i int = 42
var f float64 = float64(i)
var u uint = uint(f)

// String conversions
import "strconv"

// String to int
i, err := strconv.Atoi("42")

// Int to string
s := strconv.Itoa(42)

// String to float
f, err := strconv.ParseFloat("3.14", 64)

// Int to string (byte value to character)
var c rune = 65
s := string(c)  // "A"
```

## Type Aliases
```go
type MyInt int
type Celsius float64
type Fahrenheit float64

var temp Celsius = 25.0
```

## Key Points

✅ **Use `:=` for concise declarations** inside functions
✅ **Use `var` for package-level** variables or when zero value is needed
✅ **Constants use `const`** and cannot be changed
✅ **Type conversion is explicit** - no automatic conversion
✅ **iota is powerful** for creating enumerated constants

## Common Patterns

```go
// Swap values
a, b := 10, 20
a, b = b, a

// Ignore return values with _
value, _ := someFunction()

// Group related constants
const (
    MaxRetries    = 3
    Timeout       = 30
    RetryInterval = 5
)

// Group related variables
var (
    host string
    port int
    db   *Database
)
```

## Gotchas

❌ **Short declaration only works inside functions**
```go
// This WON'T work at package level
name := "John"  // Error!

// Use var instead
var name = "John"  // OK
```

❌ **Must use declared variables**
```go
func main() {
    x := 10
    // Error: x declared but not used
}
```

❌ **Cannot redeclare in same scope**
```go
x := 10
x := 20  // Error: no new variables on left side of :=
x = 20   // OK: assignment
```

## Best Practices

1. **Use short names for short scopes**: `i`, `j` for loop counters
2. **Use descriptive names for larger scopes**: `userCount`, `maxConnections`
3. **Prefer `const` for unchanging values**
4. **Group related declarations** with `var()` or `const()`
5. **Use type aliases** for domain clarity: `type UserID int`
