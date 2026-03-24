# Go Networking Exercises (Hard)

Build real network infrastructure from scratch. Each problem teaches protocol design, connection management, and the `net` package at a deep level.

---

## Problem 1: TCP Chat Server with Rooms

**Difficulty:** Hard

**Concepts:** net.Listener, goroutine-per-connection, buffered I/O, shared state, graceful disconnect

**Problem:**

Build a multi-room TCP chat server that clients connect to via `telnet` or `nc`.

**Requirements:**
1. Server listens on a configurable port. Each client connection is handled by its own goroutine.
2. On connect, the client is prompted for a username. Duplicate usernames are rejected.
3. Commands:
   - `/join <room>` — join/create a room (leave current room first)
   - `/list` — list all rooms and user counts
   - `/who` — list users in current room
   - `/msg <user> <text>` — private DM
   - `/quit` — disconnect
4. Messages in a room are broadcast to all members except the sender.
5. Handle abrupt disconnects (broken pipe) — clean up the user from all data structures.
6. Server tracks metrics: total connections, active connections, messages sent.

```go
type Server struct {
    mu       sync.RWMutex
    rooms    map[string]*Room
    users    map[string]*Client
    listener net.Listener
}

type Room struct {
    name    string
    members map[string]*Client
}

type Client struct {
    conn     net.Conn
    username string
    room     string
    writer   *bufio.Writer
}
```

**Expected Output:**
```
$ go run . -port=9090
[Server] Listening on :9090

[Server] New connection from 192.168.1.10
[Server] User "alice" joined
[alice] /join general
[Server] Room "general" created
[alice] joined #general

[Server] New connection from 192.168.1.11
[Server] User "bob" joined
[bob] /join general
[bob] joined #general

[alice] hello everyone!
  -> [bob] alice: hello everyone!

[bob] /msg alice hey, private msg
  -> [alice] [DM from bob]: hey, private msg

[alice] /who
  -> Users in #general: alice, bob

[Server] Connection lost: bob (broken pipe)
  -> [alice] ** bob has left #general

[Server] Stats: total_conns=2 active=1 messages=3
```

**Key learning:** goroutine-per-connection model, mutex-protected shared maps, buffered I/O for network writes, handling partial reads/writes, cleanup on disconnect.

---

## Problem 2: HTTP/1.1 Server From Scratch

**Difficulty:** Hard

**Concepts:** TCP, HTTP protocol parsing, request routing, connection keep-alive, timeouts

**Problem:**

Build a minimal HTTP/1.1 server using only `net` (no `net/http`).

**Requirements:**
1. Parse raw HTTP requests from TCP connections: method, path, headers, body.
2. Support `GET`, `POST`, `PUT`, `DELETE` methods.
3. Implement a simple router: `Register(method, path, handler)`.
4. Support `Connection: keep-alive` — reuse the same TCP connection for multiple requests.
5. Set read/write timeouts (`net.Conn.SetDeadline`) — close idle connections after 30 seconds.
6. Support `Content-Length` for request bodies and chunked responses.
7. Return proper status codes: 200, 201, 400, 404, 405, 500.
8. Add a `/stats` endpoint that returns JSON with: active connections, total requests, avg response time.

```go
type HTTPServer struct {
    listener net.Listener
    routes   map[string]map[string]HandlerFunc // method -> path -> handler
    mu       sync.RWMutex
    stats    *ServerStats
}

type Request struct {
    Method  string
    Path    string
    Headers map[string]string
    Body    []byte
    Query   map[string]string
}

type Response struct {
    StatusCode int
    Headers    map[string]string
    Body       []byte
}

type HandlerFunc func(req *Request) *Response
```

**Expected Output:**
```
$ go run . -port=8080
[HTTP] Listening on :8080

[HTTP] 192.168.1.10 -> GET /hello HTTP/1.1
[HTTP] 192.168.1.10 <- 200 OK (12 bytes, 0.3ms)

[HTTP] 192.168.1.10 -> POST /users HTTP/1.1 (Content-Length: 45)
[HTTP] 192.168.1.10 <- 201 Created (89 bytes, 1.2ms)

[HTTP] 192.168.1.10 -> GET /nonexistent HTTP/1.1
[HTTP] 192.168.1.10 <- 404 Not Found

[HTTP] 192.168.1.10 -> DELETE /users/42 HTTP/1.1
[HTTP] 192.168.1.10 <- 200 OK

[HTTP] 192.168.1.10 connection idle for 30s — closing
[HTTP] Stats: active_conns=0 total_reqs=4 avg_latency=0.7ms

# Test with curl:
$ curl http://localhost:8080/stats
{"active_connections": 1, "total_requests": 5, "avg_response_ms": 0.8}
```

**Key learning:** Deep understanding of HTTP/1.1 at the wire level, connection management, keep-alive, timeout handling, state machines for request parsing.

---

## Problem 3: Reverse Proxy / Load Balancer

**Difficulty:** Hard

**Concepts:** net.Dial, io.Copy, connection pooling, health checks, round-robin/least-connections

**Problem:**

Build a Layer 7 reverse proxy that distributes HTTP requests across multiple backend servers.

**Requirements:**
1. Accept `GET`/`POST` requests on a frontend port and forward them to backend servers.
2. Load balancing strategies (implement at least 2, switchable via config):
   - **Round-robin**: cycle through backends.
   - **Least-connections**: pick the backend with fewest active connections.
   - **Weighted round-robin**: backends have different weights.
3. **Health checking**: background goroutine pings each backend every 5 seconds. Unhealthy backends are removed from rotation and re-added when healthy.
4. **Connection pooling**: maintain a pool of reusable TCP connections to backends (avoid dialing for every request).
5. **Request/response streaming**: use `io.Copy` to stream large responses without buffering the entire body in memory.
6. Add `X-Forwarded-For` and `X-Request-ID` headers to proxied requests.
7. Timeout: if a backend doesn't respond within 5 seconds, try the next one.

```go
type Proxy struct {
    backends    []*Backend
    strategy    Strategy
    mu          sync.RWMutex
    connPool    map[string][]*pooledConn
    healthTicker *time.Ticker
}

type Backend struct {
    addr        string
    healthy     atomic.Bool
    activeConns atomic.Int64
    weight      int
}

type Strategy interface {
    Pick(backends []*Backend) *Backend
}
```

**Expected Output:**
```
$ go run . -listen=:8080 -backends=localhost:9001,localhost:9002,localhost:9003
[Proxy] Listening on :8080
[Proxy] Backends: [:9001 (w=1), :9002 (w=1), :9003 (w=1)]
[Health] :9001=OK :9002=OK :9003=OK

[Proxy] GET /api/users -> :9001 (round-robin) -> 200 OK (34ms)
[Proxy] GET /api/users -> :9002 (round-robin) -> 200 OK (28ms)
[Proxy] GET /api/users -> :9003 (round-robin) -> 200 OK (31ms)
[Proxy] POST /api/orders -> :9001 (round-robin) -> 201 Created (45ms)

[Health] :9002=UNHEALTHY (connection refused)
[Proxy] Removed :9002 from rotation

[Proxy] GET /api/users -> :9001 (round-robin) -> 200 OK (29ms)
[Proxy] GET /api/users -> :9003 (round-robin) -> 200 OK (33ms)

[Health] :9002=OK (recovered)
[Proxy] Added :9002 back to rotation

[Proxy] GET /api/slow -> :9001 TIMEOUT (5s) -> retrying :9003 -> 200 OK (120ms)

[Proxy] Stats: total=7 success=7 failed=0 avg_latency=45ms
  :9001 -> 3 reqs, 2 active conns, pool=3
  :9002 -> 1 req,  0 active conns, pool=1
  :9003 -> 3 reqs, 1 active conn,  pool=2
```

**Key learning:** Connection pooling, health check goroutines, load balancing algorithms, request forwarding, timeout + retry patterns used in real proxies (Nginx, Envoy).

---

## Problem 4: TCP Connection Multiplexer (Protocol Design)

**Difficulty:** Very Hard

**Concepts:** custom binary protocol, framing, multiplexing, net.Conn wrapping, io.Reader/io.Writer

**Problem:**

Build a multiplexer that runs multiple logical streams over a single TCP connection (like HTTP/2 or yamux).

**Requirements:**
1. **Protocol**: Design a simple binary frame format:
   ```
   [StreamID: 4 bytes][Type: 1 byte][Length: 4 bytes][Payload: N bytes]
   ```
   Types: `DATA`, `SYN` (open stream), `FIN` (close stream), `RST` (abort), `PING`/`PONG`.

2. **Server side**: Accept a single TCP connection. The client can open multiple streams over it. Each stream gets its own goroutine on the server.

3. **Client side**: `OpenStream() (Stream, error)` returns a `net.Conn`-like object. Multiple streams share one TCP conn.

4. **Flow control**: Each stream has a receive window. Sender pauses if window is full.

5. `Stream` implements `io.ReadWriteCloser` so it feels like a regular connection.

6. Demonstrate: client opens 5 streams over 1 TCP connection, sends different data on each concurrently, server echoes back on each stream independently.

```go
type Muxer struct {
    conn     net.Conn
    streams  map[uint32]*Stream
    mu       sync.RWMutex
    nextID   atomic.Uint32
    incoming chan *Stream
}

type Stream struct {
    id       uint32
    muxer    *Muxer
    readBuf  *bytes.Buffer
    readCond *sync.Cond
    closed   atomic.Bool
}

type Frame struct {
    StreamID uint32
    Type     FrameType
    Payload  []byte
}

// Stream implements io.ReadWriteCloser
func (s *Stream) Read(p []byte) (int, error)
func (s *Stream) Write(p []byte) (int, error)
func (s *Stream) Close() error
```

**Expected Output:**
```
$ go run . -mode=server -addr=:7000
[Muxer-Server] Listening on :7000
[Muxer-Server] Accepted connection from 127.0.0.1:52341

[Stream-1] Opened (SYN received)
[Stream-2] Opened (SYN received)
[Stream-3] Opened (SYN received)
[Stream-4] Opened (SYN received)
[Stream-5] Opened (SYN received)

[Stream-1] Received 1024 bytes -> echoing back
[Stream-3] Received 512 bytes -> echoing back
[Stream-2] Received 2048 bytes -> echoing back
[Stream-5] Received 256 bytes -> echoing back
[Stream-4] Received 4096 bytes -> echoing back
[Stream-1] FIN received -> closing
[Stream-3] FIN received -> closing

[Muxer-Server] Stats: 1 TCP conn, 5 streams (2 closed, 3 active)

--- client side ---
$ go run . -mode=client -addr=localhost:7000
[Muxer-Client] Connected to localhost:7000

[Client] Opening 5 streams over 1 TCP connection...
[Stream-1] Sent 1024 bytes, received 1024 bytes (echo OK ✓)  12ms
[Stream-2] Sent 2048 bytes, received 2048 bytes (echo OK ✓)  15ms
[Stream-3] Sent 512 bytes,  received 512 bytes  (echo OK ✓)  9ms
[Stream-4] Sent 4096 bytes, received 4096 bytes (echo OK ✓)  18ms
[Stream-5] Sent 256 bytes,  received 256 bytes  (echo OK ✓)  7ms

All 5 streams completed over 1 TCP connection
Total data: 7936 bytes sent, 7936 bytes received
```

**Key learning:** Binary protocol design and framing, multiplexing streams over a single connection, sync.Cond for stream-level blocking reads, understanding how HTTP/2 and gRPC work under the hood.

---

## Problem 5: DNS Resolver

**Difficulty:** Hard

**Concepts:** UDP, binary encoding/decoding, protocol implementation, caching, net.UDPConn

**Problem:**

Build a DNS resolver that can resolve domain names by querying upstream DNS servers.

**Requirements:**
1. Construct DNS query packets (binary format per RFC 1035): header + question section.
2. Send queries via UDP to an upstream server (e.g., `8.8.8.8:53`).
3. Parse DNS response packets: extract A records, CNAME records, MX records.
4. Implement a **recursive resolver**: if the answer is a CNAME, follow it. If it's a referral, query the nameserver.
5. Add a **local cache** with TTL from the DNS response. Cache hit should skip the network call.
6. Support concurrent lookups: multiple goroutines resolving different domains simultaneously.
7. Timeout: if upstream doesn't respond within 2 seconds, retry once, then return error.

```go
type Resolver struct {
    upstream string
    cache    *DNSCache
    mu       sync.RWMutex
    conn     *net.UDPConn
}

type DNSCache struct {
    mu      sync.RWMutex
    entries map[string]*CacheEntry
}

type CacheEntry struct {
    records   []DNSRecord
    expiresAt time.Time
}

type DNSRecord struct {
    Name  string
    Type  uint16 // A=1, CNAME=5, MX=15
    TTL   uint32
    Value string // IP for A, domain for CNAME
}
```

**Expected Output:**
```
$ go run .
[Resolver] Using upstream 8.8.8.8:53

> resolve google.com A
[DNS] Query: google.com A -> 8.8.8.8:53 (UDP)
[DNS] Response (28ms):
  google.com.  300  IN  A  142.250.80.46
  (cached, TTL=300s)

> resolve google.com A
[DNS] Cache HIT: google.com A (expires in 287s)
  google.com.  287  IN  A  142.250.80.46

> resolve www.github.com A
[DNS] Query: www.github.com A -> 8.8.8.8:53
[DNS] Response: CNAME -> github.com
[DNS] Following CNAME: github.com A -> 8.8.8.8:53
[DNS] Response (42ms):
  www.github.com.  3600  IN  CNAME  github.com.
  github.com.      60    IN  A      140.82.121.4

> resolve gmail.com MX
[DNS] Query: gmail.com MX -> 8.8.8.8:53
[DNS] Response (35ms):
  gmail.com.  3600  IN  MX  5  gmail-smtp-in.l.google.com.
  gmail.com.  3600  IN  MX  10 alt1.gmail-smtp-in.l.google.com.

[Resolver] Stats: queries=4 cache_hits=1 cache_misses=3 avg_latency=35ms
```

**Key learning:** UDP networking, binary protocol encoding/decoding, DNS internals, concurrent cache with TTL, real protocol implementation experience.

---

## Problem 6: Port Scanner with Service Detection

**Difficulty:** Moderate–Hard

**Concepts:** net.DialTimeout, goroutine pool, rate limiting, banner grabbing, TCP handshake

**Problem:**

Build a concurrent port scanner that detects open ports and identifies running services.

**Requirements:**
1. Scan a target host across a port range (e.g., 1–1024 or specific ports).
2. Use a goroutine pool (configurable workers, e.g., 100) to scan ports concurrently.
3. For each open port, attempt **banner grabbing**: connect, read the first response bytes, and match against known service signatures (HTTP, SSH, SMTP, FTP, MySQL, Redis, etc.).
4. Configurable timeout per port (default: 500ms).
5. Rate limiting: max N connections per second to avoid overwhelming the target.
6. Output results sorted by port number with service name.
7. Track scan progress: show percentage and estimated remaining time.

```go
type Scanner struct {
    target      string
    ports       []int
    workers     int
    timeout     time.Duration
    rateLimit   time.Duration
    results     chan ScanResult
}

type ScanResult struct {
    Port    int
    Open    bool
    Service string // detected service
    Banner  string // raw banner text
    Latency time.Duration
}

var serviceSignatures = map[string]string{
    "SSH":   "SSH-",
    "HTTP":  "HTTP/",
    "SMTP":  "220 ",
    "FTP":   "220 ",
    "MySQL": "\x00",
    "Redis": "+PONG",
}
```

**Expected Output:**
```
$ go run . -target=scanme.nmap.org -ports=1-1024 -workers=100
[Scanner] Scanning scanme.nmap.org (45.33.32.156)
[Scanner] Ports: 1-1024 | Workers: 100 | Timeout: 500ms

Scanning... ██████████████████████████████░░░░░░  72% (738/1024) ETA: 3s

PORT     STATE   SERVICE   BANNER                        LATENCY
22/tcp   open    SSH       SSH-2.0-OpenSSH_6.6.1p1       45ms
80/tcp   open    HTTP      HTTP/1.1 200 OK               38ms
443/tcp  open    HTTPS     (TLS handshake)                52ms
9929/tcp open    Unknown   (no banner)                    67ms

========== SCAN REPORT ==========
Host:         scanme.nmap.org (45.33.32.156)
Ports scanned: 1024
Open ports:    4
Closed/filtered: 1020
Scan duration: 11.2s
==================================
```

**Key learning:** High-concurrency TCP connections, banner grabbing, rate-limited goroutine pools, timeout handling, practical network reconnaissance.

---

## Problem 7: Peer-to-Peer File Transfer

**Difficulty:** Very Hard

**Concepts:** net.Listener, concurrent connections, chunked transfer, integrity verification, NAT traversal awareness

**Problem:**

Build a P2P file transfer system where two peers can send files directly to each other.

**Requirements:**
1. **Sender mode**: Listen on a port, split the file into chunks (e.g., 64KB), serve chunks to the receiver.
2. **Receiver mode**: Connect to sender, request chunks, write them to disk in order.
3. **Parallel download**: Receiver uses N goroutines to fetch different chunks simultaneously (like BitTorrent).
4. **Integrity**: SHA-256 checksum per chunk + full file checksum. Receiver verifies each chunk and retries corrupted ones.
5. **Progress reporting**: Show transfer speed (MB/s), percentage, ETA.
6. **Resume**: If transfer is interrupted, receiver can resume from last verified chunk.
7. **Protocol**: Design a request/response protocol:
   - `HANDSHAKE` — exchange file metadata (name, size, chunk count, full hash)
   - `REQUEST <chunk_id>` — request a specific chunk
   - `DATA <chunk_id> <hash> <payload>` — chunk data with hash
   - `DONE` — transfer complete

```go
type Sender struct {
    filePath   string
    listener   net.Listener
    chunks     []*Chunk
    fileHash   string
}

type Receiver struct {
    outputPath  string
    senderAddr  string
    workers     int
    progress    *Progress
    completed   map[int]bool
}

type Chunk struct {
    ID     int
    Offset int64
    Size   int
    Hash   string
    Data   []byte
}
```

**Expected Output:**
```
--- Sender ---
$ go run . send -file=ubuntu.iso -port=9000
[Sender] Serving ubuntu.iso (2.4 GB, 38400 chunks × 64KB)
[Sender] File SHA-256: a1b2c3d4...
[Sender] Listening on :9000

[Sender] Peer 192.168.1.20 connected
[Sender] HANDSHAKE sent (file metadata)
[Sender] Serving chunk 0 (64KB) to worker-1
[Sender] Serving chunk 1 (64KB) to worker-2
[Sender] Serving chunk 2 (64KB) to worker-3
...

--- Receiver ---
$ go run . recv -from=192.168.1.10:9000 -out=./ubuntu.iso -workers=5
[Receiver] Connected to 192.168.1.10:9000
[Receiver] File: ubuntu.iso (2.4 GB, 38400 chunks)
[Receiver] Expected SHA-256: a1b2c3d4...
[Receiver] Downloading with 5 parallel workers...

Downloading... ████████████░░░░░░░░░░░░░░░░░  31% (12042/38400)
Speed: 48.2 MB/s | ETA: 34s | Verified: 12042 chunks

[Worker-3] Chunk 8847 CHECKSUM MISMATCH — retrying...
[Worker-3] Chunk 8847 retry OK ✓

^C
[Receiver] Interrupted — progress saved (12042/38400 chunks)

$ go run . recv -from=192.168.1.10:9000 -out=./ubuntu.iso -workers=5 -resume
[Receiver] Resuming from chunk 12042
Downloading... ████████████████████████████░░  93% (35712/38400)

[Receiver] Download complete!
[Receiver] Verifying full file SHA-256... ✓ MATCH
[Receiver] Saved: ./ubuntu.iso (2.4 GB)
```

**Key learning:** Chunked file transfer, parallel downloads, integrity verification (crypto/sha256), resume capability, custom application protocol design, real-world systems thinking.

---

## Concept Coverage Matrix

| Problem | TCP | UDP | Goroutine-per-conn | Connection Pool | Binary Protocol | Timeouts | Broadcasting | Rate Limiting |
|---------|-----|-----|--------------------|-----------------|-----------------|----------|--------------|---------------|
| 1. Chat Server | ✓ | | ✓ | | | | ✓ | |
| 2. HTTP Server | ✓ | | ✓ | | | ✓ | | |
| 3. Reverse Proxy | ✓ | | | ✓ | | ✓ | | |
| 4. Multiplexer | ✓ | | | | ✓ | | | |
| 5. DNS Resolver | | ✓ | | | ✓ | ✓ | | |
| 6. Port Scanner | ✓ | | | | | ✓ | | ✓ |
| 7. P2P Transfer | ✓ | | ✓ | | ✓ | | | |

---

## Recommended Order

1. **Problem 1** — Chat server (TCP basics + goroutine-per-conn)
2. **Problem 6** — Port scanner (high-concurrency TCP)
3. **Problem 2** — HTTP server (protocol parsing)
4. **Problem 5** — DNS resolver (UDP + binary protocol)
5. **Problem 3** — Reverse proxy (connection pooling + health checks)
6. **Problem 4** — Multiplexer (advanced protocol design)
7. **Problem 7** — P2P transfer (full system)
