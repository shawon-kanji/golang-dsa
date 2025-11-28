# Interfaces

## Interface Basics

### Definition
```go
// Interface defines a contract
type Shape interface {
    Area() float64
    Perimeter() float64
}

// Empty interface
type Any interface{}  // or interface{}
```

### Implementation
```go
// Rectangle implements Shape
type Rectangle struct {
    Width, Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

// Circle implements Shape
type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return 3.14159 * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
    return 2 * 3.14159 * c.Radius
}
```

### Using Interfaces
```go
func printInfo(s Shape) {
    fmt.Printf("Area: %.2f\n", s.Area())
    fmt.Printf("Perimeter: %.2f\n", s.Perimeter())
}

rect := Rectangle{Width: 10, Height: 5}
circle := Circle{Radius: 7}

printInfo(rect)   // Works
printInfo(circle) // Works
```

## Empty Interface

```go
// interface{} can hold any type
var any interface{}

any = 42
any = "hello"
any = []int{1, 2, 3}
any = Rectangle{Width: 10, Height: 5}

// Function accepting any type
func print(v interface{}) {
    fmt.Println(v)
}

print(42)
print("hello")
print([]int{1, 2, 3})
```

## Type Assertions

```go
var i interface{} = "hello"

// Type assertion
s := i.(string)
fmt.Println(s)  // "hello"

// Type assertion with check
s, ok := i.(string)
if ok {
    fmt.Println(s)
} else {
    fmt.Println("not a string")
}

// Panic if wrong type
// n := i.(int)  // Panic!

// Safe version
if n, ok := i.(int); ok {
    fmt.Println(n)
} else {
    fmt.Println("not an int")
}
```

## Type Switch

```go
func describe(i interface{}) {
    switch v := i.(type) {
    case int:
        fmt.Printf("Integer: %d\n", v)
    case string:
        fmt.Printf("String: %s\n", v)
    case bool:
        fmt.Printf("Boolean: %t\n", v)
    case Rectangle:
        fmt.Printf("Rectangle: %v\n", v)
    default:
        fmt.Printf("Unknown type: %T\n", v)
    }
}

describe(42)                          // Integer: 42
describe("hello")                     // String: hello
describe(true)                        // Boolean: true
describe(Rectangle{Width: 10, Height: 5})  // Rectangle: {10 5}
```

## Interface Embedding

```go
// Reader interface
type Reader interface {
    Read(p []byte) (n int, err error)
}

// Writer interface
type Writer interface {
    Write(p []byte) (n int, err error)
}

// ReadWriter embeds both
type ReadWriter interface {
    Reader
    Writer
}

// Or inline
type ReadWriter interface {
    Read(p []byte) (n int, err error)
    Write(p []byte) (n int, err error)
}
```

## Standard Library Interfaces

### fmt.Stringer
```go
type Person struct {
    Name string
    Age  int
}

func (p Person) String() string {
    return fmt.Sprintf("%s (%d years old)", p.Name, p.Age)
}

p := Person{Name: "Alice", Age: 30}
fmt.Println(p)  // Alice (30 years old)
```

### error Interface
```go
type error interface {
    Error() string
}

// Custom error
type MyError struct {
    Code    int
    Message string
}

func (e MyError) Error() string {
    return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

func doSomething() error {
    return MyError{Code: 404, Message: "Not found"}
}
```

### io.Reader and io.Writer
```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

// Custom reader
type MyReader struct {
    data string
    pos  int
}

func (r *MyReader) Read(p []byte) (n int, err error) {
    if r.pos >= len(r.data) {
        return 0, io.EOF
    }
    n = copy(p, r.data[r.pos:])
    r.pos += n
    return n, nil
}
```

### sort.Interface
```go
type Interface interface {
    Len() int
    Less(i, j int) bool
    Swap(i, j int)
}

// Sort custom type
type Person struct {
    Name string
    Age  int
}

type ByAge []Person

func (a ByAge) Len() int           { return len(a) }
func (a ByAge) Less(i, j int) bool { return a[i].Age < a[j].Age }
func (a ByAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

people := []Person{
    {"Alice", 30},
    {"Bob", 25},
    {"Charlie", 35},
}

sort.Sort(ByAge(people))
```

## Polymorphism

```go
type Animal interface {
    Speak() string
}

type Dog struct{}
func (d Dog) Speak() string { return "Woof!" }

type Cat struct{}
func (c Cat) Speak() string { return "Meow!" }

func makeSound(a Animal) {
    fmt.Println(a.Speak())
}

makeSound(Dog{})  // Woof!
makeSound(Cat{})  // Meow!

// Slice of interfaces
animals := []Animal{Dog{}, Cat{}, Dog{}}
for _, animal := range animals {
    fmt.Println(animal.Speak())
}
```

## Interface Values

```go
// Interface value = (type, value)
var w io.Writer

w = os.Stdout  // w = (*os.File, os.Stdout)
w.Write([]byte("hello"))

w = new(bytes.Buffer)  // w = (*bytes.Buffer, buffer)
w.Write([]byte("hello"))

// Nil interface value
var w io.Writer  // w = (nil, nil)
// w.Write([]byte("hello"))  // Panic!

// Non-nil interface with nil value
var buf *bytes.Buffer  // buf = nil
w = buf                // w = (*bytes.Buffer, nil) - not nil interface!
// w.Write([]byte("hello"))  // May panic or not, depends on implementation
```

## Interface Satisfaction

```go
// Implicit satisfaction
type Shape interface {
    Area() float64
}

type Square struct {
    Side float64
}

func (s Square) Area() float64 {
    return s.Side * s.Side
}

// Square automatically satisfies Shape
var shape Shape = Square{Side: 5}

// Compile-time check
var _ Shape = (*Square)(nil)  // Compile error if Square doesn't implement Shape
var _ Shape = Square{}        // Same check
```

## Common Patterns

### Strategy Pattern
```go
type SortStrategy interface {
    Sort([]int) []int
}

type BubbleSort struct{}
func (b BubbleSort) Sort(arr []int) []int {
    // Implement bubble sort
    return arr
}

type QuickSort struct{}
func (q QuickSort) Sort(arr []int) []int {
    // Implement quick sort
    return arr
}

func sortNumbers(arr []int, strategy SortStrategy) []int {
    return strategy.Sort(arr)
}
```

### Adapter Pattern
```go
// Old interface
type OldPrinter interface {
    PrintOld(s string)
}

// New interface
type NewPrinter interface {
    Print(s string)
}

// Adapter
type PrinterAdapter struct {
    oldPrinter OldPrinter
}

func (pa PrinterAdapter) Print(s string) {
    pa.oldPrinter.PrintOld(s)
}
```

### Dependency Injection
```go
type Database interface {
    Query(query string) ([]byte, error)
}

type UserService struct {
    db Database
}

func NewUserService(db Database) *UserService {
    return &UserService{db: db}
}

func (s *UserService) GetUser(id int) (*User, error) {
    data, err := s.db.Query(fmt.Sprintf("SELECT * FROM users WHERE id = %d", id))
    // Process data
    return user, err
}

// Can inject different implementations
service := NewUserService(PostgresDB{})
service := NewUserService(MockDB{})
```

## Interface Composition

```go
type ReadCloser interface {
    io.Reader
    io.Closer
}

type WriteCloser interface {
    io.Writer
    io.Closer
}

type ReadWriteCloser interface {
    io.Reader
    io.Writer
    io.Closer
}

// Or use io.ReadWriteCloser from stdlib
```

## Gotchas

❌ **Nil interface vs interface with nil value**
```go
var i interface{}  // Nil interface
i == nil  // true

var p *int  // Nil pointer
i = p       // Interface with nil value
i == nil    // false!

// Check underlying value
if i == nil || reflect.ValueOf(i).IsNil() {
    // Handle nil
}
```

❌ **Value vs pointer receivers**
```go
type Counter interface {
    Increment()
}

type MyCounter struct {
    count int
}

func (c *MyCounter) Increment() {
    c.count++
}

// Only *MyCounter implements Counter
var c Counter
c = &MyCounter{}  // OK
// c = MyCounter{}   // Error: MyCounter doesn't implement Counter
```

❌ **Interface comparison**
```go
// Can compare interfaces
var i1, i2 interface{} = 1, 1
i1 == i2  // true

// But can panic if underlying type not comparable
i1 = []int{1}
i2 = []int{1}
// i1 == i2  // Panic! slice not comparable
```

## Best Practices

1. **Keep interfaces small** - ideally 1-3 methods
2. **Define interfaces at usage** - not implementation
3. **Accept interfaces, return structs** - flexible input, concrete output
4. **Use standard library interfaces** - io.Reader, fmt.Stringer, etc.
5. **Name single-method interfaces** with -er suffix: Reader, Writer, Closer
6. **Don't over-abstract** - interfaces add complexity
7. **Document interface contracts** - behavior expectations
8. **Test with interface** - easy to mock
9. **Check nil carefully** - understand nil interface vs nil value
10. **Use empty interface sparingly** - loses type safety

## Interface Design Guidelines

✅ **Good interface design:**
```go
// Small, focused
type Reader interface {
    Read(p []byte) (n int, err error)
}

// Composition
type ReadWriter interface {
    Reader
    Writer
}

// Descriptive names
type UserRepository interface {
    FindByID(id int) (*User, error)
    Save(user *User) error
}
```

❌ **Poor interface design:**
```go
// Too large
type Repository interface {
    Create(...)
    Read(...)
    Update(...)
    Delete(...)
    List(...)
    Search(...)
    // ... many more methods
}

// Too generic
type Service interface {
    Do(interface{}) interface{}
}
```
