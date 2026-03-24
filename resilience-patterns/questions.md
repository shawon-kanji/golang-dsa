# Go Resilience & Observability Exercises (Hard)

Build the patterns that keep production systems alive: retries, bulkheads, backpressure, tracing, and monitoring. These are the techniques behind Netflix Hystrix, Envoy, and OpenTelemetry.

---

## Problem 1: Retry Engine with Pluggable Strategies

**Difficulty:** Moderate–Hard

**Concepts:** exponential backoff, jitter, retry budgets, context integration, strategy pattern

**Problem:**

Build a composable retry engine that supports multiple backoff strategies and retry conditions.

**Requirements:**
1. **Retry function**:
```go
func Retry(ctx context.Context, op func(ctx context.Context) error, opts ...Option) error
```
2. **Backoff strategies** (pluggable):
   - **Constant**: fixed delay between retries.
   - **Exponential**: `baseDelay * 2^attempt` with configurable max delay.
   - **Exponential + jitter**: add random jitter (±25%) to avoid thundering herd.
   - **Linear**: `baseDelay * attempt`.
3. **Retry conditions**: retry only on specific error types. Don't retry on `context.Canceled` or permanent errors.
4. **Retry budget**: limit total retries per time window across all callers (e.g., max 20% of requests can be retries).
5. **Hooks**: `OnRetry(attempt int, err error, delay time.Duration)` callback for logging.
6. Demonstrate with a flaky function that fails 70% of the time, and show each strategy's behavior.

```go
type RetryConfig struct {
    MaxAttempts   int
    Backoff       BackoffStrategy
    RetryIf       func(error) bool
    OnRetry       func(attempt int, err error, delay time.Duration)
    Budget        *RetryBudget
}

type BackoffStrategy interface {
    Delay(attempt int) time.Duration
}

type RetryBudget struct {
    mu       sync.Mutex
    window   time.Duration
    maxRatio float64 // max retry/total ratio
    history  []requestRecord
}
```

**Expected Output:**
```
$ go run .

===== Exponential Backoff =====
[Attempt 1] FAILED: connection refused (retry in 100ms)
[Attempt 2] FAILED: connection refused (retry in 200ms)
[Attempt 3] FAILED: connection refused (retry in 400ms)
[Attempt 4] SUCCESS (total time: 745ms, 3 retries)

===== Exponential + Jitter =====
[Attempt 1] FAILED: connection refused (retry in 87ms)
[Attempt 2] FAILED: connection refused (retry in 223ms)
[Attempt 3] SUCCESS (total time: 341ms, 2 retries)

===== Constant Backoff =====
[Attempt 1] FAILED (retry in 500ms)
[Attempt 2] FAILED (retry in 500ms)
[Attempt 3] FAILED (retry in 500ms)
[Attempt 4] SUCCESS (total time: 1.52s, 3 retries)

===== Retry Budget =====
[Budget] Window: 10s, Max retry ratio: 20%
[t=0s] 10 requests: 7 failed, 7 retried (70% retry ratio > 20%)
[Budget] THROTTLED: retry rejected (budget exhausted: 20% cap hit)
[t=0s] 3 retries rejected by budget

===== Permanent Error (no retry) =====
[Attempt 1] FAILED: 404 not found (permanent error — not retrying)
Total time: 2ms

===== Context Cancellation =====
[Attempt 1] FAILED (retry in 100ms)
[Attempt 2] CONTEXT CANCELLED — aborting retries
Total time: 115ms
```

**Key learning:** Retry strategies used in every production system, jitter for thundering herd prevention (AWS best practice), retry budgets to prevent retry storms, composable option pattern in Go.

---

## Problem 2: Bulkhead Pattern (Resource Isolation)

**Difficulty:** Hard

**Concepts:** semaphore, resource partitioning, goroutine pools, queue overflow, adaptive concurrency

**Problem:**

Build a bulkhead that isolates different types of work into separate resource pools, preventing one slow dependency from consuming all resources.

**Requirements:**
1. **Named bulkheads**: create isolated execution contexts — each with its own goroutine limit and queue size.
   ```go
   bh := NewBulkhead("payment-service", MaxConcurrency(10), QueueSize(50))
   ```
2. `Execute(ctx, fn) error` — run `fn` within the bulkhead. If at max concurrency, queue. If queue full, reject immediately with `ErrBulkheadFull`.
3. **Metrics per bulkhead**: active count, queued count, rejected count, avg execution time.
4. **Adaptive concurrency** (stretch): dynamically adjust max concurrency based on observed latency (increase if latency is low, decrease if latency spikes — TCP congestion control style).
5. Demonstrate: 3 bulkheads (`user-api`, `payment-api`, `inventory-api`). Saturate `payment-api` → show that the other two are unaffected.

```go
type Bulkhead struct {
    name         string
    maxConc      int
    semaphore    chan struct{}
    queue        chan func()
    queueSize    int
    metrics      *BulkheadMetrics
    adaptive     bool
}

type BulkheadMetrics struct {
    Active    atomic.Int64
    Queued    atomic.Int64
    Completed atomic.Int64
    Rejected  atomic.Int64
    TotalTime atomic.Int64 // nanoseconds
}
```

**Expected Output:**
```
$ go run .
[Bulkhead] "user-api"      max=20 queue=100
[Bulkhead] "payment-api"   max=10 queue=50
[Bulkhead] "inventory-api" max=15 queue=75

--- Normal operation ---
[user-api]      active=5/20   queued=0   rejected=0  avg_latency=25ms
[payment-api]   active=3/10   queued=0   rejected=0  avg_latency=45ms
[inventory-api] active=4/15   queued=0   rejected=0  avg_latency=30ms

--- payment-api becomes slow (2s per request) ---
[payment-api]   active=10/10  queued=12  rejected=0  avg_latency=2100ms
[payment-api]   active=10/10  queued=50  rejected=0  avg_latency=2300ms ← queue filling
[payment-api]   active=10/10  queued=50  rejected=23 avg_latency=2400ms ← REJECTING

[user-api]      active=6/20   queued=0   rejected=0  avg_latency=26ms  ← UNAFFECTED ✓
[inventory-api] active=5/15   queued=0   rejected=0  avg_latency=31ms  ← UNAFFECTED ✓

Without bulkhead (shared pool of 45):
  All 45 goroutines stuck on payment-api
  user-api: STARVED (avg_latency=2300ms) ✗
  inventory-api: STARVED (avg_latency=2200ms) ✗

With bulkhead:
  payment-api: degraded (rejecting overflow)
  user-api: healthy ✓
  inventory-api: healthy ✓

--- Adaptive concurrency ---
[payment-api] Latency spike detected (2100ms > target 100ms)
[payment-api] Reducing concurrency: 10 → 5
[payment-api] Latency improving: 1200ms
[payment-api] Reducing concurrency: 5 → 3
[payment-api] Latency stabilized: 150ms
[payment-api] Gradually increasing: 3 → 4 → 5 → 6 (latency OK)

========== BULKHEAD REPORT ==========
            Completed  Rejected  Avg Latency
user-api        847       0        26ms
payment-api     312      23       1850ms
inventory-api   623       0        31ms
======================================
```

**Key learning:** Bulkhead pattern (Netflix Hystrix), resource isolation, queue-based admission control, adaptive concurrency limits (like Envoy's concurrency limiter), preventing cascading failures.

---

## Problem 3: Backpressure Controller

**Difficulty:** Hard

**Concepts:** producer-consumer balance, dynamic rate adjustment, load shedding, feedback loops

**Problem:**

Build a backpressure-aware pipeline where slow consumers signal producers to slow down.

**Requirements:**
1. **Pipeline**: Producer → Buffer → Worker Pool → Output.
2. **Backpressure signals**: when the buffer is >80% full, signal the producer to slow down. When >95%, drop new messages (load shedding).
3. **Adaptive producer rate**: producer adjusts its emission rate based on consumer throughput feedback.
4. **Consumer scaling**: dynamically add/remove workers based on queue depth (scale up when buffer grows, scale down when idle).
5. **Metrics**: track throughput (in/out), buffer utilization, dropped messages, consumer count over time.
6. Demonstrate: producer bursts at 10,000 msg/s, consumers handle 5,000 msg/s → show how backpressure stabilizes the system.

```go
type Pipeline struct {
    buffer       chan Message
    bufferSize   int
    producer     *Producer
    pool         *WorkerPool
    metrics      *PipelineMetrics
    scalingRules ScalingConfig
}

type Producer struct {
    rate      atomic.Int64 // messages per second
    baseRate  int64
    dropped   atomic.Int64
}

type WorkerPool struct {
    mu       sync.Mutex
    workers  int
    minWorkers int
    maxWorkers int
    active   atomic.Int64
}

type PipelineMetrics struct {
    Produced  atomic.Int64
    Consumed  atomic.Int64
    Dropped   atomic.Int64
    BufferPct atomic.Int64 // 0-100
}
```

**Expected Output:**
```
$ go run .
[Pipeline] Buffer: 1000 | Workers: 2-10 (min-max) | Producer: 10,000 msg/s

[t=0s]  Producer: 10000/s | Buffer:  0% | Workers: 2 | Consumed: 0/s | Dropped: 0
[t=1s]  Producer: 10000/s | Buffer: 52% | Workers: 2 | Consumed: 5000/s | Dropped: 0
  ↑ Buffer growing (in > out)

[t=2s]  Producer: 10000/s | Buffer: 81% | Workers: 2 | Consumed: 5100/s | Dropped: 0
  ⚠ BACKPRESSURE: buffer >80% → scaling up workers
  [Scale] Workers: 2 → 4

[t=3s]  Producer: 10000/s | Buffer: 72% | Workers: 4 | Consumed: 8200/s | Dropped: 0
  Buffer stabilizing...

[t=4s]  Producer: 10000/s | Buffer: 78% | Workers: 4 | Consumed: 9100/s | Dropped: 0
  ⚠ Still near threshold → scaling up
  [Scale] Workers: 4 → 6

[t=5s]  Producer: 10000/s | Buffer: 45% | Workers: 6 | Consumed: 10200/s | Dropped: 0
  ✓ Buffer draining — system balanced

--- Extreme burst: producer 30,000 msg/s ---
[t=10s] Producer: 30000/s | Buffer: 96% | Workers: 10 | Consumed: 10500/s | Dropped: 0
  🛑 LOAD SHEDDING: buffer >95% → dropping messages + signaling producer
  [Producer] Rate reduced: 30000 → 12000 msg/s

[t=11s] Producer: 12000/s | Buffer: 88% | Workers: 10 | Consumed: 10400/s | Dropped: 1450
  Buffer still high → producer slowing further
  [Producer] Rate reduced: 12000 → 10500 msg/s

[t=12s] Producer: 10500/s | Buffer: 62% | Workers: 10 | Consumed: 10300/s | Dropped: 0
  ✓ System recovered

--- Low load: scaling down ---
[t=20s] Producer:  2000/s | Buffer:  3% | Workers: 10 | Consumed: 2000/s | Dropped: 0
  Workers idle → scaling down
  [Scale] Workers: 10 → 4 → 2

========== PIPELINE REPORT ==========
Total produced:   152,000
Total consumed:   149,550
Total dropped:      2,450 (1.6%)
Max buffer util:   96%
Worker range:      2–10
Effective throughput: ~10,200 msg/s (steady state)
======================================
```

**Key learning:** Backpressure is critical for production stability. How Kafka, TCP, and reactive streams handle backpressure. Adaptive rate control. Load shedding (controlled degradation vs cascading failure).

---

## Problem 4: Distributed Tracing (OpenTelemetry-lite)

**Difficulty:** Hard

**Concepts:** trace context propagation, span trees, W3C trace context, sampling, reporting

**Problem:**

Build a simplified distributed tracing system that tracks requests across simulated microservices.

**Requirements:**
1. **Trace context**: each trace has a unique `TraceID`. Each operation within a trace is a `Span` with: `SpanID`, `ParentSpanID`, `OperationName`, `StartTime`, `Duration`, `Status`, `Tags`.
2. **Context propagation**: `SpanFromContext(ctx)` extracts the current span. When calling another service, inject the trace context into the "request" (pass via context).
3. **Automatic span creation**: `tracer.StartSpan(ctx, "operation-name")` creates a child span linked to the parent.
4. **Simulated microservices**: API Gateway → User Service → DB Query, and API Gateway → Order Service → Payment Service. Each is a function with tracing.
5. **Sampling**: only record 10% of traces (head-based sampling). Use deterministic sampling (hash of TraceID).
6. **Reporter**: collect all spans and print a trace tree after the request completes.
7. Support synchronous and asynchronous spans (goroutines that continue after the parent span ends).

```go
type Tracer struct {
    mu       sync.Mutex
    traces   map[string]*Trace
    sampler  Sampler
    reporter Reporter
}

type Trace struct {
    TraceID string
    Spans   []*Span
}

type Span struct {
    TraceID     string
    SpanID      string
    ParentID    string
    Operation   string
    StartTime   time.Time
    Duration    time.Duration
    Status      SpanStatus // OK, Error
    Tags        map[string]string
    Events      []SpanEvent // timestamped log lines within span
}

type Sampler interface {
    ShouldSample(traceID string) bool
}
```

**Expected Output:**
```
$ go run .

===== Trace: abc-123-def (sampled) =====

[12ms] api-gateway: GET /api/users/42
├─ [8ms] user-service: GetUser(42)
│  ├─ [1ms] cache: Get("user:42") → MISS
│  ├─ [5ms] postgres: SELECT * FROM users WHERE id=42
│  │  └─ tags: {db.type=postgres, db.statement="SELECT...", rows=1}
│  └─ [1ms] cache: Set("user:42") → OK
└─ [2ms] audit-service: LogAccess(user=42) [async]
   └─ tags: {async=true}

Trace Summary:
  TraceID:     abc-123-def
  Total spans: 6
  Duration:    12ms
  Services:    [api-gateway, user-service, postgres, cache, audit-service]
  Status:      OK

===== Trace: ghi-456-jkl (sampled, with error) =====

[340ms] api-gateway: POST /api/orders
├─ [15ms] order-service: CreateOrder
│  ├─ [8ms] inventory-service: CheckStock(item-99)
│  │  └─ tags: {available=true}
│  └─ [310ms] payment-service: Charge($49.99)
│     └─ STATUS: ERROR
│     └─ events: [
│          {t+200ms} "connecting to payment gateway..."
│          {t+300ms} "timeout waiting for gateway response"
│        ]
│     └─ tags: {error=true, error.message="context deadline exceeded"}
└─ tags: {http.status=504, error=true}

Trace Summary:
  TraceID:     ghi-456-jkl
  Total spans: 5
  Duration:    340ms
  Root cause:  payment-service timeout
  Status:      ERROR

===== Sampling Stats =====
Total requests:  1000
Sampled:          102 (10.2% ≈ target 10%)
Dropped:          898
Reported spans:   612 (~6 spans/trace avg)
```

**Key learning:** Distributed tracing (Jaeger, Zipkin, OpenTelemetry), context propagation pattern, span trees, sampling strategies, how to debug cross-service latency in microservices.

---

## Problem 5: Health Check Framework

**Difficulty:** Moderate–Hard

**Concepts:** health check patterns, liveness vs readiness, dependency checking, degraded state

**Problem:**

Build a comprehensive health check framework for a Go service with multiple dependencies.

**Requirements:**
1. **Health check types**:
   - **Liveness**: is the process alive? (simple, returns 200 if not hung).
   - **Readiness**: can the service accept traffic? (checks all dependencies).
   - **Startup**: has the service finished initializing? (one-time checks).
2. **Dependency checks**: register checkers for: database, cache, external API, disk space, memory usage.
3. **Check execution**: run all checks concurrently with individual timeouts.
4. **Composite status**: `UP` (all pass), `DEGRADED` (non-critical checks failing), `DOWN` (critical check failing).
5. **Caching**: cache check results for N seconds to avoid hammering dependencies.
6. **Endpoint**: `/health/live`, `/health/ready`, `/health/startup` returning JSON.
7. **History**: keep last 10 check results per dependency for trend analysis.

```go
type HealthChecker struct {
    mu        sync.RWMutex
    checks    map[string]*Check
    cache     map[string]*CachedResult
    cacheTTL  time.Duration
    history   map[string][]*CheckResult
}

type Check struct {
    Name     string
    Type     CheckType // Liveness, Readiness, Startup
    Critical bool
    Timeout  time.Duration
    Fn       func(ctx context.Context) error
}

type CheckResult struct {
    Name     string
    Status   Status // UP, DOWN
    Duration time.Duration
    Error    string
    CheckedAt time.Time
}

type HealthResponse struct {
    Status  string                   `json:"status"` // UP, DEGRADED, DOWN
    Checks  map[string]*CheckResult  `json:"checks"`
    Uptime  time.Duration            `json:"uptime"`
}
```

**Expected Output:**
```
$ go run . -port=8080
[Health] Registered checks:
  postgres   (readiness, critical, timeout=2s)
  redis      (readiness, critical, timeout=1s)
  s3         (readiness, non-critical, timeout=3s)
  disk       (liveness, critical, timeout=100ms)
  memory     (liveness, critical, timeout=100ms)

--- GET /health/ready ---
{
  "status": "UP",
  "checks": {
    "postgres": {"status": "UP", "duration_ms": 12, "cached": false},
    "redis":    {"status": "UP", "duration_ms": 3,  "cached": false},
    "s3":       {"status": "UP", "duration_ms": 45, "cached": false}
  },
  "uptime": "2h15m"
}

--- Redis goes down ---
GET /health/ready
{
  "status": "DOWN",
  "checks": {
    "postgres": {"status": "UP",   "duration_ms": 11},
    "redis":    {"status": "DOWN", "error": "connection refused", "duration_ms": 1000},
    "s3":       {"status": "UP",   "duration_ms": 48}
  }
}

--- S3 goes down (non-critical) ---
GET /health/ready
{
  "status": "DEGRADED",
  "checks": {
    "postgres": {"status": "UP",   "duration_ms": 13},
    "redis":    {"status": "UP",   "duration_ms": 4},
    "s3":       {"status": "DOWN", "error": "timeout", "duration_ms": 3000}
  }
}

--- Cached result ---
GET /health/ready (within 5s of last check)
{
  "status": "UP",
  "checks": {
    "postgres": {"status": "UP", "cached": true, "cached_at": "2026-03-21T10:30:00Z"},
    ...
  }
}

--- GET /health/live ---
{
  "status": "UP",
  "checks": {
    "disk":   {"status": "UP", "details": "85% used (15% free)"},
    "memory": {"status": "UP", "details": "1.2GB / 4GB (30%)"}
  }
}

--- History ---
GET /health/history/redis
[
  {"status": "UP",   "checked_at": "10:28:00", "duration_ms": 3},
  {"status": "UP",   "checked_at": "10:29:00", "duration_ms": 4},
  {"status": "DOWN", "checked_at": "10:30:00", "duration_ms": 1000, "error": "refused"},
  {"status": "DOWN", "checked_at": "10:31:00", "duration_ms": 1000, "error": "refused"},
  {"status": "UP",   "checked_at": "10:32:00", "duration_ms": 5}
]
```

**Key learning:** Health check patterns (Kubernetes liveness/readiness probes), degraded state handling, concurrent dependency checking with timeouts, caching to avoid probe storms, standard health check APIs.

---

## Problem 6: Metrics Collector & Aggregator

**Difficulty:** Hard

**Concepts:** counters, gauges, histograms, atomic operations, time-series, percentile calculation

**Problem:**

Build a metrics collection library (Prometheus-client-like) with aggregation and exposition.

**Requirements:**
1. **Metric types**:
   - **Counter**: monotonically increasing (requests_total, errors_total).
   - **Gauge**: can go up and down (active_connections, queue_size).
   - **Histogram**: distribution of values with configurable buckets (request_duration_seconds).
   - **Summary**: streaming percentiles (p50, p90, p99) using a sliding window.
2. **Labels**: metrics support labels for dimensionality (e.g., `http_requests_total{method="GET", path="/api/users", status="200"}`).
3. **Thread-safe**: all metric operations must be safe for concurrent use. Use atomics for counters, mutex for histograms.
4. **Exposition**: `GET /metrics` returns Prometheus-compatible text format.
5. **Aggregation**: compute rates (requests per second), ratios (error rate), and percentiles.
6. Demonstrate: instrument a simulated HTTP server with request latency, error rate, and active connections.

```go
type Registry struct {
    mu       sync.RWMutex
    counters map[string]*Counter
    gauges   map[string]*Gauge
    histograms map[string]*Histogram
    summaries  map[string]*Summary
}

type Counter struct {
    name   string
    labels map[string]*atomic.Int64 // label_combo -> count
}

type Gauge struct {
    name   string
    labels map[string]*atomic.Int64
}

type Histogram struct {
    mu      sync.RWMutex
    name    string
    buckets []float64
    counts  map[string][]atomic.Uint64 // label -> bucket counts
    sums    map[string]*atomic.Int64
}

type Summary struct {
    mu      sync.RWMutex
    name    string
    window  time.Duration
    values  map[string]*slidingWindow
}
```

**Expected Output:**
```
$ go run .
[Metrics] Registry initialized
[Server] Simulating HTTP traffic...

--- After 60s of traffic ---

GET /metrics

# HELP http_requests_total Total HTTP requests
# TYPE http_requests_total counter
http_requests_total{method="GET",path="/api/users",status="200"} 4521
http_requests_total{method="GET",path="/api/users",status="500"} 47
http_requests_total{method="POST",path="/api/orders",status="201"} 1232
http_requests_total{method="POST",path="/api/orders",status="400"} 89

# HELP http_active_connections Current active connections
# TYPE http_active_connections gauge
http_active_connections 23

# HELP http_request_duration_seconds Request latency histogram
# TYPE http_request_duration_seconds histogram
http_request_duration_seconds_bucket{path="/api/users",le="0.01"} 2341
http_request_duration_seconds_bucket{path="/api/users",le="0.05"} 3987
http_request_duration_seconds_bucket{path="/api/users",le="0.1"} 4412
http_request_duration_seconds_bucket{path="/api/users",le="0.5"} 4556
http_request_duration_seconds_bucket{path="/api/users",le="1.0"} 4568
http_request_duration_seconds_bucket{path="/api/users",le="+Inf"} 4568
http_request_duration_seconds_sum{path="/api/users"} 142.7
http_request_duration_seconds_count{path="/api/users"} 4568

# HELP http_request_duration_summary Request latency summary
# TYPE http_request_duration_summary summary
http_request_duration_summary{path="/api/users",quantile="0.5"} 0.023
http_request_duration_summary{path="/api/users",quantile="0.9"} 0.078
http_request_duration_summary{path="/api/users",quantile="0.99"} 0.312

--- Computed metrics ---
Request rate:    98.2 req/s
Error rate:      2.3% (GET /api/users)
P99 latency:     312ms (GET /api/users)
Active conns:    23 (gauge)
```

**Key learning:** Prometheus metrics internals, counter/gauge/histogram/summary semantics, label cardinality, atomic operations for high-performance metrics, percentile calculation algorithms, observability foundations.

---

## Problem 7: Chaos Testing Framework

**Difficulty:** Hard

**Concepts:** fault injection, error simulation, latency injection, goroutine control, testing patterns

**Problem:**

Build a chaos testing framework that injects faults into a running system to test resilience.

**Requirements:**
1. **Fault types**:
   - **Latency injection**: add configurable delay to function calls.
   - **Error injection**: make functions fail with configurable probability.
   - **Timeout injection**: cause functions to hang until context cancels.
   - **Panic injection**: trigger panics (with recovery verification).
   - **Resource exhaustion**: simulate memory/goroutine/fd limits.
2. **Targeting**: faults target specific functions or services by name.
3. **Scheduling**: one-shot, periodic, or probability-based triggering.
4. **Chaos controller**: enable/disable faults at runtime via API or channel.
5. **Report**: after chaos run, report which faults fired, system behavior, and whether the system recovered.
6. Demonstrate: run chaos against a simulated order processing system and verify it degrades gracefully.

```go
type ChaosController struct {
    mu      sync.RWMutex
    faults  map[string]*Fault
    active  bool
    report  *ChaosReport
}

type Fault struct {
    Name        string
    Type        FaultType
    Target      string
    Probability float64   // 0.0-1.0
    Duration    time.Duration
    LatencyMS   int
    ErrorMsg    string
    Schedule    Schedule
    Active      atomic.Bool
}

type ChaosReport struct {
    FaultsFired   map[string]int
    ErrorsCaused  int
    RecoveryTimes map[string]time.Duration
    SystemHealthy bool
}

// Wrap functions with chaos
func (cc *ChaosController) Wrap(name string, fn func() error) func() error
```

**Expected Output:**
```
$ go run .
[Chaos] Controller initialized
[System] Order processing system running...

--- Phase 1: Baseline (no faults, 30s) ---
[System] Orders processed/s: 452  Errors: 0  P99: 45ms
[Chaos] Baseline recorded ✓

--- Phase 2: Latency injection (30s) ---
[Chaos] ENABLED: payment-service +500ms latency (100% of calls)
[System] t=31s: Orders/s: 312  Errors: 0  P99: 550ms
[System] t=40s: Orders/s: 289  Errors: 12 P99: 620ms (payment timeouts)
[Chaos] System degraded but functional (errors < 5%) ✓

--- Phase 3: Error injection (30s) ---
[Chaos] DISABLED: latency injection
[Chaos] ENABLED: inventory-service 50% error rate
[System] t=61s: Orders/s: 234  Errors: 118 P99: 89ms
[System] Circuit breaker OPENED for inventory-service
[System] t=65s: Orders/s: 445  Errors: 0   P99: 42ms (circuit breaker handling it)
[Chaos] System recovered via circuit breaker ✓

--- Phase 4: Cascading failure (30s) ---
[Chaos] ENABLED: database +2s latency + 30% errors
[System] t=91s: Connection pool exhausted! Queue backing up...
[System] t=93s: Backpressure activated, shedding load
[System] t=95s: Orders/s: 89   Errors: 234  P99: 2100ms
[System] t=100s: DB recovering, circuit half-open
[System] t=105s: Orders/s: 410  Errors: 12  P99: 52ms
[Chaos] System recovered after 14s ✓

--- Phase 5: Panic injection ---
[Chaos] ENABLED: order-handler 5% panic rate
[System] t=121s: PANIC in goroutine → recovered ✓
[System] t=123s: PANIC in goroutine → recovered ✓
[System] Orders/s: 430  Errors: 22 (panics recovered) ✓

========== CHAOS REPORT ==========
Test duration:       150s
Faults injected:     4

Fault: payment-latency
  Fired: 9,360 times
  System impact: P99 +505ms
  Recovery: immediate after disable

Fault: inventory-errors
  Fired: 3,510 times
  System impact: 50% error rate → 0% (circuit breaker)
  Recovery: 4s (circuit breaker trip)

Fault: db-latency+errors
  Fired: 2,670 + 801 times
  System impact: near-total degradation
  Recovery: 14s

Fault: panic
  Fired: 22 times
  All panics recovered: ✓

Overall: System demonstrated graceful degradation ✓
===================================
```

**Key learning:** Chaos engineering (Netflix Chaos Monkey, Gremlin), fault injection techniques, verifying resilience patterns actually work, how to design systems that degrade gracefully, testing distributed systems under failure.

---

## Concept Coverage Matrix

| Problem | Retry | Bulkhead | Backpressure | Tracing | Health | Metrics | Chaos |
|---------|-------|----------|-------------|---------|--------|---------|-------|
| 1. Retry Engine | ✓ | | | | | | |
| 2. Bulkhead | | ✓ | | | | | |
| 3. Backpressure | | | ✓ | | | | |
| 4. Distributed Tracing | | | | ✓ | | | |
| 5. Health Checks | | | | | ✓ | | |
| 6. Metrics Collector | | | | | | ✓ | |
| 7. Chaos Framework | ✓ | ✓ | ✓ | | ✓ | ✓ | ✓ |

---

## Recommended Order

1. **Problem 1** — Retry engine (foundational resilience pattern)
2. **Problem 5** — Health checks (Kubernetes integration)
3. **Problem 6** — Metrics collector (Prometheus model)
4. **Problem 4** — Distributed tracing (OpenTelemetry model)
5. **Problem 2** — Bulkhead (resource isolation)
6. **Problem 3** — Backpressure (flow control)
7. **Problem 7** — Chaos framework (ties everything together)
