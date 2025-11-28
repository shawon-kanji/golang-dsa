# Go Syntax Cheatsheet

## Variables

```go
// Declaration
var name string
var age int = 30
var isActive bool

// Short declaration
name := "John"
age := 30

// Multiple
var a, b int = 1, 2
x, y := 3, 4

// Zero values: 0, false, "", nil
```

## Constants

```go
const Pi = 3.14
const (
    StatusOK = 200
    StatusNotFound = 404
)

// iota
const (
    Sunday = iota  // 0
    Monday         // 1
    Tuesday        // 2
)
```

## Types

```go
// Basic types
bool, string
int, int8, int16, int32, int64
uint, uint8, uint16, uint32, uint64
float32, float64
complex64, complex128
byte, rune

// Type conversion
i := 42
f := float64(i)
s := string(i)
```

## Control Flow

```go
// If
if x > 0 {
    // ...
} else if x < 0 {
    // ...
} else {
    // ...
}

// If with init
if v := getValue(); v > 0 {
    // v is scoped to if
}

// For
for i := 0; i < 10; i++ {
    // ...
}

// While-style
for condition {
    // ...
}

// Infinite
for {
    // break to exit
}

// Range
for i, v := range slice {
    // i: index, v: value
}

// Switch
switch x {
case 1:
    // ...
case 2, 3:
    // ...
default:
    // ...
}

// Switch with init
switch v := getValue(); v {
case 1:
    // ...
}

// Type switch
switch v := x.(type) {
case int:
    // v is int
case string:
    // v is string
}
```

## Functions

```go
// Basic
func add(a, b int) int {
    return a + b
}

// Multiple returns
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

// Named returns
func swap(a, b int) (x, y int) {
    x, y = b, a
    return  // naked return
}

// Variadic
func sum(nums ...int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}

// Closure
func adder() func(int) int {
    sum := 0
    return func(x int) int {
        sum += x
        return sum
    }
}
```

## Arrays & Slices

```go
// Array
var arr [5]int
arr := [3]int{1, 2, 3}
arr := [...]int{1, 2, 3}

// Slice
slice := []int{1, 2, 3}
slice := make([]int, 5)      // len=5, cap=5
slice := make([]int, 5, 10)  // len=5, cap=10

// Operations
slice = append(slice, 4)
copy(dest, src)
slice[1:3]  // [1, 2]
slice[:3]   // [0, 1, 2]
slice[2:]   // [2, 3, 4]
```

## Maps

```go
// Declaration
m := make(map[string]int)
m := map[string]int{"a": 1, "b": 2}

// Operations
m["key"] = value
value := m["key"]
value, ok := m["key"]  // Check existence
delete(m, "key")

// Range
for key, value := range m {
    // ...
}
```

## Structs

```go
// Definition
type Person struct {
    Name string
    Age  int
}

// Creation
p := Person{"John", 30}
p := Person{Name: "John", Age: 30}
p := Person{Name: "John"}  // Age: 0

// Access
p.Name = "Jane"

// Anonymous
point := struct{ X, Y int }{1, 2}
```

## Pointers

```go
var p *int
i := 42
p = &i
*p = 21  // Dereference

// With structs
p := &Person{Name: "John"}
p.Name = "Jane"  // Auto-dereference
```

## Interfaces

```go
type Reader interface {
    Read() ([]byte, error)
}

// Implementation (implicit)
type File struct{}
func (f File) Read() ([]byte, error) {
    return nil, nil
}

// Empty interface
var x interface{} = 42
```

## Methods

```go
type Person struct {
    Name string
}

// Value receiver
func (p Person) Greet() string {
    return "Hello, " + p.Name
}

// Pointer receiver
func (p *Person) SetName(name string) {
    p.Name = name
}
```

## Error Handling

```go
// Check error
result, err := doSomething()
if err != nil {
    return err
}

// Create error
errors.New("error message")
fmt.Errorf("error: %v", value)
```

## Goroutines & Channels

```go
// Goroutine
go func() {
    // runs concurrently
}()

// Channel
ch := make(chan int)
ch := make(chan int, 10)  // Buffered

// Send/Receive
ch <- value      // Send
value := <-ch    // Receive
value, ok := <-ch  // Check if closed

// Close
close(ch)

// Select
select {
case v := <-ch1:
    // ...
case ch2 <- v:
    // ...
default:
    // ...
}
```

## Defer, Panic, Recover

```go
// Defer (LIFO order)
defer fmt.Println("last")

// Panic
panic("something went wrong")

// Recover
func safe() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered:", r)
        }
    }()
    panic("error")
}
```

## Packages & Imports

```go
package main

import (
    "fmt"
    "time"
    alias "package/path"
    . "package"  // Import all
    _ "package"  // Side effects only
)
```
