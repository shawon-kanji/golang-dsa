# Go Distributed Systems Exercises (Hard)

Build core distributed systems primitives from scratch. These teach consensus, replication, partitioning, and failure handling — the concepts behind databases, message queues, and orchestration systems.

---

## Problem 1: Consistent Hashing Ring

**Difficulty:** Hard

**Concepts:** hash functions, virtual nodes, ring topology, rebalancing, replication

**Problem:**

Build a consistent hashing ring that distributes keys across nodes with minimal reshuffling when nodes join or leave.

**Requirements:**
1. `AddNode(id string)` — add a node with V virtual nodes (configurable, e.g., 150) placed on the ring using `crc32` or `fnv` hash.
2. `RemoveNode(id string)` — remove a node and all its virtual nodes.
3. `GetNode(key string) string` — find the responsible node for a given key (clockwise walk).
4. `GetNodes(key string, n int) []string` — get N distinct physical nodes for replication (skip duplicate physicals from virtual nodes).
5. Track key distribution across nodes. Print standard deviation to show balance.
6. Demonstrate: start with 3 nodes, insert 100,000 keys, add a 4th node → show that only ~25% of keys are remapped (vs ~75% with naive modulo hashing).

```go
type Ring struct {
    mu           sync.RWMutex
    nodes        map[string]*Node
    sortedHashes []uint32
    hashMap      map[uint32]string // hash -> physical node ID
    vnodeCount   int
}

type Node struct {
    ID     string
    Addr   string
    VNodes []uint32
}
```

**Expected Output:**
```
$ go run .
[Ring] Created with vnodes=150

[Ring] Added node-A (150 vnodes)
[Ring] Added node-B (150 vnodes)
[Ring] Added node-C (150 vnodes)

[Distribution] 100,000 keys across 3 nodes:
  node-A: 33,412 keys (33.4%)
  node-B: 33,187 keys (33.2%)
  node-C: 33,401 keys (33.4%)
  Std deviation: 0.12% (excellent balance)

[Ring] Adding node-D...
[Rebalance] Keys remapped: 24,891 / 100,000 (24.9%)
  node-A: 33,412 → 25,102 (moved 8,310 keys)
  node-B: 33,187 → 24,956 (moved 8,231 keys)
  node-C: 33,401 → 25,091 (moved 8,310 keys)
  node-D: 0      → 24,851 (received 24,891 keys)

[Comparison] Naive modulo hash: 74,892 keys remapped (74.9%)
[Comparison] Consistent hash:   24,891 keys remapped (24.9%) ← 3x better

[Replication] GetNodes("order-12345", 3) = [node-B, node-D, node-A]
```

**Key learning:** Understanding why consistent hashing is used in DynamoDB, Cassandra, and CDNs. Virtual nodes for balance. Replication placement on the ring.

---

## Problem 2: Vector Clocks & Conflict Detection

**Difficulty:** Hard

**Concepts:** logical clocks, partial ordering, causal consistency, conflict resolution

**Problem:**

Implement vector clocks for a simulated distributed key-value store with 3 replicas that can accept writes independently.

**Requirements:**
1. Each replica maintains a vector clock `map[string]uint64` (replica ID → counter).
2. On local write, increment own counter. Attach vector clock to the value.
3. On receiving a write from another replica, merge vector clocks (`max(local[i], remote[i])` for each entry).
4. Detect **conflicts**: two versions conflict if neither vector clock dominates the other (concurrent writes).
5. Implement conflict resolution strategies:
   - **Last-writer-wins** (use wall-clock timestamp as tiebreaker).
   - **Multi-value** (keep both versions like Dynamo — "siblings").
6. Simulate: 3 replicas, network partitions cause concurrent writes to the same key, then partition heals and replicas sync.

```go
type VectorClock map[string]uint64

func (vc VectorClock) Increment(nodeID string)
func (vc VectorClock) Merge(other VectorClock)
func (vc VectorClock) Compare(other VectorClock) Ordering // Before, After, Concurrent

type Ordering int
const (
    Before     Ordering = iota // vc happened before other
    After                      // vc happened after other
    Concurrent                 // neither dominates — CONFLICT
    Equal
)

type Replica struct {
    id     string
    store  map[string]*VersionedValue
    clock  VectorClock
    peers  []*Replica
}

type VersionedValue struct {
    Value   string
    Clock   VectorClock
    Timestamp time.Time
}
```

**Expected Output:**
```
$ go run .
[Cluster] 3 replicas: A, B, C

--- Normal operation ---
[A] PUT "user:1" = "Alice"  clock={A:1}
[A] Replicate to B, C
[B] Received "user:1"="Alice" clock={A:1} — accepted (A:1 > local ∅)
[C] Received "user:1"="Alice" clock={A:1} — accepted

--- Network partition: A isolated from B,C ---
[A] PUT "user:1" = "Alice-v2"  clock={A:2}
[B] PUT "user:1" = "Bob-edit"  clock={A:1, B:1}

[A] GET "user:1" → "Alice-v2"  clock={A:2}
[B] GET "user:1" → "Bob-edit"  clock={A:1, B:1}
  ⚠ Compare {A:2} vs {A:1, B:1} → CONCURRENT (conflict!)

--- Partition healed: A syncs with B ---
[Sync] A sends "user:1"="Alice-v2" {A:2} to B
[B] Compare {A:2} vs {A:1, B:1}:
  A: 2 > 1 (A wins this dimension)
  B: 0 < 1 (B wins this dimension)
  → CONCURRENT — conflict detected!

[Resolve: LWW] A wrote at 10:15:30, B wrote at 10:15:28
  → A wins (newer timestamp): "Alice-v2"
  → Merged clock: {A:2, B:1}

[Resolve: Siblings] Keeping both versions:
  → Sibling 1: "Alice-v2" {A:2}
  → Sibling 2: "Bob-edit" {A:1, B:1}
  → Client must resolve on next read

--- After resolution ---
[A] GET "user:1" → "Alice-v2"  clock={A:2, B:1}
[B] GET "user:1" → "Alice-v2"  clock={A:2, B:1}
[C] GET "user:1" → "Alice-v2"  clock={A:2, B:1}
  ✓ All replicas consistent
```

**Key learning:** Why vector clocks exist (Lamport clocks can't detect concurrency), how DynamoDB/Riak handle conflicts, partial ordering vs total ordering, the CAP theorem in practice.

---

## Problem 3: Gossip Protocol (Membership & Failure Detection)

**Difficulty:** Hard

**Concepts:** UDP, randomized protocols, failure detection, convergence, membership management

**Problem:**

Implement a gossip protocol for cluster membership and failure detection (like SWIM / HashiCorp Memberlist).

**Requirements:**
1. N nodes (goroutines) communicate via channels (simulating UDP).
2. Each node maintains a **membership list** with states: `Alive`, `Suspect`, `Dead`.
3. **Gossip cycle** (every 500ms): pick a random peer, send your membership list, merge their response.
4. **Failure detection (SWIM-style)**:
   - Ping a random member. If no ack within 200ms, ask K other members to ping it (indirect ping).
   - If still no ack, mark as `Suspect`. After 2 seconds of sustained suspicion, declare `Dead`.
5. **Join**: new node contacts a seed node and gets the full membership list via gossip.
6. **Leave**: node broadcasts a "leave" message that propagates via gossip.
7. Measure **convergence time**: how many gossip rounds until all N nodes agree on membership.

```go
type Node struct {
    id       string
    addr     string
    members  map[string]*Member
    mu       sync.RWMutex
    inbox    chan *GossipMessage
    sequence uint64
}

type Member struct {
    ID        string
    State     MemberState // Alive, Suspect, Dead
    Heartbeat uint64
    LastSeen  time.Time
}

type GossipMessage struct {
    From    string
    Type    MsgType // Ping, Ack, IndirectPing, Sync
    Members []*Member
}
```

**Expected Output:**
```
$ go run . -nodes=10
[Cluster] Starting 10-node gossip simulation

[t=0.0s] Node-0 started (seed)
[t=0.0s] Node-1 joining via Node-0
[t=0.0s] Node-2 joining via Node-0
... (all 10 nodes join)

[t=0.5s] Gossip round 1:
  Node-0 → Node-7: synced 10 members
  Node-3 → Node-1: synced 10 members
  Node-5 → Node-9: synced 10 members
  ...
[t=0.5s] Convergence: 8/10 nodes have full membership

[t=1.0s] Gossip round 2:
  Convergence: 10/10 nodes have full membership ✓ (2 rounds)

--- Simulating Node-4 crash ---
[t=3.0s] Node-7 pings Node-4 → no ack (200ms timeout)
[t=3.2s] Node-7 indirect ping via [Node-2, Node-8] → no ack
[t=3.2s] Node-7 marks Node-4 as SUSPECT
[t=3.5s] Gossip propagates: 4/9 nodes know Node-4 is suspect
[t=4.0s] Gossip propagates: 9/9 nodes know Node-4 is suspect
[t=5.2s] Suspicion timeout → Node-4 declared DEAD

[t=5.5s] All nodes agree: 9 alive, 1 dead

--- Node-4 recovers ---
[t=8.0s] Node-4 rejoining via Node-0
[t=8.5s] Gossip: 5/9 nodes know Node-4 is alive again
[t=9.0s] Gossip: 10/10 nodes have consistent membership ✓

========== GOSSIP STATS ==========
Total gossip rounds:    18
Messages sent:          180
Avg convergence time:   2 rounds (1.0s)
False positives:        0
Failure detection time: 2.2s
===================================
```

**Key learning:** Epidemic/gossip protocols, probabilistic failure detection, convergence properties, how HashiCorp Consul and Serf work, designing protocols that tolerate partial failures.

---

## Problem 4: Distributed Key-Value Store with Replication

**Difficulty:** Very Hard

**Concepts:** quorum reads/writes, replication, anti-entropy, hinted handoff, consistent hashing

**Problem:**

Build a distributed KV store (simplified Dynamo/Cassandra) across N nodes.

**Requirements:**
1. Data is partitioned using consistent hashing (reuse Problem 1).
2. Each key is replicated to N=3 nodes (the key's primary + 2 successors on the ring).
3. **Quorum protocol**: configurable W (write quorum) and R (read quorum). Default: W=2, R=2.
   - `Put(key, value)` succeeds if W replicas ack.
   - `Get(key)` reads from R replicas, returns the value with the highest version.
4. **Hinted handoff**: if a replica is down during write, the coordinator stores a "hint" and forwards it when the node recovers.
5. **Anti-entropy (read repair)**: during a read, if replicas disagree, update the stale replica.
6. Simulate with 5 nodes (goroutines), channels as network, and inject node failures.

```go
type KVNode struct {
    id         string
    store      map[string]*Record
    ring       *Ring // consistent hash ring
    peers      map[string]*KVNode
    hints      map[string][]*HintedWrite // target node -> pending writes
    mu         sync.RWMutex
}

type Record struct {
    Key     string
    Value   string
    Version VectorClock
}

type Coordinator struct {
    ring   *Ring
    nodes  map[string]*KVNode
    W, R   int
}
```

**Expected Output:**
```
$ go run .
[KV] Starting 5-node cluster (N=3, W=2, R=2)
[Ring] node-A, node-B, node-C, node-D, node-E

--- Write ---
[Coord] PUT "product:99" → primary=node-B, replicas=[node-C, node-D]
[node-B] Stored "product:99"="Widget" v={B:1} ✓
[node-C] Stored "product:99"="Widget" v={B:1} ✓
[node-D] Stored "product:99"="Widget" v={B:1} ✓
[Coord] Write quorum met (3/3 acks, needed 2) ✓

--- Read ---
[Coord] GET "product:99" → reading from [node-B, node-C]
[node-B] "product:99" = "Widget" v={B:1}
[node-C] "product:99" = "Widget" v={B:1}
[Coord] Read quorum met (2/2), versions agree ✓
→ "Widget"

--- Node failure + hinted handoff ---
[Simulate] node-C goes DOWN

[Coord] PUT "product:99" = "Gadget" → replicas=[node-B, node-C(down), node-D]
[node-B] Stored v={B:2} ✓
[node-C] UNREACHABLE — storing hint at coordinator
[node-D] Stored v={B:2} ✓
[Coord] Write quorum met (2/3 acks) ✓

[Simulate] node-C comes BACK UP
[Coord] Delivering 1 hinted write to node-C
[node-C] Applied hint: "product:99"="Gadget" v={B:2} ✓

--- Read repair ---
[Coord] GET "product:99" → reading from [node-B, node-D]
[node-B] v={B:2} "Gadget"
[node-D] v={B:2} "Gadget"
✓ Consistent

[Coord] GET (later read from different quorum) → [node-B, node-C]
[node-B] v={B:2} "Gadget"
[node-C] v={B:2} "Gadget"  (hinted write applied earlier)
✓ Consistent — read repair not needed

========== CLUSTER STATS ==========
Total PUTs:        50
Total GETs:        50
Quorum failures:    0
Hinted handoffs:    3
Read repairs:       1
====================================
```

**Key learning:** Quorum math (W+R > N ensures consistency), the mechanics behind DynamoDB, hinted handoff for availability during failures, anti-entropy for eventual consistency.

---

## Problem 5: Distributed Lock Manager (Lease-Based)

**Difficulty:** Hard

**Concepts:** TTL-based leases, fencing tokens, lock contention, split-brain protection

**Problem:**

Build a distributed lock service with TTL-based leases and fencing tokens (like a simplified ZooKeeper/etcd lock).

**Requirements:**
1. A central lock server manages named locks.
2. `Acquire(name, ttl) (Lease, error)` — acquire a lock with a TTL. Returns a unique fencing token (monotonically increasing).
3. `Renew(lease) error` — extend the TTL before it expires.
4. `Release(lease) error` — release the lock.
5. If a client dies (doesn't renew), the lock auto-expires after TTL and another client can acquire it.
6. **Fencing tokens**: each lock acquisition returns an incrementing token. Backend resources should reject operations with a stale token.
7. Simulate: 10 clients competing for 3 named locks. Show lock acquisition, renewal, expiry, and fencing token usage.
8. Handle split-brain: if a client thinks it still holds the lock after TTL expires, the fencing token prevents stale writes.

```go
type LockServer struct {
    mu     sync.Mutex
    locks  map[string]*Lock
    nextToken atomic.Uint64
}

type Lock struct {
    Name      string
    Owner     string
    Token     uint64
    ExpiresAt time.Time
    mu        sync.Mutex
}

type Lease struct {
    Name  string
    Token uint64
    TTL   time.Duration
}
```

**Expected Output:**
```
$ go run .
[LockServer] Started

[Client-1] Acquire("db-migration") → token=1, TTL=5s ✓
[Client-2] Acquire("db-migration") → LOCKED by Client-1 (waiting...)
[Client-3] Acquire("cache-rebuild") → token=2, TTL=5s ✓

[Client-1] Renew("db-migration") → extended to 10s ✓
[Client-2] Still waiting for "db-migration"...

[Client-1] Release("db-migration") ✓
[Client-2] Acquire("db-migration") → token=3, TTL=5s ✓

--- Simulating Client-4 crash (holds lock, never releases) ---
[Client-4] Acquire("worker-lock") → token=4, TTL=3s ✓
[Client-4] CRASHED (no renew, no release)
[t+3s] Lock "worker-lock" EXPIRED (Client-4 failed to renew)
[Client-5] Acquire("worker-lock") → token=5, TTL=5s ✓

--- Fencing token demo ---
[Client-4] (zombie) Tries to write to DB with token=4
[DB] REJECTED: stale token 4 < current token 5
[Client-5] Writes to DB with token=5 ✓

========== LOCK STATS ==========
Total acquisitions:  12
Total releases:       8
Expired leases:       2
Fencing rejections:   1
Avg wait time:       1.2s
Max contention:      3 clients on "db-migration"
=================================
```

**Key learning:** Lease-based distributed locking, fencing tokens for correctness, TTL expiry for fault tolerance, understanding why distributed locks are hard (see Martin Kleppmann's analysis), etcd/ZooKeeper locking patterns.

---

## Problem 6: Write-Ahead Log (WAL) with Log Shipping

**Difficulty:** Hard

**Concepts:** sequential I/O, fsync, log-structured storage, replication via log shipping

**Problem:**

Build a write-ahead log that supports log shipping to replicas for replication.

**Requirements:**
1. **WAL Writer**: append-only log with binary format: `[Length:4][CRC32:4][Sequence:8][Data:N]`.
2. All writes go to WAL first, then to an in-memory state machine.
3. `fsync` after each write (or batch of writes for performance).
4. **Log reader**: read entries sequentially from a WAL file for recovery.
5. **Truncation**: compact the WAL by removing entries already applied to a snapshot.
6. **Log shipping**: a background goroutine sends new WAL entries to replica nodes via TCP. Replicas apply entries to their own state machine.
7. Track: last applied sequence number, replication lag per replica.

```go
type WAL struct {
    mu        sync.Mutex
    file      *os.File
    sequence  uint64
    buffer    *bufio.Writer
    syncMode  SyncMode // SyncEvery, SyncBatch, SyncNone
}

type WALEntry struct {
    Sequence uint64
    CRC      uint32
    Data     []byte
}

type LogShipper struct {
    wal       *WAL
    replicas  []*Replica
    mu        sync.RWMutex
    positions map[string]uint64 // replica -> last shipped seq
}
```

**Expected Output:**
```
$ go run .
[WAL] Opened wal-0001.log (sync=batch)
[Shipper] Replicating to: replica-1(:9001), replica-2(:9002)

[WAL] Append seq=1 "SET user:1 Alice" (CRC=0xA3B2C1D0) ✓
[WAL] Append seq=2 "SET user:2 Bob"   (CRC=0xF1E2D3C4) ✓
[WAL] Batch fsync (2 entries, 187 bytes)
[Shipper] Sent seq=1..2 to replica-1 ✓ (lag=0)
[Shipper] Sent seq=1..2 to replica-2 ✓ (lag=0)

[WAL] Append seq=3..10 (8 entries, batch fsync)
[Shipper] Sent seq=3..10 to replica-1 ✓ (lag=0)
[Shipper] replica-2 unreachable — buffering (lag=8)

[WAL] Append seq=11..20 (10 entries)
[Shipper] replica-2 reconnected — catching up seq=3..20
[Shipper] replica-2 caught up ✓ (lag=0)

--- Recovery test ---
[WAL] Simulating crash...
[WAL] Recovering from wal-0001.log
[WAL] Read seq=1: CRC ✓
[WAL] Read seq=2: CRC ✓
...
[WAL] Read seq=20: CRC ✓
[WAL] Recovery complete — 20 entries, state machine rebuilt ✓

--- Truncation ---
[WAL] Snapshot at seq=15
[WAL] Truncating entries seq=1..15
[WAL] WAL size: 2.1KB → 0.7KB

========== WAL STATS ==========
Total entries:      20
Bytes written:      2.1KB
Fsyncs:             5 (batched)
Replication lag:
  replica-1: 0 entries
  replica-2: 0 entries
Recovery time:      3ms
================================
```

**Key learning:** Write-ahead logging (foundation of all databases), CRC for data integrity, fsync semantics, log-based replication (used by PostgreSQL, MySQL, Kafka), recovery from crash scenarios.

---

## Problem 7: Conflict-Free Replicated Data Types (CRDTs)

**Difficulty:** Very Hard

**Concepts:** eventual consistency, commutativity, idempotency, merge functions, replicated state

**Problem:**

Implement several CRDTs that can be independently updated on different replicas and merged without conflicts.

**Requirements:**

Implement these CRDT types:
1. **G-Counter** (grow-only counter): each node has its own counter; value = sum of all nodes.
2. **PN-Counter** (positive-negative counter): a G-Counter for increments + a G-Counter for decrements.
3. **G-Set** (grow-only set): union merge.
4. **OR-Set** (observed-remove set): add and remove with unique tags to resolve conflicts.
5. **LWW-Register** (last-writer-wins register): timestamped values, latest wins.

Each CRDT must implement:
```go
type CRDT[T any] interface {
    Value() T            // current value
    Merge(other CRDT[T]) // merge remote state (commutative, idempotent)
    State() []byte       // serialize for replication
}
```

Simulate 3 replicas updating independently, then merging. Prove commutativity (A.Merge(B) == B.Merge(A)) and idempotency (A.Merge(B).Merge(B) == A.Merge(B)).

```go
type GCounter struct {
    counts map[string]uint64 // nodeID -> count
}

type PNCounter struct {
    pos *GCounter
    neg *GCounter
}

type ORSet struct {
    elements map[string]map[string]bool // element -> {unique_tag -> exists}
    tombstones map[string]map[string]bool
}
```

**Expected Output:**
```
$ go run .

===== G-Counter =====
[node-A] Increment 3 times → local={A:3}
[node-B] Increment 5 times → local={B:5}
[node-C] Increment 2 times → local={C:2}

[Merge] A.Merge(B) → {A:3, B:5}          Value=8
[Merge] A.Merge(C) → {A:3, B:5, C:2}     Value=10
[Verify] Commutative: B.Merge(A) → {A:3, B:5} Value=8 ✓
[Verify] Idempotent: A.Merge(B).Merge(B) → Value=8 ✓

===== PN-Counter =====
[node-A] +5, -2 → local_value=3
[node-B] +3, -1 → local_value=2
[Merge] → Value=5 (total +8, total -3) ✓

===== OR-Set =====
[node-A] Add("apple"), Add("banana")
[node-B] Add("banana"), Remove("banana"), Add("cherry")

[Merge] A + B:
  "apple"  — in A, not removed anywhere → PRESENT ✓
  "banana" — added by A (tag-1), removed by B (tag-2), added by A (tag-1 survives) → PRESENT ✓
  "cherry" — added by B → PRESENT ✓
  Result: {"apple", "banana", "cherry"}

[Verify] Commutative: A.Merge(B) == B.Merge(A) ✓
[Verify] Idempotent ✓

===== LWW-Register =====
[node-A] Set("hello") at t=100
[node-B] Set("world") at t=200
[Merge] → "world" (t=200 > t=100) ✓
```

**Key learning:** CRDTs are the foundation of collaborative editing (Google Docs), distributed caches, and any system needing strong eventual consistency without coordination. Understanding merge semantics, commutativity, and idempotency.

---

## Concept Coverage Matrix

| Problem | Consistent Hashing | Vector Clocks | Gossip | Quorum | Leases | WAL/Replication | CRDTs |
|---------|-------------------|---------------|--------|--------|--------|-----------------|-------|
| 1. Consistent Hashing | ✓ | | | | | | |
| 2. Vector Clocks | | ✓ | | | | | |
| 3. Gossip Protocol | | | ✓ | | | | |
| 4. Distributed KV | ✓ | ✓ | | ✓ | | | |
| 5. Distributed Lock | | | | | ✓ | | |
| 6. WAL + Log Shipping | | | | | | ✓ | |
| 7. CRDTs | | | | | | | ✓ |

---

## Recommended Order

1. **Problem 1** — Consistent hashing (foundation for partitioning)
2. **Problem 2** — Vector clocks (foundation for ordering)
3. **Problem 6** — WAL (foundation for durability)
4. **Problem 5** — Distributed locks (coordination)
5. **Problem 3** — Gossip protocol (failure detection)
6. **Problem 7** — CRDTs (conflict-free replication)
7. **Problem 4** — Distributed KV (combines everything)
