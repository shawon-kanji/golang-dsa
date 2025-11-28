# Reflection

## Reflection Basics

### reflect.TypeOf and reflect.ValueOf
```go
import "reflect"

var x float64 = 3.4

// Get type
t := reflect.TypeOf(x)
fmt.Println("Type:", t)  // float64

// Get value
v := reflect.ValueOf(x)
fmt.Println("Value:", v)  // 3.4
fmt.Println("Kind:", v.Kind())  // float64
```

### Type and Value Methods
```go
v := reflect.ValueOf(42)

fmt.Println(v.Type())     // int
fmt.Println(v.Kind())     // int
fmt.Println(v.Int())      // 42
fmt.Println(v.Interface()) // 42 as interface{}
```

## Inspecting Types

### Struct Fields
```go
type Person struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

p := Person{"Alice", 30}
t := reflect.TypeOf(p)

for i := 0; i < t.NumField(); i++ {
    field := t.Field(i)
    fmt.Printf("Field: %s, Type: %s, Tag: %s\n",
        field.Name, field.Type, field.Tag.Get("json"))
}
```

### Method Inspection
```go
v := reflect.ValueOf(p)
t := v.Type()

for i := 0; i < t.NumMethod(); i++ {
    method := t.Method(i)
    fmt.Printf("Method: %s\n", method.Name)
}
```

## Modifying Values

### Setting Values
```go
var x float64 = 3.4
v := reflect.ValueOf(&x)  // Pointer to x
v = v.Elem()              // Dereference
v.SetFloat(7.1)           // Set new value
fmt.Println(x)            // 7.1
```

### Settability
```go
v := reflect.ValueOf(x)
// v.SetFloat(7.1)  // Panic! v is not settable

v := reflect.ValueOf(&x).Elem()
v.SetFloat(7.1)  // OK, v is settable
```

## Common Use Cases

### JSON Marshal/Unmarshal
```go
func structToMap(obj interface{}) map[string]interface{} {
    result := make(map[string]interface{})
    v := reflect.ValueOf(obj)
    t := v.Type()

    for i := 0; i < v.NumField(); i++ {
        field := t.Field(i)
        value := v.Field(i)
        result[field.Name] = value.Interface()
    }

    return result
}
```

### Function Invocation
```go
func callMethod(obj interface{}, method string, args ...interface{}) []interface{} {
    v := reflect.ValueOf(obj)
    m := v.MethodByName(method)

    in := make([]reflect.Value, len(args))
    for i, arg := range args {
        in[i] = reflect.ValueOf(arg)
    }

    results := m.Call(in)

    out := make([]interface{}, len(results))
    for i, result := range results {
        out[i] = result.Interface()
    }

    return out
}
```

### Deep Copy
```go
func deepCopy(src interface{}) interface{} {
    v := reflect.ValueOf(src)
    t := v.Type()

    if v.Kind() == reflect.Ptr {
        v = v.Elem()
        t = v.Type()
    }

    copy := reflect.New(t).Elem()

    for i := 0; i < v.NumField(); i++ {
        copy.Field(i).Set(v.Field(i))
    }

    return copy.Interface()
}
```

## Type Switches with Reflection
```go
func printType(v interface{}) {
    rv := reflect.ValueOf(v)

    switch rv.Kind() {
    case reflect.Int, reflect.Int64:
        fmt.Printf("Integer: %d\n", rv.Int())
    case reflect.String:
        fmt.Printf("String: %s\n", rv.String())
    case reflect.Slice:
        fmt.Printf("Slice of %d elements\n", rv.Len())
    default:
        fmt.Printf("Unknown type: %v\n", rv.Kind())
    }
}
```

## Best Practices

1. **Avoid reflection when possible** - it's slow and loses type safety
2. **Use interfaces** instead of reflection
3. **Cache reflect.Type** if used repeatedly
4. **Check CanSet()** before setting values
5. **Handle panics** - reflection can panic easily
6. **Use type assertions** for known types
7. **Document when using reflection** - explain why it's needed

## Performance Considerations
- Reflection is ~10-100x slower than direct access
- Use sparingly in hot paths
- Consider code generation alternatives
