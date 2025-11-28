# Structs

## Struct Basics

### Declaration
```go
// Define a struct type
type Person struct {
    Name string
    Age  int
}

// Anonymous struct
var person struct {
    Name string
    Age  int
}
```

### Creating Instances
```go
// Using var (zero values)
var p Person  // Name: "", Age: 0

// Using struct literal
p := Person{Name: "Alice", Age: 30}

// Positional (not recommended)
p := Person{"Alice", 30}

// Partial initialization
p := Person{Name: "Alice"}  // Age: 0

// New (returns pointer)
p := new(Person)  // &Person{Name: "", Age: 0}

// Pointer with &
p := &Person{Name: "Alice", Age: 30}
```

## Accessing Fields

```go
type Person struct {
    Name string
    Age  int
}

p := Person{Name: "Alice", Age: 30}

// Access fields
name := p.Name
age := p.Age

// Modify fields
p.Name = "Bob"
p.Age = 25

// Pointer automatic dereferencing
ptr := &p
ptr.Name = "Charlie"  // Same as (*ptr).Name = "Charlie"
```

## Nested Structs

```go
type Address struct {
    Street  string
    City    string
    ZipCode string
}

type Person struct {
    Name    string
    Age     int
    Address Address  // Nested struct
}

// Create with nested struct
p := Person{
    Name: "Alice",
    Age:  30,
    Address: Address{
        Street:  "123 Main St",
        City:    "Springfield",
        ZipCode: "12345",
    },
}

// Access nested fields
city := p.Address.City
```

## Anonymous Fields (Embedding)

```go
type Address struct {
    Street  string
    City    string
    ZipCode string
}

type Person struct {
    Name    string
    Age     int
    Address  // Anonymous field (embedded)
}

p := Person{
    Name: "Alice",
    Age:  30,
    Address: Address{
        Street:  "123 Main St",
        City:    "Springfield",
        ZipCode: "12345",
    },
}

// Direct access to embedded fields
city := p.City  // Instead of p.Address.City
```

## Struct Tags

```go
type User struct {
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"-"`                // Omit from JSON
    Age      int    `json:"age,omitempty"`    // Omit if zero value
}

// Multiple tags
type Product struct {
    ID    int    `json:"id" xml:"id" db:"product_id"`
    Name  string `json:"name" xml:"name" db:"product_name"`
    Price float64 `json:"price" xml:"price" db:"price"`
}

// Validation tags
type Request struct {
    Email string `validate:"required,email"`
    Age   int    `validate:"gte=0,lte=130"`
}
```

## Methods on Structs

See Methods section for details. Quick example:

```go
type Rectangle struct {
    Width  float64
    Height float64
}

// Value receiver
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

// Pointer receiver
func (r *Rectangle) Scale(factor float64) {
    r.Width *= factor
    r.Height *= factor
}

rect := Rectangle{Width: 10, Height: 5}
area := rect.Area()    // 50
rect.Scale(2)          // Modifies rect
```

## Struct Comparison

```go
type Point struct {
    X int
    Y int
}

p1 := Point{1, 2}
p2 := Point{1, 2}
p3 := Point{3, 4}

p1 == p2  // true
p1 == p3  // false

// Structs with uncomparable fields cannot be compared
type Container struct {
    Data []int  // Slice is not comparable
}

// c1 == c2  // Error! Cannot compare
```

## Struct with Interfaces

```go
type Shape interface {
    Area() float64
}

type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return 3.14159 * c.Radius * c.Radius
}

// Circle implements Shape
var s Shape = Circle{Radius: 5}
area := s.Area()
```

## Common Patterns

### Constructor Function
```go
type Person struct {
    name string  // Unexported
    age  int     // Unexported
}

// Constructor
func NewPerson(name string, age int) *Person {
    return &Person{
        name: name,
        age:  age,
    }
}

// Getters
func (p *Person) Name() string {
    return p.name
}

// Setters
func (p *Person) SetName(name string) {
    p.name = name
}
```

### Builder Pattern
```go
type RequestBuilder struct {
    method  string
    url     string
    headers map[string]string
    body    []byte
}

func NewRequestBuilder() *RequestBuilder {
    return &RequestBuilder{
        headers: make(map[string]string),
    }
}

func (rb *RequestBuilder) Method(method string) *RequestBuilder {
    rb.method = method
    return rb
}

func (rb *RequestBuilder) URL(url string) *RequestBuilder {
    rb.url = url
    return rb
}

func (rb *RequestBuilder) Build() Request {
    return Request{
        method:  rb.method,
        url:     rb.url,
        headers: rb.headers,
    }
}

// Usage
req := NewRequestBuilder().
    Method("GET").
    URL("https://api.example.com").
    Build()
```

### Option Pattern
```go
type Server struct {
    host string
    port int
    timeout time.Duration
}

type ServerOption func(*Server)

func WithHost(host string) ServerOption {
    return func(s *Server) {
        s.host = host
    }
}

func WithPort(port int) ServerOption {
    return func(s *Server) {
        s.port = port
    }
}

func NewServer(opts ...ServerOption) *Server {
    s := &Server{
        host:    "localhost",
        port:    8080,
        timeout: 30 * time.Second,
    }
    for _, opt := range opts {
        opt(s)
    }
    return s
}

// Usage
server := NewServer(
    WithHost("0.0.0.0"),
    WithPort(9000),
)
```

### Composition Over Inheritance
```go
type Engine struct {
    Power int
}

func (e Engine) Start() {
    fmt.Println("Engine started")
}

type Car struct {
    Engine  // Composition
    Brand string
}

car := Car{
    Engine: Engine{Power: 200},
    Brand:  "Toyota",
}
car.Start()  // Calls Engine.Start()
```

## Struct Copying

```go
// Shallow copy
p1 := Person{Name: "Alice", Age: 30}
p2 := p1  // p2 is a copy

// Deep copy with slices/maps
type Data struct {
    Values []int
}

d1 := Data{Values: []int{1, 2, 3}}
d2 := Data{Values: make([]int, len(d1.Values))}
copy(d2.Values, d1.Values)
```

## Empty Structs

```go
// Empty struct takes zero bytes
var empty struct{}
fmt.Println(unsafe.Sizeof(empty))  // 0

// Useful for signaling channels
done := make(chan struct{})

// Signal completion
close(done)

// Or send signal
done <- struct{}{}

// Sets
set := make(map[string]struct{})
set["key"] = struct{}{}
```

## Struct JSON Encoding

```go
import "encoding/json"

type User struct {
    ID       int    `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email,omitempty"`
    Password string `json:"-"`  // Never in JSON
}

user := User{
    ID:       1,
    Name:     "Alice",
    Email:    "alice@example.com",
    Password: "secret",
}

// Marshal to JSON
data, err := json.Marshal(user)
// {"id":1,"name":"Alice","email":"alice@example.com"}

// Unmarshal from JSON
var u User
err := json.Unmarshal(data, &u)
```

## Anonymous Structs

```go
// Define and initialize inline
person := struct {
    Name string
    Age  int
}{
    Name: "Alice",
    Age:  30,
}

// Common for test data
tests := []struct {
    input    string
    expected int
}{
    {"test1", 1},
    {"test2", 2},
}

// JSON mapping
data := struct {
    Status  string `json:"status"`
    Message string `json:"message"`
}{
    Status:  "success",
    Message: "Operation completed",
}
```

## Struct Alignment and Padding

```go
// Poor alignment (24 bytes on 64-bit)
type Bad struct {
    a bool   // 1 byte + 7 padding
    b int64  // 8 bytes
    c bool   // 1 byte + 7 padding
}

// Good alignment (16 bytes on 64-bit)
type Good struct {
    b int64  // 8 bytes
    a bool   // 1 byte
    c bool   // 1 byte + 6 padding
}

// Check size
fmt.Println(unsafe.Sizeof(Bad{}))   // 24
fmt.Println(unsafe.Sizeof(Good{}))  // 16
```

## Gotchas

❌ **Struct comparison with uncomparable fields**
```go
type Container struct {
    Data []int
}
// c1 == c2  // Error!
```

❌ **Method receivers: value vs pointer**
```go
type Counter struct {
    count int
}

func (c Counter) Increment() {
    c.count++  // Modifies copy!
}

// Use pointer receiver
func (c *Counter) Increment() {
    c.count++  // Modifies original
}
```

❌ **Nil pointer dereference**
```go
var p *Person  // nil
// p.Name = "Alice"  // Panic!

// Check for nil
if p != nil {
    p.Name = "Alice"
}
```

❌ **Unintended field shadowing**
```go
type Base struct {
    Name string
}

type Derived struct {
    Base
    Name string  // Shadows Base.Name
}

d := Derived{}
d.Name = "derived"      // Sets Derived.Name
d.Base.Name = "base"    // Sets Base.Name
```

## Best Practices

1. **Use pointer receivers** for methods that modify state
2. **Use value receivers** for small, immutable structs
3. **Export fields** that need to be accessed externally
4. **Use struct tags** for encoding/decoding
5. **Provide constructors** for complex initialization
6. **Consider field order** for memory efficiency
7. **Use embedding** for composition
8. **Document exported types** and fields
9. **Validate in constructors** - ensure valid state
10. **Use zero values wisely** - design for default initialization

## Performance Tips

```go
// Pass large structs by pointer
func ProcessLarge(data *LargeStruct) {
    // Avoids copying large struct
}

// Small structs can be passed by value
func ProcessSmall(p Point) {
    // Point is small (2 ints)
}

// Pool structs if creating many
var pool = sync.Pool{
    New: func() interface{} {
        return &MyStruct{}
    },
}

obj := pool.Get().(*MyStruct)
// Use obj
pool.Put(obj)
```
