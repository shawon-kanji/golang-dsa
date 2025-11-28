# Control Flow

## If Statement

### Basic If
```go
if x > 10 {
    fmt.Println("x is greater than 10")
}

// If-else
if x > 10 {
    fmt.Println("Greater")
} else {
    fmt.Println("Not greater")
}

// If-else if-else
if x > 10 {
    fmt.Println("Greater than 10")
} else if x > 5 {
    fmt.Println("Greater than 5")
} else {
    fmt.Println("5 or less")
}
```

### If with Short Statement
```go
// Initialize variable in if statement
if num := 9; num < 0 {
    fmt.Println(num, "is negative")
} else if num < 10 {
    fmt.Println(num, "has 1 digit")
} else {
    fmt.Println(num, "has multiple digits")
}
// num is only in scope within if/else blocks

// Common pattern: error checking
if err := doSomething(); err != nil {
    return err
}

// File operations
if file, err := os.Open("file.txt"); err == nil {
    defer file.Close()
    // Use file
}
```

## For Loop

Go has only `for` - no `while` or `do-while`

### Basic For
```go
// Traditional for loop
for i := 0; i < 10; i++ {
    fmt.Println(i)
}

// Multiple variables
for i, j := 0, 10; i < j; i, j = i+1, j-1 {
    fmt.Println(i, j)
}
```

### For as While
```go
i := 0
for i < 10 {
    fmt.Println(i)
    i++
}
```

### Infinite Loop
```go
for {
    // Runs forever until break
    if condition {
        break
    }
}
```

### For Range

#### Range over Slice/Array
```go
numbers := []int{1, 2, 3, 4, 5}

// Index and value
for i, num := range numbers {
    fmt.Printf("Index: %d, Value: %d\n", i, num)
}

// Just index
for i := range numbers {
    fmt.Println(i)
}

// Just value
for _, num := range numbers {
    fmt.Println(num)
}
```

#### Range over Map
```go
m := map[string]int{"a": 1, "b": 2, "c": 3}

// Key and value
for key, value := range m {
    fmt.Printf("%s: %d\n", key, value)
}

// Just keys
for key := range m {
    fmt.Println(key)
}
```

#### Range over String
```go
s := "Hello 世界"

// Iterates over runes (Unicode code points)
for i, r := range s {
    fmt.Printf("Index: %d, Rune: %c\n", i, r)
}
```

#### Range over Channel
```go
ch := make(chan int)

// Receives values until channel is closed
for value := range ch {
    fmt.Println(value)
}
```

## Break and Continue

```go
// Break - exit loop
for i := 0; i < 10; i++ {
    if i == 5 {
        break  // Exits the loop
    }
    fmt.Println(i)
}

// Continue - skip to next iteration
for i := 0; i < 10; i++ {
    if i%2 == 0 {
        continue  // Skip even numbers
    }
    fmt.Println(i)  // Only prints odd numbers
}

// Break with label (break from nested loop)
outer:
for i := 0; i < 3; i++ {
    for j := 0; j < 3; j++ {
        if i == 1 && j == 1 {
            break outer  // Breaks from outer loop
        }
        fmt.Println(i, j)
    }
}

// Continue with label
outer:
for i := 0; i < 3; i++ {
    for j := 0; j < 3; j++ {
        if i == 1 {
            continue outer  // Continue outer loop
        }
        fmt.Println(i, j)
    }
}
```

## Switch Statement

### Basic Switch
```go
day := "Monday"

switch day {
case "Monday":
    fmt.Println("Start of work week")
case "Friday":
    fmt.Println("End of work week")
case "Saturday", "Sunday":
    fmt.Println("Weekend!")
default:
    fmt.Println("Midweek day")
}
```

### Switch with Short Statement
```go
switch day := time.Now().Weekday(); day {
case time.Saturday, time.Sunday:
    fmt.Println("Weekend")
default:
    fmt.Println("Weekday")
}
```

### Switch without Expression (like if-else chain)
```go
hour := 15

switch {
case hour < 12:
    fmt.Println("Good morning")
case hour < 17:
    fmt.Println("Good afternoon")
default:
    fmt.Println("Good evening")
}
```

### Switch with Fallthrough
```go
switch num := 2; num {
case 1:
    fmt.Println("One")
case 2:
    fmt.Println("Two")
    fallthrough  // Executes next case regardless
case 3:
    fmt.Println("Three")  // This will also print
default:
    fmt.Println("Other")
}
// Output: Two, Three
```

### Type Switch
```go
func describe(i interface{}) {
    switch v := i.(type) {
    case int:
        fmt.Printf("Integer: %d\n", v)
    case string:
        fmt.Printf("String: %s\n", v)
    case bool:
        fmt.Printf("Boolean: %t\n", v)
    default:
        fmt.Printf("Unknown type: %T\n", v)
    }
}

describe(42)      // Integer: 42
describe("hello") // String: hello
describe(true)    // Boolean: true
```

## Goto Statement

Use sparingly - can make code hard to follow:

```go
func main() {
    i := 0
loop:
    fmt.Println(i)
    i++
    if i < 5 {
        goto loop
    }
}

// Error handling pattern
func doWork() error {
    // ... code ...
    if err != nil {
        goto cleanup
    }

    // ... more code ...
    if err != nil {
        goto cleanup
    }

    return nil

cleanup:
    // Cleanup code
    return err
}
```

## Common Patterns

### Loop until condition
```go
for {
    item, err := getNext()
    if err != nil {
        break
    }
    process(item)
}
```

### Iterate with index and value
```go
for i, v := range slice {
    fmt.Printf("slice[%d] = %v\n", i, v)
}
```

### Reverse iteration
```go
for i := len(slice) - 1; i >= 0; i-- {
    fmt.Println(slice[i])
}
```

### Iterate every nth element
```go
for i := 0; i < len(slice); i += 3 {
    fmt.Println(slice[i])
}
```

### Early return pattern
```go
func process(items []Item) error {
    for _, item := range items {
        if !item.IsValid() {
            return fmt.Errorf("invalid item: %v", item)
        }
        // Process item
    }
    return nil
}
```

### Loop with timeout
```go
timeout := time.After(5 * time.Second)
for {
    select {
    case <-timeout:
        return fmt.Errorf("timeout")
    case result := <-resultCh:
        return result
    }
}
```

## Gotchas

❌ **No parentheses required, but braces are mandatory**
```go
if x > 10 {  // OK
    // code
}

if (x > 10) {  // Parentheses optional
    // code
}

if x > 10   // Error! Missing braces
    fmt.Println(x)
```

❌ **Variable scope in if statement**
```go
if x := 10; x > 5 {
    fmt.Println(x)  // OK
}
fmt.Println(x)  // Error! x is out of scope
```

❌ **Range creates copies, not references**
```go
items := []Item{{value: 1}, {value: 2}}
for _, item := range items {
    item.value = 0  // Doesn't modify original!
}

// Use index instead
for i := range items {
    items[i].value = 0  // Modifies original
}

// Or use pointers
for _, item := range items {
    item := &items[i]  // Get pointer
    item.value = 0
}
```

❌ **Switch doesn't fall through by default**
```go
switch x {
case 1:
    fmt.Println("One")
case 2:
    fmt.Println("Two")  // Won't execute if x == 1
}

// Use fallthrough explicitly
switch x {
case 1:
    fmt.Println("One")
    fallthrough
case 2:
    fmt.Println("Two")  // Executes if x == 1
}
```

❌ **Break in switch only breaks switch, not loop**
```go
for i := 0; i < 10; i++ {
    switch i {
    case 5:
        break  // Only breaks switch, not loop
    }
}

// Use label to break loop
loop:
for i := 0; i < 10; i++ {
    switch i {
    case 5:
        break loop  // Breaks the loop
    }
}
```

## Best Practices

1. **Use short variable declarations in if** for error checking
2. **Prefer switch over long if-else chains**
3. **Use range for iterating** over slices, maps, and channels
4. **Break complex conditions** into separate if statements for clarity
5. **Avoid goto** unless necessary for cleanup patterns
6. **Use labels sparingly** for breaking nested loops
7. **Keep loop bodies short** - extract to functions if needed
