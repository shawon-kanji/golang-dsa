package main

import (
	"fmt"
	"sync"
	"time"
)

// =====================================================
// Rate Limiter using Token Bucket Algorithm
// =====================================================

type TokenBucketLimiter struct {
	tokens         int           // Current number of tokens
	maxTokens      int           // Maximum tokens in bucket
	refillRate     int           // Tokens added per refill interval
	refillInterval time.Duration // How often to add tokens
	mu             sync.Mutex
	stopChan       chan struct{}
}

// NewTokenBucketLimiter creates a new token bucket rate limiter
func NewTokenBucketLimiter(maxTokens, refillRate int, refillInterval time.Duration) *TokenBucketLimiter {
	limiter := &TokenBucketLimiter{
		tokens:         maxTokens, // Start with full bucket
		maxTokens:      maxTokens,
		refillRate:     refillRate,
		refillInterval: refillInterval,
		stopChan:       make(chan struct{}),
	}

	// Start the refill goroutine
	go limiter.refillTokens()

	return limiter
}

// refillTokens periodically adds tokens to the bucket
func (t *TokenBucketLimiter) refillTokens() {
	ticker := time.NewTicker(t.refillInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			t.mu.Lock()
			t.tokens += t.refillRate
			if t.tokens > t.maxTokens {
				t.tokens = t.maxTokens
			}
			t.mu.Unlock()
		case <-t.stopChan:
			return
		}
	}
}

func (t *TokenBucketLimiter) isAllowed() int {
	return 5
}

// Allow checks if a request is allowed (non-blocking)
func (t *TokenBucketLimiter) Allow() bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.tokens > 0 {
		t.tokens--
		return true
	}
	return false
}

// Wait blocks until a token is available
func (t *TokenBucketLimiter) Wait() {
	for {
		if t.Allow() {
			return
		}
		time.Sleep(10 * time.Millisecond) // Small sleep before retry
	}
}

// Stop stops the refill goroutine
func (t *TokenBucketLimiter) Stop() {
	close(t.stopChan)
}

// =====================================================
// Rate Limiter using Fixed Window Algorithm
// =====================================================

type FixedWindowLimiter struct {
	requests    int           // Current requests in window
	maxRequests int           // Max requests allowed per window
	windowSize  time.Duration // Size of the time window
	windowStart time.Time     // Start of current window
	mu          sync.Mutex
}

// NewFixedWindowLimiter creates a new fixed window rate limiter
func NewFixedWindowLimiter(maxRequests int, windowSize time.Duration) *FixedWindowLimiter {
	return &FixedWindowLimiter{
		requests:    0,
		maxRequests: maxRequests,
		windowSize:  windowSize,
		windowStart: time.Now(),
	}
}

// Allow checks if a request is allowed
func (f *FixedWindowLimiter) Allow() bool {
	f.mu.Lock()
	defer f.mu.Unlock()

	now := time.Now()

	// Check if we're in a new window
	if now.Sub(f.windowStart) >= f.windowSize {
		f.requests = 0
		f.windowStart = now
	}

	if f.requests < f.maxRequests {
		f.requests++
		return true
	}
	return false
}

// =====================================================
// Rate Limiter using time.Ticker (Simple approach)
// =====================================================

type TickerLimiter struct {
	ticker *time.Ticker
	tokens chan struct{}
}

// NewTickerLimiter creates a rate limiter that allows 1 request per interval
func NewTickerLimiter(interval time.Duration, burst int) *TickerLimiter {
	limiter := &TickerLimiter{
		ticker: time.NewTicker(interval),
		tokens: make(chan struct{}, burst),
	}

	// Fill initial burst
	for i := 0; i < burst; i++ {
		limiter.tokens <- struct{}{}
	}

	// Refill tokens at regular intervals
	go func() {
		for range limiter.ticker.C {
			select {
			case limiter.tokens <- struct{}{}:
			default: // Bucket full, discard token
			}
		}
	}()

	return limiter
}

// Wait blocks until allowed to proceed
func (t *TickerLimiter) Wait() {
	<-t.tokens
}

// Stop stops the limiter
func (t *TickerLimiter) Stop() {
	t.ticker.Stop()
}

// =====================================================
// Demo
// =====================================================

func main() {
	fmt.Println("=== Rate Limiter Demo ===\n")

	// Demo 1: Token Bucket Limiter
	fmt.Println("--- Token Bucket Limiter ---")
	fmt.Println("Config: 3 max tokens, refill 1 token every 500ms")

	tokenLimiter := NewTokenBucketLimiter(3, 1, 500*time.Millisecond)
	defer tokenLimiter.Stop()

	for i := 1; i <= 10; i++ {
		start := time.Now()
		if tokenLimiter.Allow() {
			fmt.Printf("Request %d: ALLOWED at %v\n", i, start.Format("15:04:05.000"))
		} else {
			fmt.Printf("Request %d: DENIED at %v\n", i, start.Format("15:04:05.000"))
		}
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("\n--- Fixed Window Limiter ---")
	fmt.Println("Config: 3 requests per 1 second window")

	windowLimiter := NewFixedWindowLimiter(3, 1*time.Second)

	for i := 1; i <= 10; i++ {
		start := time.Now()
		if windowLimiter.Allow() {
			fmt.Printf("Request %d: ALLOWED at %v\n", i, start.Format("15:04:05.000"))
		} else {
			fmt.Printf("Request %d: DENIED at %v\n", i, start.Format("15:04:05.000"))
		}
		time.Sleep(200 * time.Millisecond)
	}

	fmt.Println("\n--- Ticker Limiter (Blocking) ---")
	fmt.Println("Config: 1 request per 300ms, burst of 2")

	tickerLimiter := NewTickerLimiter(300*time.Millisecond, 2)
	defer tickerLimiter.Stop()

	for i := 1; i <= 6; i++ {
		tickerLimiter.Wait() // Blocks until allowed
		fmt.Printf("Request %d: PROCESSED at %v\n", i, time.Now().Format("15:04:05.000"))
	}

	fmt.Println("\n=== Rate Limiter Demo Complete ===")
}
