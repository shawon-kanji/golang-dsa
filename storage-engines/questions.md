# Go Storage Engine Exercises (Hard)

Build database internals from scratch. These are the data structures and algorithms behind PostgreSQL, SQLite, RocksDB, and Redis.

---

## Problem 1: LSM Tree (Log-Structured Merge Tree)

**Difficulty:** Very Hard

**Concepts:** memtable, SSTables, compaction, bloom filters, write amplification

**Problem:**

Build an LSM tree storage engine — the architecture behind RocksDB, Cassandra, and LevelDB.

**Requirements:**
1. **MemTable**: in-memory sorted structure (use a skip list or red-black tree) that accepts writes. When it reaches a size threshold (e.g., 4MB), flush it to disk as an SSTable.
2. **SSTable (Sorted String Table)**: immutable, sorted on-disk file. Format:
   ```
   [DataBlock][DataBlock]...[IndexBlock][BloomFilter][Footer]
   ```
   Each data block: sorted key-value pairs. Index block: maps key ranges to block offsets.
3. **Write path**: Write → WAL → MemTable → (when full) → flush to L0 SSTable.
4. **Read path**: MemTable → L0 SSTables (newest first) → L1 → L2 → ... Return first match.
5. **Bloom filters**: per-SSTable bloom filter to skip SSTables that definitely don't contain a key.
6. **Compaction**: background goroutine merges overlapping SSTables from L0 into L1, deduplicating and removing tombstones.
7. **Tombstones**: `Delete(key)` writes a tombstone marker. Compaction removes it after it passes through all levels.

```go
type LSMTree struct {
    mu         sync.RWMutex
    memtable   *MemTable
    immutable  *MemTable        // being flushed
    levels     [][]*SSTable     // levels[0] = L0, etc.
    wal        *WAL
    compactCh  chan struct{}
    opts       Options
}

type MemTable struct {
    data    *SkipList
    size    int64
    maxSize int64
}

type SSTable struct {
    path       string
    index      []*IndexEntry  // key -> offset
    bloom      *BloomFilter
    minKey     []byte
    maxKey     []byte
    level      int
}
```

**Expected Output:**
```
$ go run .
[LSM] Initialized (memtable=4MB, L0 threshold=4 files)

--- Write phase (100,000 keys) ---
[WAL] Append seq=1..25000
[MemTable] Size: 4.0MB → flushing to L0
[Flush] L0/sst-0001.db written (25,000 entries, 4.1MB, bloom=48KB)
[WAL] Append seq=25001..50000
[Flush] L0/sst-0002.db written
[WAL] Append seq=50001..75000
[Flush] L0/sst-0003.db written
[WAL] Append seq=75001..100000
[Flush] L0/sst-0004.db written

[Compaction] L0 has 4 files → merging into L1
[Compaction] Merging sst-0001..0004 → L1/sst-0005.db (100,000 entries)
[Compaction] Removed 342 duplicate keys, 28 tombstones
[Compaction] L0: 4 files → 0 files | L1: 0 → 1 file

--- Read phase ---
[GET] "user:5000"
  MemTable: MISS
  Bloom(L0): empty
  Bloom(L1/sst-0005): MAYBE → scanning block 12
  Found: "user:5000" = "value-5000" (2 lookups, 0.3ms)

[GET] "nonexistent"
  MemTable: MISS
  Bloom(L1/sst-0005): DEFINITELY NOT → skipped
  Not found (1 lookup, 0.04ms)  ← bloom filter saved a disk read

--- Delete ---
[DELETE] "user:5000" → tombstone written to memtable
[GET] "user:5000" → tombstone found in memtable → NOT FOUND

========== LSM STATS ==========
Writes:         100,000
Reads:           10,000
Bloom filter saves: 3,412 (34.1% of reads skipped disk)
Compactions:     1
Write amplification: 2.1x
Disk usage:      12.3MB
================================
```

**Key learning:** How RocksDB/LevelDB work internally, write amplification trade-offs, bloom filters for read optimization, compaction strategies, the LSM vs B-tree trade-off.

---

## Problem 2: B+ Tree with Disk Persistence

**Difficulty:** Very Hard

**Concepts:** B+ tree, page-based storage, buffer pool, node splitting/merging, disk I/O

**Problem:**

Build a B+ tree that stores data on disk using fixed-size pages — the index structure behind PostgreSQL and MySQL (InnoDB).

**Requirements:**
1. **Page-based storage**: each B+ tree node is stored in a fixed-size page (4KB). Nodes are read/written by page number.
2. **Internal nodes**: store keys + child page pointers. No values.
3. **Leaf nodes**: store key-value pairs + a "next leaf" pointer for range scans.
4. `Insert(key, value)` — insert with node splitting when full.
5. `Get(key) (value, bool)` — point lookup traversing from root to leaf.
6. `Range(startKey, endKey) []KeyValue` — range scan using leaf node linked list.
7. `Delete(key)` — delete with node merging/redistribution when underflow.
8. **Buffer pool**: cache recently accessed pages in memory (LRU eviction). Track cache hit ratio.
9. **Pager**: reads/writes 4KB pages to a file. Pages are addressed by page number (offset = pageNum × 4096).

```go
type BPlusTree struct {
    pager      *Pager
    bufferPool *BufferPool
    rootPage   uint32
    order      int // max keys per node
    mu         sync.RWMutex
}

type Node struct {
    pageNum   uint32
    isLeaf    bool
    keys      [][]byte
    values    [][]byte   // leaf only
    children  []uint32   // internal only (page numbers)
    nextLeaf  uint32     // leaf only
}

type Pager struct {
    file     *os.File
    pageSize int
    numPages uint32
}

type BufferPool struct {
    pages   map[uint32]*CachedPage
    maxSize int
    lru     *list.List
    mu      sync.RWMutex
    hits    int64
    misses  int64
}
```

**Expected Output:**
```
$ go run .
[B+Tree] Initialized (order=128, pageSize=4KB, bufferPool=256 pages)

--- Insert 50,000 keys ---
[Insert] key="aaa:00001" → leaf page 2
[Insert] key="aaa:00002" → leaf page 2
...
[Split] Leaf page 2 full → split into pages 2, 5
[Split] Internal page 1 full → split, new root at page 8
...
[Insert] Complete: 50,000 keys, tree height=3

--- Point lookup ---
[GET] "user:25000"
  Root (page 1) → internal (page 14) → leaf (page 2847)
  Found: "user:25000" = "value-25000"
  Pages read: 3 | Buffer pool: 2 hits, 1 miss
  Latency: 0.08ms

--- Range scan ---
[RANGE] "user:10000" to "user:10100"
  Seek to leaf → scan 101 entries across 2 leaf pages
  Results: 101 key-value pairs
  Latency: 0.2ms

--- Delete ---
[DELETE] "user:25000"
  Leaf page 2847: removed key
  Underflow → redistributing from sibling page 2848
  Latency: 0.1ms

========== B+TREE STATS ==========
Total keys:       49,999
Tree height:      3
Total pages:      412
Buffer pool hits: 147,823 (94.2%)
Buffer pool misses: 9,112 (5.8%)
Avg lookup I/Os:  1.2 pages (with buffer pool)
Range scan speed: 505 entries/page
===================================
```

**Key learning:** How database indexes actually work on disk, page-based storage, buffer pool management (the #1 optimization in databases), node splitting/merging algorithms, why B+ trees are preferred over B-trees for databases (leaf linked list → efficient range scans).

---

## Problem 3: Bloom Filter with Counting Support

**Difficulty:** Moderate–Hard

**Concepts:** probabilistic data structures, hash functions, false positive rates, bit manipulation

**Problem:**

Build a bloom filter and a counting bloom filter, with configurable false positive rate.

**Requirements:**
1. **Standard Bloom Filter**:
   - `New(expectedItems int, fpRate float64)` — auto-calculate optimal bit array size and hash count.
   - `Add(item []byte)` — set K bits using K hash functions.
   - `MayContain(item []byte) bool` — check if item might be present.
   - Formulas: `m = -n*ln(p) / (ln2)^2`, `k = (m/n) * ln2`.

2. **Counting Bloom Filter**:
   - Same API + `Remove(item []byte)` — decrement counters instead of bits.
   - Uses 4-bit counters instead of single bits (still compact).

3. **Benchmark**: insert 1M items, test false positive rate with 1M items NOT in the set. Compare actual vs theoretical FP rate.

4. **Serialization**: `Marshal() []byte` and `Unmarshal([]byte)` for saving/loading the filter.

```go
type BloomFilter struct {
    bits      []uint64 // bit array packed into uint64s
    numBits   uint
    numHashes uint
    count     uint
}

type CountingBloom struct {
    counters  []uint8 // 4-bit counters (packed 2 per byte)
    numSlots  uint
    numHashes uint
    count     uint
}

func NewBloom(n int, fpRate float64) *BloomFilter
func (bf *BloomFilter) Add(item []byte)
func (bf *BloomFilter) MayContain(item []byte) bool
```

**Expected Output:**
```
$ go run .

===== Standard Bloom Filter =====
Config: 1,000,000 expected items, 1% target FP rate
  Optimal bits: 9,585,059 (1.14 MB)
  Optimal hashes: 7

[Phase 1] Inserting 1,000,000 items...
[Phase 2] Testing membership (should all be true):
  True positives: 1,000,000 / 1,000,000 (100%) ✓

[Phase 3] Testing 1,000,000 items NOT in set:
  False positives: 10,127 / 1,000,000 (1.01%)
  Target FP rate:  1.00%
  Actual FP rate:  1.01% ← very close ✓

===== Counting Bloom Filter =====
[Insert] Added 100,000 items
[Remove] Removed 10,000 items
[Verify] Removed items: 0 / 10,000 found (0%) ✓
[Verify] Remaining items: 90,000 / 90,000 found (100%) ✓

===== Serialization =====
[Bloom] Serialized: 1.14 MB
[Bloom] Deserialized and verified: all 1M items still present ✓

===== Memory Comparison =====
HashSet (1M items):    ~48 MB
Bloom (1% FP):          1.14 MB (42x smaller)
Bloom (0.1% FP):        1.71 MB (28x smaller)
Bloom (0.01% FP):       2.28 MB (21x smaller)
```

**Key learning:** Probabilistic data structures, optimal parameter calculation, bit manipulation, how databases use bloom filters (LSM trees, Cassandra, HBase) to skip unnecessary I/O.

---

## Problem 4: Skip List (Concurrent)

**Difficulty:** Hard

**Concepts:** probabilistic data structure, lock-free or fine-grained locking, ordered map

**Problem:**

Build a concurrent skip list — the data structure used as the memtable in many LSM tree implementations (RocksDB, LevelDB).

**Requirements:**
1. **Operations**: `Insert(key, value)`, `Get(key) (value, bool)`, `Delete(key) bool`, `Range(start, end) []KV`.
2. **Probabilistic levels**: each node has a random height (p=0.5 for each additional level, max 32 levels).
3. **Concurrency**: safe for concurrent reads and writes. Implement one of:
   - Fine-grained locking (lock per node).
   - Lock-free using `sync/atomic` and CAS (compare-and-swap).
4. **Iterator**: `NewIterator()` returns an iterator for sequential access.
5. Benchmark: compare throughput vs `sync.Map` and `sync.RWMutex + map` for mixed read/write workloads.

```go
type SkipList struct {
    head    *skipNode
    height  atomic.Int32
    length  atomic.Int64
    maxLevel int
}

type skipNode struct {
    key   []byte
    value []byte
    next  []*atomic.Pointer[skipNode] // one per level
    marked atomic.Bool // for lock-free delete
}

type Iterator struct {
    current *skipNode
    list    *SkipList
}
```

**Expected Output:**
```
$ go run .
[SkipList] Initialized (maxLevel=32, p=0.5)

--- Basic Operations ---
[Insert] "alpha"=1 (height=1) ✓
[Insert] "beta"=2  (height=3) ✓
[Insert] "gamma"=3 (height=1) ✓
[Insert] "delta"=4 (height=2) ✓

Level 3: head -----> beta -----> nil
Level 2: head -----> beta -> delta -> nil
Level 1: head -> alpha -> beta -> delta -> gamma -> nil

[Get] "beta" → 2 ✓ (traversed 1 level)
[Get] "missing" → not found (traversed 3 levels)
[Range] "alpha".."delta" → [alpha=1, beta=2, delta=4]
[Delete] "beta" ✓

--- Concurrent Benchmark (10s) ---
  100 goroutines (80% reads, 20% writes)

  SkipList:
    Ops/sec:     2,450,000
    Read  lat:   0.12μs avg
    Write lat:   0.34μs avg

  sync.Map:
    Ops/sec:     3,100,000 (faster for pure KV)
    No range scan support ✗

  RWMutex+map:
    Ops/sec:     1,800,000
    Lock contention visible at high concurrency

  Conclusion: SkipList wins for ordered access + range scans
```

**Key learning:** Skip list internals (used in Redis sorted sets, LevelDB memtable), probabilistic balancing, fine-grained concurrent data structures, lock-free programming patterns.

---

## Problem 5: Write-Optimized Buffer Pool Manager

**Difficulty:** Very Hard

**Concepts:** page management, LRU-K eviction, dirty page tracking, write-back vs write-through, concurrency

**Problem:**

Build a buffer pool manager that caches disk pages in memory — the core component of any disk-based database.

**Requirements:**
1. Fixed number of frames (e.g., 256). Each frame holds one 4KB page.
2. `FetchPage(pageID) (*Page, error)` — load page from disk into a buffer frame. If no free frame, evict using LRU-K (K=2).
3. **Pin counting**: pages have a pin count. Pinned pages cannot be evicted. `Unpin(pageID)` decrements the count.
4. **Dirty tracking**: mark pages dirty on write. Dirty pages must be flushed to disk before eviction.
5. **LRU-K eviction**: track the K-th most recent access timestamp. Evict the page with the oldest K-th access (better than LRU for scan resistance).
6. `FlushPage(pageID)` — write dirty page to disk.
7. `FlushAll()` — flush all dirty pages (for checkpoint).
8. Thread-safe: multiple goroutines can pin/unpin/fetch pages concurrently.

```go
type BufferPoolManager struct {
    mu        sync.Mutex
    pool      []*Frame
    pageTable map[uint32]int // pageID -> frameIndex
    disk      *DiskManager
    replacer  *LRUKReplacer
    freeList  []int
}

type Frame struct {
    page     *Page
    pageID   uint32
    dirty    bool
    pinCount int
}

type LRUKReplacer struct {
    k          int
    history    map[uint32][]time.Time // pageID -> access timestamps
    evictable  map[uint32]bool
}

type Page struct {
    data [4096]byte
}
```

**Expected Output:**
```
$ go run .
[BufferPool] 256 frames, LRU-K (K=2), page_size=4KB

--- Sequential scan (512 pages through 256 frames) ---
[Fetch] Page 0 → frame 0 (disk read) ✓
[Fetch] Page 1 → frame 1 (disk read) ✓
...
[Fetch] Page 255 → frame 255 (disk read) ✓
[Fetch] Page 256 → frame 0 (evict page 0, disk read)

After sequential scan:
  Cache hits:   0 / 512 (0%) — expected for cold start
  Pages evicted: 256

--- Random access pattern (10,000 reads, 80% from hot set of 50 pages) ---
[Fetch] Page 42 → frame 12 (cache HIT) ✓
[Fetch] Page 301 → frame 47 (cache MISS → evict page 489)
...

After random access:
  Cache hits: 7,847 / 10,000 (78.5%)
  LRU-K kept hot pages, evicted scan pages ✓

--- Dirty page handling ---
[Write] Page 42 → marked dirty
[Write] Page 42 → already in buffer, just mark dirty
[Evict] Page 42 → dirty! flushing to disk first...
[Evict] Page 42 → flushed and evicted ✓

--- Checkpoint ---
[FlushAll] 47 dirty pages flushed to disk (188KB written)

========== BUFFER POOL STATS ==========
Total fetches:     10,512
Cache hits:         7,847 (74.6%)
Cache misses:       2,665 (25.4%)
Pages evicted:      2,665
Dirty flushes:      1,234
Disk reads:         2,665
Disk writes:        1,234
Avg pin duration:   0.05ms
========================================
```

**Key learning:** Buffer pool management (the single most important DB optimization), LRU-K eviction policy (better than plain LRU), dirty page tracking, pin counting, understanding how PostgreSQL and MySQL manage memory.

---

## Problem 6: Transaction Manager (MVCC)

**Difficulty:** Very Hard

**Concepts:** multi-version concurrency control, snapshots, read/write conflicts, serialization

**Problem:**

Build a simplified MVCC transaction manager — the concurrency control mechanism behind PostgreSQL, MySQL InnoDB, and CockroachDB.

**Requirements:**
1. Each transaction gets a unique, monotonically increasing `txnID`.
2. **Writes** create a new version of the key tagged with `txnID`. Old versions are kept.
3. **Reads** see the latest version that was committed before the transaction's start time (snapshot isolation).
4. **Commit**: mark all the transaction's writes as committed.
5. **Abort/Rollback**: discard all the transaction's writes.
6. **Write-write conflict detection**: if two concurrent transactions write the same key, the second one to commit gets an error (`ErrWriteConflict`).
7. **Garbage collection**: background goroutine removes versions that are no longer visible to any active transaction.
8. Demonstrate: concurrent transactions with reads, writes, conflicts, and rollbacks.

```go
type MVCCStore struct {
    mu          sync.RWMutex
    versions    map[string][]*Version // key -> versions (newest first)
    activeTxns  map[uint64]*Transaction
    nextTxnID   atomic.Uint64
    gcInterval  time.Duration
}

type Version struct {
    TxnID     uint64
    Value     []byte
    Committed bool
    Deleted   bool
    CreatedAt time.Time
}

type Transaction struct {
    ID        uint64
    StartID   uint64    // snapshot: see commits before this ID
    Writes    map[string]*Version
    Status    TxnStatus // Active, Committed, Aborted
}
```

**Expected Output:**
```
$ go run .
[MVCC] Initialized

--- Snapshot Isolation ---
[Txn-1] BEGIN (snapshot at txnID=1)
[Txn-1] PUT "balance:alice" = 1000
[Txn-1] COMMIT ✓

[Txn-2] BEGIN (snapshot at txnID=2)
[Txn-3] BEGIN (snapshot at txnID=3)

[Txn-2] READ "balance:alice" → 1000 (sees txn-1's commit)
[Txn-2] PUT "balance:alice" = 800 (debit 200)

[Txn-3] READ "balance:alice" → 1000 (snapshot: still sees original!)
[Txn-3] PUT "balance:alice" = 1200 (credit 200)

[Txn-2] COMMIT ✓ (writes version with txnID=2)
[Txn-3] COMMIT → ERROR: ErrWriteConflict on "balance:alice"
  (txn-2 already committed a write to this key after txn-3 started)

[Txn-3] ABORT — all writes discarded

--- After conflict resolution ---
[Txn-4] BEGIN
[Txn-4] READ "balance:alice" → 800 (sees txn-2's commit)
[Txn-4] PUT "balance:alice" = 1000 (credit 200, retry)
[Txn-4] COMMIT ✓

--- Version history for "balance:alice" ---
  txnID=4: 1000 (committed)  ← current
  txnID=2: 800  (committed)  ← visible to txns started before 4
  txnID=1: 1000 (committed)  ← old, GC candidate

--- Garbage Collection ---
[GC] Active txns: [5, 6] (oldest snapshot = 5)
[GC] Removed 1 old version of "balance:alice" (txnID=1, not visible to anyone)
[GC] Cleaned 47 obsolete versions across 23 keys

========== MVCC STATS ==========
Transactions:       20
Committed:          17
Aborted (conflict):  3
Versions created:   45
Versions GC'd:      12
Avg versions/key:   1.8
=================================
```

**Key learning:** How PostgreSQL/MySQL implement concurrent transactions, snapshot isolation vs serializable, write-write conflict detection, version garbage collection, why MVCC is preferred over locking for read-heavy workloads.

---

## Concept Coverage Matrix

| Problem | Disk I/O | In-Memory | Compaction | Bloom Filter | B-Tree | Buffer Pool | MVCC | Concurrency |
|---------|----------|-----------|------------|--------------|--------|-------------|------|-------------|
| 1. LSM Tree | ✓ | ✓ | ✓ | ✓ | | | | ✓ |
| 2. B+ Tree | ✓ | | | | ✓ | ✓ | | ✓ |
| 3. Bloom Filter | | ✓ | | ✓ | | | | |
| 4. Skip List | | ✓ | | | | | | ✓ |
| 5. Buffer Pool | ✓ | ✓ | | | | ✓ | | ✓ |
| 6. MVCC | | ✓ | | | | | ✓ | ✓ |

---

## Recommended Order

1. **Problem 3** — Bloom filter (standalone, used by later problems)
2. **Problem 4** — Skip list (memtable for LSM tree)
3. **Problem 5** — Buffer pool (page cache for B+ tree)
4. **Problem 2** — B+ tree (uses buffer pool)
5. **Problem 1** — LSM tree (uses skip list, bloom filter, WAL)
6. **Problem 6** — MVCC (transaction layer on top of any storage)
