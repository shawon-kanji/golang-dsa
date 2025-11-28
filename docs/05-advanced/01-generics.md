# Generics (Go 1.18+)

## Generic Basics

### Generic Functions
```go
// Generic function with type parameter
func Print[T any](s []T) {
    for _, v := range s {
        fmt.Println(v)
    }
}

// Usage
Print([]int{1, 2, 3})
Print([]string{"a", "b", "c"})
```

### Type Constraints
```go
// Using comparable constraint
func Index[T comparable](s []T, x T) int {
    for i, v := range s {
        if v == x {
            return i
        }
    }
    return -1
}

// Using any constraint
func FirstElement[T any](s []T) T {
    return s[0]
}
```

## Generic Types

### Generic Struct
```go
type Stack[T any] struct {
    items []T
}

func (s *Stack[T]) Push(item T) {
    s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
    if len(s.items) == 0 {
        var zero T
        return zero, false
    }
    item := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return item, true
}

// Usage
stack := Stack[int]{}
stack.Push(1)
stack.Push(2)
val, ok := stack.Pop()
```

### Generic Interface
```go
type Container[T any] interface {
    Add(T)
    Get() T
}

type Box[T any] struct {
    value T
}

func (b *Box[T]) Add(v T) {
    b.value = v
}

func (b *Box[T]) Get() T {
    return b.value
}
```

## Custom Constraints

### Interface Constraints
```go
type Number interface {
    int | int64 | float64
}

func Sum[T Number](nums []T) T {
    var sum T
    for _, n := range nums {
        sum += n
    }
    return sum
}
```

### Combining Constraints
```go
type Ordered interface {
    ~int | ~int64 | ~float64 | ~string
}

func Min[T Ordered](a, b T) T {
    if a < b {
        return a
    }
    return b
}
```

### Method Constraints
```go
type Stringer interface {
    String() string
}

func Concat[T Stringer](items []T) string {
    var result string
    for _, item := range items {
        result += item.String()
    }
    return result
}
```

## Common Generic Patterns

### Map Function
```go
func Map[T, U any](s []T, f func(T) U) []U {
    result := make([]U, len(s))
    for i, v := range s {
        result[i] = f(v)
    }
    return result
}

// Usage
nums := []int{1, 2, 3}
doubled := Map(nums, func(n int) int { return n * 2 })
```

### Filter Function
```go
func Filter[T any](s []T, f func(T) bool) []T {
    var result []T
    for _, v := range s {
        if f(v) {
            result = append(result, v)
        }
    }
    return result
}

// Usage
evens := Filter(nums, func(n int) bool { return n%2 == 0 })
```

### Reduce Function
```go
func Reduce[T, U any](s []T, init U, f func(U, T) U) U {
    result := init
    for _, v := range s {
        result = f(result, v)
    }
    return result
}

// Usage
sum := Reduce(nums, 0, func(acc, n int) int { return acc + n })
```

### Contains Function
```go
func Contains[T comparable](s []T, v T) bool {
    for _, item := range s {
        if item == v {
            return true
        }
    }
    return false
}
```

## Generic Data Structures

### Linked List
```go
type Node[T any] struct {
    Value T
    Next  *Node[T]
}

type LinkedList[T any] struct {
    head *Node[T]
}

func (l *LinkedList[T]) Add(value T) {
    node := &Node[T]{Value: value, Next: l.head}
    l.head = node
}
```

### Binary Tree
```go
type TreeNode[T any] struct {
    Value T
    Left  *TreeNode[T]
    Right *TreeNode[T]
}

func (n *TreeNode[T]) Insert(value T, compare func(T, T) int) {
    if compare(value, n.Value) < 0 {
        if n.Left == nil {
            n.Left = &TreeNode[T]{Value: value}
        } else {
            n.Left.Insert(value, compare)
        }
    } else {
        if n.Right == nil {
            n.Right = &TreeNode[T]{Value: value}
        } else {
            n.Right.Insert(value, compare)
        }
    }
}
```

### Generic Queue
```go
type Queue[T any] struct {
    items []T
}

func (q *Queue[T]) Enqueue(item T) {
    q.items = append(q.items, item)
}

func (q *Queue[T]) Dequeue() (T, bool) {
    if len(q.items) == 0 {
        var zero T
        return zero, false
    }
    item := q.items[0]
    q.items = q.items[1:]
    return item, true
}
```

## Type Inference
```go
// Type parameters can be inferred
func Print[T any](v T) {
    fmt.Println(v)
}

Print(42)       // T inferred as int
Print("hello")  // T inferred as string

// Explicit type parameters
Print[int](42)
Print[string]("hello")
```

## Best Practices

1. **Use generics sparingly** - don't over-engineer
2. **Prefer specific types** when possible
3. **Use meaningful constraints** - not just `any`
4. **Type parameter names** - single letter (T, U, V) or descriptive
5. **Document constraints** - explain what types are valid
6. **Avoid complex type parameters** - keep it simple
7. **Test with multiple types** - ensure generics work as expected
