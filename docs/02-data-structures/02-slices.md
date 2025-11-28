# Slices

## Slice Basics

### Declaration and Initialization
```go
// Declare nil slice
var slice []int  // nil, length 0, capacity 0

// Empty slice literal
slice := []int{}  // not nil, length 0, capacity 0

// Initialize with values
slice := []int{1, 2, 3, 4, 5}

// Create slice with make
slice := make([]int, 5)       // length 5, capacity 5, values are zeros
slice := make([]int, 5, 10)   // length 5, capacity 10

// From array
arr := [5]int{1, 2, 3, 4, 5}
slice := arr[:]     // All elements
slice := arr[1:4]   // Elements at index 1, 2, 3
slice := arr[:3]    // Elements at index 0, 1, 2
slice := arr[2:]    // Elements from index 2 to end
```

### Slice Properties
```go
slice := []int{1, 2, 3, 4, 5}

len(slice)   // 5 - number of elements
cap(slice)   // 5 - capacity of underlying array

// Nil slice
var nilSlice []int
len(nilSlice)  // 0
cap(nilSlice)  // 0
nilSlice == nil  // true

// Empty slice
emptySlice := []int{}
len(emptySlice)  // 0
cap(emptySlice)  // 0
emptySlice == nil  // false
```

## Understanding Slices

### Slice Internals
A slice is a descriptor containing:
- Pointer to underlying array
- Length (current number of elements)
- Capacity (size of underlying array)

```go
// Slice structure (conceptually)
type slice struct {
    ptr *array
    len int
    cap int
}
```

### Slice vs Array
```go
// Array - fixed size, value type
arr := [5]int{1, 2, 3, 4, 5}

// Slice - dynamic size, reference type
slice := []int{1, 2, 3, 4, 5}

// Slicing creates views into same array
s1 := slice[0:3]  // [1, 2, 3]
s2 := slice[2:5]  // [3, 4, 5]
s1[2] = 100       // Changes slice and s2
// slice is now [1, 2, 100, 4, 5]
// s2 is now [100, 4, 5]
```

## Accessing and Modifying

```go
slice := []int{10, 20, 30, 40, 50}

// Access elements
first := slice[0]    // 10
last := slice[len(slice)-1]  // 50

// Modify elements
slice[0] = 100
slice[2] = 300

// Slicing
sub := slice[1:4]    // [20, 300, 40]
sub := slice[:3]     // [100, 20, 300]
sub := slice[2:]     // [300, 40, 50]
sub := slice[:]      // [100, 20, 300, 40, 50] (full slice)
```

## Append Operation

```go
// Append one element
slice := []int{1, 2, 3}
slice = append(slice, 4)  // [1, 2, 3, 4]

// Append multiple elements
slice = append(slice, 5, 6, 7)  // [1, 2, 3, 4, 5, 6, 7]

// Append another slice
other := []int{8, 9, 10}
slice = append(slice, other...)  // [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]

// Append to nil slice (creates new slice)
var slice []int
slice = append(slice, 1, 2, 3)  // [1, 2, 3]
```

### Append and Capacity
```go
slice := make([]int, 0, 5)
fmt.Printf("len=%d cap=%d\n", len(slice), cap(slice))  // len=0 cap=5

slice = append(slice, 1)  // len=1 cap=5
slice = append(slice, 2)  // len=2 cap=5
slice = append(slice, 3, 4, 5)  // len=5 cap=5

// Exceeding capacity allocates new array (typically doubles)
slice = append(slice, 6)  // len=6 cap=10 (capacity doubled)
```

## Copy Operation

```go
// copy(dst, src) copies elements, returns number copied
src := []int{1, 2, 3, 4, 5}
dst := make([]int, 3)

n := copy(dst, src)  // dst is [1, 2, 3], n is 3

// Copy full slice
dst := make([]int, len(src))
copy(dst, src)

// Overlapping copy is safe
slice := []int{1, 2, 3, 4, 5}
copy(slice[2:], slice[:])  // [1, 2, 1, 2, 3]
```

## Iterating Slices

```go
slice := []int{10, 20, 30, 40, 50}

// Traditional for loop
for i := 0; i < len(slice); i++ {
    fmt.Printf("slice[%d] = %d\n", i, slice[i])
}

// Range with index and value
for i, v := range slice {
    fmt.Printf("Index: %d, Value: %d\n", i, v)
}

// Range with just value
for _, v := range slice {
    fmt.Println(v)
}

// Range with just index
for i := range slice {
    fmt.Println(i)
}
```

## Multidimensional Slices

```go
// 2D slice (slice of slices)
matrix := [][]int{
    {1, 2, 3},
    {4, 5, 6},
    {7, 8, 9},
}

// Create 2D slice with make
rows, cols := 3, 4
matrix := make([][]int, rows)
for i := range matrix {
    matrix[i] = make([]int, cols)
}

// Jagged arrays (rows can have different lengths)
jagged := [][]int{
    {1, 2},
    {3, 4, 5},
    {6},
}
```

## Common Operations

### Insert Element
```go
// Insert at position
func insert(slice []int, index int, value int) []int {
    // Grow slice by one
    slice = append(slice, 0)
    // Shift elements right
    copy(slice[index+1:], slice[index:])
    // Insert value
    slice[index] = value
    return slice
}

slice := []int{1, 2, 4, 5}
slice = insert(slice, 2, 3)  // [1, 2, 3, 4, 5]
```

### Remove Element
```go
// Remove at index
func remove(slice []int, index int) []int {
    return append(slice[:index], slice[index+1:]...)
}

slice := []int{1, 2, 3, 4, 5}
slice = remove(slice, 2)  // [1, 2, 4, 5]

// Remove preserving order
slice = append(slice[:i], slice[i+1:]...)

// Remove without preserving order (faster)
slice[i] = slice[len(slice)-1]
slice = slice[:len(slice)-1]
```

### Filter
```go
func filter(slice []int, predicate func(int) bool) []int {
    result := []int{}
    for _, v := range slice {
        if predicate(v) {
            result = append(result, v)
        }
    }
    return result
}

slice := []int{1, 2, 3, 4, 5, 6}
evens := filter(slice, func(n int) bool { return n%2 == 0 })
// [2, 4, 6]

// In-place filter (more efficient)
func filterInPlace(slice []int, predicate func(int) bool) []int {
    n := 0
    for _, v := range slice {
        if predicate(v) {
            slice[n] = v
            n++
        }
    }
    return slice[:n]
}
```

### Map
```go
func mapSlice(slice []int, transform func(int) int) []int {
    result := make([]int, len(slice))
    for i, v := range slice {
        result[i] = transform(v)
    }
    return result
}

slice := []int{1, 2, 3, 4, 5}
doubled := mapSlice(slice, func(n int) int { return n * 2 })
// [2, 4, 6, 8, 10]
```

### Reduce
```go
func reduce(slice []int, initial int, reducer func(int, int) int) int {
    result := initial
    for _, v := range slice {
        result = reducer(result, v)
    }
    return result
}

slice := []int{1, 2, 3, 4, 5}
sum := reduce(slice, 0, func(a, b int) int { return a + b })
// 15
```

### Contains
```go
func contains(slice []int, target int) bool {
    for _, v := range slice {
        if v == target {
            return true
        }
    }
    return false
}

slice := []int{1, 2, 3, 4, 5}
exists := contains(slice, 3)  // true
```

### Reverse
```go
func reverse(slice []int) {
    for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
        slice[i], slice[j] = slice[j], slice[i]
    }
}

slice := []int{1, 2, 3, 4, 5}
reverse(slice)  // [5, 4, 3, 2, 1]
```

### Deduplicate
```go
func deduplicate(slice []int) []int {
    seen := make(map[int]bool)
    result := []int{}
    for _, v := range slice {
        if !seen[v] {
            seen[v] = true
            result = append(result, v)
        }
    }
    return result
}

slice := []int{1, 2, 2, 3, 3, 3, 4, 5, 5}
unique := deduplicate(slice)  // [1, 2, 3, 4, 5]
```

## Slice Tricks

### Expand
```go
// Expand slice by inserting empty space
slice = append(slice[:i], append(make([]int, j), slice[i:]...)...)
```

### Prepend
```go
slice = append([]int{x}, slice...)
```

### Pop
```go
// Pop last
x, slice := slice[len(slice)-1], slice[:len(slice)-1]

// Pop first
x, slice := slice[0], slice[1:]
```

### Clear
```go
// Clear slice (keep capacity)
slice = slice[:0]

// Clear and release memory
slice = nil
```

### Compact (Remove consecutive duplicates)
```go
func compact(slice []int) []int {
    if len(slice) == 0 {
        return slice
    }
    j := 0
    for i := 1; i < len(slice); i++ {
        if slice[i] != slice[j] {
            j++
            slice[j] = slice[i]
        }
    }
    return slice[:j+1]
}
```

## Performance Tips

### Pre-allocate Capacity
```go
// Bad: Multiple allocations
var slice []int
for i := 0; i < 1000; i++ {
    slice = append(slice, i)
}

// Good: Single allocation
slice := make([]int, 0, 1000)
for i := 0; i < 1000; i++ {
    slice = append(slice, i)
}

// Best: If you know the final length
slice := make([]int, 1000)
for i := 0; i < 1000; i++ {
    slice[i] = i
}
```

### Avoid Slice Leaks
```go
// Memory leak: slice keeps reference to large array
var data []byte = getHugeData()  // 1GB
smallSlice := data[:5]  // Keeps entire 1GB in memory!

// Fix: Copy to new slice
smallSlice := make([]byte, 5)
copy(smallSlice, data[:5])
data = nil  // Allow GC to collect
```

## Common Patterns

### Append if not exists
```go
if !contains(slice, value) {
    slice = append(slice, value)
}
```

### Batch append
```go
slice = append(slice, item1, item2, item3)
```

### Slice as stack
```go
// Push
stack = append(stack, item)

// Pop
item, stack := stack[len(stack)-1], stack[:len(stack)-1]

// Peek
item := stack[len(stack)-1]
```

### Slice as queue
```go
// Enqueue
queue = append(queue, item)

// Dequeue
item, queue := queue[0], queue[1:]
```

## Gotchas

❌ **Append may reallocate**
```go
slice1 := make([]int, 0, 5)
slice2 := slice1

slice1 = append(slice1, 1, 2, 3, 4, 5, 6)  // Reallocates!
// slice1 and slice2 now point to different arrays
```

❌ **Slice expressions share storage**
```go
original := []int{1, 2, 3, 4, 5}
sub := original[1:4]  // [2, 3, 4]
sub[0] = 100
// original is now [1, 100, 3, 4, 5]
```

❌ **Append result must be assigned**
```go
slice := []int{1, 2, 3}
append(slice, 4)  // Wrong! Result is lost
slice = append(slice, 4)  // Correct
```

❌ **Range creates copies of values**
```go
type Point struct{ x, y int }
points := []Point{{1, 2}, {3, 4}}

for _, p := range points {
    p.x = 0  // Modifies copy, not original!
}

// Use index
for i := range points {
    points[i].x = 0  // Modifies original
}
```

## Best Practices

1. **Pre-allocate with make** when size is known
2. **Always assign append result** back to slice
3. **Use copy for full slices** to avoid sharing
4. **Use slices instead of arrays** for most cases
5. **Check for nil vs empty** when it matters
6. **Avoid memory leaks** by copying small portions
7. **Use capacity wisely** - understand growth patterns
8. **Prefer index loops** when modifying elements
9. **Use range for read-only** iteration
10. **Document if function modifies** the input slice
