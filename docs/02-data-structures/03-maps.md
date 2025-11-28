# Maps

## Map Basics

### Declaration and Initialization
```go
// Declare nil map
var m map[string]int  // nil map, cannot add elements

// Make map
m := make(map[string]int)

// Map literal
m := map[string]int{
    "apple":  5,
    "banana": 3,
    "orange": 7,
}

// Empty map literal
m := map[string]int{}

// With capacity hint (optimization)
m := make(map[string]int, 100)
```

### Map Properties
```go
m := map[string]int{"a": 1, "b": 2, "c": 3}

// Length
len(m)  // 3

// Maps are reference types
m2 := m
m2["d"] = 4
// m also has "d": 4

// Nil map
var nilMap map[string]int
len(nilMap)  // 0
nilMap == nil  // true
// nilMap["key"] = 1  // Panic! Cannot assign to nil map
```

## Adding and Updating

```go
m := make(map[string]int)

// Add or update
m["apple"] = 5
m["banana"] = 3

// Update existing key
m["apple"] = 10

// Conditional update
if val, exists := m["apple"]; exists {
    m["apple"] = val + 1
}
```

## Accessing Elements

```go
m := map[string]int{"apple": 5, "banana": 3}

// Get value
value := m["apple"]  // 5

// Get non-existent key returns zero value
value := m["notfound"]  // 0

// Check if key exists
value, exists := m["apple"]
if exists {
    fmt.Println("Found:", value)
}

// Common pattern
if value, ok := m["apple"]; ok {
    // Key exists, use value
    fmt.Println(value)
} else {
    // Key doesn't exist
    fmt.Println("Not found")
}
```

## Deleting Elements

```go
m := map[string]int{"apple": 5, "banana": 3, "orange": 7}

// Delete key
delete(m, "banana")

// Delete non-existent key is safe (no-op)
delete(m, "notfound")

// Check before delete
if _, exists := m["apple"]; exists {
    delete(m, "apple")
}

// Clear all elements
for key := range m {
    delete(m, key)
}

// Or just create new map
m = make(map[string]int)
```

## Iterating Maps

```go
m := map[string]int{"apple": 5, "banana": 3, "orange": 7}

// Iterate key and value
for key, value := range m {
    fmt.Printf("%s: %d\n", key, value)
}

// Iterate keys only
for key := range m {
    fmt.Println(key)
}

// Iterate values only
for _, value := range m {
    fmt.Println(value)
}

// Note: Iteration order is random!
```

## Maps with Different Types

### String Keys
```go
m := map[string]string{
    "name":  "John",
    "email": "john@example.com",
}
```

### Integer Keys
```go
m := map[int]string{
    1: "one",
    2: "two",
    3: "three",
}
```

### Struct Keys (must be comparable)
```go
type Point struct {
    x, y int
}

m := map[Point]string{
    {0, 0}: "origin",
    {1, 2}: "point A",
}
```

### Complex Values
```go
// Map of slices
m := map[string][]int{
    "even": {2, 4, 6, 8},
    "odd":  {1, 3, 5, 7},
}

// Map of maps
m := map[string]map[string]int{
    "fruits": {"apple": 5, "banana": 3},
    "veggies": {"carrot": 7, "broccoli": 2},
}

// Map of structs
type User struct {
    Name  string
    Age   int
    Email string
}

users := map[int]User{
    1: {"Alice", 30, "alice@example.com"},
    2: {"Bob", 25, "bob@example.com"},
}
```

## Common Operations

### Check if Key Exists
```go
_, exists := m["key"]

// Or
if value, ok := m["key"]; ok {
    // Key exists
}
```

### Get with Default Value
```go
func getOrDefault(m map[string]int, key string, defaultValue int) int {
    if value, ok := m[key]; ok {
        return value
    }
    return defaultValue
}

value := getOrDefault(m, "key", 0)
```

### Merge Maps
```go
func merge(m1, m2 map[string]int) map[string]int {
    result := make(map[string]int)
    for k, v := range m1 {
        result[k] = v
    }
    for k, v := range m2 {
        result[k] = v  // Overwrites if key exists
    }
    return result
}
```

### Copy Map
```go
func copyMap(m map[string]int) map[string]int {
    result := make(map[string]int, len(m))
    for k, v := range m {
        result[k] = v
    }
    return result
}
```

### Invert Map
```go
func invert(m map[string]int) map[int]string {
    result := make(map[int]string)
    for k, v := range m {
        result[v] = k
    }
    return result
}

original := map[string]int{"a": 1, "b": 2}
inverted := invert(original)  // map[int]string{1: "a", 2: "b"}
```

### Filter Map
```go
func filter(m map[string]int, predicate func(string, int) bool) map[string]int {
    result := make(map[string]int)
    for k, v := range m {
        if predicate(k, v) {
            result[k] = v
        }
    }
    return result
}

m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
evens := filter(m, func(k string, v int) bool { return v%2 == 0 })
// map[string]int{"b": 2, "d": 4}
```

### Keys and Values
```go
func keys(m map[string]int) []string {
    result := make([]string, 0, len(m))
    for k := range m {
        result = append(result, k)
    }
    return result
}

func values(m map[string]int) []int {
    result := make([]int, 0, len(m))
    for _, v := range m {
        result = append(result, v)
    }
    return result
}
```

## Map as Set

```go
// Set using map[T]bool
set := make(map[string]bool)

// Add elements
set["apple"] = true
set["banana"] = true

// Check membership
if set["apple"] {
    fmt.Println("apple is in set")
}

// Remove element
delete(set, "banana")

// Set using map[T]struct{} (saves memory)
type Set map[string]struct{}

set := make(Set)
set["apple"] = struct{}{}
set["banana"] = struct{}{}

// Check membership
if _, exists := set["apple"]; exists {
    fmt.Println("apple is in set")
}

// Set operations
func union(a, b Set) Set {
    result := make(Set)
    for k := range a {
        result[k] = struct{}{}
    }
    for k := range b {
        result[k] = struct{}{}
    }
    return result
}

func intersection(a, b Set) Set {
    result := make(Set)
    for k := range a {
        if _, ok := b[k]; ok {
            result[k] = struct{}{}
        }
    }
    return result
}

func difference(a, b Set) Set {
    result := make(Set)
    for k := range a {
        if _, ok := b[k]; !ok {
            result[k] = struct{}{}
        }
    }
    return result
}
```

## Concurrent Map Access

```go
// Maps are NOT safe for concurrent use!

// Use sync.RWMutex
type SafeMap struct {
    mu sync.RWMutex
    m  map[string]int
}

func (sm *SafeMap) Get(key string) (int, bool) {
    sm.mu.RLock()
    defer sm.mu.RUnlock()
    value, ok := sm.m[key]
    return value, ok
}

func (sm *SafeMap) Set(key string, value int) {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    sm.m[key] = value
}

// Or use sync.Map (for specific use cases)
var m sync.Map

m.Store("key", "value")
value, ok := m.Load("key")
m.Delete("key")
```

## Ordering Maps

### Sort by Keys
```go
m := map[string]int{"c": 3, "a": 1, "b": 2}

// Get keys and sort
keys := make([]string, 0, len(m))
for k := range m {
    keys = append(keys, k)
}
sort.Strings(keys)

// Iterate in sorted order
for _, k := range keys {
    fmt.Printf("%s: %d\n", k, m[k])
}
```

### Sort by Values
```go
type kv struct {
    key   string
    value int
}

// Convert to slice
var pairs []kv
for k, v := range m {
    pairs = append(pairs, kv{k, v})
}

// Sort by value
sort.Slice(pairs, func(i, j int) bool {
    return pairs[i].value < pairs[j].value
})

for _, pair := range pairs {
    fmt.Printf("%s: %d\n", pair.key, pair.value)
}
```

## Advanced Patterns

### Frequency Counter
```go
func frequency(items []string) map[string]int {
    freq := make(map[string]int)
    for _, item := range items {
        freq[item]++
    }
    return freq
}

items := []string{"apple", "banana", "apple", "orange", "banana", "apple"}
freq := frequency(items)
// map[string]int{"apple": 3, "banana": 2, "orange": 1}
```

### Grouping
```go
type Person struct {
    Name string
    Age  int
}

func groupByAge(people []Person) map[int][]Person {
    groups := make(map[int][]Person)
    for _, p := range people {
        groups[p.Age] = append(groups[p.Age], p)
    }
    return groups
}
```

### Memoization
```go
var cache = make(map[int]int)

func fibonacci(n int) int {
    if n <= 1 {
        return n
    }
    if val, ok := cache[n]; ok {
        return val
    }
    result := fibonacci(n-1) + fibonacci(n-2)
    cache[n] = result
    return result
}
```

### Nested Maps
```go
// Auto-initialize nested maps
m := make(map[string]map[string]int)

// Check and initialize
if m["group1"] == nil {
    m["group1"] = make(map[string]int)
}
m["group1"]["key1"] = 10

// Or use helper function
func getOrCreateMap(m map[string]map[string]int, key string) map[string]int {
    if m[key] == nil {
        m[key] = make(map[string]int)
    }
    return m[key]
}

getOrCreateMap(m, "group2")["key2"] = 20
```

## Gotchas

❌ **Map iteration order is random**
```go
m := map[string]int{"a": 1, "b": 2, "c": 3}
for k := range m {
    fmt.Println(k)  // Order is not guaranteed!
}
```

❌ **Nil map cannot be assigned to**
```go
var m map[string]int  // nil map
// m["key"] = 1  // Panic!

// Must initialize first
m = make(map[string]int)
m["key"] = 1  // OK
```

❌ **Maps are not comparable**
```go
m1 := map[string]int{"a": 1}
m2 := map[string]int{"a": 1}
// if m1 == m2 {}  // Error! Cannot compare maps

// Can only compare to nil
if m1 == nil {}  // OK
```

❌ **Not safe for concurrent access**
```go
// This will cause race condition
go func() {
    m["key"] = 1
}()
go func() {
    m["key"] = 2
}()

// Use sync.RWMutex or sync.Map instead
```

❌ **Modifying map during iteration**
```go
// Safe to delete during iteration
for k := range m {
    if someCondition {
        delete(m, k)  // OK
    }
}

// Adding during iteration may or may not be seen
for k := range m {
    m["new"+k] = 1  // May or may not iterate over new keys
}
```

## Best Practices

1. **Initialize before use** - don't use nil maps
2. **Use make for empty maps** - not map literal if no initial values
3. **Check existence** with `value, ok := m[key]`
4. **Use struct{} for sets** - saves memory compared to bool
5. **Protect concurrent access** - use mutexes or sync.Map
6. **Don't rely on iteration order** - sort keys if order matters
7. **Consider capacity** - use `make(map[K]V, capacity)` for large maps
8. **Document side effects** - if function modifies input map
9. **Clear large maps** - set to nil to allow GC
10. **Use meaningful key types** - that reflect domain

## Performance Considerations

```go
// Pre-allocate capacity if size is known
m := make(map[string]int, 1000)

// Map lookup is O(1) average case
value := m[key]

// Deletion is O(1)
delete(m, key)

// Iteration is O(n)
for k, v := range m {}

// Memory: maps grow but don't shrink
// Create new map if many deletions occur
```
