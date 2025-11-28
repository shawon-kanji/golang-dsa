# Methods

## Method Basics

### Defining Methods
```go
type Rectangle struct {
    Width, Height float64
}

// Method with value receiver
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

// Method with pointer receiver
func (r *Rectangle) Scale(factor float64) {
    r.Width *= factor
    r.Height *= factor
}

// Usage
rect := Rectangle{Width: 10, Height: 5}
area := rect.Area()    // 50
rect.Scale(2)          // Modifies rect
fmt.Println(rect.Area())  // 200
```

## Value vs Pointer Receivers

### Value Receiver
```go
type Point struct {
    X, Y int
}

// Value receiver - receives copy
func (p Point) Distance() float64 {
    return math.Sqrt(float64(p.X*p.X + p.Y*p.Y))
}

// Doesn't modify original
func (p Point) Move(dx, dy int) {
    p.X += dx  // Modifies copy
    p.Y += dy
}

p := Point{3, 4}
p.Move(1, 1)
fmt.Println(p)  // {3 4} - unchanged
```

### Pointer Receiver
```go
// Pointer receiver - receives address
func (p *Point) MovePtr(dx, dy int) {
    p.X += dx  // Modifies original
    p.Y += dy
}

p := Point{3, 4}
p.MovePtr(1, 1)
fmt.Println(p)  // {4 5} - changed
```

### When to Use Each

**Use pointer receivers when:**
- Method needs to modify the receiver
- Receiver is a large struct (avoid copying)
- Consistency (if some methods use pointer, all should)

**Use value receivers when:**
- Receiver is small and copyable
- Receiver is immutable
- Receiver is a basic type or small struct

```go
type Counter struct {
    count int
}

// Pointer receiver - modifies state
func (c *Counter) Increment() {
    c.count++
}

// Value receiver - read-only
func (c Counter) Value() int {
    return c.count
}
```

## Methods on Different Types

### Methods on Built-in Types
```go
// Can't define methods on built-in types directly
// func (i int) Double() int { }  // Error!

// Use custom type
type MyInt int

func (i MyInt) Double() MyInt {
    return i * 2
}

var x MyInt = 5
fmt.Println(x.Double())  // 10
```

### Methods on Slices
```go
type IntSlice []int

func (s IntSlice) Sum() int {
    total := 0
    for _, v := range s {
        total += v
    }
    return total
}

func (s IntSlice) Average() float64 {
    if len(s) == 0 {
        return 0
    }
    return float64(s.Sum()) / float64(len(s))
}

nums := IntSlice{1, 2, 3, 4, 5}
fmt.Println(nums.Sum())      // 15
fmt.Println(nums.Average())  // 3.0
```

### Methods on Maps
```go
type Counter map[string]int

func (c Counter) Increment(key string) {
    c[key]++
}

func (c Counter) Get(key string) int {
    return c[key]
}

counter := make(Counter)
counter.Increment("apple")
counter.Increment("apple")
fmt.Println(counter.Get("apple"))  // 2
```

## Method Sets and Interfaces

### Interface Implementation
```go
type Shape interface {
    Area() float64
    Perimeter() float64
}

type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
    return 2 * math.Pi * c.Radius
}

// Circle implements Shape
var s Shape = Circle{Radius: 5}
```

### Pointer Receivers and Interfaces
```go
type Incrementer interface {
    Increment()
}

type Counter struct {
    count int
}

func (c *Counter) Increment() {
    c.count++
}

// Only *Counter implements Incrementer
var i Incrementer
i = &Counter{}  // OK
// i = Counter{}   // Error: Counter doesn't implement Incrementer
```

### Method Set Rules
```go
type T struct{}

func (t T) ValueMethod() {}
func (t *T) PointerMethod() {}

// T's method set: ValueMethod
// *T's method set: ValueMethod, PointerMethod

var t T
t.ValueMethod()    // OK
t.PointerMethod()  // OK (Go automatically takes address)

var p *T = &T{}
p.ValueMethod()    // OK (Go automatically dereferences)
p.PointerMethod()  // OK
```

## Method Chaining

```go
type Builder struct {
    data string
}

func (b *Builder) Append(s string) *Builder {
    b.data += s
    return b
}

func (b *Builder) Prepend(s string) *Builder {
    b.data = s + b.data
    return b
}

func (b *Builder) String() string {
    return b.data
}

// Chaining
result := new(Builder).
    Append("World").
    Prepend("Hello ").
    String()
// "Hello World"
```

## Method Expressions

```go
type Point struct {
    X, Y int
}

func (p Point) Distance() float64 {
    return math.Sqrt(float64(p.X*p.X + p.Y*p.Y))
}

// Method expression
distance := Point.Distance  // func(Point) float64

p := Point{3, 4}
d := distance(p)  // 5
```

## Common Patterns

### Constructor Method
```go
type Person struct {
    name string
    age  int
}

func NewPerson(name string, age int) *Person {
    return &Person{
        name: name,
        age:  age,
    }
}

func (p *Person) Name() string {
    return p.name
}

func (p *Person) SetName(name string) {
    p.name = name
}
```

### Validation Method
```go
type Email string

func (e Email) Validate() error {
    if !strings.Contains(string(e), "@") {
        return errors.New("invalid email")
    }
    return nil
}

email := Email("user@example.com")
if err := email.Validate(); err != nil {
    fmt.Println("Invalid email")
}
```

### String Representation
```go
type Person struct {
    Name string
    Age  int
}

func (p Person) String() string {
    return fmt.Sprintf("%s (%d years)", p.Name, p.Age)
}

p := Person{"Alice", 30}
fmt.Println(p)  // Alice (30 years)
```

### Comparison Methods
```go
type Version struct {
    Major, Minor, Patch int
}

func (v Version) Equal(other Version) bool {
    return v.Major == other.Major &&
           v.Minor == other.Minor &&
           v.Patch == other.Patch
}

func (v Version) Less(other Version) bool {
    if v.Major != other.Major {
        return v.Major < other.Major
    }
    if v.Minor != other.Minor {
        return v.Minor < other.Minor
    }
    return v.Patch < other.Patch
}
```

## Embedded Types and Methods

```go
type Engine struct {
    Power int
}

func (e Engine) Start() {
    fmt.Println("Engine started")
}

type Car struct {
    Engine  // Embedded
    Brand string
}

car := Car{
    Engine: Engine{Power: 200},
    Brand:  "Toyota",
}

car.Start()  // Calls Engine.Start()
// Equivalent to car.Engine.Start()
```

### Method Promotion
```go
type Logger struct{}

func (l Logger) Log(msg string) {
    fmt.Println("LOG:", msg)
}

type Service struct {
    Logger  // Logger methods promoted to Service
    Name string
}

s := Service{Name: "MyService"}
s.Log("Started")  // Calls Logger.Log
```

### Overriding Promoted Methods
```go
type Base struct{}

func (b Base) Method() {
    fmt.Println("Base method")
}

type Derived struct {
    Base
}

func (d Derived) Method() {
    fmt.Println("Derived method")
}

d := Derived{}
d.Method()       // "Derived method"
d.Base.Method()  // "Base method"
```

## Best Practices

1. **Consistent receiver names** - use short, consistent names (often first letter of type)
2. **Consistent receiver type** - all methods should use same receiver type
3. **Pointer receivers for modification** - or for large structs
4. **Value receivers for small, immutable** data
5. **Don't mix receivers** unless necessary
6. **Method names should be verbs** - describe actions
7. **Getter methods don't use Get prefix** - `Name()` not `GetName()`
8. **Setter methods use Set prefix** - `SetName(string)`
9. **Document exported methods** - especially complex behavior
10. **Keep methods focused** - single responsibility

## Gotchas

❌ **Mixing receiver types**
```go
type Counter struct {
    count int
}

func (c Counter) Value() int {
    return c.count
}

func (c *Counter) Increment() {
    c.count++
}

// Inconsistent - prefer all pointer receivers if any modify state
```

❌ **Forgetting pointer receiver**
```go
func (c Counter) Increment() {
    c.count++  // Modifies copy!
}

c := Counter{}
c.Increment()
fmt.Println(c.count)  // 0, not 1
```

❌ **Calling method on nil receiver**
```go
var p *Point
// p.Move(1, 1)  // Panic if method dereferences p!

// Safe implementation
func (p *Point) Move(dx, dy int) {
    if p == nil {
        return
    }
    p.X += dx
    p.Y += dy
}
```

## Testing Methods

```go
func TestCounterIncrement(t *testing.T) {
    c := Counter{}
    c.Increment()
    if c.Value() != 1 {
        t.Errorf("expected 1, got %d", c.Value())
    }
}

func TestCircleArea(t *testing.T) {
    c := Circle{Radius: 5}
    expected := math.Pi * 25
    if math.Abs(c.Area()-expected) > 0.0001 {
        t.Errorf("expected %.2f, got %.2f", expected, c.Area())
    }
}
```
