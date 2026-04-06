# Recommended Solving Order

A structured path through all 39 exercises across 5 topic folders. Each phase builds on the previous — start with Go concurrency primitives, progress through networking and storage, layer on resilience and system design, and finish with distributed systems and chaos testing.

---

## Phase 1: Core Building Blocks

| # | Problem | Folder | Key Concepts |
|---|---------|--------|-------------|
| 1 | Worker Pool | concurrency/#1 | goroutines, channels, fan-out/fan-in |
| 2 | Pipeline Pattern | concurrency/#2 | channel chaining, stage composition |
| 3 | TTL Cache | concurrency/#4 | sync.RWMutex, expiration, concurrent map |
| 4 | Retry Engine | resilience-patterns/#1 | backoff strategies, context, jitter |

## Phase 2: Networking & Storage Foundations

| # | Problem | Folder | Key Concepts |
|---|---------|--------|-------------|
| 5 | TCP Chat Server | networking/#1 | net package, goroutine-per-connection |
| 6 | HTTP/1.1 Server from Scratch | networking/#2 | protocol parsing, request/response |
| 7 | Bloom Filter | storage-engines/#3 | hashing, probabilistic data structures |
| 8 | Concurrent Skip List | storage-engines/#4 | lock-free concurrency, ordered map |

## Phase 3: Resilience & Observability

| # | Problem | Folder | Key Concepts |
|---|---------|--------|-------------|
| 9 | Health Check Framework | resilience-patterns/#5 | Kubernetes liveness/readiness probes |
| 10 | Metrics Collector | resilience-patterns/#6 | Prometheus model, counters/histograms, atomics |
| 11 | Circuit Breaker | concurrency/#7 | state machines, failure detection |
| 12 | Bulkhead Pattern | resilience-patterns/#2 | resource isolation, semaphores, adaptive concurrency |

## Phase 4: System Design Patterns

| # | Problem | Folder | Key Concepts |
|---|---------|--------|-------------|
| 13 | In-Memory Message Queue | system-design/#1 | Kafka-lite, consumer groups, offsets |
| 14 | URL Shortener | system-design/#5 | base62 encoding, analytics, TTL |
| 15 | Pub/Sub Message Broker | concurrency/#6 | topic routing, fan-out delivery |
| 16 | Distributed Task Queue | system-design/#2 | Celery-like, priority, dead letter |
| 17 | API Gateway | system-design/#3 | middleware chain, routing, auth |

## Phase 5: Storage Engines Deep Dive

| # | Problem | Folder | Key Concepts |
|---|---------|--------|-------------|
| 18 | LSM Tree | storage-engines/#1 | memtable, SSTables, compaction |
| 19 | B+ Tree with Disk Persistence | storage-engines/#2 | page-based I/O, node splitting |
| 20 | Buffer Pool Manager | storage-engines/#5 | LRU-K eviction, page pinning |
| 21 | MVCC Transaction Manager | storage-engines/#6 | snapshot isolation, version chains |

## Phase 6: Distributed Systems

| # | Problem | Folder | Key Concepts |
|---|---------|--------|-------------|
| 22 | Consistent Hashing Ring | distributed-systems/#1 | virtual nodes, rebalancing |
| 23 | Vector Clocks & Conflict Detection | distributed-systems/#2 | causality, happens-before |
| 24 | Write-Ahead Log with Log Shipping | distributed-systems/#6 | durability, replication |
| 25 | Distributed KV Store | distributed-systems/#4 | replication, quorum reads/writes |
| 26 | Distributed Lock Manager | distributed-systems/#5 | fencing tokens, TTL, deadlock |
| 27 | Gossip Protocol | distributed-systems/#3 | epidemic broadcast, failure detection |

## Phase 7: Advanced Composition

| # | Problem | Folder | Key Concepts |
|---|---------|--------|-------------|
| 28 | Reverse Proxy / Load Balancer | networking/#3 | round-robin, health checks, buffering |
| 29 | Backpressure Controller | resilience-patterns/#3 | flow control, load shedding, scaling |
| 30 | Distributed Tracing (OTel-lite) | resilience-patterns/#4 | span trees, context propagation, sampling |
| 31 | Distributed Rate Limiter | system-design/#6 | sliding window, token bucket, Redis-like |
| 32 | Log Aggregation Service | system-design/#7 | log collection, indexing, search |
| 33 | Order Processing System | concurrency/#9 | full pipeline, saga pattern |
| 34 | P2P File Transfer | networking/#7 | chunking, peer discovery, integrity |

## Phase 8: Capstone

| # | Problem | Folder | Key Concepts |
|---|---------|--------|-------------|
| 35 | CRDTs | distributed-systems/#7 | eventual consistency, merge semantics |
| 36 | Cron Scheduler with Persistence | system-design/#8 | scheduling, WAL, leader election |
| 37 | Service Registry & Discovery | system-design/#4 | heartbeats, DNS, load balancing |
| 38 | Chaos Testing Framework | resilience-patterns/#7 | fault injection, recovery verification |
| 39 | Map-Reduce Framework | concurrency/#11 | distributed computation, shuffle |

---

## Remaining Exercises (solve anytime)

These don't have strict ordering dependencies — fit them in whenever you want extra practice:

| Problem | Folder | Key Concepts |
|---------|--------|-------------|
| Singleflight | concurrency/#3 | deduplication, sync.Once variant |
| Task Scheduler | concurrency/#5 | priority queue, cron expressions |
| Web Crawler | concurrency/#8 | BFS, politeness, dedup |
| Rate Limiter | concurrency/#10 | token bucket, sliding window |
| Connection Pool | concurrency/#12 | resource management, idle timeout |
| Leader Election | concurrency/#13 | consensus, bully algorithm |
| Graceful Server | concurrency/#14 | signal handling, drain connections |
| Merge Sort | concurrency/#15 | parallel divide-and-conquer |
| TCP Multiplexer | networking/#4 | stream multiplexing, framing |
| DNS Resolver | networking/#5 | UDP, recursive resolution |
| Port Scanner | networking/#6 | concurrent scanning, SYN detection |
| Service Registry | system-design/#4 | heartbeats, service mesh |
| Rate Limiter (distributed) | system-design/#6 | consensus, sliding window |

---

## Progress Tracker

- [x] 1. Worker Pool
- [x] 2. Pipeline Pattern
- [ ] 3. TTL Cache
- [ ] 4. Retry Engine
- [ ] 5. TCP Chat Server
- [ ] 6. HTTP/1.1 Server
- [ ] 7. Bloom Filter
- [ ] 8. Concurrent Skip List
- [ ] 9. Health Check Framework
- [ ] 10. Metrics Collector
- [ ] 11. Circuit Breaker
- [ ] 12. Bulkhead Pattern
- [ ] 13. In-Memory Message Queue
- [ ] 14. URL Shortener
- [ ] 15. Pub/Sub Message Broker
- [ ] 16. Distributed Task Queue
- [ ] 17. API Gateway
- [ ] 18. LSM Tree
- [ ] 19. B+ Tree
- [ ] 20. Buffer Pool Manager
- [ ] 21. MVCC Transaction Manager
- [ ] 22. Consistent Hashing Ring
- [ ] 23. Vector Clocks
- [ ] 24. Write-Ahead Log
- [ ] 25. Distributed KV Store
- [ ] 26. Distributed Lock Manager
- [ ] 27. Gossip Protocol
- [ ] 28. Reverse Proxy / Load Balancer
- [ ] 29. Backpressure Controller
- [ ] 30. Distributed Tracing
- [ ] 31. Distributed Rate Limiter
- [ ] 32. Log Aggregation Service
- [ ] 33. Order Processing System
- [ ] 34. P2P File Transfer
- [ ] 35. CRDTs
- [ ] 36. Cron Scheduler
- [ ] 37. Service Registry & Discovery
- [ ] 38. Chaos Testing Framework
- [ ] 39. Map-Reduce Framework
