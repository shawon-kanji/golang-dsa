# Go Concurrency Exercises (Moderate → Hard)

Practice these problems in order. Each builds on concepts from previous ones. Implement each as a standalone `main.go` in its own subfolder under `concurrency/`.

---

## Problem 1: Concurrent File Downloader (Worker Pool + WaitGroup)

**Difficulty:** Moderate

**Concepts:** goroutines, channels, sync.WaitGroup, worker pool pattern

**Problem:**

Build a concurrent file downloader that processes a list of URLs using a fixed-size worker pool.

**Requirements:**
1. Accept a list of URLs (can be simulated with `[]string`) and a configurable number of workers (e.g., 5).
2. Each worker picks a URL from a shared job channel, simulates an HTTP download (`time.Sleep` for random 100–500ms), and sends back a `DownloadResult{URL, BytesDownloaded, Error, Duration}`.
3. A collector goroutine reads results and prints a summary after all downloads finish.
4. If a "download" takes > 300ms, simulate a random failure (return error 30% of the time).
5. Print a final report: total files downloaded, total bytes, total failures, and total wall-clock time.

**Stretch:** Add retry logic — failed downloads go back into the job queue (max 3 retries per URL).

```go
type DownloadResult struct {
    URL            string
    BytesDownloaded int
    Err            error
    Duration       time.Duration
    Retries        int
}
```

**Expected Output:**
```
$ go run .
[Worker 3] Downloading https://example.com/file1.zip ... 247ms OK (14523 bytes)
[Worker 1] Downloading https://example.com/file2.zip ... 312ms FAILED (timeout)
[Worker 2] Downloading https://example.com/file3.zip ... 189ms OK (8291 bytes)
[Worker 5] Downloading https://example.com/file4.zip ... 401ms OK (22017 bytes)
[Worker 1] Retrying https://example.com/file2.zip (attempt 2/3) ... 156ms OK (9102 bytes)
[Worker 4] Downloading https://example.com/file5.zip ... 334ms FAILED (timeout)
[Worker 3] Downloading https://example.com/file6.zip ... 112ms OK (5530 bytes)
[Worker 4] Retrying https://example.com/file5.zip (attempt 2/3) ... 278ms OK (11204 bytes)
[Worker 2] Downloading https://example.com/file7.zip ... 203ms OK (7788 bytes)
[Worker 5] Downloading https://example.com/file8.zip ... 451ms FAILED (timeout)
[Worker 5] Retrying https://example.com/file8.zip (attempt 2/3) ... 389ms FAILED (timeout)
[Worker 5] Retrying https://example.com/file8.zip (attempt 3/3) ... 198ms OK (12300 bytes)

========== DOWNLOAD REPORT ==========
Total files:       8
Successful:        8 (70455 bytes)
Failed:            0
Retries used:      4
Wall-clock time:   1.247s
======================================
```

**Key learning:** Worker pool lifecycle, channel direction types (`chan<-`, `<-chan`), WaitGroup coordination, graceful shutdown of producers and consumers.

---

## Problem 2: Pipeline Data Processor (Pipeline + Fan-Out/Fan-In)

**Difficulty:** Moderate

**Concepts:** pipeline pattern, fan-out, fan-in, channel chaining

**Problem:**

Build a multi-stage data processing pipeline that reads raw log lines, parses them, filters, enriches, and aggregates results — all concurrently.

**Requirements:**
1. **Stage 1 (Generator):** Emit 10,000 simulated log lines like `"2026-03-20T10:15:30Z|INFO|user-service|user_login|user_id=123"`.
2. **Stage 2 (Parser):** Parse each line into a `LogEntry` struct. Fan out to 3 parser workers.
3. **Stage 3 (Filter):** Drop entries older than 1 hour or with level `DEBUG`. Single goroutine.
4. **Stage 4 (Enricher):** Fan out to 2 workers. Each looks up `user_id` from a simulated in-memory map and adds `username` field.
5. **Stage 5 (Aggregator):** Count events per service, per log level. Print final stats.

Each stage is connected by channels. The pipeline must shut down cleanly when the generator is done (close channels, propagate shutdown).

```go
type LogEntry struct {
    Timestamp time.Time
    Level     string
    Service   string
    Action    string
    UserID    int
    Username  string // filled by enricher
}
```

**Expected Output:**
```
$ go run .
[Generator] Emitting 10000 log lines...
[Parser-1] Started
[Parser-2] Started
[Parser-3] Started
[Filter] Started — dropping DEBUG and stale entries
[Enricher-1] Started
[Enricher-2] Started
[Aggregator] Collecting results...

[Generator] Done — 10000 lines emitted
[Parser-1] Parsed 3342 entries
[Parser-2] Parsed 3331 entries
[Parser-3] Parsed 3327 entries
[Filter] Passed 7214 / 10000 (filtered 2786: 1902 DEBUG, 884 stale)
[Enricher-1] Enriched 3611 entries
[Enricher-2] Enriched 3603 entries

========== PIPELINE RESULTS ==========
Events per service:
  user-service:     2408
  order-service:    2397
  payment-service:  2409

Events per log level:
  INFO:    3624
  WARN:    2156
  ERROR:   1434

Pipeline duration: 312ms
=======================================
```

**Key learning:** Channel chaining, close propagation through pipeline stages, fan-out/fan-in merging, back-pressure behavior.

---

## Problem 3: Request Coalescer / Singleflight (sync.Mutex + sync.Cond + channels)

**Difficulty:** Moderate–Hard

**Concepts:** sync.Mutex, sync.Cond (or channel-based signaling), deduplication, shared state

**Problem:**

Implement a `Singleflight`-style request coalescer from scratch (don't use `golang.org/x/sync/singleflight`).

**Requirements:**
1. Multiple goroutines may request the same key simultaneously (`Call(key string, fn func() (interface{}, error))`).
2. Only **one** goroutine actually executes `fn` for a given key. All other callers for the same key block and receive the same result.
3. Once the in-flight call completes, subsequent calls for the same key should trigger a new execution (not cache forever).
4. Demonstrate with 50 goroutines all calling `FetchUserProfile("user-42")` simultaneously — only 1 actual "fetch" should execute.

```go
type Singleflight struct {
    mu    sync.Mutex
    calls map[string]*call
}

type call struct {
    wg  sync.WaitGroup
    val interface{}
    err error
}

func (sf *Singleflight) Do(key string, fn func() (interface{}, error)) (interface{}, error)
```

**Expected Output:**
```
$ go run .
Launching 50 goroutines to fetch "user-42" simultaneously...

[Singleflight] Key="user-42" — executing fetch (goroutine #7 won the race)
[Fetch] Calling external API for user-42... (200ms simulated)
[Fetch] Got profile: {ID:42 Name:"Alice" Email:"alice@example.com"}

[Goroutine  1] result={ID:42 Name:"Alice"} err=<nil> (shared)
[Goroutine  2] result={ID:42 Name:"Alice"} err=<nil> (shared)
[Goroutine  3] result={ID:42 Name:"Alice"} err=<nil> (shared)
... (47 more identical lines)
[Goroutine 50] result={ID:42 Name:"Alice"} err=<nil> (shared)

Actual fetch calls made: 1
Total goroutines served: 50
De-duplication ratio:    50:1
```

**Key learning:** Protecting shared maps with mutexes, WaitGroup as a broadcast signal, understanding the difference between Mutex (exclusive access) and Cond/WaitGroup (coordination).

---

## Problem 4: Bounded Concurrent Cache with TTL (sync.RWMutex + time.AfterFunc)

**Difficulty:** Moderate–Hard

**Concepts:** sync.RWMutex, goroutine-based expiry, bounded concurrency, eviction

**Problem:**

Build a thread-safe in-memory cache with TTL expiration and max capacity.

**Requirements:**
1. `Get(key)` — returns value if present and not expired, otherwise `nil, false`.
2. `Set(key, value, ttl)` — stores value. If at max capacity, evict the oldest entry (LRU or FIFO).
3. TTL-based expiration: expired entries are lazily cleaned on access AND proactively cleaned by a background sweeper goroutine running every N seconds.
4. The background sweeper must be stoppable via a `Close()` method (use a `done` channel or `context.Context`).
5. Must be safe for concurrent reads (many goroutines calling Get) and writes (goroutines calling Set) — use `sync.RWMutex`.
6. Write a benchmark: spin up 100 goroutines doing random Gets and 20 goroutines doing Sets. Measure throughput.

```go
type Cache struct {
    mu       sync.RWMutex
    items    map[string]*entry
    maxSize  int
    done     chan struct{}
}

type entry struct {
    value     interface{}
    expiresAt time.Time
    accessedAt time.Time
}
```

**Expected Output:**
```
$ go run .
[Cache] Started (maxSize=100, sweepInterval=2s)
[Benchmark] Spawning 100 readers + 20 writers...

[Sweeper] Tick — removed 12 expired entries (88 remaining)
[Sweeper] Tick — removed 8 expired entries (97 remaining)
[Sweeper] Tick — removed 15 expired entries (100 remaining)
[Eviction] Cache full — evicted key="item-23" (oldest access)
[Eviction] Cache full — evicted key="item-07" (oldest access)
[Sweeper] Tick — removed 11 expired entries (96 remaining)

========== CACHE BENCHMARK (5s) ==========
Total Get ops:    247,832
Total Set ops:     49,211
Cache hits:       189,104 (76.3%)
Cache misses:      58,728 (23.7%)
Evictions:            412
Expired (swept):      198
Expired (lazy):       304
Avg Get latency:     ~2μs
Avg Set latency:     ~5μs
===========================================

[Cache] Closed — sweeper stopped
```

**Key learning:** RWMutex for read-heavy workloads, background goroutine lifecycle management, graceful shutdown patterns, measuring contention.

---

## Problem 5: Distributed Task Scheduler (Context + Select + Timeout)

**Difficulty:** Hard

**Concepts:** context.WithTimeout, context.WithCancel, select, graceful shutdown, error groups

**Problem:**

Build a task scheduler that distributes tasks to workers with deadlines, cancellation, and graceful shutdown.

**Requirements:**
1. Tasks have different priorities (high/medium/low) and deadlines.
2. A scheduler goroutine reads tasks from a submission channel and assigns them to available workers based on priority (high-priority tasks jump the queue).
3. Each worker receives a `context.Context` with the task's deadline. If the task exceeds its deadline, it's cancelled and reported as timed out.
4. Support a global shutdown signal: calling `scheduler.Shutdown(ctx)` should:
   - Stop accepting new tasks.
   - Wait for in-flight tasks to finish (up to the shutdown context's deadline).
   - Force-cancel remaining tasks if shutdown deadline is exceeded.
5. Use `errgroup.Group` (from `golang.org/x/sync/errgroup`) for worker management.

```go
type Task struct {
    ID       string
    Priority Priority // High=0, Medium=1, Low=2
    Deadline time.Duration
    Execute  func(ctx context.Context) error
}

type Scheduler struct {
    workers    int
    tasks      chan Task
    ctx        context.Context
    cancel     context.CancelFunc
}

func (s *Scheduler) Submit(task Task) error       // returns error if shut down
func (s *Scheduler) Shutdown(ctx context.Context)  // graceful shutdown
```

**Expected Output:**
```
$ go run .
[Scheduler] Started with 4 workers
[Scheduler] Submitted task-001 (priority=HIGH,  deadline=200ms)
[Scheduler] Submitted task-002 (priority=LOW,   deadline=500ms)
[Scheduler] Submitted task-003 (priority=HIGH,  deadline=100ms)
[Scheduler] Submitted task-004 (priority=MED,   deadline=300ms)
[Worker-1] Executing task-001 (HIGH) ...
[Worker-2] Executing task-003 (HIGH) ... ← high priority jumped queue
[Worker-1] task-001 completed in 87ms ✓
[Worker-3] Executing task-004 (MED) ...
[Worker-2] task-003 TIMEOUT after 100ms ✗ (context deadline exceeded)
[Worker-4] Executing task-002 (LOW) ...
[Worker-3] task-004 completed in 210ms ✓
[Worker-4] task-002 completed in 312ms ✓

[Scheduler] Shutdown initiated (grace period: 5s)
[Scheduler] Waiting for 0 in-flight tasks...
[Scheduler] All tasks drained

========== SCHEDULER REPORT ==========
Tasks submitted:  4
Tasks completed:  3
Tasks timed out:  1
Tasks cancelled:  0
Avg exec time:    203ms
=======================================
```

**Key learning:** Context propagation, nested contexts (parent cancel → child cancel), select with multiple channels, graceful shutdown patterns used in production systems.

---

## Problem 6: Pub/Sub Message Broker (Channels + sync.RWMutex)

**Difficulty:** Hard

**Concepts:** pub/sub pattern, dynamic subscriber management, buffered channels, goroutine leak prevention

**Problem:**

Build an in-process pub/sub message broker supporting topics, subscribers, and backpressure.

**Requirements:**
1. `Subscribe(topic string, bufferSize int) (<-chan Message, func())` — returns a receive-only channel and an unsubscribe function.
2. `Publish(topic string, data interface{})` — delivers message to all current subscribers of that topic.
3. If a subscriber's buffer is full, the publisher should **not block** — drop the message for that slow subscriber and increment a dropped-message counter.
4. Unsubscribe must: close the subscriber's channel, remove it from the topic's subscriber list, and not leak goroutines.
5. Support multiple topics, each with multiple subscribers.
6. Demonstrate with:
   - 3 topics: `"orders"`, `"payments"`, `"notifications"`
   - Fast and slow subscribers on each topic
   - A publisher blasting 1000 messages per topic
   - Print per-subscriber stats: received, dropped

```go
type Broker struct {
    mu          sync.RWMutex
    subscribers map[string][]*subscriber
}

type subscriber struct {
    ch      chan Message
    dropped int64  // atomic counter
}

type Message struct {
    Topic     string
    Data      interface{}
    Timestamp time.Time
}
```

**Expected Output:**
```
$ go run .
[Broker] Started
[Sub] fast-orders    subscribed to "orders"        (buffer=100)
[Sub] slow-orders    subscribed to "orders"        (buffer=5)
[Sub] fast-payments  subscribed to "payments"      (buffer=100)
[Sub] slow-payments  subscribed to "payments"      (buffer=5)
[Sub] fast-notifs    subscribed to "notifications" (buffer=100)
[Sub] slow-notifs    subscribed to "notifications" (buffer=5)

[Pub] Publishing 1000 messages to each topic...

[slow-orders]   ⚠ buffer full — dropping msg #47
[slow-orders]   ⚠ buffer full — dropping msg #48
[slow-payments]  ⚠ buffer full — dropping msg #51
... (more drops for slow subscribers)

[Pub] Done publishing

========== SUBSCRIBER STATS ==========
Topic: orders
  fast-orders:     received=1000  dropped=0
  slow-orders:     received=127   dropped=873

Topic: payments
  fast-payments:   received=1000  dropped=0
  slow-payments:   received=134   dropped=866

Topic: notifications
  fast-notifs:     received=1000  dropped=0
  slow-notifs:     received=119   dropped=881
=======================================

[Broker] All subscribers unsubscribed — no goroutine leaks
```

**Key learning:** Managing dynamic goroutine lifecycles, non-blocking sends with select/default, preventing goroutine leaks with unsubscribe, safe concurrent map access.

---

## Problem 7: Circuit Breaker (State Machine + sync/atomic + time.Ticker)

**Difficulty:** Hard

**Concepts:** sync/atomic, state machines, time.Ticker, goroutines for background monitoring

**Problem:**

Implement a circuit breaker that protects calls to an unreliable external service.

**Requirements:**
1. Three states: **Closed** (normal), **Open** (failing — reject all calls), **Half-Open** (testing recovery).
2. Transitions:
   - Closed → Open: when failure count in a rolling window exceeds threshold (e.g., 5 failures in 10 seconds).
   - Open → Half-Open: after a cooldown period (e.g., 5 seconds).
   - Half-Open → Closed: if a probe call succeeds.
   - Half-Open → Open: if a probe call fails.
3. `Execute(fn func() error) error` — runs `fn` if circuit is closed/half-open, returns `ErrCircuitOpen` if open.
4. Use `sync/atomic` for the state variable and failure counters (no mutex for the hot path).
5. A background goroutine manages the rolling window, resetting old failure counts.
6. Demonstrate with a simulated service that fails 80% of the time for 10 seconds, then recovers.

```go
type State int32
const (
    StateClosed   State = iota
    StateOpen
    StateHalfOpen
)

type CircuitBreaker struct {
    state        atomic.Int32
    failures     atomic.Int64
    successes    atomic.Int64
    threshold    int64
    cooldown     time.Duration
    lastFailure  atomic.Int64  // unix nano
}

func (cb *CircuitBreaker) Execute(fn func() error) error
```

**Expected Output:**
```
$ go run .
[CircuitBreaker] Initialized (threshold=5, cooldown=5s, window=10s)
[Simulator] Service will FAIL 80% for first 10s, then recover

[t=0.0s]  CLOSED   call #1  → FAIL   (failures: 1/5)
[t=0.1s]  CLOSED   call #2  → FAIL   (failures: 2/5)
[t=0.2s]  CLOSED   call #3  → OK     (failures: 2/5)
[t=0.3s]  CLOSED   call #4  → FAIL   (failures: 3/5)
[t=0.4s]  CLOSED   call #5  → FAIL   (failures: 4/5)
[t=0.5s]  CLOSED   call #6  → FAIL   (failures: 5/5) ← THRESHOLD REACHED
[t=0.5s]  ⚡ State: CLOSED → OPEN
[t=0.6s]  OPEN     call #7  → REJECTED (circuit open)
[t=0.7s]  OPEN     call #8  → REJECTED (circuit open)
... (all calls rejected while open)
[t=5.5s]  ⚡ State: OPEN → HALF-OPEN (cooldown elapsed)
[t=5.5s]  HALF-OPEN call #57 → FAIL (probe failed)
[t=5.5s]  ⚡ State: HALF-OPEN → OPEN
... (another cooldown cycle)
[t=10.5s] ⚡ State: OPEN → HALF-OPEN
[t=10.5s] HALF-OPEN call #112 → OK (probe succeeded!)
[t=10.5s] ⚡ State: HALF-OPEN → CLOSED
[t=10.6s] CLOSED   call #113 → OK
[t=10.7s] CLOSED   call #114 → OK

========== CIRCUIT BREAKER STATS ==========
Total calls:      150
Successful:        43
Failed (real):     22
Rejected (open):   85
State transitions:  6
============================================
```

**Key learning:** Lock-free programming with atomics, state machine design, rolling windows, background goroutine coordination, building resilience patterns used in microservices.

---

## Problem 8: Parallel Web Crawler with Deduplication (sync.Map + Semaphore + Context)

**Difficulty:** Hard

**Concepts:** sync.Map (or mutex+map), semaphore pattern (buffered channel), context cancellation, recursion with goroutines

**Problem:**

Build a parallel web crawler that starts from a seed URL, follows links up to a configurable depth, and avoids revisiting URLs.

**Requirements:**
1. Start from a seed URL. Simulate `FetchPage(url string) ([]string, error)` which returns a list of discovered links (generate random links from a fixed domain pool).
2. Maximum concurrency: at most N pages fetched in parallel (use a **semaphore** via buffered channel of size N).
3. Track visited URLs using `sync.Map` or a `mutex`-protected `map` — never fetch the same URL twice.
4. Respect a maximum crawl depth (e.g., 3 levels from seed).
5. Accept a `context.Context` — if cancelled, all in-flight and future fetches should abort.
6. Return a sitemap: `map[string][]string` (URL → list of links found on that page).
7. Print crawl stats: total pages fetched, duplicates skipped, errors, time elapsed.

```go
type Crawler struct {
    maxDepth    int
    concurrency int
    visited     sync.Map
    semaphore   chan struct{}
    results     chan CrawlResult
}

type CrawlResult struct {
    URL    string
    Links  []string
    Depth  int
    Err    error
}
```

**Expected Output:**
```
$ go run .
[Crawler] Starting from https://example.com (maxDepth=3, concurrency=5)

[depth=0] Fetching https://example.com ... 3 links found
[depth=1] Fetching https://example.com/about ... 2 links found
[depth=1] Fetching https://example.com/products ... 4 links found
[depth=1] Fetching https://example.com/blog ... 3 links found
[depth=2] Fetching https://example.com/products/item1 ... 1 link found
[depth=2] Skipping https://example.com/about (already visited)
[depth=2] Fetching https://example.com/blog/post1 ... 2 links found
[depth=2] Fetching https://example.com/products/item2 ... 0 links found
[depth=2] Fetching https://example.com/blog/post2 ... 1 link found
[depth=2] Fetching https://example.com/contact ... 0 links found
[depth=2] Fetching https://example.com/careers ... 1 link found
[depth=3] Skipping https://example.com/blog (already visited)
[depth=3] Fetching https://example.com/blog/post3 ... 0 links found
[depth=3] Skipping https://example.com/products (already visited)
[depth=3] Fetching https://example.com/team ... 0 links found

========== CRAWL REPORT ==========
Pages fetched:     12
Duplicates skipped: 3
Errors:             0
Max depth reached:  3
Time elapsed:       1.82s

Sitemap:
  https://example.com → [/about, /products, /blog]
  https://example.com/about → [/contact, /team]
  https://example.com/products → [/products/item1, /products/item2, /about, /blog]
  https://example.com/blog → [/blog/post1, /blog/post2, /products]
  ... (8 more entries)
===================================
```

**Key learning:** Semaphore pattern for bounded concurrency, sync.Map vs mutex+map trade-offs, recursive concurrency with depth tracking, context-driven cancellation of a tree of goroutines.

---

## Problem 9: Event-Driven Order Processing System (Multiple Patterns Combined)

**Difficulty:** Hard

**Concepts:** channels as event buses, pipeline, worker pool, context, sync.WaitGroup, graceful shutdown, state machines

**Problem:**

Build an event-driven order processing system that simulates an e-commerce backend.

**Requirements:**

1. **Order Submission:** A producer goroutine generates orders at random intervals (10–100ms). Each order has: `ID`, `CustomerID`, `Items[]`, `Total`.

2. **Validation Stage (Worker Pool, 3 workers):** Validate each order:
   - Reject if total > $10,000 (fraud check).
   - Reject if customer is in a blocklist.
   - Simulate validation delay (50–200ms).

3. **Payment Stage (Worker Pool, 2 workers):** Process payment:
   - Simulate payment gateway call with `context.WithTimeout` (500ms deadline).
   - 10% chance of payment failure.
   - On timeout/failure, move order to a **retry queue** (max 2 retries).

4. **Fulfillment Stage (Pipeline):** Accepted orders go through:
   - Inventory check → Packaging → Shipping
   - Each sub-stage is a pipeline stage with its own goroutine.

5. **Event Log:** Every state change (`submitted → validating → validated → payment_processing → paid → fulfilling → shipped` or `→ rejected/failed`) is emitted as an event on a shared event channel. A logger goroutine persists events (print to stdout with timestamp).

6. **Metrics Collector:** A goroutine collects metrics using `sync/atomic` counters:
   - Orders submitted, validated, rejected, paid, failed, shipped
   - Average processing time per stage

7. **Graceful Shutdown:** On `SIGINT` (or context cancellation):
   - Stop accepting new orders.
   - Drain all in-flight orders through the pipeline.
   - Print final metrics.
   - Exit cleanly.

```go
type Order struct {
    ID         string
    CustomerID string
    Items      []Item
    Total      float64
    Status     OrderStatus
    CreatedAt  time.Time
    Retries    int
}

type OrderEvent struct {
    OrderID   string
    FromState OrderStatus
    ToState   OrderStatus
    Timestamp time.Time
    Error     error
}
```

**Expected Output:**
```
$ go run .
[System] Starting order processing system...
[Validator-1] Ready
[Validator-2] Ready
[Validator-3] Ready
[Payment-1] Ready
[Payment-2] Ready
[Fulfillment] Pipeline ready (inventory → packaging → shipping)
[EventLog] Listening...
[Metrics] Collecting...

[Producer] Order ORD-001 submitted ($245.50, customer=C-12)
[Event] ORD-001: submitted → validating
[Validator-2] ORD-001: checking... OK
[Event] ORD-001: validating → validated
[Event] ORD-001: validated → payment_processing
[Payment-1] ORD-001: charging $245.50... OK
[Event] ORD-001: payment_processing → paid
[Event] ORD-001: paid → fulfilling
[Fulfillment] ORD-001: inventory ✓ → packaging ✓ → shipping ✓
[Event] ORD-001: fulfilling → shipped

[Producer] Order ORD-002 submitted ($15200.00, customer=C-05)
[Event] ORD-002: submitted → validating
[Validator-1] ORD-002: REJECTED (fraud: total > $10,000)
[Event] ORD-002: validating → rejected

[Producer] Order ORD-003 submitted ($89.99, customer=C-08)
[Event] ORD-003: submitted → validating → validated → payment_processing
[Payment-2] ORD-003: charging $89.99... FAILED (gateway error)
[Event] ORD-003: payment_processing → payment_retry (attempt 1/2)
[Payment-1] ORD-003: retry charging $89.99... OK
[Event] ORD-003: payment_retry → paid → fulfilling → shipped

... (more orders flowing through)

^C
[System] SIGINT received — shutting down gracefully
[Producer] Stopped (no new orders)
[System] Draining in-flight orders...
[System] All orders drained

========== FINAL METRICS ==========
Orders submitted:   47
Orders validated:   42
Orders rejected:     5 (3 fraud, 2 blocklist)
Payments OK:        38
Payments failed:     4 (2 recovered via retry)
Orders shipped:     38
Orders failed:       2 (payment — max retries exceeded)

Avg time per stage:
  Validation:    127ms
  Payment:       203ms
  Fulfillment:   340ms
  End-to-end:    670ms
====================================
```

**Key learning:** Combining all major concurrency patterns into a realistic system: worker pools, pipelines, event buses, context-driven timeouts, graceful shutdown, atomic metrics — mirrors real microservice architectures.

---

## Problem 10: Rate Limiter with Multiple Algorithms (Token Bucket + Sliding Window)

**Difficulty:** Hard

**Concepts:** sync.Mutex, time.Ticker, channels, goroutines, interface design

**Problem:**

Implement two rate-limiting algorithms behind a common interface, and build a middleware that applies them.

**Requirements:**

1. **Interface:**
```go
type RateLimiter interface {
    Allow() bool                   // non-blocking: can I proceed?
    Wait(ctx context.Context) error // blocking: wait until allowed or ctx cancelled
}
```

2. **Token Bucket Limiter:**
   - Configurable: rate (tokens/sec), burst (max tokens).
   - A background goroutine refills tokens at the specified rate.
   - `Allow()` takes a token if available (atomic or mutex).
   - `Wait()` blocks until a token is available or context is done.

3. **Sliding Window Limiter:**
   - Configurable: max requests per window duration.
   - Track timestamps of recent requests.
   - Clean up old entries in a background goroutine.
   - `Allow()` checks if current window count < max.

4. **Per-Client Rate Limiting:**
   - `RateLimiterGroup` maps client IDs to individual limiters.
   - Lazy initialization with `sync.Once` or mutex.
   - Auto-expire idle client limiters after a configurable duration.

5. **Demonstrate:** Simulate 20 clients, each making 100 requests. Measure and print per-client accept/reject rates.

```go
type TokenBucket struct {
    mu       sync.Mutex
    tokens   float64
    maxTokens float64
    rate     float64
    lastTime time.Time
    notify   chan struct{}
}

type SlidingWindow struct {
    mu       sync.Mutex
    requests []time.Time
    max      int
    window   time.Duration
}

type RateLimiterGroup struct {
    mu       sync.RWMutex
    limiters map[string]*clientLimiter
    factory  func() RateLimiter
}
```

**Expected Output:**
```
$ go run .
===== Token Bucket Limiter (10 tokens/sec, burst=20) =====
[Client-01] 100 requests → 62 accepted, 38 rejected (62.0% pass)
[Client-02] 100 requests → 59 accepted, 41 rejected (59.0% pass)
[Client-03] 100 requests → 61 accepted, 39 rejected (61.0% pass)
[Client-04] 100 requests → 58 accepted, 42 rejected (58.0% pass)
... (16 more clients)
[Client-20] 100 requests → 60 accepted, 40 rejected (60.0% pass)
Total: 1193/2000 accepted (59.7%)

===== Sliding Window Limiter (50 req/5s window) =====
[Client-01] 100 requests → 50 accepted, 50 rejected (50.0% pass)
[Client-02] 100 requests → 50 accepted, 50 rejected (50.0% pass)
... (18 more clients)
[Client-20] 100 requests → 50 accepted, 50 rejected (50.0% pass)
Total: 1000/2000 accepted (50.0%)

===== Wait() mode (Token Bucket, 3s timeout) =====
[Client-05] Wait() OK — proceeded after 120ms
[Client-12] Wait() OK — proceeded after 340ms
[Client-08] Wait() FAILED — context deadline exceeded

===== Per-Client Limiter Group =====
[LimiterGroup] Active limiters: 20
[LimiterGroup] Idle cleanup: removed 5 expired client limiters
[LimiterGroup] Active limiters: 15
```

**Key learning:** Implementing concurrency-safe algorithms, condition-variable-like signaling with channels, per-entity state management, background cleanup goroutines, `sync.Once` for lazy init.

---

## Problem 11: Concurrent Map-Reduce Engine (Generics + Goroutines)

**Difficulty:** Hard

**Concepts:** generics, fan-out/fan-in, channels, sync.WaitGroup, functional patterns

**Problem:**

Build a generic concurrent map-reduce framework.

**Requirements:**

1. Generic API:
```go
func MapReduce[T any, M any, R any](
    input []T,
    mapFn func(T) M,
    reduceFn func([]M) R,
    concurrency int,
) R
```

2. The `map` phase fans out to `concurrency` workers, each processing a chunk of the input.
3. The `reduce` phase collects all mapped results and applies `reduceFn`.
4. Support context cancellation — if cancelled, map workers stop early.
5. If any map worker panics, recover it and report the error without crashing the whole pipeline.

**Demonstrate with 3 use cases:**
- **Word count:** Input = []string (files), Map = count words per file, Reduce = merge counts.
- **Sum of squares:** Input = []int, Map = square, Reduce = sum.
- **Log analysis:** Input = []LogEntry, Map = extract error counts per service, Reduce = merge into totals.

**Expected Output:**
```
$ go run .
===== Use Case 1: Word Count =====
[MapReduce] Input: 10 files, Concurrency: 4 workers
[Worker-1] Processing chunk 1 (3 files)...
[Worker-2] Processing chunk 2 (3 files)...
[Worker-3] Processing chunk 3 (2 files)...
[Worker-4] Processing chunk 4 (2 files)...
[Reduce] Merging 4 partial word counts...
Top 5 words:
  "the"    → 1247
  "error"  → 892
  "func"   → 634
  "return" → 521
  "nil"    → 498
Duration: 45ms (sequential estimate: 156ms)

===== Use Case 2: Sum of Squares =====
[MapReduce] Input: 1,000,000 ints, Concurrency: 8 workers
[Worker-1..8] Processing 125,000 elements each...
[Reduce] Summing 8 partial results...
Result: 333,333,833,333,500,000
Duration: 12ms (sequential estimate: 87ms)

===== Use Case 3: Log Analysis =====
[MapReduce] Input: 50,000 log entries, Concurrency: 4 workers
[Worker-3] ⚠ PANIC recovered: "index out of range" — skipped 1 entry
[Reduce] Merging error counts...
Errors per service:
  auth-service:    234
  order-service:   187
  payment-service: 312
  user-service:    156
Total errors: 889
Duration: 67ms
[Note] 1 worker panic recovered without crashing pipeline
```

**Key learning:** Generics with concurrency, panic recovery in goroutines, chunked parallel processing, building reusable concurrent abstractions.

---

## Problem 12: Database Connection Pool (Semaphore + sync.Cond + Context)

**Difficulty:** Hard

**Concepts:** sync.Cond, semaphore, resource pooling, context deadline, goroutine lifecycle

**Problem:**

Build a database connection pool from scratch.

**Requirements:**

1. `NewPool(maxConns int, factory func() (Conn, error)) *Pool`
2. `Acquire(ctx context.Context) (Conn, error)` — get a connection. If none available:
   - If pool size < max, create a new one.
   - Otherwise, **wait** until one is released or `ctx` is cancelled/timed out.
3. `Release(conn Conn)` — return connection to pool. Wake up one waiting acquirer.
4. `Close()` — close all connections, reject new Acquire calls.
5. Health check: a background goroutine pings idle connections every 30 seconds, removing dead ones.
6. Track metrics: active connections, idle connections, total wait time, timeouts.

```go
type Pool struct {
    mu          sync.Mutex
    cond        *sync.Cond
    idle        []*poolConn
    active      int
    maxSize     int
    closed      bool
    factory     func() (Conn, error)
}

type Conn interface {
    Exec(query string) error
    Ping() error
    Close() error
}
```

**Demonstrate:** 50 goroutines competing for 10 connections, each running random queries (50–200ms). Print pool stats every second.

**Expected Output:**
```
$ go run .
[Pool] Initialized (maxConns=10)
[HealthCheck] Started (interval=30s)
[Benchmark] Spawning 50 goroutines competing for 10 connections...

[t=0.0s] Pool stats: active=0  idle=0  waiting=0
[Goroutine-01] Acquired conn-1 (new)       pool: active=1  idle=0
[Goroutine-02] Acquired conn-2 (new)       pool: active=2  idle=0
[Goroutine-03] Acquired conn-3 (new)       pool: active=3  idle=0
... (7 more new connections created)
[Goroutine-11] Waiting for connection...   pool: active=10 idle=0  waiting=1
[Goroutine-12] Waiting for connection...   pool: active=10 idle=0  waiting=2
...

[t=1.0s] Pool stats: active=10 idle=0  waiting=12  avg_wait=234ms
[Goroutine-01] Released conn-1             pool: active=9  idle=0  waiting=11
[Goroutine-11] Acquired conn-1 (reused)    pool: active=10 idle=0  waiting=10
...

[t=2.0s] Pool stats: active=10 idle=0  waiting=4   avg_wait=312ms
[t=3.0s] Pool stats: active=8  idle=2  waiting=0   avg_wait=287ms
[t=4.0s] Pool stats: active=3  idle=7  waiting=0

[Goroutine-47] Acquire TIMEOUT (context deadline exceeded after 500ms)

[t=5.0s] Pool stats: active=0  idle=10 waiting=0

========== POOL REPORT ==========
Total Acquire calls:   250
Successful:            248
Timeouts:                2
Connections created:    10
Avg wait time:         287ms
Max wait time:         498ms
==================================

[Pool] Closed — 10 connections released
[HealthCheck] Stopped
```

**Key learning:** sync.Cond for efficient wait/notify (vs busy-waiting or channel), resource pool design used in real DB drivers (`database/sql`), context integration with blocking operations.

---

## Problem 13: Leader Election Simulator (Channels + Goroutines + Timeouts)

**Difficulty:** Hard

**Concepts:** distributed systems simulation, select + timeout, channels for message passing, state machines

**Problem:**

Simulate a simplified Raft-like leader election among N nodes communicating via channels.

**Requirements:**

1. Each node is a goroutine with states: `Follower`, `Candidate`, `Leader`.
2. Nodes communicate via directional channels (simulating network).
3. **Follower behavior:** Wait for heartbeats. If no heartbeat received within a random timeout (150–300ms), become Candidate.
4. **Candidate behavior:** Increment term, vote for self, send `RequestVote` RPCs to all other nodes. If majority votes received before timeout, become Leader. If heartbeat from valid leader received, revert to Follower.
5. **Leader behavior:** Send periodic heartbeats (every 50ms) to all nodes. Step down if higher term is seen.
6. Simulate network partitions: randomly block channels between certain nodes for a period, then restore. The system should elect a new leader during partition and resolve when healed.
7. Print a timeline showing state transitions and elections.

```go
type Node struct {
    id         int
    state      NodeState
    term       int
    votedFor   int
    peers      map[int]*peerConn
    inbox      chan Message
}

type Message struct {
    Type    MsgType  // RequestVote, VoteGranted, Heartbeat
    From    int
    Term    int
}
```

**Expected Output:**
```
$ go run .
[Cluster] Starting 5-node Raft simulation...

[t=0.000s] Node-0: Follower (term=0)
[t=0.000s] Node-1: Follower (term=0)
[t=0.000s] Node-2: Follower (term=0)
[t=0.000s] Node-3: Follower (term=0)
[t=0.000s] Node-4: Follower (term=0)

[t=0.217s] Node-2: election timeout — becoming Candidate (term=1)
[t=0.217s] Node-2: requesting votes from [0,1,3,4]
[t=0.218s] Node-0: voted for Node-2 (term=1)
[t=0.219s] Node-1: voted for Node-2 (term=1)
[t=0.219s] Node-4: voted for Node-2 (term=1)
[t=0.219s] Node-2: received 4/5 votes — becoming LEADER (term=1) ★
[t=0.270s] Node-2: heartbeat → [0,1,3,4]
[t=0.320s] Node-2: heartbeat → [0,1,3,4]
...

[t=3.000s] ⚡ Network partition: [Node-0, Node-1] isolated from [Node-2, Node-3, Node-4]
[t=3.000s] Node-2: heartbeat to Node-0 BLOCKED
[t=3.000s] Node-2: heartbeat to Node-1 BLOCKED
[t=3.250s] Node-0: no heartbeat — becoming Candidate (term=2)
[t=3.251s] Node-1: voted for Node-0 (term=2)
[t=3.251s] Node-0: received 2/5 votes — NOT enough (need 3)
[t=3.490s] Node-1: election timeout — becoming Candidate (term=3)
[t=3.491s] Node-0: voted for Node-1 (term=3)
[t=3.491s] Node-1: received 2/5 votes — NOT enough

[t=6.000s] ⚡ Network partition HEALED
[t=6.050s] Node-2: heartbeat → all nodes
[t=6.051s] Node-0: saw term=1 < my term=3 — rejecting heartbeat
[t=6.051s] Node-2: saw higher term=3 — stepping down to Follower
[t=6.280s] Node-4: election timeout — becoming Candidate (term=4)
[t=6.282s] Node-4: received 3/5 votes — becoming LEADER (term=4) ★

========== ELECTION TIMELINE ==========
Term 1: Node-2 elected leader (0.219s)
Term 4: Node-4 elected leader (6.282s)
Total elections:  5 (2 successful, 3 failed—no quorum)
Partition events: 1 (duration: 3.0s)
========================================
```

**Key learning:** Using goroutines and channels to model distributed systems, randomized timeouts for contention avoidance, select for handling multiple event sources, simulating real consensus protocols.

---

## Problem 14: Graceful HTTP Server with Dependency Lifecycle (Real-World Pattern)

**Difficulty:** Moderate–Hard

**Concepts:** context, os/signal, sync.WaitGroup, errgroup, goroutine lifecycle management

**Problem:**

Build a simulated HTTP-like server that manages multiple background services with proper startup and shutdown ordering.

**Requirements:**

1. **Services to manage** (each is a long-running goroutine):
   - Cache warmer (runs every 5 seconds)
   - Metrics exporter (runs every 10 seconds)
   - Background job processor (worker pool of 3)
   - Health checker (pings dependencies every 3 seconds)

2. **Startup order:** Cache warmer → Health checker → Job processor → Metrics exporter. Each service signals when ready (use channels).

3. **The "server"** accepts simulated requests (generate them in a goroutine) and processes them using all services.

4. **Shutdown sequence** on `SIGINT`:
   - Stop accepting new requests.
   - Wait for in-flight requests (max 30 seconds).
   - Shut down services in reverse startup order.
   - Each service has a `Shutdown(ctx context.Context) error` with its own timeout.
   - If a service doesn't stop in time, force-kill and log a warning.
   - Print final status of each service.

```go
type Service interface {
    Name() string
    Start(ctx context.Context) error  // blocks until stopped
    Ready() <-chan struct{}           // signals when ready to serve
    Shutdown(ctx context.Context) error
}

type Server struct {
    services []Service
    wg       sync.WaitGroup
}
```

**Expected Output:**
```
$ go run .
[Server] Starting services...

[CacheWarmer] Starting...
[CacheWarmer] Warming cache: 1247 entries loaded
[CacheWarmer] ✓ Ready
[HealthChecker] Starting...
[HealthChecker] Pinging dependencies: DB=OK, Redis=OK, API=OK
[HealthChecker] ✓ Ready
[JobProcessor] Starting 3 workers...
[JobProcessor-1] ✓ Ready
[JobProcessor-2] ✓ Ready
[JobProcessor-3] ✓ Ready
[JobProcessor] ✓ Ready
[MetricsExporter] Starting...
[MetricsExporter] ✓ Ready

[Server] All services ready — accepting requests

[Server] Request #1 → 200 OK (12ms)
[CacheWarmer] Refresh: 23 entries updated
[Server] Request #2 → 200 OK (8ms)
[Server] Request #3 → 200 OK (15ms)
[HealthChecker] Ping: DB=OK, Redis=OK, API=OK
[JobProcessor-2] Processed background job J-042
[MetricsExporter] Exported: reqs=3, latency_p99=15ms, jobs=1
[Server] Request #4 → 200 OK (9ms)
...

^C
[Server] SIGINT received — initiating graceful shutdown
[Server] Stopped accepting new requests
[Server] Waiting for 2 in-flight requests (max 30s)...
[Server] In-flight requests drained

[Server] Shutting down services (reverse order)...
[MetricsExporter] Shutdown... flushing final metrics... ✓ (120ms)
[JobProcessor] Shutdown... waiting for 1 active job... ✓ (2.1s)
[HealthChecker] Shutdown... ✓ (5ms)
[CacheWarmer] Shutdown... ✓ (3ms)

========== SERVER STATUS ==========
Uptime:             47s
Requests served:    123
Services shutdown:  4/4 OK
  MetricsExporter:  stopped (120ms)
  JobProcessor:     stopped (2.1s)
  HealthChecker:    stopped (5ms)
  CacheWarmer:      stopped (3ms)
====================================
```

**Key learning:** Production-grade goroutine lifecycle management, ordered startup/shutdown, context propagation through a service tree, matching patterns used in Kubernetes controllers and production Go servers.

---

## Problem 15: Concurrent Merge Sort with Goroutine Budget (Bounded Parallelism)

**Difficulty:** Moderate

**Concepts:** goroutines, channels, semaphore, recursion, bounded concurrency

**Problem:**

Implement merge sort that runs sub-sorts concurrently, but limits the total number of concurrent goroutines.

**Requirements:**

1. Sort a large slice (1 million elements) using merge sort.
2. Each recursive split **may** spawn a goroutine for the left and right halves.
3. Use a **goroutine budget** (semaphore): at most N goroutines can run in parallel. If the budget is exhausted, fall back to sequential sort for that subtree.
4. Compare performance (with `time.Now()`) for budgets of: 1 (sequential), 4, 8, 16, runtime.NumCPU(), 100, 1000.
5. Show that unbounded goroutines (no budget) is actually slower due to overhead.

```go
func ConcurrentMergeSort(data []int, budget chan struct{}) {
    if len(data) <= 1 {
        return
    }
    mid := len(data) / 2

    select {
    case budget <- struct{}{}:
        // Got budget — sort left half concurrently
        var wg sync.WaitGroup
        wg.Add(1)
        go func() {
            defer wg.Done()
            ConcurrentMergeSort(data[:mid], budget)
            <-budget
        }()
        ConcurrentMergeSort(data[mid:], budget)
        wg.Wait()
    default:
        // No budget — sort sequentially
        ConcurrentMergeSort(data[:mid], budget)
        ConcurrentMergeSort(data[mid:], budget)
    }

    merge(data, mid)
}
```

**Expected Output:**
```
$ go run .
[MergeSort] Sorting 1,000,000 elements with varying goroutine budgets...

 Budget=1  (sequential):       312ms  goroutines_spawned=0
 Budget=4:                     142ms  goroutines_spawned=4       (2.2x speedup)
 Budget=8:                      89ms  goroutines_spawned=8       (3.5x speedup)
 Budget=16:                     71ms  goroutines_spawned=16      (4.4x speedup)
 Budget=NumCPU(10):             78ms  goroutines_spawned=10      (4.0x speedup)
 Budget=100:                    68ms  goroutines_spawned=100     (4.6x speedup)
 Budget=1000:                   74ms  goroutines_spawned=1000    (4.2x speedup)
 Budget=UNLIMITED:             187ms  goroutines_spawned=524287  (1.7x speedup) ← SLOWER!

[MergeSort] Verifying sort correctness... ✓ all sorted

========== ANALYSIS ==========
Best budget:        100 (68ms, 4.6x speedup)
Diminishing returns after ~NumCPU goroutines
Unlimited spawns 524K goroutines → scheduling overhead dominates
Sweet spot: budget = 2*NumCPU to 4*NumCPU
===============================
```

**Key learning:** Understanding goroutine overhead, semaphore pattern for bounding concurrency, adaptive parallelism, practical performance tuning.

---

## Concept Coverage Matrix

| Problem | goroutines | channels | WaitGroup | Mutex/RWMutex | atomic | Context | select | Pipeline | Fan-out/in | Semaphore | sync.Cond | errgroup |
|---------|-----------|----------|-----------|---------------|--------|---------|--------|----------|------------|-----------|-----------|----------|
| 1. File Downloader | ✓ | ✓ | ✓ | | | | | | | | | |
| 2. Pipeline Processor | ✓ | ✓ | ✓ | | | | | ✓ | ✓ | | | |
| 3. Singleflight | ✓ | | ✓ | ✓ | | | | | | | | |
| 4. TTL Cache | ✓ | ✓ | | ✓ | | | ✓ | | | | | |
| 5. Task Scheduler | ✓ | ✓ | | | | ✓ | ✓ | | | | | ✓ |
| 6. Pub/Sub Broker | ✓ | ✓ | | ✓ | ✓ | | ✓ | | | | | |
| 7. Circuit Breaker | ✓ | | | | ✓ | | | | | | | |
| 8. Web Crawler | ✓ | ✓ | ✓ | | | ✓ | | | | ✓ | | |
| 9. Order System | ✓ | ✓ | ✓ | | ✓ | ✓ | ✓ | ✓ | | | | |
| 10. Rate Limiter | ✓ | ✓ | | ✓ | | ✓ | ✓ | | | | | |
| 11. Map-Reduce | ✓ | ✓ | ✓ | | | ✓ | | | ✓ | | | |
| 12. Conn Pool | ✓ | | | ✓ | | ✓ | | | | ✓ | ✓ | |
| 13. Leader Election | ✓ | ✓ | | | | | ✓ | | | | | |
| 14. Graceful Server | ✓ | ✓ | ✓ | | | ✓ | ✓ | | | | | ✓ |
| 15. Merge Sort | ✓ | | ✓ | | | | ✓ | | | ✓ | | |

---

## Recommended Order

1. **Problem 1** — Warm up: worker pool basics
2. **Problem 2** — Pipelines and fan-out/fan-in
3. **Problem 15** — Bounded parallelism intuition
4. **Problem 3** — Mutex + coordination primitives
5. **Problem 4** — RWMutex + background goroutines
6. **Problem 10** — Rate limiting algorithms
7. **Problem 7** — Atomics + state machines
8. **Problem 6** — Pub/Sub + goroutine leak prevention
9. **Problem 8** — Semaphore + sync.Map + context
10. **Problem 5** — Context-heavy scheduling
11. **Problem 12** — sync.Cond + resource pooling
12. **Problem 11** — Generics + map-reduce
13. **Problem 14** — Production server lifecycle
14. **Problem 9** — Full system integration
15. **Problem 13** — Distributed systems simulation

---

## Tips

- **Always** use `go vet` and `go run -race .` (race detector) on every solution.
- If a test hangs, you have a deadlock — check channel/mutex ordering.
- Draw the goroutine topology (who sends to whom) on paper before coding.
- Prefer `context.Context` over raw `done` channels for cancellation in real code.
- Every goroutine you spawn must have a clear termination condition — ask "what makes this goroutine stop?"
