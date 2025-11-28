# Standard Library Reference

## fmt - Formatting

```go
// Print
fmt.Print("text")
fmt.Println("text")
fmt.Printf("format: %v", value)

// Format to string
s := fmt.Sprintf("format: %v", value)

// Scan input
var name string
fmt.Scan(&name)
fmt.Scanf("%s", &name)

// Format verbs
%v    // Default format
%+v   // With field names
%#v   // Go syntax
%T    // Type
%t    // Boolean
%d    // Integer
%f    // Float
%s    // String
%q    // Quoted string
%p    // Pointer
```

## strings - String Operations

```go
import "strings"

// Contains/Search
strings.Contains("hello", "ll")  // true
strings.HasPrefix("hello", "he") // true
strings.HasSuffix("hello", "lo") // true
strings.Index("hello", "ll")     // 2
strings.Count("hello", "l")      // 2

// Modify
strings.ToUpper("hello")         // HELLO
strings.ToLower("HELLO")         // hello
strings.TrimSpace(" hello ")     // hello
strings.Trim("--hello--", "-")   // hello
strings.Replace("hello", "l", "L", 2)  // heLLo
strings.ReplaceAll("hello", "l", "L")  // heLLo

// Split/Join
strings.Split("a,b,c", ",")      // []string{"a", "b", "c"}
strings.Join([]string{"a", "b"}, ",")  // "a,b"
strings.Fields("a b  c")         // []string{"a", "b", "c"}

// Builder
var b strings.Builder
b.WriteString("hello")
b.String()  // "hello"
```

## strconv - String Conversions

```go
import "strconv"

// String to int
i, err := strconv.Atoi("42")
i, err := strconv.ParseInt("42", 10, 64)

// Int to string
s := strconv.Itoa(42)
s := strconv.FormatInt(42, 10)

// String to float
f, err := strconv.ParseFloat("3.14", 64)

// Float to string
s := strconv.FormatFloat(3.14, 'f', 2, 64)

// Bool conversions
b, err := strconv.ParseBool("true")
s := strconv.FormatBool(true)
```

## time - Time and Duration

```go
import "time"

// Current time
now := time.Now()

// Create time
t := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)

// Parse/Format
t, err := time.Parse("2006-01-02", "2024-01-01")
s := time.Now().Format("2006-01-02 15:04:05")

// Duration
d := 5 * time.Second
time.Sleep(d)

// Arithmetic
future := now.Add(24 * time.Hour)
diff := future.Sub(now)

// Comparison
if t1.Before(t2) { }
if t1.After(t2) { }
if t1.Equal(t2) { }

// Timers
timer := time.NewTimer(5 * time.Second)
<-timer.C

ticker := time.NewTicker(1 * time.Second)
for t := range ticker.C {
    // Every second
}
```

## io - Input/Output

```go
import "io"

// Readers
r := strings.NewReader("hello")
bytes, _ := io.ReadAll(r)
io.Copy(dst, src)

// Writers
w := &bytes.Buffer{}
io.WriteString(w, "hello")

// Closer
defer r.Close()
```

## os - Operating System

```go
import "os"

// Files
f, err := os.Open("file.txt")
f, err := os.Create("file.txt")
f.Write([]byte("data"))
f.Close()

// Read file
data, err := os.ReadFile("file.txt")

// Write file
err := os.WriteFile("file.txt", data, 0644)

// Environment
val := os.Getenv("PATH")
os.Setenv("KEY", "value")

// Arguments
args := os.Args  // [program, arg1, arg2]

// Working directory
dir, _ := os.Getwd()
os.Chdir("/path")

// File info
info, _ := os.Stat("file.txt")
info.Size()
info.ModTime()
info.IsDir()

// Remove
os.Remove("file.txt")
os.RemoveAll("dir")
```

## path/filepath - File Paths

```go
import "path/filepath"

// Join
path := filepath.Join("dir", "subdir", "file.txt")

// Base/Dir
base := filepath.Base("/path/to/file.txt")  // file.txt
dir := filepath.Dir("/path/to/file.txt")    // /path/to

// Ext
ext := filepath.Ext("file.txt")  // .txt

// Absolute
abs, _ := filepath.Abs("file.txt")

// Walk directory
filepath.Walk("dir", func(path string, info os.FileInfo, err error) error {
    fmt.Println(path)
    return nil
})
```

## encoding/json - JSON

```go
import "encoding/json"

// Marshal
data, err := json.Marshal(obj)
data, err := json.MarshalIndent(obj, "", "  ")

// Unmarshal
var obj Type
err := json.Unmarshal(data, &obj)

// Encoder/Decoder
encoder := json.NewEncoder(w)
encoder.Encode(obj)

decoder := json.NewDecoder(r)
decoder.Decode(&obj)

// Struct tags
type Person struct {
    Name string `json:"name"`
    Age  int    `json:"age,omitempty"`
}
```

## net/http - HTTP

```go
import "net/http"

// Client
resp, err := http.Get("http://example.com")
defer resp.Body.Close()
body, _ := io.ReadAll(resp.Body)

// POST
resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))

// Custom request
req, _ := http.NewRequest("GET", url, nil)
req.Header.Set("Authorization", "Bearer token")
client := &http.Client{Timeout: 10 * time.Second}
resp, _ := client.Do(req)

// Server
http.HandleFunc("/", handler)
http.ListenAndServe(":8080", nil)

// Handler
func handler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Hello"))
}
```

## sync - Synchronization

```go
import "sync"

// Mutex
var mu sync.Mutex
mu.Lock()
defer mu.Unlock()

// RWMutex
var rw sync.RWMutex
rw.RLock()
defer rw.RUnlock()

// WaitGroup
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    // work
}()
wg.Wait()

// Once
var once sync.Once
once.Do(func() {
    // Initialize
})

// Pool
pool := &sync.Pool{
    New: func() interface{} {
        return &Object{}
    },
}
obj := pool.Get().(*Object)
defer pool.Put(obj)
```

## context - Context

```go
import "context"

// Background
ctx := context.Background()
ctx := context.TODO()

// Cancel
ctx, cancel := context.WithCancel(parent)
defer cancel()

// Timeout
ctx, cancel := context.WithTimeout(parent, 5*time.Second)
defer cancel()

// Deadline
ctx, cancel := context.WithDeadline(parent, time.Now().Add(5*time.Second))
defer cancel()

// Values
ctx := context.WithValue(parent, key, value)
val := ctx.Value(key)
```

## log - Logging

```go
import "log"

// Print
log.Print("message")
log.Printf("format: %v", value)
log.Println("message")

// Fatal (exit)
log.Fatal("error")
log.Fatalf("error: %v", err)

// Panic
log.Panic("error")

// Custom logger
logger := log.New(os.Stdout, "PREFIX: ", log.LstdFlags)
logger.Println("message")
```

## errors - Error Handling

```go
import "errors"

// New
err := errors.New("error message")

// Format
err := fmt.Errorf("error: %w", originalErr)

// Is
if errors.Is(err, ErrNotFound) { }

// As
var target *MyError
if errors.As(err, &target) { }

// Unwrap
unwrapped := errors.Unwrap(err)
```

## sort - Sorting

```go
import "sort"

// Slices
ints := []int{3, 1, 2}
sort.Ints(ints)
sort.Strings([]string{"c", "a", "b"})

// Custom
sort.Slice(items, func(i, j int) bool {
    return items[i] < items[j]
})

// Search
i := sort.SearchInts(sortedInts, target)
```

## regexp - Regular Expressions

```go
import "regexp"

// Compile
re := regexp.MustCompile(`\d+`)

// Match
matched := re.MatchString("abc123")

// Find
result := re.FindString("abc123")      // "123"
results := re.FindAllString("a1b2", -1)  // ["1", "2"]

// Replace
result := re.ReplaceAllString("a1b2", "X")  // "aXbX"
```
