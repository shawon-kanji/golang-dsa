# Pointers

## Pointer Basics

### What is a Pointer?
A pointer holds the memory address of a value.

```go
// Declaration
var p *int  // Pointer to int, initially nil

// & operator: get address
x := 10
p = &x  // p points to x

// * operator: dereference (get value)
value := *p  // value is 10

// Modify through pointer
*p = 20  // Changes x to 20
```

### Nil Pointers
```go
var p *int  // nil pointer

if p == nil {
    fmt.Println("p is nil")
}

// Dereferencing nil pointer causes panic!
// *p = 10  // Panic!

// Must initialize before use
x := 10
p = &x
*p = 20  // OK
```

## Creating Pointers

### Using & Operator
```go
x := 42
p := &x  // p is *int pointing to x

fmt.Println(*p)  // 42
*p = 100
fmt.Println(x)   // 100
```

### Using new
```go
// new allocates memory and returns pointer
p := new(int)  // p is *int, value is 0
*p = 42

// Equivalent to
var x int
p := &x
```

### Pointer to Struct
```go
type Person struct {
    Name string
    Age  int
}

// Method 1: & operator
p := &Person{Name: "Alice", Age: 30}

// Method 2: new
p := new(Person)
p.Name = "Alice"
p.Age = 30

// Method 3: value then address
person := Person{Name: "Alice", Age: 30}
p := &person
```

## Pointer Operations

### Access Struct Fields
```go
type Person struct {
    Name string
    Age  int
}

p := &Person{Name: "Alice", Age: 30}

// Go automatically dereferences
name := p.Name  // Same as (*p).Name

// Modify through pointer
p.Age = 31
```

### Pointer Arithmetic
```go
// Go does NOT support pointer arithmetic
// This is NOT allowed:
// p++
// p += 1
// p = p + 1

// Use slices for array-like operations
```

### Pointer Comparison
```go
x := 10
p1 := &x
p2 := &x
p3 := &x

p1 == p2  // true (point to same address)
p1 == p3  // true

y := 10
p4 := &y
p1 == p4  // false (different addresses)

// Can compare to nil
var p *int
p == nil  // true
```

## Pass by Value vs Reference

### Pass by Value
```go
func modify(x int) {
    x = 100  // Modifies copy
}

x := 10
modify(x)
fmt.Println(x)  // 10 (unchanged)
```

### Pass by Pointer
```go
func modify(p *int) {
    *p = 100  // Modifies original
}

x := 10
modify(&x)
fmt.Println(x)  // 100 (changed)
```

### Struct Example
```go
type Person struct {
    Name string
    Age  int
}

// Pass by value - no modification
func updateAge(p Person, newAge int) {
    p.Age = newAge  // Modifies copy
}

// Pass by pointer - modifies original
func updateAgePtr(p *Person, newAge int) {
    p.Age = newAge  // Modifies original
}

person := Person{Name: "Alice", Age: 30}
updateAge(person, 31)
fmt.Println(person.Age)  // 30 (unchanged)

updateAgePtr(&person, 31)
fmt.Println(person.Age)  // 31 (changed)
```

## Pointer Receivers

```go
type Counter struct {
    count int
}

// Value receiver - no modification
func (c Counter) Increment() {
    c.count++  // Modifies copy
}

// Pointer receiver - modifies original
func (c *Counter) IncrementPtr() {
    c.count++  // Modifies original
}

c := Counter{count: 0}
c.Increment()
fmt.Println(c.count)  // 0

c.IncrementPtr()
fmt.Println(c.count)  // 1
```

## Pointers to Arrays and Slices

### Array Pointers
```go
arr := [3]int{1, 2, 3}
p := &arr

// Access elements
(*p)[0] = 10  // Explicit dereference
p[0] = 10     // Automatic dereference

fmt.Println(arr)  // [10, 2, 3]
```

### Slice Internals
```go
// Slices are already references!
func modify(s []int) {
    s[0] = 100  // Modifies original
}

slice := []int{1, 2, 3}
modify(slice)
fmt.Println(slice)  // [100, 2, 3]

// But append needs pointer if slice grows
func appendValue(s *[]int, val int) {
    *s = append(*s, val)
}
```

## Pointer to Pointer

```go
x := 10
p := &x      // *int
pp := &p     // **int (pointer to pointer)

fmt.Println(**pp)  // 10

**pp = 20
fmt.Println(x)  // 20

// Modify pointer itself
y := 30
*pp = &y  // Changes p to point to y
fmt.Println(*p)  // 30
```

## Common Patterns

### Optional Parameters
```go
func process(data string, timeout *int) {
    if timeout != nil {
        // Use provided timeout
        fmt.Println("Timeout:", *timeout)
    } else {
        // Use default
        fmt.Println("Using default timeout")
    }
}

// Call with timeout
timeout := 30
process("data", &timeout)

// Call without timeout
process("data", nil)
```

### Factory Functions
```go
func NewPerson(name string, age int) *Person {
    return &Person{
        Name: name,
        Age:  age,
    }
}

p := NewPerson("Alice", 30)
```

### Linked Structures
```go
type Node struct {
    Value int
    Next  *Node
}

// Build linked list
head := &Node{Value: 1}
head.Next = &Node{Value: 2}
head.Next.Next = &Node{Value: 3}

// Traverse
for node := head; node != nil; node = node.Next {
    fmt.Println(node.Value)
}
```

### Swap Function
```go
func swap(a, b *int) {
    *a, *b = *b, *a
}

x, y := 10, 20
swap(&x, &y)
fmt.Println(x, y)  // 20, 10
```

### Pointer Receiver for Large Structs
```go
type LargeStruct struct {
    data [1000]int
}

// Avoid copying large struct
func (ls *LargeStruct) Process() {
    // Process data
}
```

## Pointer Methods

### Method Set Rules
```go
type T struct{}

func (t T) ValueMethod() {}
func (t *T) PointerMethod() {}

// Value can call both
var t T
t.ValueMethod()    // OK
t.PointerMethod()  // OK (Go takes address automatically)

// Pointer can call both
var p *T = &T{}
p.ValueMethod()    // OK (Go dereferences automatically)
p.PointerMethod()  // OK

// Interface considerations
type I interface {
    PointerMethod()
}

var i I
i = &T{}  // OK: *T implements I
// i = T{}   // Error: T doesn't implement I
```

## Pointers and Interfaces

```go
type Shape interface {
    Area() float64
}

type Circle struct {
    Radius float64
}

// Pointer receiver
func (c *Circle) Area() float64 {
    return 3.14159 * c.Radius * c.Radius
}

// Only *Circle implements Shape
var s Shape
s = &Circle{Radius: 5}  // OK
// s = Circle{Radius: 5}   // Error!
```

## Memory Management

### Stack vs Heap
```go
// Stack allocation (if possible)
func stackAlloc() int {
    x := 10  // Allocated on stack
    return x  // Value returned
}

// Heap allocation (escapes to heap)
func heapAlloc() *int {
    x := 10  // Allocated on heap (escapes)
    return &x  // Pointer returned
}

// Compiler decides based on escape analysis
```

### Garbage Collection
```go
// Go has automatic garbage collection
// No manual memory management needed

func createData() {
    p := new(LargeStruct)
    // Use p
    // No need to free - GC handles it
}
```

## Gotchas

❌ **Nil pointer dereference**
```go
var p *int
// *p = 10  // Panic: nil pointer dereference
```

❌ **Returning pointer to local variable**
```go
// This is actually OK in Go!
func createInt() *int {
    x := 10
    return &x  // x escapes to heap, safe to return
}
```

❌ **Pointer loop variable**
```go
// Common mistake
slice := []int{1, 2, 3}
var pointers []*int

for _, v := range slice {
    pointers = append(pointers, &v)  // Wrong! v is reused
}
// All pointers point to same address

// Fix: create new variable
for _, v := range slice {
    v := v  // Create new variable
    pointers = append(pointers, &v)
}

// Or use index
for i := range slice {
    pointers = append(pointers, &slice[i])
}
```

❌ **Comparing pointers to different types**
```go
var p1 *int
var p2 *float64
// p1 == p2  // Error: mismatched types
```

❌ **Modifying through nil pointer**
```go
var p *Person
// p.Name = "Alice"  // Panic: nil pointer dereference

// Check first
if p != nil {
    p.Name = "Alice"
}
```

## Best Practices

1. **Use pointers for large structs** - avoid copying
2. **Use pointer receivers** when modifying state
3. **Check for nil** before dereferencing
4. **Return pointers from constructors** - idiomatic
5. **Use values for small, immutable data** - simpler
6. **Avoid pointer to slice/map** - already references
7. **Be careful with loop variables** - don't take their address
8. **Document nil behavior** - clarify if nil is valid
9. **Use pointers for optional fields** - nil means absent
10. **Trust the compiler** - escape analysis is smart

## When to Use Pointers

✅ **Use pointers when:**
- Modifying the receiver in methods
- Large structs (avoid copying)
- Shared state between functions
- Implementing interfaces with pointer receivers
- Optional/absent values (nil)
- Building linked data structures

❌ **Don't use pointers when:**
- Small, immutable data
- Simple value types (int, bool, etc.)
- Read-only operations
- Slices and maps (already references)
- Unnecessary complexity
- Not sure (start with values)

## Performance Considerations

```go
// Small structs - value is fine
type Point struct {
    X, Y int
}

func distance(p Point) float64 {
    // No performance issue
}

// Large structs - use pointer
type LargeData struct {
    data [10000]int
}

func process(ld *LargeData) {
    // Avoids copying 10000 ints
}

// Benchmark to decide threshold
// Generally: > 64 bytes, consider pointer
```

## Pointer Safety

Go provides pointer safety:
- No pointer arithmetic
- No manual memory management
- Garbage collection
- Type safety
- Nil pointer checks

```go
// These are NOT allowed in Go:
// p++
// p += 1
// free(p)
// (int*)p

// Making Go pointers safer than C/C++
```
