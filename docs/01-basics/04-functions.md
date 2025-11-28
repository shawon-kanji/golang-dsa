# Functions

## Basic Function Syntax

```go
// Basic function
func functionName(param1 type1, param2 type2) returnType {
    // function body
    return value
}

// Example
func add(a int, b int) int {
    return a + b
}

// Shorthand when parameters share type
func add(a, b int) int {
    return a + b
}

// Function with no parameters
func greet() {
    fmt.Println("Hello!")
}

// Function with no return value
func printSum(a, b int) {
    fmt.Println(a + b)
}
```

## Multiple Return Values

```go
// Return multiple values
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
}

// Usage
result, err := divide(10, 2)
if err != nil {
    fmt.Println("Error:", err)
    return
}
fmt.Println("Result:", result)

// Ignore return values with _
result, _ := divide(10, 2)
```

## Named Return Values

```go
// Named returns
func split(sum int) (x, y int) {
    x = sum * 4 / 9
    y = sum - x
    return  // naked return
}

// Explicit return (recommended for clarity)
func split(sum int) (x, y int) {
    x = sum * 4 / 9
    y = sum - x
    return x, y
}

// Named return with early return
func divide(a, b float64) (result float64, err error) {
    if b == 0 {
        err = fmt.Errorf("division by zero")
        return  // Returns zero value for result, error for err
    }
    result = a / b
    return
}
```

## Variadic Functions

```go
// Variadic parameter (must be last)
func sum(nums ...int) int {
    total := 0
    for _, num := range nums {
        total += num
    }
    return total
}

// Usage
sum(1, 2, 3)           // 6
sum(1, 2, 3, 4, 5)     // 15

// Pass slice as variadic
numbers := []int{1, 2, 3, 4}
sum(numbers...)        // Unpacks slice

// Mix regular and variadic parameters
func printf(format string, args ...interface{}) {
    fmt.Printf(format, args...)
}
```

## Anonymous Functions

```go
// Anonymous function
func() {
    fmt.Println("Anonymous function")
}()  // Immediate execution

// Assign to variable
add := func(a, b int) int {
    return a + b
}
result := add(3, 4)  // 7

// Return function from function
func makeAdder(x int) func(int) int {
    return func(y int) int {
        return x + y
    }
}

add5 := makeAdder(5)
fmt.Println(add5(3))   // 8
fmt.Println(add5(10))  // 15
```

## Closures

```go
// Closure captures outer variable
func counter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}

c1 := counter()
fmt.Println(c1())  // 1
fmt.Println(c1())  // 2
fmt.Println(c1())  // 3

c2 := counter()
fmt.Println(c2())  // 1 (new counter)

// Closure in loops (common gotcha)
funcs := []func(){}
for i := 0; i < 3; i++ {
    funcs = append(funcs, func() {
        fmt.Println(i)  // Captures reference to i
    })
}
for _, f := range funcs {
    f()  // Prints: 3, 3, 3
}

// Fix: capture value
funcs := []func(){}
for i := 0; i < 3; i++ {
    i := i  // Create new variable
    funcs = append(funcs, func() {
        fmt.Println(i)
    })
}
for _, f := range funcs {
    f()  // Prints: 0, 1, 2
}
```

## Higher-Order Functions

```go
// Function as parameter
func apply(f func(int) int, value int) int {
    return f(value)
}

double := func(x int) int { return x * 2 }
result := apply(double, 5)  // 10

// Map function
func mapInts(nums []int, f func(int) int) []int {
    result := make([]int, len(nums))
    for i, num := range nums {
        result[i] = f(num)
    }
    return result
}

nums := []int{1, 2, 3, 4}
doubled := mapInts(nums, func(x int) int { return x * 2 })
// [2, 4, 6, 8]

// Filter function
func filter(nums []int, f func(int) bool) []int {
    result := []int{}
    for _, num := range nums {
        if f(num) {
            result = append(result, num)
        }
    }
    return result
}

evens := filter(nums, func(x int) bool { return x%2 == 0 })
// [2, 4]
```

## Defer Statement

```go
// Defer executes when surrounding function returns
func main() {
    defer fmt.Println("World")
    fmt.Println("Hello")
}
// Output: Hello, World

// Multiple defers execute in LIFO order
func main() {
    defer fmt.Println("1")
    defer fmt.Println("2")
    defer fmt.Println("3")
}
// Output: 3, 2, 1

// Common use: cleanup
func processFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()  // Ensures file is closed

    // Process file
    return nil
}

// Defer with variables (evaluated immediately)
func main() {
    x := 10
    defer fmt.Println(x)  // Will print 10
    x = 20
    fmt.Println(x)        // Will print 20
}
// Output: 20, 10
```

## Methods vs Functions

Methods are covered in detail in the Methods section, but here's a quick comparison:

```go
// Function
func Add(a, b int) int {
    return a + b
}

// Method (function with receiver)
type Calculator struct{}

func (c Calculator) Add(a, b int) int {
    return a + b
}

// Usage
result := Add(2, 3)          // Function

calc := Calculator{}
result := calc.Add(2, 3)     // Method
```

## Recursive Functions

```go
// Factorial
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n-1)
}

// Fibonacci
func fibonacci(n int) int {
    if n <= 1 {
        return n
    }
    return fibonacci(n-1) + fibonacci(n-2)
}

// Binary search
func binarySearch(arr []int, target, low, high int) int {
    if low > high {
        return -1
    }
    mid := (low + high) / 2
    if arr[mid] == target {
        return mid
    }
    if arr[mid] > target {
        return binarySearch(arr, target, low, mid-1)
    }
    return binarySearch(arr, target, mid+1, high)
}
```

## Function Types

```go
// Define function type
type Operation func(int, int) int

// Use as parameter
func calculate(a, b int, op Operation) int {
    return op(a, b)
}

// Define operations
var add Operation = func(a, b int) int { return a + b }
var multiply Operation = func(a, b int) int { return a * b }

// Use
result := calculate(10, 5, add)       // 15
result = calculate(10, 5, multiply)   // 50
```

## Common Patterns

### Error handling wrapper
```go
func must(err error) {
    if err != nil {
        panic(err)
    }
}

file := must(os.Open("file.txt"))
```

### Builder pattern
```go
type Request struct {
    url     string
    method  string
    headers map[string]string
}

func NewRequest(url string) *Request {
    return &Request{url: url, headers: make(map[string]string)}
}

func (r *Request) Method(method string) *Request {
    r.method = method
    return r
}

func (r *Request) Header(key, value string) *Request {
    r.headers[key] = value
    return r
}

// Usage
req := NewRequest("https://api.com").
    Method("GET").
    Header("Auth", "token")
```

### Option pattern
```go
type Server struct {
    host string
    port int
}

type Option func(*Server)

func WithHost(host string) Option {
    return func(s *Server) {
        s.host = host
    }
}

func WithPort(port int) Option {
    return func(s *Server) {
        s.port = port
    }
}

func NewServer(opts ...Option) *Server {
    s := &Server{host: "localhost", port: 8080}
    for _, opt := range opts {
        opt(s)
    }
    return s
}

// Usage
server := NewServer(WithHost("0.0.0.0"), WithPort(9090))
```

## Gotchas

❌ **Defer order and evaluation**
```go
func main() {
    i := 0
    defer fmt.Println(i)  // Prints 0 (value captured)
    i++
}
```

❌ **Closure loop variable capture**
```go
for i := 0; i < 3; i++ {
    defer func() {
        fmt.Println(i)  // Captures reference, prints 3, 3, 3
    }()
}

// Fix
for i := 0; i < 3; i++ {
    i := i  // Shadow variable
    defer func() {
        fmt.Println(i)  // Prints 2, 1, 0
    }()
}
```

❌ **Named returns can be confusing**
```go
func confusing() (i int) {
    defer func() {
        i++  // Modifies return value!
    }()
    return 5  // Actually returns 6
}
```

## Best Practices

1. **Keep functions short and focused** - single responsibility
2. **Use descriptive names** - function name should describe what it does
3. **Return errors, don't panic** - except for unrecoverable errors
4. **Use defer for cleanup** - file closes, unlocks, etc.
5. **Naked returns only for short functions** - be explicit in longer functions
6. **Document exported functions** - use doc comments
7. **Variadic functions for flexibility** - but don't overuse
8. **Use closures judiciously** - they capture state
9. **Prefer passing values** - unless data is large
10. **Check errors immediately** - don't ignore them

## Function Documentation

```go
// Add returns the sum of a and b.
// This is a simple addition function for demonstration.
func Add(a, b int) int {
    return a + b
}

// Process handles data processing with the following steps:
//   1. Validates input data
//   2. Transforms data
//   3. Stores result
//
// Returns error if validation fails or storage operation fails.
func Process(data []byte) error {
    // implementation
}
```
