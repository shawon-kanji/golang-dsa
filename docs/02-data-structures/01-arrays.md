# Arrays

## Array Basics

### Declaration
```go
// Declare array with size
var arr [5]int  // Array of 5 integers, initialized to zeros

// Declare and initialize
var arr = [5]int{1, 2, 3, 4, 5}

// Short declaration
arr := [5]int{1, 2, 3, 4, 5}

// Let compiler count size
arr := [...]int{1, 2, 3, 4, 5}  // Size is 5

// Partial initialization (rest are zero values)
arr := [5]int{1, 2}  // [1, 2, 0, 0, 0]

// Initialize specific indices
arr := [5]int{0: 10, 2: 30, 4: 50}  // [10, 0, 30, 0, 50]
```

### Array Properties
```go
arr := [5]int{1, 2, 3, 4, 5}

// Length (fixed at compile time)
length := len(arr)  // 5

// Arrays are values (not references)
arr2 := arr  // Creates a copy
arr2[0] = 100
// arr[0] is still 1, arr2[0] is 100
```

## Accessing Elements

```go
arr := [5]int{10, 20, 30, 40, 50}

// Access by index
first := arr[0]   // 10
last := arr[4]    // 50

// Modify elements
arr[0] = 100
arr[2] = 300

// Out of bounds causes panic
// value := arr[10]  // Panic: index out of range
```

## Iterating Arrays

### For Loop
```go
arr := [5]int{1, 2, 3, 4, 5}

// Traditional for loop
for i := 0; i < len(arr); i++ {
    fmt.Printf("arr[%d] = %d\n", i, arr[i])
}

// Range with index and value
for i, v := range arr {
    fmt.Printf("Index: %d, Value: %d\n", i, v)
}

// Range with just index
for i := range arr {
    fmt.Println(i)
}

// Range with just value
for _, v := range arr {
    fmt.Println(v)
}
```

## Multidimensional Arrays

```go
// 2D array
var matrix [3][3]int

// Initialize 2D array
matrix := [3][3]int{
    {1, 2, 3},
    {4, 5, 6},
    {7, 8, 9},
}

// Access elements
value := matrix[1][2]  // 6

// Iterate 2D array
for i := 0; i < len(matrix); i++ {
    for j := 0; j < len(matrix[i]); j++ {
        fmt.Printf("%d ", matrix[i][j])
    }
    fmt.Println()
}

// Range over 2D array
for i, row := range matrix {
    for j, val := range row {
        fmt.Printf("matrix[%d][%d] = %d\n", i, j, val)
    }
}

// 3D array
cube := [2][2][2]int{
    {
        {1, 2},
        {3, 4},
    },
    {
        {5, 6},
        {7, 8},
    },
}
```

## Array Operations

### Copy Array
```go
// Arrays are copied by value
arr1 := [3]int{1, 2, 3}
arr2 := arr1  // Creates a complete copy

// Manual copy
arr1 := [3]int{1, 2, 3}
var arr2 [3]int
for i, v := range arr1 {
    arr2[i] = v
}
```

### Compare Arrays
```go
// Arrays of same type and size can be compared
arr1 := [3]int{1, 2, 3}
arr2 := [3]int{1, 2, 3}
arr3 := [3]int{4, 5, 6}

fmt.Println(arr1 == arr2)  // true
fmt.Println(arr1 == arr3)  // false
fmt.Println(arr1 != arr3)  // true
```

### Pass Array to Function
```go
// Arrays are passed by value (copied)
func modify(arr [5]int) {
    arr[0] = 100  // Modifies the copy, not original
}

func main() {
    arr := [5]int{1, 2, 3, 4, 5}
    modify(arr)
    fmt.Println(arr[0])  // Still 1
}

// Pass by reference using pointer
func modifyByRef(arr *[5]int) {
    arr[0] = 100  // Modifies original
}

func main() {
    arr := [5]int{1, 2, 3, 4, 5}
    modifyByRef(&arr)
    fmt.Println(arr[0])  // 100
}
```

## Common Operations

### Find Element
```go
func find(arr [5]int, target int) int {
    for i, v := range arr {
        if v == target {
            return i
        }
    }
    return -1  // Not found
}

arr := [5]int{10, 20, 30, 40, 50}
index := find(arr, 30)  // 2
```

### Sum of Elements
```go
func sum(arr [5]int) int {
    total := 0
    for _, v := range arr {
        total += v
    }
    return total
}

arr := [5]int{1, 2, 3, 4, 5}
result := sum(arr)  // 15
```

### Max/Min
```go
func max(arr [5]int) int {
    maxVal := arr[0]
    for _, v := range arr {
        if v > maxVal {
            maxVal = v
        }
    }
    return maxVal
}

func min(arr [5]int) int {
    minVal := arr[0]
    for _, v := range arr {
        if v < minVal {
            minVal = v
        }
    }
    return minVal
}
```

### Reverse Array
```go
func reverse(arr [5]int) [5]int {
    for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
        arr[i], arr[j] = arr[j], arr[i]
    }
    return arr
}

arr := [5]int{1, 2, 3, 4, 5}
reversed := reverse(arr)  // [5, 4, 3, 2, 1]
```

## Arrays vs Slices

| Feature | Array | Slice |
|---------|-------|-------|
| Size | Fixed | Dynamic |
| Memory | Value type | Reference type |
| Pass to function | Copied | Reference passed |
| Common usage | Less common | Very common |
| Type | `[n]T` | `[]T` |

```go
// Array - fixed size
var arr [5]int = [5]int{1, 2, 3, 4, 5}

// Slice - dynamic size
var slice []int = []int{1, 2, 3, 4, 5}

// Array size is part of type
var a1 [5]int
var a2 [10]int
// a1 and a2 are different types!

// Slices of same element type are compatible
var s1 []int
var s2 []int
// s1 and s2 are the same type
```

## Performance Considerations

```go
// Arrays are efficient for:
// - Fixed-size collections
// - Stack allocation
// - Passing small amounts of data by value

// Example: Small fixed buffer
var buffer [64]byte

// Example: Matrix operations
var matrix [100][100]float64

// For most use cases, prefer slices:
// - Dynamic sizing
// - Reference semantics
// - More idiomatic in Go
```

## Common Patterns

### Fixed-size buffer
```go
const BufferSize = 1024
var buffer [BufferSize]byte

n, err := reader.Read(buffer[:])
```

### Lookup tables
```go
var daysInMonth = [12]int{
    31, 28, 31, 30, 31, 30,
    31, 31, 30, 31, 30, 31,
}

days := daysInMonth[monthIndex-1]
```

### State machines
```go
const (
    StateIdle = iota
    StateRunning
    StateComplete
)

var stateNames = [3]string{
    "Idle",
    "Running",
    "Complete",
}

fmt.Println(stateNames[StateRunning])  // "Running"
```

## Gotchas

❌ **Array size is part of type**
```go
var a1 [5]int
var a2 [10]int
// Cannot assign a1 to a2 - different types!

// Must use same size
var a3 [5]int = a1  // OK
```

❌ **Arrays are passed by value**
```go
func modify(arr [100]int) {
    // This copies all 100 integers!
    arr[0] = 1
}

// Use slice or pointer instead
func modify(arr *[100]int) {
    arr[0] = 1  // Modifies original
}

func modifySlice(arr []int) {
    arr[0] = 1  // Modifies original (slice references array)
}
```

❌ **Out of bounds not checked at compile time**
```go
arr := [5]int{1, 2, 3, 4, 5}
// This compiles but panics at runtime
// value := arr[10]  // Panic!
```

❌ **Range creates copies**
```go
arr := [3]Point{{1, 2}, {3, 4}, {5, 6}}

for _, p := range arr {
    p.x = 0  // Modifies copy, not original!
}

// Use index instead
for i := range arr {
    arr[i].x = 0  // Modifies original
}
```

## Best Practices

1. **Prefer slices over arrays** - more flexible and idiomatic
2. **Use arrays for fixed-size data** - constants, lookup tables
3. **Use `[...]` for array literals** - compiler determines size
4. **Pass large arrays by pointer** - avoid expensive copies
5. **Use range for iteration** - more readable than index loops
6. **Check bounds manually** - when accessing with variables
7. **Use arrays for stack allocation** - when size is known and small
8. **Consider memory layout** - arrays are contiguous in memory

## When to Use Arrays

✅ **Use arrays when:**
- Size is known at compile time and won't change
- Need value semantics (independent copies)
- Building lookup tables or constants
- Working with fixed-size buffers
- Stack allocation is preferred

❌ **Don't use arrays when:**
- Size needs to be dynamic
- Passing to functions frequently
- Need reference semantics
- Working with collections
- Size is large (use slices with pointers)

## Array to Slice Conversion

```go
arr := [5]int{1, 2, 3, 4, 5}

// Create slice from array
slice := arr[:]        // All elements
slice := arr[1:4]      // Elements 1, 2, 3
slice := arr[:3]       // Elements 0, 1, 2
slice := arr[2:]       // Elements 2, 3, 4

// Slice references the array
slice[0] = 100
fmt.Println(arr[0])    // 100 (original array modified)
```
