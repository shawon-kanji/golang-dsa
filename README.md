# Golang DSA Practice

A collection of Data Structures, Algorithms, and Go programming exercises organized by difficulty and topic.

## Table of Contents

- [Fundamentals](#fundamentals)
- [Data Structures](#data-structures)
- [Intermediate: Interfaces & Composition](#intermediate-interfaces--composition)
- [Arrays, Slices & Maps](#arrays-slices--maps)
- [Concurrency Deep Dive](#concurrency-deep-dive)
- [Creative & Real-World Problems](#creative--real-world-problems)
- [Advanced Challenges](#advanced-challenges)

---

## Fundamentals

### 1. Reverse a String In-Place

**Description:** Reverse a string containing Unicode characters, handling multi-byte runes correctly.

**Problem Statement:**
Given a string that may contain Unicode characters (emojis, Chinese, Japanese, etc.), reverse it in-place. Unlike ASCII strings, Unicode strings require special handling because characters can be multi-byte.

**Key Concepts:**
- Rune vs byte in Go
- Two-pointer technique
- Unicode/UTF-8 encoding
- Slice manipulation

**Edge Cases:**
- Empty string
- Single character
- Strings with emojis and mixed Unicode
- Palindromes

**Example:**
```
Input: "Hello, 世界"
Output: "界世 ,olleH"

Input: "golang"
Output: "gnalog"

Input: "👋🌍"
Output: "🌍👋"
```

---

### 2. Implement a Stack with Generics

**Description:** Create a generic stack data structure supporting any type with standard stack operations.

**Problem Statement:**
Build a type-safe stack using Go 1.18+ generics. The stack should support Push, Pop, Peek, and IsEmpty operations with proper error handling for edge cases like underflow.

**Key Concepts:**
- Go generics (type parameters)
- Method receivers
- Error handling patterns
- Slice as backing storage

**Operations:**
- `Push(item T)` - Add item to top
- `Pop() (T, error)` - Remove and return top item
- `Peek() (T, error)` - View top item without removing
- `IsEmpty() bool` - Check if stack is empty
- `Size() int` - Return number of elements

**Example:**
```go
stack := NewStack[int]()
stack.Push(1)
stack.Push(2)
stack.Push(3)
val, err := stack.Pop()     // val=3, err=nil
peek, _ := stack.Peek()     // peek=2
isEmpty := stack.IsEmpty()  // false
size := stack.Size()        // 2
```

---

## Data Structures

### 3. Custom Linked List

**Description:** Implement a singly linked list with comprehensive operations for insertion, deletion, search, and reversal.

**Problem Statement:**
Create a linked list from scratch with operations to manipulate nodes. Focus on pointer management and edge case handling.

**Key Concepts:**
- Pointer manipulation
- Node struct design
- Head/tail pointer management
- In-place reversal algorithm

**Operations:**
- `InsertAtHead(value)` - Add node at beginning
- `InsertAtTail(value)` - Add node at end
- `InsertAt(index, value)` - Insert at specific position
- `DeleteNode(value)` - Remove first occurrence
- `DeleteAt(index)` - Remove at specific position
- `Find(value) bool` - Search for value
- `Reverse()` - Reverse the list in-place
- `GetMiddle()` - Find middle node (useful for palindrome check)

**Example:**
```
Operations:
InsertAtHead(3)
InsertAtHead(2)
InsertAtTail(4)
InsertAtHead(1)

List: 1 -> 2 -> 3 -> 4 -> nil

Delete(2)
List: 1 -> 3 -> 4 -> nil

Reverse()
List: 4 -> 3 -> 1 -> nil

GetMiddle()
Output: 3
```

---

### 4. Binary Search Tree

**Description:** Implement a Binary Search Tree with insertion, deletion, search, and traversal operations.

**Problem Statement:**
Build a BST that maintains the property: left child < parent < right child. Implement recursive and iterative approaches for different operations.

**Key Concepts:**
- Tree node structure
- Recursive algorithms
- BST property maintenance
- Three deletion cases
- Tree traversals

**Operations:**
- `Insert(value)` - Add new node
- `Search(value) bool` - Find if value exists
- `Delete(value)` - Remove node (handle 3 cases: leaf, one child, two children)
- `FindMin()` / `FindMax()` - Find minimum/maximum values
- `InOrder()` - Left, Root, Right (sorted output)
- `PreOrder()` - Root, Left, Right
- `PostOrder()` - Left, Right, Root
- `Height()` - Calculate tree height
- `IsValid()` - Verify BST property

**Example:**
```
Insert: 50, 30, 70, 20, 40, 60, 80

        50
       /  \
      30   70
     / \   / \
    20 40 60 80

InOrder: [20, 30, 40, 50, 60, 70, 80]
Search(40): found
Delete(30)  // Node with two children
InOrder: [20, 40, 50, 60, 70, 80]
Height: 2
```

---

### 5. LRU Cache

**Description:** Implement a Least Recently Used cache with O(1) Get and Put operations.

**Problem Statement:**
Design a data structure that supports getting and putting key-value pairs with a fixed capacity. When capacity is exceeded, evict the least recently used item. Both operations must run in constant time.

**Key Concepts:**
- HashMap + Doubly-Linked List combination
- O(1) access and eviction
- Node movement for access order
- Capacity management

**Data Structure:**
- Map for O(1) key lookup
- Doubly-linked list for O(1) insertion/deletion and LRU tracking
- Head pointer = most recently used
- Tail pointer = least recently used

**Operations:**
- `Get(key)` - Return value and mark as recently used
- `Put(key, value)` - Insert/update and mark as recently used
- Eviction happens automatically when at capacity

**Example:**
```go
cache := NewLRU(3) // capacity 3
cache.Put("a", 1)
cache.Put("b", 2)
cache.Put("c", 3)
// Order: c -> b -> a

cache.Get("a")     // returns 1
// Order: a -> c -> b (a becomes most recent)

cache.Put("d", 4)  // evicts "b" (least recently used)
// Order: d -> a -> c

cache.Get("b")     // returns -1 (not found)
cache.Get("c")     // returns 3
// Order: c -> d -> a
```

---

## Intermediate: Interfaces & Composition

### 6. Shape Calculator

**Description:** Define a Shape interface and implement it for different geometric shapes to demonstrate polymorphism.

**Problem Statement:**
Create an interface-based design where different shapes implement the same contract. Calculate areas and perimeters for a collection of shapes without knowing their specific types.

**Key Concepts:**
- Interface definition and implementation
- Polymorphism in Go
- Struct embedding (optional)
- Type assertions

**Shapes to Implement:**
- Circle (radius)
- Rectangle (width, height)
- Triangle (three sides, use Heron's formula)
- Square (side, can embed Rectangle)

**Interface:**
```go
type Shape interface {
    Area() float64
    Perimeter() float64
    Name() string
}
```

**Example:**
```go
shapes := []Shape{
    Circle{radius: 5},
    Rectangle{width: 4, height: 6},
    Triangle{a: 3, b: 4, c: 5},
}

Output:
Circle (r=5.0) - Area: 78.54, Perimeter: 31.42
Rectangle (4x6) - Area: 24.00, Perimeter: 20.00
Triangle (3,4,5) - Area: 6.00, Perimeter: 12.00
Total Area: 108.54
```

---

### 7. Plugin System

**Description:** Create a flexible plugin architecture where plugins can be dynamically registered and executed.

**Problem Statement:**
Design a system that allows different plugins to be loaded at runtime without modifying the core code. Each plugin implements a common interface but performs different operations.

**Key Concepts:**
- Interface-based design
- Registry pattern
- Type switches
- Plugin lifecycle

**Plugin Types to Implement:**
- LoggerPlugin - Logs data to console/file
- ValidatorPlugin - Validates input data
- TransformerPlugin - Transforms data (uppercase, lowercase, etc.)
- FilterPlugin - Filters data based on criteria

**Interface:**
```go
type Plugin interface {
    Name() string
    Execute(data interface{}) (interface{}, error)
}
```

**Example:**
```go
registry := NewRegistry()
registry.Register(&LoggerPlugin{})
registry.Register(&UpperCasePlugin{})
registry.Register(&ValidatorPlugin{min: 5})

Input: "hello world"

LoggerPlugin: Logged: hello world
UpperCasePlugin: HELLO WORLD
ValidatorPlugin: Valid (length: 11 >= 5)
```

---

### 8. Sort Custom Types

**Description:** Sort slices of custom structs by implementing the `sort.Interface`.

**Problem Statement:**
Demonstrate how to make custom types sortable by implementing Len(), Less(), and Swap() methods. Show sorting by multiple fields and stable vs unstable sorting.

**Key Concepts:**
- `sort.Interface` implementation
- Comparison logic
- Stable sorting for preserving order
- Multi-field sorting

**Sort Variations:**
- By single field (age)
- By multiple fields (age then name)
- By custom criteria (reverse order)
- Using `sort.Slice()` with anonymous functions

**Example:**
```go
type Person struct {
    Name string
    Age  int
    City string
}

people := []Person{
    {"Alice", 30, "NYC"},
    {"Bob", 25, "LA"},
    {"Charlie", 30, "NYC"},
    {"Alice", 25, "SF"},
}

sort.Sort(ByAge(people))
Output: [Bob(25), Alice(25), Charlie(30), Alice(30)]

sort.Sort(ByNameThenAge(people))
Output: [Alice(25,SF), Alice(30,NYC), Bob(25,LA), Charlie(30,NYC)]

sort.Sort(ByAgeDesc(people))
Output: [Alice(30,NYC), Charlie(30,NYC), Bob(25,LA), Alice(25,SF)]
```

---

### 9. Reader/Writer Wrapper

**Description:** Create custom types that implement `io.Reader` and `io.Writer` interfaces with additional functionality.

**Problem Statement:**
Wrap existing readers/writers to add features like counting bytes, transformation, buffering, or logging without modifying the original source.

**Key Concepts:**
- `io.Reader` and `io.Writer` interfaces
- Composition over inheritance
- Decorator pattern
- Byte stream manipulation

**Implementations:**
- CountingReader - Tracks bytes read
- UpperCaseReader - Converts to uppercase while reading
- ROT13Reader - Applies ROT13 cipher
- TeeReader - Writes to multiple destinations
- LimitReader - Limits bytes read

**Example:**
```go
input := strings.NewReader("Hello, World!")
counter := &CountingReader{reader: input}

data, _ := io.ReadAll(counter)
fmt.Printf("Read %d bytes: %s\n", counter.Count, data)
Output: Read 13 bytes: Hello, World!

// With transformer
transformer := &UpperCaseReader{
    reader: strings.NewReader("hello"),
}
data, _ = io.ReadAll(transformer)
Output: HELLO

// Chaining wrappers
chain := &CountingReader{
    reader: &UpperCaseReader{
        reader: strings.NewReader("test"),
    },
}
Output: TEST (4 bytes counted)
```

---

### 10. Error Handler Chain

**Description:** Create custom error types with context and implement error wrapping compatible with Go 1.13+ error handling.

**Problem Statement:**
Design a hierarchy of custom errors that carry additional context and can be inspected using `errors.Is()` and `errors.As()`. Implement error wrapping to maintain error chains.

**Key Concepts:**
- Custom error types
- Error wrapping with `fmt.Errorf("%w", err)`
- `errors.Is()` for error comparison
- `errors.As()` for type assertion
- Error context and metadata

**Custom Errors to Implement:**
- ValidationError (field, value, reason)
- NetworkError (status code, endpoint)
- DatabaseError (query, operation)
- AuthenticationError (user, reason)

**Example:**
```go
type ValidationError struct {
    Field string
    Value interface{}
    Err   error
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error: field '%s' with value '%v': %v",
        e.Field, e.Value, e.Err)
}

func (e *ValidationError) Unwrap() error {
    return e.Err
}

// Usage
err := ValidateUser(User{Age: -5})
Output: validation error: field 'age' with value '-5': age must be positive

if errors.Is(err, ErrInvalidAge) {
    // Handle specific error
}

var valErr *ValidationError
if errors.As(err, &valErr) {
    fmt.Printf("Invalid field: %s\n", valErr.Field)
}
```

---

## Arrays, Slices & Maps

### 11. Two Sum / Three Sum

**Description:** Find pairs or triplets in an array that sum to a target value.

**Problem Statement:**
Given an array of integers and a target sum, find all unique combinations of elements that sum to the target. For Two Sum, find pairs; for Three Sum, find triplets.

**Key Concepts:**
- Hash map for O(n) Two Sum
- Sorting + two pointers for Three Sum
- Handling duplicates
- Multiple solutions handling

**Variations:**
- Two Sum - return indices
- Two Sum - return values
- Three Sum - find all unique triplets
- Four Sum - extend to four numbers

**Example:**
```
Two Sum:
Input: nums = [2, 7, 11, 15], target = 9
Output: [0, 1] (nums[0] + nums[1] = 2 + 7 = 9)

Input: nums = [3, 2, 4], target = 6
Output: [1, 2]

Three Sum:
Input: nums = [-1, 0, 1, 2, -1, -4], target = 0
Output: [[-1, -1, 2], [-1, 0, 1]]

Input: nums = [0, 0, 0, 0], target = 0
Output: [[0, 0, 0]]
```

---

### 12. Merge Intervals

**Description:** Merge overlapping intervals in a collection.

**Problem Statement:**
Given a collection of intervals, merge all overlapping intervals and return non-overlapping intervals that cover the same range.

**Key Concepts:**
- Sorting by start time
- Greedy algorithm
- Interval comparison logic
- Edge case handling

**Applications:**
- Meeting room scheduling
- Time range consolidation
- Calendar event merging

**Example:**
```
Input: [[1,3], [2,6], [8,10], [15,18]]
Output: [[1,6], [8,10], [15,18]]
Explanation: [1,3] and [2,6] overlap, merged to [1,6]

Input: [[1,4], [4,5]]
Output: [[1,5]]
Explanation: Adjacent intervals merged

Input: [[1,4], [0,2], [3,5]]
Output: [[0,5]]
Explanation: After sorting: [0,2], [1,4], [3,5] all merge
```

---

### 13. Group Anagrams

**Description:** Group strings that are anagrams of each other.

**Problem Statement:**
Given a list of strings, group together all strings that are anagrams. Anagrams are words formed by rearranging letters of another word.

**Key Concepts:**
- Hash map with sorted string as key
- String sorting
- Map of slices
- Character frequency counting (alternative approach)

**Optimizations:**
- Use character count array instead of sorting
- Case-insensitive grouping option
- Unicode support

**Example:**
```
Input: ["eat", "tea", "tan", "ate", "nat", "bat"]
Output:
[
  ["eat", "tea", "ate"],
  ["tan", "nat"],
  ["bat"]
]

Input: ["listen", "silent", "hello", "world"]
Output:
[
  ["listen", "silent"],
  ["hello"],
  ["world"]
]

Input: ["", ""]
Output: [["", ""]]
```

---

### 14. Sliding Window Maximum

**Description:** Find the maximum value in each sliding window of size k.

**Problem Statement:**
Given an array and a window size k, slide the window from left to right, one position at a time. For each window position, find the maximum value.

**Key Concepts:**
- Deque (double-ended queue) data structure
- Monotonic decreasing queue
- O(n) time complexity
- Amortized constant time per operation

**Approach:**
- Maintain indices in deque, not values
- Keep deque in decreasing order
- Remove indices outside window
- Front of deque is always current maximum

**Example:**
```
Input: nums = [1,3,-1,-3,5,3,6,7], k = 3
Output: [3,3,5,5,6,7]

Explanation:
Window [1,3,-1]: max = 3
Window [3,-1,-3]: max = 3
Window [-1,-3,5]: max = 5
Window [-3,5,3]: max = 5
Window [5,3,6]: max = 6
Window [3,6,7]: max = 7

Input: nums = [1], k = 1
Output: [1]

Input: nums = [1,-1], k = 1
Output: [1,-1]
```

---

### 15. Word Frequency Counter

**Description:** Count word frequencies in text and provide various analysis features.

**Problem Statement:**
Read text from a file or string, count how many times each word appears, and provide options for sorting and filtering results.

**Key Concepts:**
- Text parsing and tokenization
- Map for frequency counting
- Sorting by value (not just key)
- File I/O and string manipulation
- Unicode word boundaries

**Features:**
- Ignore punctuation and special characters
- Case-insensitive counting
- Stop word filtering
- Top N most frequent words
- Export to different formats

**Example:**
```
Input file (text.txt):
"Hello world! Hello Go. Go is awesome. Go Go Go!"

Output (default):
go: 5
hello: 2
world: 1
is: 1
awesome: 1

Top 3 words:
1. go (5 occurrences)
2. hello (2 occurrences)
3. world (1 occurrence)

With min frequency filter (>= 2):
go: 5
hello: 2
```

---

## Concurrency Deep Dive

### 16. Worker Pool Pattern

**Description:** Create a fixed number of worker goroutines to process jobs from a shared queue.

**Problem Statement:**
Implement a thread pool pattern where a limited number of workers process jobs concurrently. This prevents spawning too many goroutines and manages resource usage.

**Key Concepts:**
- Buffered channels for job queue
- Worker goroutines
- `sync.WaitGroup` for synchronization
- Graceful shutdown
- Result collection

**Components:**
- Job queue channel
- Result channel
- Worker pool
- Dispatcher

**Example:**
```go
jobs := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
numWorkers := 3

Input: Process each job (square the number)
Output:
Worker 1 processing job 1 -> result: 1
Worker 2 processing job 2 -> result: 4
Worker 3 processing job 3 -> result: 9
Worker 1 processing job 4 -> result: 16
Worker 2 processing job 5 -> result: 25
...

All results: [1, 4, 9, 16, 25, 36, 49, 64, 81, 100]
Processing time: 2.5s (vs 5s sequential)
Speedup: 2x with 3 workers
```

---

### 17. Pipeline Pattern

**Description:** Build a data processing pipeline where each stage is a goroutine connected by channels.

**Problem Statement:**
Create a series of processing stages where data flows through channels from one stage to the next, allowing concurrent processing at each stage.

**Key Concepts:**
- Channel chaining
- Generator pattern
- Fan-out for parallelism
- Closing channels properly
- Range over channels

**Pipeline Stages:**
1. Generator - produces data
2. Processor - transforms data
3. Filter - selects data
4. Aggregator - combines results

**Example:**
```go
// Stage 1: Generate numbers 1-10
// Stage 2: Square each number
// Stage 3: Filter even results
// Stage 4: Sum all values

Input: 1, 2, 3, 4, 5
Stage 1 (generate): 1, 2, 3, 4, 5
Stage 2 (square): 1, 4, 9, 16, 25
Stage 3 (filter evens): 4, 16
Stage 4 (sum): 20

Final output: 20

// Another example with strings
Input: ["hello", "world", "go"]
Stage 1 (uppercase): ["HELLO", "WORLD", "GO"]
Stage 2 (add prefix): ["MSG: HELLO", "MSG: WORLD", "MSG: GO"]
Stage 3 (filter length > 10): ["MSG: HELLO", "MSG: WORLD"]
```

---

### 18. Fan-Out/Fan-In

**Description:** Distribute work across multiple workers (fan-out) and merge results from all workers (fan-in).

**Problem Statement:**
Take a single input stream, distribute work to multiple workers for parallel processing, then merge all results back into a single output stream.

**Key Concepts:**
- Multiple worker goroutines from single source
- Merge channels from multiple sources
- Load balancing across workers
- Result synchronization

**Pattern:**
```
        [Input Channel]
             |
      Fan-Out (distribute)
       /     |     \
    Worker1 Worker2 Worker3
       \     |     /
       Fan-In (merge)
             |
       [Output Channel]
```

**Example:**
```go
// Check if numbers are prime (CPU-intensive)
Input: [2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20]

Fan-out to 4 workers:
Worker 1: checking [2, 6, 10, 14, 18]
Worker 2: checking [3, 7, 11, 15, 19]
Worker 3: checking [4, 8, 12, 16, 20]
Worker 4: checking [5, 9, 13, 17]

Results (fan-in):
Primes found: [2, 3, 5, 7, 11, 13, 17, 19]
Time with 4 workers: 0.5s
Time with 1 worker: 2.0s
Speedup: 4x
```

---

### 19. Rate Limiter

**Description:** Limit the number of operations per time period using channels and timers.

**Problem Statement:**
Implement a rate limiter that allows only N operations per time window, queuing or blocking additional requests until capacity is available.

**Key Concepts:**
- `time.Ticker` for rate control
- Token bucket algorithm
- Buffered channels as semaphore
- Non-blocking try operations

**Algorithms:**
- Token bucket (burst allowed)
- Fixed window
- Sliding window
- Leaky bucket

**Example:**
```go
limiter := NewRateLimiter(3, time.Second) // 3 requests per second

for i := 1; i <= 10; i++ {
    limiter.Wait()
    fmt.Printf("Request %d at %s\n", i, time.Now().Format("15:04:05.000"))
}

Output:
Request 1 at 10:00:00.000
Request 2 at 10:00:00.000
Request 3 at 10:00:00.000
Request 4 at 10:00:01.000  // waited 1s
Request 5 at 10:00:01.000
Request 6 at 10:00:01.000
Request 7 at 10:00:02.000  // waited 1s
Request 8 at 10:00:02.000
Request 9 at 10:00:02.000
Request 10 at 10:00:03.000

// Non-blocking variant
if limiter.TryAcquire() {
    // Process request
} else {
    // Reject request
}
```

---

### 20. Concurrent Web Scraper

**Description:** Fetch multiple URLs concurrently with error handling, timeouts, and rate limiting.

**Problem Statement:**
Build a web scraper that fetches multiple URLs in parallel while respecting concurrency limits, handling errors gracefully, and implementing timeouts.

**Key Concepts:**
- HTTP client usage
- Goroutine per request
- Semaphore pattern for limiting concurrency
- Context for timeout and cancellation
- Error collection

**Features:**
- Concurrent fetching with max limit
- Timeout per request
- Retry logic
- Response size tracking
- Error aggregation

**Example:**
```go
urls := []string{
    "https://example.com",
    "https://golang.org",
    "https://github.com",
    "https://invalid-url-12345.com",
}

maxConcurrent := 2
timeout := 5 * time.Second

Output:
[Worker 1] Fetching https://example.com...
[Worker 2] Fetching https://golang.org...
✓ https://example.com (1.2s, 1256 bytes)
[Worker 1] Fetching https://github.com...
✓ https://golang.org (0.8s, 4589 bytes)
[Worker 2] Fetching https://invalid-url-12345.com...
✓ https://github.com (1.5s, 12456 bytes)
✗ https://invalid-url-12345.com (error: no such host)

Summary:
- Succeeded: 3
- Failed: 1
- Total time: 2.3s
- Total bytes: 18301
```

---

### 21. Pub/Sub System

**Description:** Build a publish-subscribe message broker for topic-based messaging.

**Problem Statement:**
Create a message broker where publishers send messages to topics, and multiple subscribers receive those messages. Each subscriber runs in its own goroutine.

**Key Concepts:**
- Topic-based routing
- Channel per subscriber
- Goroutine per subscriber
- Subscribe/Unsubscribe operations
- Broadcast to all subscribers

**Features:**
- Multiple topics
- Dynamic subscription
- Buffered subscriber channels
- Graceful shutdown
- Message persistence (optional)

**Example:**
```go
broker := NewBroker()
broker.CreateTopic("news")
broker.CreateTopic("sports")

sub1 := broker.Subscribe("news")
sub2 := broker.Subscribe("news")
sub3 := broker.Subscribe("sports")

broker.Publish("news", "Breaking: Go 1.23 released!")
broker.Publish("news", "Tutorial: Concurrency patterns")
broker.Publish("sports", "Team wins championship")

Output:
[Sub1 - news] Received: "Breaking: Go 1.23 released!"
[Sub2 - news] Received: "Breaking: Go 1.23 released!"
[Sub1 - news] Received: "Tutorial: Concurrency patterns"
[Sub2 - news] Received: "Tutorial: Concurrency patterns"
[Sub3 - sports] Received: "Team wins championship"

broker.Unsubscribe("news", sub1)
broker.Publish("news", "Another news item")
// Only sub2 receives this message
```

---

### 22. Dining Philosophers

**Description:** Solve the classic concurrency problem of 5 philosophers sharing 5 forks without deadlock.

**Problem Statement:**
Five philosophers sit at a round table with five forks. Each philosopher needs two forks to eat. Implement a solution that prevents deadlock and starvation.

**Key Concepts:**
- Mutex for resource locking
- Deadlock prevention strategies
- Resource ordering
- Timeout-based acquisition
- Arbitrator pattern

**Deadlock Prevention Strategies:**
1. Resource ordering (always pick lower-numbered fork first)
2. Timeout on lock acquisition
3. Arbitrator (limit to 4 eating at once)
4. Chandy-Misra algorithm

**Example:**
```
5 Philosophers, 5 forks arranged in circle

Solution: Resource ordering (pick lower-numbered fork first)

Output:
Philosopher 0 is thinking
Philosopher 1 is thinking
Philosopher 2 is thinking
Philosopher 0 picked up fork 0
Philosopher 0 picked up fork 1
Philosopher 0 is eating
Philosopher 2 picked up fork 2
Philosopher 0 finished eating
Philosopher 0 released forks 0, 1
Philosopher 1 picked up fork 1
Philosopher 1 picked up fork 2  // Wait, fork 2 is held by Phil 2!
Philosopher 2 picked up fork 3
Philosopher 2 is eating
...

Statistics after 100 cycles:
- Total meals: 500 (100 per philosopher)
- No deadlocks detected
- Average wait time: 0.3s
```

---

### 23. Parallel Merge Sort

**Description:** Implement merge sort using goroutines to parallelize the divide step.

**Problem Statement:**
Implement merge sort that spawns goroutines for recursive calls, with a depth limit to avoid excessive overhead.

**Key Concepts:**
- Recursive parallelism
- Goroutine spawning with depth limit
- Channel synchronization
- Comparison with sequential version

**Optimization:**
- Only spawn goroutines above certain array size
- Limit recursion depth for goroutines
- Use sequential sort for small subarrays

**Example:**
```go
Input: [38, 27, 43, 3, 9, 82, 10, 50, 12, 7]

Sequential merge sort: 15ms
Parallel merge sort (max 4 goroutines): 6ms
Speedup: 2.5x

Output: [3, 7, 9, 10, 12, 27, 38, 43, 50, 82]

Goroutine tree (depth=2):
            [38,27,43,3,9,82,10,50,12,7]
               /                    \
    [38,27,43,3,9]             [82,10,50,12,7]
       /        \                 /         \
  [38,27,43]  [3,9]        [82,10,50]   [12,7]
  (sequential) (sequential) (sequential) (sequential)

Goroutines created: 3
Max concurrent: 4
```

---

### 24. Context Cancellation

**Description:** Demonstrate proper use of `context.Context` for cancellation, timeouts, and deadlines.

**Problem Statement:**
Show how to propagate cancellation signals through a call stack using context, implement timeout handling, and properly clean up resources.

**Key Concepts:**
- `context.WithTimeout`
- `context.WithCancel`
- `context.WithDeadline`
- `context.WithValue`
- Cancellation propagation
- Cleanup with defer

**Use Cases:**
- HTTP request timeout
- Database query cancellation
- Multi-tier service calls
- Background job cancellation

**Example:**
```go
// Scenario 1: Timeout
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()

result := LongRunningTask(ctx) // Takes 5 seconds
Output: "Task cancelled: context deadline exceeded (after 2s)"

// Scenario 2: Manual cancellation
ctx, cancel := context.WithCancel(context.Background())
go func() {
    time.Sleep(1 * time.Second)
    cancel() // Cancel after 1 second
}()

result := LongRunningTask(ctx)
Output: "Task cancelled after 1s: context canceled"

// Scenario 3: Successful completion
ctx := context.Background()
result := QuickTask(ctx) // Takes 0.5 seconds
Output: "Task completed: result data"

// Scenario 4: Context with values
ctx := context.WithValue(context.Background(), "requestID", "abc123")
ProcessRequest(ctx)
// Inside ProcessRequest:
requestID := ctx.Value("requestID").(string)
Output: "Processing request abc123"
```

---

### 25. Semaphore Implementation

**Description:** Implement a counting semaphore using channels to limit concurrent resource access.

**Problem Statement:**
Create a semaphore that allows up to N concurrent operations, blocking additional attempts until slots become available.

**Key Concepts:**
- Buffered channel as semaphore
- Acquire/Release pattern
- Try-acquire (non-blocking)
- Timeout-based acquire

**Operations:**
- `Acquire()` - Block until slot available
- `Release()` - Free up a slot
- `TryAcquire() bool` - Non-blocking attempt
- `AcquireWithTimeout(duration) bool` - Timeout variant

**Example:**
```go
sem := NewSemaphore(3) // Allow 3 concurrent operations

for i := 1; i <= 10; i++ {
    go func(id int) {
        sem.Acquire()
        defer sem.Release()

        fmt.Printf("Task %d started\n", id)
        time.Sleep(1 * time.Second) // Simulate work
        fmt.Printf("Task %d finished\n", id)
    }(i)
}

Output:
Task 1 started
Task 2 started
Task 3 started
[After 1s]
Task 1 finished
Task 4 started  // waited for slot
Task 2 finished
Task 5 started
Task 3 finished
Task 6 started
[After 1s]
Task 4 finished
Task 7 started
...

// Using TryAcquire
if sem.TryAcquire() {
    defer sem.Release()
    // Do work
} else {
    fmt.Println("Resource busy, skipping")
}
```

---

## Creative & Real-World Problems

### 26. CLI Tool with Flags

**Description:** Build a command-line application with subcommands, flags, and help text.

**Problem Statement:**
Create a CLI tool similar to `git` or `docker` with multiple subcommands, required and optional flags, and comprehensive help documentation.

**Key Concepts:**
- `flag` package or `cobra` library
- Subcommand routing
- Flag parsing and validation
- Help text generation
- Exit codes

**Suggested Tool:** Resource manager (create, list, delete, update resources)

**Example:**
```bash
$ myapp --help
Usage: myapp [command] [flags]

A tool for managing resources

Commands:
  create    Create a new resource
  list      List all resources
  delete    Delete a resource
  update    Update a resource

Flags:
  --config string    Config file path
  --verbose          Enable verbose output
  --version          Show version

$ myapp create --name "database" --type "postgres" --size 100
✓ Created postgres resource: database (100GB)

$ myapp list --format json --filter "type=postgres"
[
  {
    "name": "database",
    "type": "postgres",
    "size": 100,
    "created": "2025-01-15T10:00:00Z"
  }
]

$ myapp delete --name "database" --force
Warning: This will permanently delete the resource
✓ Deleted resource: database

$ myapp update --name "database" --size 200
✓ Updated database size to 200GB
```

---

### 27. JSON/CSV Parser

**Description:** Parse, validate, and convert between JSON and CSV formats.

**Problem Statement:**
Build a tool that can read JSON and CSV files, validate data against schemas, and convert between formats with proper type handling.

**Key Concepts:**
- `encoding/json` package
- `encoding/csv` package
- Struct tags for field mapping
- Custom unmarshalers
- Data validation

**Features:**
- CSV to JSON conversion
- JSON to CSV conversion
- Schema validation
- Pretty printing
- Nested object handling

**Example:**
```go
// Input CSV (users.csv):
name,age,email,active
Alice,30,alice@example.com,true
Bob,25,bob@example.com,false
Charlie,35,charlie@example.com,true

// Parse to structs
users := ParseCSV("users.csv")

// Convert to JSON
Output JSON:
[
  {
    "name": "Alice",
    "age": 30,
    "email": "alice@example.com",
    "active": true
  },
  {
    "name": "Bob",
    "age": 25,
    "email": "bob@example.com",
    "active": false
  },
  {
    "name": "Charlie",
    "age": 35,
    "email": "charlie@example.com",
    "active": true
  }
]

// Validation
Validation results:
✓ All records have required fields
✓ All ages are positive integers
✓ All emails match pattern
✗ 1 invalid email: "charlie@example" (missing domain)

// Convert JSON to CSV
Input JSON -> Output CSV with headers
```

---

### 28. Simple HTTP Server

**Description:** Build a RESTful API server with CRUD operations using `net/http`.

**Problem Statement:**
Create an HTTP server that handles REST API endpoints for managing a resource (e.g., users, todos, products) with proper routing, JSON handling, and middleware.

**Key Concepts:**
- HTTP handlers and handler functions
- Routing (using `net/http` or `gorilla/mux`)
- JSON encoding/decoding
- Middleware pattern
- In-memory data store or database

**Features:**
- CRUD endpoints (Create, Read, Update, Delete)
- Request validation
- Error handling and status codes
- Logging middleware
- Authentication middleware (basic)
- CORS support

**Example:**
```bash
# Start server
$ go run main.go
Server listening on :8080

# Create user (POST /users)
$ curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","age":30,"email":"alice@example.com"}'

Response: {
  "id": "1",
  "name": "Alice",
  "age": 30,
  "email": "alice@example.com",
  "created_at": "2025-01-15T10:00:00Z"
}

# Get all users (GET /users)
$ curl http://localhost:8080/users

Response: [
  {"id":"1","name":"Alice","age":30,"email":"alice@example.com"}
]

# Get specific user (GET /users/:id)
$ curl http://localhost:8080/users/1

Response: {"id":"1","name":"Alice","age":30}

# Update user (PUT /users/:id)
$ curl -X PUT http://localhost:8080/users/1 \
  -d '{"name":"Alice Smith","age":31}'

Response: {"id":"1","name":"Alice Smith","age":31}

# Delete user (DELETE /users/:id)
$ curl -X DELETE http://localhost:8080/users/1

Response: {"message":"User deleted successfully"}

# Error handling
$ curl http://localhost:8080/users/999

Response: 404 {"error":"User not found"}
```

---

### 29. Markdown to HTML Converter

**Description:** Parse markdown syntax and convert it to HTML.

**Problem Statement:**
Build a markdown parser that supports common markdown syntax and converts it to valid HTML output.

**Key Concepts:**
- String parsing and regex
- State machine for parsing
- HTML escaping
- Nested structure handling

**Supported Syntax:**
- Headers (`#`, `##`, `###`)
- Bold (`**text**`) and italic (`*text*`)
- Links (`[text](url)`)
- Images (`![alt](url)`)
- Code blocks (` ``` `)
- Inline code (`` `code` ``)
- Lists (ordered and unordered)
- Blockquotes (`>`)
- Horizontal rules (`---`)

**Example:**
```markdown
Input:
# My Blog Post

This is **bold** and this is *italic*.
You can also have ***bold italic***.

## Features

- Item 1
- Item 2
  - Nested item
- Item 3

1. First
2. Second
3. Third

Here's a [link to Go](https://golang.org).

Inline `code` looks like this.

```go
func main() {
    fmt.Println("Hello, World!")
}
```

> This is a blockquote
> It can span multiple lines

---

Output HTML:
<h1>My Blog Post</h1>

<p>This is <strong>bold</strong> and this is <em>italic</em>.</p>
<p>You can also have <strong><em>bold italic</em></strong>.</p>

<h2>Features</h2>

<ul>
  <li>Item 1</li>
  <li>Item 2
    <ul>
      <li>Nested item</li>
    </ul>
  </li>
  <li>Item 3</li>
</ul>

<ol>
  <li>First</li>
  <li>Second</li>
  <li>Third</li>
</ol>

<p>Here's a <a href="https://golang.org">link to Go</a>.</p>

<p>Inline <code>code</code> looks like this.</p>

<pre><code class="language-go">func main() {
    fmt.Println("Hello, World!")
}
</code></pre>

<blockquote>
  <p>This is a blockquote</p>
  <p>It can span multiple lines</p>
</blockquote>

<hr>
```

---

### 30. Log Aggregator

**Description:** Read logs from multiple sources concurrently, merge by timestamp, and filter by criteria.

**Problem Statement:**
Build a tool that reads log files from multiple sources in parallel, parses timestamps and log levels, and merges them into a single chronologically-sorted stream.

**Key Concepts:**
- Concurrent file reading
- Timestamp parsing
- Merge-sort algorithm
- Log level filtering
- Stream processing

**Features:**
- Read from multiple files concurrently
- Parse various timestamp formats
- Merge by timestamp
- Filter by log level (DEBUG, INFO, WARN, ERROR)
- Search by keyword
- Output to file or stdout

**Example:**
```
Input files:
- app1.log
- app2.log
- app3.log

app1.log:
2025-01-15 10:00:01 INFO Server started on port 8080
2025-01-15 10:00:05 ERROR Connection to database failed
2025-01-15 10:00:08 INFO Retrying connection

app2.log:
2025-01-15 10:00:02 INFO Request received: GET /api/users
2025-01-15 10:00:04 WARN High memory usage: 85%
2025-01-15 10:00:07 INFO Request completed: 200 OK

app3.log:
2025-01-15 10:00:03 INFO Background job started
2025-01-15 10:00:06 DEBUG Processing 1000 records
2025-01-15 10:00:09 INFO Background job completed

Merged output (sorted by time):
2025-01-15 10:00:01 [app1] INFO Server started on port 8080
2025-01-15 10:00:02 [app2] INFO Request received: GET /api/users
2025-01-15 10:00:03 [app3] INFO Background job started
2025-01-15 10:00:04 [app2] WARN High memory usage: 85%
2025-01-15 10:00:05 [app1] ERROR Connection to database failed
2025-01-15 10:00:06 [app3] DEBUG Processing 1000 records
2025-01-15 10:00:07 [app2] INFO Request completed: 200 OK
2025-01-15 10:00:08 [app1] INFO Retrying connection
2025-01-15 10:00:09 [app3] INFO Background job completed

Filter (ERROR only):
2025-01-15 10:00:05 [app1] ERROR Connection to database failed

Filter (WARN and ERROR):
2025-01-15 10:00:04 [app2] WARN High memory usage: 85%
2025-01-15 10:00:05 [app1] ERROR Connection to database failed
```

---

### 31. Memory Cache with TTL

**Description:** In-memory key-value cache with automatic expiration based on time-to-live.

**Problem Statement:**
Implement a thread-safe cache where each entry has a TTL. Automatically clean up expired entries in the background.

**Key Concepts:**
- `sync.RWMutex` for concurrent access
- Time-based expiration with `time.Timer`
- Background cleanup goroutine
- Map for O(1) access

**Operations:**
- `Set(key, value, ttl)` - Store with expiration
- `Get(key) (value, found)` - Retrieve if not expired
- `Delete(key)` - Manual removal
- `Clear()` - Remove all entries
- `Stats()` - Get statistics

**Features:**
- Thread-safe operations
- Automatic background cleanup
- Lazy expiration on access
- Memory usage tracking

**Example:**
```go
cache := NewCache()

cache.Set("user:1", "Alice", 5*time.Second)
cache.Set("user:2", "Bob", 10*time.Second)
cache.Set("session:abc", "data", 3*time.Second)

val, found := cache.Get("user:1") // "Alice", true

time.Sleep(4 * time.Second)
val, found = cache.Get("session:abc") // "", false (expired)
val, found = cache.Get("user:1")       // "", false (expired)
val, found = cache.Get("user:2")       // "Bob", true (still valid)

time.Sleep(7 * time.Second)
val, found = cache.Get("user:2") // "", false (expired)

Output log:
[00:00] Set user:1 -> Alice (TTL: 5s)
[00:00] Set user:2 -> Bob (TTL: 10s)
[00:00] Set session:abc -> data (TTL: 3s)
[00:03] session:abc expired and removed
[00:05] user:1 expired and removed
[00:10] user:2 expired and removed

Stats:
- Total keys: 0
- Total sets: 3
- Total gets: 5
- Hits: 3
- Misses: 2
- Expired: 3
- Active: 0
```

---

### 32. File Watcher

**Description:** Monitor a directory for file system events (create, modify, delete).

**Problem Statement:**
Build a file watcher that detects changes in a directory tree and triggers callbacks for different event types.

**Key Concepts:**
- `fsnotify` library or polling
- Recursive directory watching
- Event debouncing
- Callback registration

**Events:**
- Create
- Modify
- Delete
- Rename/Move

**Use Cases:**
- Auto-recompile on file change
- Auto-run tests
- Live reload for web development
- Backup trigger

**Example:**
```go
watcher := NewFileWatcher("/path/to/watch")
watcher.Start()

// Register callbacks
watcher.OnCreate(func(path string) {
    fmt.Printf("[CREATED] %s\n", path)
})

watcher.OnModify(func(path string) {
    fmt.Printf("[MODIFIED] %s\n", path)
})

watcher.OnDelete(func(path string) {
    fmt.Printf("[DELETED] %s\n", path)
})

Output (as files change):
[10:00:01] CREATED: /path/to/watch/newfile.txt
[10:00:03] MODIFIED: /path/to/watch/newfile.txt
[10:00:05] MODIFIED: /path/to/watch/newfile.txt (debounced)
[10:00:07] DELETED: /path/to/watch/newfile.txt
[10:00:10] CREATED: /path/to/watch/subdir/
[10:00:12] CREATED: /path/to/watch/subdir/file.go

// Specific use case: auto-run tests on .go file change
watcher.OnModify(func(path string) {
    if strings.HasSuffix(path, ".go") {
        fmt.Println("Go file changed, running tests...")
        exec.Command("go", "test", "./...").Run()
    }
})
```

---

### 33. URL Shortener

**Description:** Build a URL shortening service with redirection and statistics tracking.

**Problem Statement:**
Create a service that generates short codes for long URLs, handles redirection, and tracks usage statistics.

**Key Concepts:**
- Base62 encoding for short codes
- HTTP redirection
- URL validation
- Counter synchronization for stats

**Features:**
- Generate short codes (auto or custom)
- Redirect to original URL
- Click tracking
- Expiration dates
- QR code generation (bonus)

**Example:**
```bash
# Shorten URL (POST /shorten)
$ curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url":"https://github.com/golang/go"}'

Response: {
  "short_url": "http://localhost:8080/abc123",
  "original_url": "https://github.com/golang/go",
  "created_at": "2025-01-15T10:00:00Z"
}

# Visit short URL
$ curl -L http://localhost:8080/abc123
→ 302 Redirect to: https://github.com/golang/go

# Get statistics
$ curl http://localhost:8080/stats/abc123

Response: {
  "short_code": "abc123",
  "original_url": "https://github.com/golang/go",
  "created_at": "2025-01-15T10:00:00Z",
  "click_count": 42,
  "last_accessed": "2025-01-15T14:30:00Z"
}

# Custom short code
$ curl -X POST http://localhost:8080/shorten \
  -d '{"url":"https://golang.org","custom_code":"go"}'

Response: {
  "short_url": "http://localhost:8080/go",
  "original_url": "https://golang.org"
}

# List all URLs
$ curl http://localhost:8080/urls

Response: [
  {"code":"abc123","url":"https://github.com/golang/go","clicks":42},
  {"code":"go","url":"https://golang.org","clicks":5}
]
```

---

### 34. WebSocket Chat Server

**Description:** Real-time chat application using WebSockets with rooms and private messaging.

**Problem Statement:**
Build a WebSocket server that supports multiple chat rooms, broadcasting messages, and private messages between users.

**Key Concepts:**
- `gorilla/websocket` library
- Connection management
- Room-based messaging
- Broadcast patterns

**Features:**
- Multiple chat rooms
- Join/leave notifications
- Broadcast to room
- Private messages
- User nicknames
- Online user list

**Example:**
```
Server started on :8080

Client 1 connects:
→ User1 joined the server

Client 1: /join general
→ User1 joined room "general"

Client 2 connects and joins "general":
→ User2 joined the server
→ User2 joined room "general"
Broadcast to general: "User2 joined the room"

User1: Hello everyone!
→ Broadcast to all in "general"
User2 sees: [User1] Hello everyone!

User2: Hi User1!
→ Broadcast to all in "general"
User1 sees: [User2] Hi User1!

User1: /pm User2 Secret message
→ Private message sent
Only User2 sees: [Private from User1] Secret message

Client 3 connects and joins "random":
→ User3 joined the server
→ User3 joined room "random"
(User3 doesn't see "general" room messages)

User1: /users
→ Shows online users in current room
Response: Users in "general": User1, User2

User1 disconnects:
→ User1 left the server
Broadcast to general: "User1 left the room"
```

---

### 35. Image Processing Pipeline

**Description:** Load images, apply transformations concurrently, and save results.

**Problem Statement:**
Build an image processing pipeline that applies multiple transformations to images in sequence, processing multiple images concurrently.

**Key Concepts:**
- `image` package for decoding/encoding
- Image manipulation algorithms
- Pipeline pattern
- Concurrent processing

**Transformations:**
- Grayscale conversion
- Resize
- Blur
- Rotate
- Flip (horizontal/vertical)
- Brightness/contrast adjustment
- Crop

**Example:**
```go
pipeline := NewImagePipeline()
pipeline.AddStage(Grayscale)
pipeline.AddStage(Resize(800, 600))
pipeline.AddStage(Blur(5))
pipeline.AddStage(Rotate(90))

images := []string{
    "photo1.jpg",
    "photo2.jpg",
    "photo3.jpg",
}

pipeline.ProcessConcurrent(images, "output/", 3) // 3 workers

Output:
[Worker 1] Processing photo1.jpg
  → Grayscale applied
  → Resized to 800x600
  → Blur (radius=5) applied
  → Rotated 90° clockwise
  → Saved to output/photo1_processed.jpg (1.2s)

[Worker 2] Processing photo2.jpg
  → Grayscale applied
  → Resized to 800x600
  → Blur (radius=5) applied
  → Rotated 90° clockwise
  → Saved to output/photo2_processed.jpg (1.1s)

[Worker 3] Processing photo3.jpg
  → Grayscale applied
  → Resized to 800x600
  → Blur (radius=5) applied
  → Rotated 90° clockwise
  → Saved to output/photo3_processed.jpg (1.3s)

Summary:
- Processed: 3 images
- Total time: 1.3s (concurrent) vs 3.6s (sequential)
- Speedup: 2.8x
- Total output size: 2.4 MB
```

---

## Advanced Challenges

### 36. Implement sync.Pool

**Description:** Create your own object pool implementation for reducing allocations.

**Problem Statement:**
Build an object pool that reuses objects to minimize garbage collection pressure. Implement per-P (processor) pools for reduced contention.

**Key Concepts:**
- Object pooling pattern
- Per-P local pools
- `runtime.GOMAXPROCS`
- Benchmark and profiling

**Features:**
- Get/Put operations
- Automatic object creation with New function
- Per-P pools to reduce lock contention
- GC integration (optional cleanup)

**Example:**
```go
type Buffer struct {
    data []byte
}

pool := &CustomPool{
    New: func() interface{} {
        return &Buffer{data: make([]byte, 0, 1024)}
    },
}

// Use case: HTTP handlers reusing buffers
func handler(w http.ResponseWriter, r *http.Request) {
    buf := pool.Get().(*Buffer)
    defer pool.Put(buf)

    // Reset before use
    buf.data = buf.data[:0]

    // Use buffer for processing
    buf.data = append(buf.data, []byte("Hello")...)
    w.Write(buf.data)
}

Benchmark results:
BenchmarkWithoutPool-8    10000    502847 ns/op    10485760 B/op    10000 allocs/op
BenchmarkWithPool-8      500000      3421 ns/op      524288 B/op      100 allocs/op

Improvement:
- Speed: 147x faster
- Memory: 20x less allocated
- Allocations: 100x fewer
```

---

### 37. Custom Context Implementation

**Description:** Build your own implementation of `context.Context` interface.

**Problem Statement:**
Recreate the context.Context functionality from scratch to understand how cancellation propagation and deadlines work internally.

**Key Concepts:**
- Interface implementation
- Channel-based cancellation
- Deadline tracking
- Value storage and propagation
- Parent-child relationship

**Methods to Implement:**
- `Deadline() (deadline time.Time, ok bool)`
- `Done() <-chan struct{}`
- `Err() error`
- `Value(key interface{}) interface{}`
- Helper functions: `WithCancel`, `WithTimeout`, `WithDeadline`, `WithValue`

**Example:**
```go
type MyContext struct {
    parent   Context
    done     chan struct{}
    err      error
    deadline time.Time
    mu       sync.Mutex
    values   map[interface{}]interface{}
}

// Usage
ctx := NewMyContext()
ctx = WithTimeout(ctx, 5*time.Second)
ctx = WithValue(ctx, "userID", "123")
ctx = WithValue(ctx, "requestID", "abc")

// Simulate work with timeout
select {
case <-ctx.Done():
    fmt.Println("Context cancelled:", ctx.Err())
case <-time.After(3 * time.Second):
    fmt.Println("Work completed successfully")
}
Output: Work completed successfully

// Extract values
userID := ctx.Value("userID").(string)    // "123"
reqID := ctx.Value("requestID").(string)  // "abc"

// Test cancellation propagation
parent, cancel := WithCancel(Background())
child1 := WithValue(parent, "key", "child1")
child2 := WithValue(parent, "key", "child2")

cancel() // Cancels parent

<-child1.Done() // Also cancelled
<-child2.Done() // Also cancelled
```

---

### 38. Goroutine Leak Detector

**Description:** Build a tool to detect goroutines that never terminate.

**Problem Statement:**
Create a debugging tool that tracks goroutine creation, monitors their lifecycle, and reports leaks with stack traces.

**Key Concepts:**
- `runtime.Stack()` for stack traces
- `runtime.NumGoroutine()` for counting
- Goroutine profiling
- Test integration

**Features:**
- Track goroutine creation with stack traces
- Periodic leak checks
- Compare before/after snapshots
- Integration with tests
- Configurable threshold

**Example:**
```go
detector := NewLeakDetector()
detector.Start()

// Take snapshot before test
before := detector.Snapshot()

// Intentional leak #1: blocking channel
go func() {
    <-make(chan struct{}) // Blocks forever
}()

// Intentional leak #2: infinite loop
go func() {
    for {
        time.Sleep(1 * time.Second)
    }
}()

// Proper goroutine (exits)
go func() {
    time.Sleep(100 * time.Millisecond)
}()

time.Sleep(5 * time.Second)

// Check for leaks
leaks := detector.CheckLeaks(before)

Output:
⚠️  Goroutine leak detected!

Found 2 leaked goroutines:

Leak #1 (age: 5s):
Created at: main.go:45
Stack trace:
  main.main.func1()
    main.go:45
  blocking on: chan receive (nil channel)

Leak #2 (age: 5s):
Created at: main.go:50
Stack trace:
  main.main.func2()
    main.go:50
  status: running (infinite loop)

Summary:
- Total goroutines: 5
- Leaked: 2
- Exited properly: 1

// Usage in tests
func TestNoLeaks(t *testing.T) {
    defer detector.AssertNoLeaks(t)

    // Test code that should not leak
}
```

---

### 39. Distributed Rate Limiter

**Description:** Rate limiter that works across multiple service instances using Redis.

**Problem Statement:**
Implement a rate limiter that shares state across distributed services using Redis, ensuring accurate rate limiting even with multiple instances.

**Key Concepts:**
- Redis for shared state
- Sliding window algorithm
- Atomic operations with Lua scripts
- Fallback on Redis failure

**Algorithms:**
- Fixed window counter
- Sliding window log
- Sliding window counter
- Token bucket (distributed)

**Example:**
```go
limiter := NewDistributedRateLimiter(
    redisClient,
    100,           // 100 requests
    time.Minute,   // per minute
)

// Instance 1 (server1):
for i := 0; i < 60; i++ {
    allowed := limiter.Allow("user:123")
    if allowed {
        fmt.Printf("Instance 1: Request %d allowed\n", i)
    } else {
        fmt.Printf("Instance 1: Request %d rejected\n", i)
    }
}

// Instance 2 (server2), simultaneously:
for i := 0; i < 60; i++ {
    allowed := limiter.Allow("user:123")
    if allowed {
        fmt.Printf("Instance 2: Request %d allowed\n", i)
    } else {
        fmt.Printf("Instance 2: Request %d rejected\n", i)
    }
}

Combined Output:
Total requests attempted: 120
Allowed: 100 (split between instances)
Rejected: 20

Redis state:
Key: ratelimit:user:123
Type: ZSET (sorted set with timestamps)
Members: 100 (sliding window entries)
Expiry: 60s

After 1 minute:
Window slides, allowing 100 more requests

// Sliding window visualization
Time: [====|====|====|====] (60s window)
Req:   [30  |25  |25  |20 ] = 100 total
       ^current time

After 15s:
Time: [    |====|====|====|====] (window slides)
Req:       [25  |25  |20  |30 ] = 100 total
                           ^current
```

---

### 40. Build a Simple Database

**Description:** Create a key-value database with indexing, persistence, and concurrency control.

**Problem Statement:**
Build a simple database engine from scratch with B-tree indexing, write-ahead logging, and support for concurrent access.

**Key Concepts:**
- B-tree data structure
- Write-ahead log (WAL)
- `sync.RWMutex` for concurrent access
- File I/O and serialization
- Transaction support

**Features:**
- CRUD operations
- B-tree index for efficient lookups
- Write-ahead log for durability
- Concurrent reads, exclusive writes
- Range scans
- Crash recovery
- Simple transactions

**Example:**
```go
db := NewDatabase("mydb.dat")
db.Open()

// Insert
db.Put("user:1", `{"name":"Alice","age":30}`)
db.Put("user:2", `{"name":"Bob","age":25}`)
db.Put("user:3", `{"name":"Charlie","age":35}`)

// Get
val, err := db.Get("user:1")
Output: {"name":"Alice","age":30}, nil

val, err = db.Get("user:999")
Output: "", ErrNotFound

// Range scan
results := db.Range("user:1", "user:3")
Output: [
  {key: "user:1", value: `{"name":"Alice","age":30}`},
  {key: "user:2", value: `{"name":"Bob","age":25}`},
  {key: "user:3", value: `{"name":"Charlie","age":35}`}
]

// Delete
db.Delete("user:2")
val, _ = db.Get("user:2")
Output: "", ErrNotFound

// Transaction
tx := db.Begin()
tx.Put("user:4", `{"name":"Diana","age":28}`)
tx.Put("user:5", `{"name":"Eve","age":32}`)
tx.Commit()

// Rollback on error
tx2 := db.Begin()
tx2.Put("user:6", `{"name":"Frank","age":40}`)
tx2.Rollback()
val, _ = db.Get("user:6")
Output: "", ErrNotFound

// Close and reopen (test persistence)
db.Close()

db2 := NewDatabase("mydb.dat")
db2.Open()
val, _ = db2.Get("user:1")
Output: {"name":"Alice","age":30} (persisted!)

// Statistics
stats := db.Stats()
Output:
{
  "total_keys": 4,
  "db_size_bytes": 512,
  "btree_height": 2,
  "btree_nodes": 5,
  "wal_entries": 156,
  "wal_size_bytes": 2048
}

// Concurrent access
for i := 0; i < 100; i++ {
    go db.Put(fmt.Sprintf("key:%d", i), fmt.Sprintf("value:%d", i))
}
// All writes are safely serialized
```

---

## Tips for Practice

1. **Start Simple**: Begin with basic implementations, then add optimizations
2. **Write Tests**: Practice writing unit tests for each implementation
3. **Benchmark**: Use Go's benchmarking tools to measure performance
4. **Use Go Idioms**: Follow Go best practices and idiomatic patterns
5. **Error Handling**: Always handle errors explicitly
6. **Documentation**: Write clear comments explaining your approach
7. **Iterate**: Refactor and improve your initial implementations

## Resources

- [Go Documentation](https://go.dev/doc/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go by Example](https://gobyexample.com/)
- [Go Blog](https://go.dev/blog/)

---

Happy coding! Focus on understanding the concepts rather than rushing through implementations.
