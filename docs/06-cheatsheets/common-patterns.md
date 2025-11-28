# Common Go Patterns

## Functional Options Pattern

```go
type Server struct {
    host    string
    port    int
    timeout time.Duration
}

type Option func(*Server)

func WithHost(host string) Option {
    return func(s *Server) {
        s.host = host
    }
}

func WithPort(port int) Option {
    return func(s *Server) {
        s.port = port
    }
}

func NewServer(opts ...Option) *Server {
    s := &Server{
        host: "localhost",
        port: 8080,
        timeout: 30 * time.Second,
    }
    for _, opt := range opts {
        opt(s)
    }
    return s
}

// Usage
server := NewServer(
    WithHost("example.com"),
    WithPort(9000),
)
```

## Builder Pattern

```go
type QueryBuilder struct {
    query  strings.Builder
    params []interface{}
}

func NewQueryBuilder() *QueryBuilder {
    return &QueryBuilder{}
}

func (qb *QueryBuilder) Select(fields ...string) *QueryBuilder {
    qb.query.WriteString("SELECT ")
    qb.query.WriteString(strings.Join(fields, ", "))
    return qb
}

func (qb *QueryBuilder) From(table string) *QueryBuilder {
    qb.query.WriteString(" FROM ")
    qb.query.WriteString(table)
    return qb
}

func (qb *QueryBuilder) Where(condition string, params ...interface{}) *QueryBuilder {
    qb.query.WriteString(" WHERE ")
    qb.query.WriteString(condition)
    qb.params = append(qb.params, params...)
    return qb
}

func (qb *QueryBuilder) Build() (string, []interface{}) {
    return qb.query.String(), qb.params
}

// Usage
query, params := NewQueryBuilder().
    Select("id", "name").
    From("users").
    Where("age > ?", 18).
    Build()
```

## Factory Pattern

```go
type Database interface {
    Connect() error
}

type MySQL struct{}
func (m *MySQL) Connect() error { return nil }

type PostgreSQL struct{}
func (p *PostgreSQL) Connect() error { return nil }

func NewDatabase(dbType string) Database {
    switch dbType {
    case "mysql":
        return &MySQL{}
    case "postgres":
        return &PostgreSQL{}
    default:
        return nil
    }
}
```

## Singleton Pattern

```go
type Database struct {
    conn *sql.DB
}

var (
    instance *Database
    once     sync.Once
)

func GetInstance() *Database {
    once.Do(func() {
        instance = &Database{}
        instance.conn, _ = sql.Open("driver", "dsn")
    })
    return instance
}
```

## Worker Pool Pattern

```go
func workerPool(jobs <-chan Job, results chan<- Result, numWorkers int) {
    var wg sync.WaitGroup

    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for job := range jobs {
                results <- process(job)
            }
        }()
    }

    wg.Wait()
    close(results)
}

// Usage
jobs := make(chan Job, 100)
results := make(chan Result, 100)

go workerPool(jobs, results, 5)

// Send jobs
for _, job := range allJobs {
    jobs <- job
}
close(jobs)

// Collect results
for result := range results {
    // Handle result
}
```

## Pipeline Pattern

```go
func generator(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, n := range nums {
            out <- n
        }
    }()
    return out
}

func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            out <- n * n
        }
    }()
    return out
}

func filter(in <-chan int, f func(int) bool) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            if f(n) {
                out <- n
            }
        }
    }()
    return out
}

// Usage
numbers := generator(1, 2, 3, 4, 5)
squared := square(numbers)
filtered := filter(squared, func(n int) bool { return n > 10 })

for n := range filtered {
    fmt.Println(n)
}
```

## Fan-Out/Fan-In Pattern

```go
func fanOut(in <-chan int, n int) []<-chan int {
    channels := make([]<-chan int, n)
    for i := 0; i < n; i++ {
        ch := make(chan int)
        channels[i] = ch
        go func(out chan<- int) {
            defer close(out)
            for val := range in {
                out <- process(val)
            }
        }(ch)
    }
    return channels
}

func fanIn(channels ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup

    for _, ch := range channels {
        wg.Add(1)
        go func(c <-chan int) {
            defer wg.Done()
            for val := range c {
                out <- val
            }
        }(ch)
    }

    go func() {
        wg.Wait()
        close(out)
    }()

    return out
}
```

## Rate Limiting Pattern

```go
// Token bucket
type RateLimiter struct {
    tokens chan struct{}
}

func NewRateLimiter(rate int, burst int) *RateLimiter {
    rl := &RateLimiter{
        tokens: make(chan struct{}, burst),
    }

    // Fill tokens
    for i := 0; i < burst; i++ {
        rl.tokens <- struct{}{}
    }

    // Refill tokens
    go func() {
        ticker := time.NewTicker(time.Second / time.Duration(rate))
        for range ticker.C {
            select {
            case rl.tokens <- struct{}{}:
            default:
            }
        }
    }()

    return rl
}

func (rl *RateLimiter) Allow() bool {
    select {
    case <-rl.tokens:
        return true
    default:
        return false
    }
}

// Usage
limiter := NewRateLimiter(10, 20)  // 10/sec, burst 20
if limiter.Allow() {
    // Process request
}
```

## Circuit Breaker Pattern

```go
type CircuitBreaker struct {
    maxFailures int
    timeout     time.Duration
    failures    int
    lastAttempt time.Time
    state       string  // "closed", "open", "half-open"
    mu          sync.Mutex
}

func (cb *CircuitBreaker) Call(fn func() error) error {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    if cb.state == "open" {
        if time.Since(cb.lastAttempt) > cb.timeout {
            cb.state = "half-open"
        } else {
            return errors.New("circuit breaker open")
        }
    }

    err := fn()

    if err != nil {
        cb.failures++
        cb.lastAttempt = time.Now()
        if cb.failures >= cb.maxFailures {
            cb.state = "open"
        }
        return err
    }

    cb.failures = 0
    cb.state = "closed"
    return nil
}
```

## Retry Pattern

```go
func retry(attempts int, delay time.Duration, fn func() error) error {
    var err error
    for i := 0; i < attempts; i++ {
        err = fn()
        if err == nil {
            return nil
        }
        time.Sleep(delay)
        delay *= 2  // Exponential backoff
    }
    return err
}

// Usage
err := retry(3, time.Second, func() error {
    return makeAPICall()
})
```

## Context Propagation Pattern

```go
type contextKey string

const userIDKey contextKey = "userID"

func WithUserID(ctx context.Context, userID int) context.Context {
    return context.WithValue(ctx, userIDKey, userID)
}

func GetUserID(ctx context.Context) (int, bool) {
    userID, ok := ctx.Value(userIDKey).(int)
    return userID, ok
}

// Middleware
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        userID := getUserFromToken(r)
        ctx := WithUserID(r.Context(), userID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

## Table-Driven Testing Pattern

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive", 2, 3, 5},
        {"negative", -1, -2, -3},
        {"zero", 0, 0, 0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Add(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("got %d, want %d", result, tt.expected)
            }
        })
    }
}
```
