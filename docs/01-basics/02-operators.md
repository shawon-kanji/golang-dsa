# Operators

## Arithmetic Operators

```go
+   // Addition
-   // Subtraction
*   // Multiplication
/   // Division
%   // Modulus (remainder)

// Examples
a := 10 + 5   // 15
b := 10 - 5   // 5
c := 10 * 5   // 50
d := 10 / 5   // 2
e := 10 % 3   // 1

// Increment/Decrement (statement, not expression)
i++  // i = i + 1
i--  // i = i - 1

// Note: ++i and --i are NOT valid in Go
```

## Comparison Operators

```go
==  // Equal to
!=  // Not equal to
<   // Less than
<=  // Less than or equal to
>   // Greater than
>=  // Greater than or equal to

// Examples
5 == 5   // true
5 != 3   // true
5 < 10   // true
5 <= 5   // true
10 > 5   // true
10 >= 10 // true
```

## Logical Operators

```go
&&  // Logical AND
||  // Logical OR
!   // Logical NOT

// Examples
true && false   // false
true || false   // true
!true          // false

// Short-circuit evaluation
if user != nil && user.IsActive() {
    // user.IsActive() only called if user != nil
}
```

## Bitwise Operators

```go
&   // Bitwise AND
|   // Bitwise OR
^   // Bitwise XOR
&^  // Bit clear (AND NOT)
<<  // Left shift
>>  // Right shift

// Examples
a := 12        // 1100 in binary
b := 10        // 1010 in binary

c := a & b     // 1000 = 8
d := a | b     // 1110 = 14
e := a ^ b     // 0110 = 6
f := a &^ b    // 0100 = 4

g := 1 << 3    // 1000 = 8 (multiply by 2^3)
h := 8 >> 2    // 0010 = 2 (divide by 2^2)
```

## Assignment Operators

```go
=   // Simple assignment
+=  // Add and assign
-=  // Subtract and assign
*=  // Multiply and assign
/=  // Divide and assign
%=  // Modulus and assign
&=  // Bitwise AND and assign
|=  // Bitwise OR and assign
^=  // Bitwise XOR and assign
<<= // Left shift and assign
>>= // Right shift and assign

// Examples
x := 10
x += 5     // x = x + 5 = 15
x -= 3     // x = x - 3 = 12
x *= 2     // x = x * 2 = 24
x /= 4     // x = x / 4 = 6
x %= 4     // x = x % 4 = 2
```

## Address and Pointer Operators

```go
&   // Address-of operator
*   // Dereference operator

// Examples
x := 10
ptr := &x      // ptr holds the address of x
value := *ptr  // value = 10 (dereference ptr)

*ptr = 20      // Changes x to 20
```

## Operator Precedence (Highest to Lowest)

```go
1. * / % << >> & &^
2. + - | ^
3. == != < <= > >=
4. &&
5. ||

// Examples
result := 2 + 3 * 4      // 14 (not 20)
result := (2 + 3) * 4    // 20 (use parentheses)

result := 10 > 5 && 3 < 7  // true
result := 10 > 5 || 3 > 7  // true
```

## String Operators

```go
+   // Concatenation
+=  // Append

// Examples
s1 := "Hello"
s2 := "World"
s3 := s1 + " " + s2  // "Hello World"

s1 += " Go"          // "Hello Go"
```

## Special Operators

### Type Operators
```go
// Type assertion
value := interfaceVar.(Type)

// Type assertion with check
value, ok := interfaceVar.(Type)

// Type switch
switch v := i.(type) {
case int:
    // v is int
case string:
    // v is string
}
```

### Channel Operators
```go
<-   // Send/Receive

// Send to channel
ch <- value

// Receive from channel
value := <-ch

// Receive and ignore
<-ch
```

## Operator Examples

### Swap without temp variable
```go
a, b := 10, 20
a, b = b, a  // a=20, b=10
```

### Check if number is even/odd
```go
isEven := num % 2 == 0
isOdd := num % 2 != 0
// Or using bitwise
isEven := num & 1 == 0
```

### Check if power of 2
```go
isPowerOfTwo := n > 0 && (n & (n - 1)) == 0
```

### Toggle boolean
```go
flag := true
flag = !flag  // false
```

### Set/Clear/Toggle bit
```go
// Set bit at position n
num |= (1 << n)

// Clear bit at position n
num &^= (1 << n)

// Toggle bit at position n
num ^= (1 << n)

// Check if bit is set
isSet := (num & (1 << n)) != 0
```

## Common Patterns

### Min/Max without if
```go
// Not idiomatic, but possible
min := a
if b < min {
    min = b
}

// Or use a function
func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}
```

### Absolute value
```go
func abs(n int) int {
    if n < 0 {
        return -n
    }
    return n
}
```

### Clamp value
```go
// Clamp value between min and max
if value < min {
    value = min
} else if value > max {
    value = max
}
```

## Gotchas

❌ **Integer division truncates**
```go
result := 5 / 2  // 2, not 2.5
// Use float for decimal
result := 5.0 / 2.0  // 2.5
```

❌ **No automatic type coercion**
```go
var i int = 10
var f float64 = 3.14
result := i + f  // Error! Cannot mix types
result := float64(i) + f  // OK
```

❌ **++ and -- are statements, not expressions**
```go
x := i++  // Error!
x := ++i  // Error!

// Correct way
i++
x := i
```

❌ **Bitwise XOR vs logical XOR**
```go
// Bitwise XOR
a ^ b

// Logical XOR (not built-in)
(a || b) && !(a && b)
```

## Best Practices

1. **Use parentheses for clarity** when mixing operators
2. **Prefer `+=` style** for compound operations
3. **Use bit operations** for flags and low-level optimization
4. **Short-circuit evaluation** - order matters in logical operations
5. **Type conversions should be explicit** and intentional
