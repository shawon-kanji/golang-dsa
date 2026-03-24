# Go System Design Exercises (Hard)

Build production-grade systems from scratch. Each problem is a real backend service you'd encounter in system design interviews and production environments.

---

## Problem 1: In-Memory Message Queue (Kafka-lite)

**Difficulty:** Very Hard

**Concepts:** partitioned logs, consumer groups, offset management, retention, concurrent producers/consumers

**Problem:**

Build a simplified in-memory message queue with Kafka-like semantics.

**Requirements:**
1. **Topics**: each topic has N partitions (configurable). Messages are ordered within a partition.
2. **Producer**: `Publish(topic, key, value)`. Key determines partition (hash(key) % numPartitions). If no key, round-robin.
3. **Consumer groups**: consumers in a group split partitions among themselves. A partition is assigned to exactly one consumer per group.
4. **Offset tracking**: each consumer tracks its offset per partition. Support `CommitOffset` and `SeekToOffset`.
5. **At-least-once delivery**: if a consumer crashes without committing, messages are re-delivered from the last committed offset.
6. **Retention**: messages older than a configurable TTL are garbage collected.
7. **Rebalancing**: when a consumer joins or leaves a group, partitions are reassigned.
8. Multiple consumer groups can independently consume the same topic.

```go
type Broker struct {
    mu     sync.RWMutex
    topics map[string]*Topic
    groups map[string]*ConsumerGroup
}

type Topic struct {
    name       string
    partitions []*Partition
}

type Partition struct {
    mu       sync.RWMutex
    id       int
    messages []*Message
    offset   int64 // next write offset
}

type Message struct {
    Key       []byte
    Value     []byte
    Offset    int64
    Timestamp time.Time
    Partition int
}

type ConsumerGroup struct {
    name        string
    consumers   map[string]*Consumer
    assignments map[int]string // partition -> consumerID
    offsets     map[int]int64  // partition -> committed offset
}
```

**Expected Output:**
```
$ go run .
[Broker] Started
[Topic] "orders" created (partitions=3)

--- Producing ---
[Producer] Publish("orders", key="user-1", "order-A") → partition=0, offset=0
[Producer] Publish("orders", key="user-2", "order-B") → partition=1, offset=0
[Producer] Publish("orders", key="user-1", "order-C") → partition=0, offset=1
[Producer] Published 1000 messages across 3 partitions

--- Consumer Group "processors" (2 consumers) ---
[Rebalance] Group "processors":
  consumer-1 → [partition-0, partition-1]
  consumer-2 → [partition-2]

[consumer-1] partition-0, offset=0: "order-A" ✓
[consumer-1] partition-0, offset=1: "order-C" ✓
[consumer-1] CommitOffset(partition=0, offset=2)
[consumer-2] partition-2, offset=0: "order-X" ✓

--- Consumer crash + rebalance ---
[consumer-2] CRASHED
[Rebalance] Group "processors":
  consumer-1 → [partition-0, partition-1, partition-2]
[consumer-1] partition-2 resuming from committed offset=0
[consumer-1] partition-2, offset=0: "order-X" (re-delivered!) ← at-least-once

--- Second consumer group (independent) ---
[Consumer Group "analytics"] consuming same topic from offset=0
  All 1000 messages consumed independently ✓

--- Retention ---
[GC] Removing messages older than 1h from "orders"
[GC] Removed 342 messages from partition-0

========== BROKER STATS ==========
Topics:            1
Total messages:    1000
Consumer groups:   2
Active consumers:  3
Rebalances:        2
Messages retained: 658
===================================
```

**Key learning:** Kafka internals (partitioned log, consumer groups, offset management), at-least-once delivery semantics, rebalancing strategies, retention policies.

---

## Problem 2: Distributed Task Queue (Celery/Sidekiq-like)

**Difficulty:** Hard

**Concepts:** task serialization, priority queues, delayed execution, retries with backoff, dead letter queue

**Problem:**

Build a task queue system where producers enqueue tasks and workers process them.

**Requirements:**
1. **Task definition**: tasks have a name, payload (JSON), priority, and optional scheduled time (for delayed execution).
2. **Enqueue**: `Enqueue(task)` — immediate execution. `EnqueueAt(task, time)` — delayed execution.
3. **Workers**: configurable pool of workers that pull tasks by priority (highest first).
4. **Retries with exponential backoff**: on failure, retry up to N times with increasing delays (1s, 2s, 4s, 8s...).
5. **Dead Letter Queue (DLQ)**: tasks that exhaust all retries go to a DLQ for manual inspection.
6. **Task middleware/hooks**: `BeforeExecute`, `AfterExecute`, `OnFailure` callbacks.
7. **Unique tasks**: optional deduplication — don't enqueue a duplicate of an in-flight task with the same key.
8. **Progress tracking**: `GetTaskStatus(id) → (pending|running|completed|failed|dead)`.

```go
type Queue struct {
    mu          sync.Mutex
    pending     *PriorityQueue
    delayed     *DelayedQueue
    inFlight    map[string]*Task
    dlq         []*Task
    workers     int
    middleware  []Middleware
}

type Task struct {
    ID         string
    Name       string
    Payload    json.RawMessage
    Priority   int
    MaxRetries int
    Retries    int
    Status     TaskStatus
    ScheduleAt *time.Time
    UniqueKey  string
    CreatedAt  time.Time
}

type Middleware func(task *Task, next func() error) error
```

**Expected Output:**
```
$ go run .
[TaskQueue] Started (5 workers)
[Middleware] Registered: logging, metrics, timeout(30s)

--- Immediate tasks ---
[Enqueue] task-001 "send_email" (priority=HIGH)
[Enqueue] task-002 "resize_image" (priority=LOW)
[Enqueue] task-003 "send_email" (priority=HIGH)
[Worker-2] Executing task-001 "send_email" ... ✓ (120ms)
[Worker-1] Executing task-003 "send_email" ... ✓ (95ms)
[Worker-3] Executing task-002 "resize_image" ... ✓ (340ms)

--- Delayed task ---
[Enqueue] task-004 "daily_report" scheduled for +5s
[Scheduler] task-004 ready → moving to pending queue
[Worker-1] Executing task-004 "daily_report" ... ✓ (2.1s)

--- Retry with backoff ---
[Worker-2] Executing task-005 "webhook_call" ... FAILED (connection refused)
[Retry] task-005 retry 1/3, backoff=1s
[Worker-4] Executing task-005 "webhook_call" ... FAILED
[Retry] task-005 retry 2/3, backoff=2s
[Worker-1] Executing task-005 "webhook_call" ... FAILED
[Retry] task-005 retry 3/3, backoff=4s
[Worker-3] Executing task-005 "webhook_call" ... FAILED
[DLQ] task-005 exhausted retries → moved to dead letter queue

--- Unique task dedup ---
[Enqueue] task-007 "sync_user:42" (unique_key="sync:42")
[Enqueue] task-008 "sync_user:42" (unique_key="sync:42") → SKIPPED (duplicate in-flight)

--- Status ---
[Status] task-001: completed
[Status] task-005: dead (in DLQ)
[Status] task-007: running

========== QUEUE STATS ==========
Enqueued:      15
Completed:     11
Failed:         3 (all retried)
Dead (DLQ):     1
Duplicates:     1 (skipped)
Avg exec time:  450ms
Avg queue time: 120ms
==================================
```

**Key learning:** Task queue architecture (Celery, Sidekiq, Bull), priority scheduling, exponential backoff, dead letter queue pattern, middleware/interceptor pattern, at-least-once execution.

---

## Problem 3: API Gateway with Middleware Chain

**Difficulty:** Hard

**Concepts:** middleware pattern, authentication, rate limiting, circuit breaking, request routing, metrics

**Problem:**

Build an API gateway that sits in front of backend services and handles cross-cutting concerns.

**Requirements:**
1. **Request routing**: route requests to different backends based on path prefix (e.g., `/api/users/*` → user-service, `/api/orders/*` → order-service).
2. **Middleware chain** (executed in order for every request):
   - **Request ID**: inject a unique `X-Request-ID` header.
   - **Logging**: log method, path, status code, latency.
   - **Authentication**: validate a bearer token (simulated). Reject with 401 if invalid.
   - **Rate limiting**: per-client token bucket (from concurrency Problem 10). Return 429 if exceeded.
   - **Circuit breaker**: per-backend circuit breaker (from concurrency Problem 7). Return 503 if circuit is open.
   - **Timeout**: wrap backend call with `context.WithTimeout`. Return 504 if exceeded.
   - **Retry**: retry failed backend calls up to 2 times with backoff.
3. **Health check**: `/health` endpoint that checks all backends.
4. **Metrics endpoint**: `/metrics` returns JSON with per-route request count, error rate, p50/p99 latency.

```go
type Gateway struct {
    routes     map[string]*Route
    middleware []Middleware
    metrics    *MetricsCollector
}

type Route struct {
    PathPrefix string
    Backend    string
    Breaker    *CircuitBreaker
    Limiter    *RateLimiter
}

type Middleware func(Handler) Handler
type Handler func(ctx context.Context, req *Request) *Response

type MetricsCollector struct {
    mu        sync.RWMutex
    routes    map[string]*RouteMetrics
}

type RouteMetrics struct {
    Requests    int64
    Errors      int64
    Latencies   []time.Duration // for percentile calculation
}
```

**Expected Output:**
```
$ go run . -port=8080
[Gateway] Routes:
  /api/users/*  → localhost:9001 (user-service)
  /api/orders/* → localhost:9002 (order-service)
[Gateway] Middleware: [request-id, logging, auth, rate-limit, circuit-breaker, timeout, retry]
[Gateway] Listening on :8080

[req-abc123] GET /api/users/42
  [auth] Token valid (user=admin) ✓
  [rate-limit] client=admin: 5/100 tokens remaining ✓
  [circuit-breaker] user-service: CLOSED ✓
  [proxy] → localhost:9001/users/42 (45ms)
  [log] 200 OK 45ms

[req-def456] GET /api/orders/99
  [auth] Token valid (user=admin) ✓
  [rate-limit] client=admin: 4/100 tokens remaining ✓
  [circuit-breaker] order-service: CLOSED ✓
  [proxy] → localhost:9002/orders/99 TIMEOUT (5s deadline)
  [retry] attempt 2/3 → localhost:9002/orders/99 (120ms)
  [log] 200 OK 5120ms (with retry)

[req-ghi789] POST /api/users
  [auth] INVALID TOKEN → 401 Unauthorized
  [log] 401 Unauthorized 0.1ms

[req-jkl012] GET /api/orders/1
  [circuit-breaker] order-service: OPEN (5 recent failures) → 503 Service Unavailable
  [log] 503 0.02ms

--- Metrics ---
$ curl localhost:8080/metrics
{
  "/api/users": {
    "requests": 847,
    "errors": 12,
    "error_rate": "1.4%",
    "latency_p50_ms": 23,
    "latency_p99_ms": 189
  },
  "/api/orders": {
    "requests": 523,
    "errors": 67,
    "error_rate": "12.8%",
    "latency_p50_ms": 45,
    "latency_p99_ms": 4200
  }
}
```

**Key learning:** API gateway architecture (Kong, Envoy), middleware chain pattern, composing resilience patterns (rate limit + circuit breaker + retry + timeout), metrics collection, understanding what API gateways like Kong/Envoy actually do.

---

## Problem 4: Service Registry & Discovery

**Difficulty:** Hard

**Concepts:** heartbeats, TTL, health status, watch/notify pattern, consistent reads

**Problem:**

Build a service registry where microservices register themselves and discover each other (like Consul/etcd/Eureka).

**Requirements:**
1. **Register**: service instances register with: name, address, port, metadata, health check URL.
2. **Heartbeat**: instances send heartbeats every N seconds. If no heartbeat within TTL, mark as `critical` then deregister.
3. **Health states**: `passing`, `warning`, `critical`. Background health checker calls each service's health endpoint.
4. **Discovery**: `Discover(serviceName) []Instance` — returns healthy instances. Optional: filter by tags/metadata.
5. **Watch**: `Watch(serviceName, callback)` — get notified when instances join, leave, or change health status. Uses long-polling or channels.
6. **DNS interface**: `Resolve(serviceName) []string` — returns IP:port list (simulated DNS-style lookup).
7. Multiple clients registering, querying, and watching simultaneously.

```go
type Registry struct {
    mu        sync.RWMutex
    services  map[string][]*Instance
    watchers  map[string][]chan Event
    ttl       time.Duration
}

type Instance struct {
    ID         string
    Service    string
    Address    string
    Port       int
    Tags       []string
    Metadata   map[string]string
    Health     HealthStatus
    LastSeen   time.Time
    RegisteredAt time.Time
}

type Event struct {
    Type     EventType // Registered, Deregistered, HealthChanged
    Instance *Instance
}
```

**Expected Output:**
```
$ go run .
[Registry] Started (heartbeat TTL=10s, health check interval=5s)

--- Registration ---
[Register] user-service/inst-1 (192.168.1.10:8080) tags=["v2","primary"]
[Register] user-service/inst-2 (192.168.1.11:8080) tags=["v2","secondary"]
[Register] order-service/inst-1 (192.168.1.20:8080)

--- Discovery ---
[Discover] "user-service" → 2 healthy instances:
  inst-1 192.168.1.10:8080 [passing] tags=[v2,primary]
  inst-2 192.168.1.11:8080 [passing] tags=[v2,secondary]

[Discover] "user-service" (tag="primary") → 1 instance:
  inst-1 192.168.1.10:8080 [passing]

--- Health degradation ---
[HealthCheck] user-service/inst-2: /health → 503 (warning)
[Event] user-service: HealthChanged inst-2 passing → warning
[Watcher-1] Notified: inst-2 health changed to WARNING

[HealthCheck] user-service/inst-2: /health → connection refused (critical)
[Event] user-service: HealthChanged inst-2 warning → critical

--- TTL expiry ---
[Heartbeat] user-service/inst-2: last heartbeat 12s ago > TTL 10s
[Deregister] user-service/inst-2 (TTL expired)
[Event] user-service: Deregistered inst-2
[Watcher-1] Notified: inst-2 deregistered

[Discover] "user-service" → 1 healthy instance:
  inst-1 192.168.1.10:8080 [passing]

--- Re-registration ---
[Register] user-service/inst-2 (192.168.1.11:8080) — recovered
[Event] user-service: Registered inst-2
[Watcher-1] Notified: inst-2 registered (health=passing)

========== REGISTRY STATS ==========
Services:        2
Total instances: 3
Healthy:         3
Warning:         0
Critical:        0
Deregistered (TTL): 1
Watcher callbacks:  4
=====================================
```

**Key learning:** Service discovery patterns (used in Consul, etcd, Eureka), heartbeat/TTL mechanisms, watch/notification pattern, health check design, understanding how Kubernetes service discovery works.

---

## Problem 5: URL Shortener with Analytics

**Difficulty:** Moderate–Hard

**Concepts:** base62 encoding, concurrent map, TTL, analytics pipeline, atomic counters

**Problem:**

Build a URL shortener service with real-time click analytics.

**Requirements:**
1. **Shorten**: `POST /shorten {"url": "https://...", "ttl": "24h", "custom_alias": "mylink"}` → returns short code.
2. **Redirect**: `GET /:code` → 302 redirect to original URL. If expired or not found, 404.
3. **ID generation**: base62 encoding of an atomic counter (or random for collision avoidance).
4. **Custom aliases**: optional user-defined short codes. Reject duplicates.
5. **Analytics pipeline**: every click is asynchronously logged to an analytics channel:
   - Analytics worker aggregates: clicks per URL, clicks per hour, unique visitors (by simulated IP), referrers.
6. **Stats endpoint**: `GET /stats/:code` → returns JSON with click count, unique visitors, top referrers, click timeline.
7. **TTL**: short URLs expire after configurable duration. Background cleanup goroutine.
8. Thread-safe: handle concurrent shortening + redirects.

```go
type URLShortener struct {
    mu        sync.RWMutex
    urls      map[string]*ShortURL
    counter   atomic.Uint64
    analytics chan ClickEvent
}

type ShortURL struct {
    Code      string
    Original  string
    CreatedAt time.Time
    ExpiresAt time.Time
    Clicks    atomic.Int64
}

type ClickEvent struct {
    Code      string
    IP        string
    Referrer  string
    Timestamp time.Time
}

type URLStats struct {
    TotalClicks   int64
    UniqueVisitors int
    TopReferrers  map[string]int
    ClicksPerHour map[string]int // "2026-03-21T10:00" -> count
}
```

**Expected Output:**
```
$ go run . -port=8080
[Shortener] Listening on :8080

--- Shorten ---
POST /shorten {"url":"https://golang.org/doc/effective_go","ttl":"24h"}
→ {"short_url":"http://localhost:8080/aB3x9","code":"aB3x9","expires":"2026-03-22T10:00:00Z"}

POST /shorten {"url":"https://go.dev","custom_alias":"godev"}
→ {"short_url":"http://localhost:8080/godev","code":"godev"}

POST /shorten {"custom_alias":"godev"} → 409 Conflict (alias taken)

--- Redirect ---
GET /aB3x9 (IP=1.2.3.4, Referrer=twitter.com)
→ 302 Location: https://golang.org/doc/effective_go

GET /aB3x9 (IP=5.6.7.8, Referrer=reddit.com)
→ 302 Location: https://golang.org/doc/effective_go

GET /expired123 → 404 Not Found (expired)

--- Analytics ---
[Analytics] Processing click event for "aB3x9"
[Analytics] Processing click event for "aB3x9"

GET /stats/aB3x9
→ {
    "code": "aB3x9",
    "original_url": "https://golang.org/doc/effective_go",
    "total_clicks": 247,
    "unique_visitors": 189,
    "top_referrers": {
      "twitter.com": 98,
      "reddit.com": 72,
      "google.com": 45,
      "direct": 32
    },
    "clicks_per_hour": {
      "2026-03-21T09:00": 45,
      "2026-03-21T10:00": 123,
      "2026-03-21T11:00": 79
    }
  }

========== SERVICE STATS ==========
URLs created:      124
URLs expired:       12
Total redirects:   2,847
Analytics events:  2,847
Active short URLs: 112
====================================
```

**Key learning:** ID generation strategies (sequential vs random), base62 encoding, async analytics pipeline with buffered channels, TTL management, building a complete CRUD service with Go.

---

## Problem 6: Distributed Rate Limiter (Redis-like)

**Difficulty:** Hard

**Concepts:** sliding window log, token bucket with persistence, lua-script-like atomics, cluster coordination

**Problem:**

Build a centralized rate limiter service that multiple API servers call to check rate limits — like what you'd build with Redis in production.

**Requirements:**
1. **Server**: listens on TCP, accepts rate limit check requests from multiple clients.
2. **Protocol**: simple text protocol:
   - `ALLOW <client_id> <limit> <window_seconds>` → `OK <remaining>` or `DENIED <retry_after_ms>`
   - `STATUS <client_id>` → `LIMIT=100 REMAINING=42 RESET=1234567890`
3. **Sliding window log algorithm**: store exact timestamps of each request. Count requests in the last N seconds.
4. **Global rate limits**: optional limit applied across all clients (e.g., 10,000 req/s total).
5. **Burst handling**: allow short bursts above the sustained rate (configurable burst multiplier).
6. **Persistence**: periodically snapshot rate limit state to disk. Restore on restart so counters survive restarts.
7. **Multi-server coordination**: if running multiple limiter instances, they share state via periodic sync (eventual consistency for rate limits is usually acceptable).
8. Demonstrate: 10 simulated API servers checking limits for 100 different client IDs.

```go
type RateLimitServer struct {
    listener   net.Listener
    mu         sync.RWMutex
    windows    map[string]*SlidingWindowLog
    globalRate *TokenBucket
    config     Config
}

type SlidingWindowLog struct {
    mu         sync.Mutex
    timestamps []int64 // unix milliseconds
    limit      int
    window     time.Duration
}

type Config struct {
    DefaultLimit    int
    DefaultWindow   time.Duration
    BurstMultiplier float64
    GlobalLimit     int
    SnapshotInterval time.Duration
}
```

**Expected Output:**
```
$ go run . -port=6380
[RateLimiter] Listening on :6380 (snapshot every 30s)

[Client] ALLOW user:alice 100 60
[Server] user:alice → OK remaining=99 (1/100 in 60s window)

[Client] ALLOW user:alice 100 60
[Server] user:alice → OK remaining=98

... (98 more requests from alice in quick succession)

[Client] ALLOW user:alice 100 60
[Server] user:alice → DENIED retry_after_ms=42380

[Client] STATUS user:alice
[Server] LIMIT=100 REMAINING=0 RESET=1679400000

--- Burst handling ---
[Client] ALLOW user:bob 10 1 (10 req/sec)
[Server] user:bob → OK (burst allows up to 15 in first 100ms)
... (14 rapid requests)
[Server] user:bob → OK remaining=0 (burst exhausted too)

--- Global limit ---
[Server] Global rate: 8,450 / 10,000 req/s
[Client] ALLOW user:charlie 100 60
[Server] Global limit approaching → OK (warning: 85% global capacity)

--- Snapshot ---
[Snapshot] Saved 100 client states to rate_limits.snapshot (4.2KB)

--- Restart recovery ---
[RateLimiter] Restoring from rate_limits.snapshot
[RateLimiter] Restored 100 client states ✓
[Client] STATUS user:alice
[Server] LIMIT=100 REMAINING=0 RESET=1679400000 (persisted!)

========== RATE LIMITER STATS ==========
Active clients:    100
Requests checked:  45,230
Allowed:           38,901 (86%)
Denied:             6,329 (14%)
Global utilization: 72%
Snapshots taken:    5
=========================================
```

**Key learning:** Centralized rate limiting (how Stripe, GitHub, Cloudflare implement it), sliding window log algorithm, state persistence/recovery, burst handling, building a TCP service from scratch.

---

## Problem 7: Log Aggregation Service (ELK-lite)

**Difficulty:** Hard

**Concepts:** structured logging, indexing, full-text search, retention, pipeline processing

**Problem:**

Build a log aggregation service that ingests, indexes, and queries structured logs from multiple services.

**Requirements:**
1. **Ingest API**: accept log entries via TCP or channel from multiple producer goroutines (simulating microservices).
2. **Log format**: structured JSON with fields: `timestamp`, `level`, `service`, `message`, `trace_id`, plus arbitrary fields.
3. **Indexing**: build inverted indexes for efficient querying:
   - Field-value index: `service=user-service` → list of matching log IDs.
   - Full-text index on `message` field (simple word tokenization).
4. **Query language** (simple):
   - `service=order-service AND level=ERROR` → filter by fields.
   - `message CONTAINS "timeout"` → full-text search.
   - `timestamp > 2026-03-21T10:00:00Z` → time range.
   - `| count by service` → aggregation.
5. **Retention**: delete logs older than configurable duration. Rebuild indexes.
6. **Tail mode**: `Tail(filter) <-chan LogEntry` — real-time streaming of matching logs.
7. Benchmark: ingest 1M log entries, measure query latency.

```go
type LogStore struct {
    mu          sync.RWMutex
    entries     []*LogEntry
    fieldIndex  map[string]map[string][]int // field -> value -> entry IDs
    textIndex   map[string][]int            // word -> entry IDs
    tailers     []*Tailer
}

type LogEntry struct {
    ID        int
    Timestamp time.Time
    Level     string
    Service   string
    Message   string
    TraceID   string
    Fields    map[string]string
}

type Query struct {
    Filters     []Filter
    TextSearch  string
    TimeRange   *TimeRange
    Aggregation *Aggregation
    Limit       int
}
```

**Expected Output:**
```
$ go run .
[LogStore] Initialized (retention=24h)
[Ingest] Listening for log entries...

--- Ingesting from 5 simulated services ---
[Ingest] user-service: 200,000 entries
[Ingest] order-service: 200,000 entries
[Ingest] payment-service: 200,000 entries
[Ingest] auth-service: 200,000 entries
[Ingest] notification-service: 200,000 entries
[Index] Built field index: 5 fields, 23 unique values
[Index] Built text index: 12,847 unique words

--- Queries ---
> service=order-service AND level=ERROR
  Results: 1,247 entries (query time: 2ms)
  [2026-03-21T10:15:30] ERROR order-service "Payment timeout for order ORD-4521"
  [2026-03-21T10:15:31] ERROR order-service "Inventory check failed: item-789"
  ... (1,245 more)

> message CONTAINS "timeout" AND level=ERROR
  Results: 3,891 entries (query time: 5ms)

> timestamp > 2026-03-21T10:00:00Z AND timestamp < 2026-03-21T11:00:00Z | count by service
  Results:
    user-service:         8,412
    order-service:        7,891
    payment-service:      8,203
    auth-service:         9,104
    notification-service: 6,722
  Query time: 12ms

> service=payment-service | count by level
  Results:
    INFO:  178,234
    WARN:   15,891
    ERROR:   5,875
  Query time: 8ms

--- Tail mode ---
> TAIL service=auth-service AND level=ERROR
  [streaming] 2026-03-21T10:30:01 ERROR "Invalid token for user-789"
  [streaming] 2026-03-21T10:30:03 ERROR "Rate limit exceeded for IP 1.2.3.4"
  ... (real-time)

--- Retention ---
[GC] Removed 45,123 entries older than 24h
[GC] Rebuilt indexes (32ms)

========== LOG STORE STATS ==========
Total entries:    954,877
Index size:       12.3MB
Unique services:  5
Unique words:     12,847
Avg query time:   7ms
Tail subscribers: 2
======================================
```

**Key learning:** Log aggregation architecture (Elasticsearch, Loki), inverted index construction, query parsing and execution, real-time streaming with goroutine-based tailers, retention management.

---

## Problem 8: Cron Scheduler with Persistence

**Difficulty:** Moderate–Hard

**Concepts:** time parsing, heap-based scheduling, persistence, concurrent job execution, timezone handling

**Problem:**

Build a cron-like job scheduler that persists job definitions and execution history.

**Requirements:**
1. **Cron expression parser**: parse standard cron expressions (`* * * * *` — minute, hour, day, month, weekday). Support `*/5`, `1-15`, `1,15,30`.
2. **Job registry**: `Register(name, cronExpr, handler)`. Jobs can be added/removed at runtime.
3. **Scheduler**: compute next run time for each job. Use a min-heap ordered by next run time. Sleep until the next job is due.
4. **Concurrent execution**: each job runs in its own goroutine. Configurable max concurrent jobs.
5. **Overlap prevention**: option to skip a job run if the previous run hasn't finished yet.
6. **Persistence**: save job definitions and last run times to a JSON file. On restart, resume scheduling.
7. **Execution history**: keep last N runs per job with status, duration, and error message.
8. **Timezone support**: jobs run in a configurable timezone.

```go
type Scheduler struct {
    mu       sync.Mutex
    jobs     map[string]*Job
    heap     *JobHeap // min-heap by nextRun
    running  map[string]bool
    maxConc  int
    sem      chan struct{}
    history  map[string][]*RunRecord
    tz       *time.Location
    persist  string // file path
}

type Job struct {
    Name       string
    CronExpr   string
    Schedule   *CronSchedule
    Handler    func(ctx context.Context) error
    NextRun    time.Time
    NoOverlap  bool
    Timeout    time.Duration
}

type CronSchedule struct {
    Minutes  []int
    Hours    []int
    Days     []int
    Months   []int
    Weekdays []int
}
```

**Expected Output:**
```
$ go run .
[Scheduler] Started (timezone=America/New_York, max_concurrent=5)
[Scheduler] Loaded 3 jobs from scheduler.json

[Register] "cleanup"      "0 * * * *"    (every hour at :00)
[Register] "daily-report"  "30 9 * * 1-5" (9:30 AM weekdays)
[Register] "health-check"  "*/5 * * * *"  (every 5 minutes)
[Register] "db-backup"     "0 2 * * *"    (2:00 AM daily)

[Schedule] Next runs:
  health-check:  2026-03-21T10:15:00 (in 2m30s)
  cleanup:       2026-03-21T11:00:00 (in 47m30s)
  daily-report:  2026-03-22T09:30:00 (in 23h17m)
  db-backup:     2026-03-22T02:00:00 (in 15h47m)

[t=10:15:00] Executing "health-check"... ✓ (45ms)
[t=10:20:00] Executing "health-check"... ✓ (52ms)
[t=10:25:00] Executing "health-check"... running (slow! 4.8s so far)
[t=10:30:00] "health-check" due but previous still running → SKIPPED (no-overlap)
[t=10:25:05] "health-check" completed (5.1s)
[t=10:35:00] Executing "health-check"... ✓ (48ms)

[History] health-check:
  10:15:00  ✓  45ms
  10:20:00  ✓  52ms
  10:25:00  ✓  5100ms  (slow)
  10:30:00  ⏭  skipped (overlap)
  10:35:00  ✓  48ms

[Persist] Saved scheduler state to scheduler.json

========== SCHEDULER STATS ==========
Registered jobs:  4
Executions today: 12
Successful:       11
Skipped (overlap): 1
Failed:            0
======================================
```

**Key learning:** Cron expression parsing, heap-based scheduling (how systemd timers and Kubernetes CronJobs work), overlap prevention, state persistence, timezone-aware scheduling.

---

## Concept Coverage Matrix

| Problem | Networking | Persistence | Concurrency | Queuing | Indexing | Scheduling | Middleware |
|---------|-----------|-------------|-------------|---------|---------|------------|-----------|
| 1. Message Queue | | | ✓ | ✓ | | | |
| 2. Task Queue | | | ✓ | ✓ | | ✓ | ✓ |
| 3. API Gateway | ✓ | | ✓ | | | | ✓ |
| 4. Service Registry | ✓ | | ✓ | | | | |
| 5. URL Shortener | ✓ | | ✓ | | | | |
| 6. Rate Limiter | ✓ | ✓ | ✓ | | | | |
| 7. Log Aggregation | ✓ | | ✓ | | ✓ | | |
| 8. Cron Scheduler | | ✓ | ✓ | | | ✓ | |

---

## Recommended Order

1. **Problem 5** — URL shortener (warm up: complete service)
2. **Problem 8** — Cron scheduler (scheduling fundamentals)
3. **Problem 2** — Task queue (queuing + retries)
4. **Problem 1** — Message queue (Kafka internals)
5. **Problem 4** — Service registry (discovery patterns)
6. **Problem 6** — Rate limiter service (stateful TCP service)
7. **Problem 3** — API gateway (combines multiple patterns)
8. **Problem 7** — Log aggregation (indexing + querying)
