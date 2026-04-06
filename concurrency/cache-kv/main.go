package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type Cache struct {
	mu            sync.RWMutex
	items         map[string]*Entry
	maxSize       int
	defaultTTL    time.Duration
	sweepInterval time.Duration
	done          chan struct{}
	closeOnce     sync.Once
	sweeperWG     sync.WaitGroup
	stats         CacheStats
}

type Entry struct {
	value     int
	expiresAt time.Time
	accessAt  time.Time
}

type CacheStats struct {
	getOps       int64
	setOps       int64
	hits         int64
	misses       int64
	evictions    int64
	expiredSwept int64
	expiredLazy  int64
	getLatencyNs int64
	setLatencyNs int64
}

func NewCache(maxSize int, defaultTTL, sweepInterval time.Duration) *Cache {
	c := &Cache{
		items:         make(map[string]*Entry),
		maxSize:       maxSize,
		defaultTTL:    defaultTTL,
		sweepInterval: sweepInterval,
		done:          make(chan struct{}),
	}

	fmt.Printf("[Cache] Started (maxSize=%d, sweepInterval=%s)\n", maxSize, sweepInterval)

	c.sweeperWG.Add(1)
	go c.sweeper()

	return c
}

func (c *Cache) sweeper() {
	defer c.sweeperWG.Done()
	ticker := time.NewTicker(c.sweepInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			now := time.Now()
			removed := 0
			remaining := 0

			c.mu.Lock()
			for k, v := range c.items {
				if now.After(v.expiresAt) {
					delete(c.items, k)
					removed++
				}
			}
			remaining = len(c.items)
			c.mu.Unlock()

			if removed > 0 {
				atomic.AddInt64(&c.stats.expiredSwept, int64(removed))
			}
			fmt.Printf("[Sweeper] Tick - removed %d expired entries (%d remaining)\n", removed, remaining)
		case <-c.done:
			return
		}
	}
}

func (c *Cache) Close() {
	c.closeOnce.Do(func() {
		close(c.done)
		c.sweeperWG.Wait()
		fmt.Println("[Cache] Closed - sweeper stopped")
	})
}

func (c *Cache) Get(key string) (int, bool) {
	start := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()
	defer func() {
		atomic.AddInt64(&c.stats.getOps, 1)
		atomic.AddInt64(&c.stats.getLatencyNs, time.Since(start).Nanoseconds())
	}()

	if v, ok := c.items[key]; ok && time.Now().Before(v.expiresAt) {
		v.accessAt = time.Now()
		atomic.AddInt64(&c.stats.hits, 1)
		return v.value, true
	}

	if v, ok := c.items[key]; ok && time.Now().After(v.expiresAt) {
		delete(c.items, key)
		atomic.AddInt64(&c.stats.expiredLazy, 1)
	}

	atomic.AddInt64(&c.stats.misses, 1)

	return 0, false
}

func (c *Cache) Set(key string, val int, ttl ...time.Duration) *Entry {
	start := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()
	defer func() {
		atomic.AddInt64(&c.stats.setOps, 1)
		atomic.AddInt64(&c.stats.setLatencyNs, time.Since(start).Nanoseconds())
	}()

	effectiveTTL := c.defaultTTL
	if len(ttl) > 0 {
		effectiveTTL = ttl[0]
	}

	if _, exists := c.items[key]; !exists && len(c.items) >= c.maxSize {
		var lru *Entry
		var lruKey string
		for k, v := range c.items {
			if lru == nil || v.accessAt.Before(lru.accessAt) {
				lru = v
				lruKey = k
			}
		}
		if lru != nil {
			delete(c.items, lruKey)
			evictions := atomic.AddInt64(&c.stats.evictions, 1)
			if evictions <= 5 || evictions%5000 == 0 {
				fmt.Printf("[Eviction] Cache full - evicted key=\"%s\" (oldest access)\n", lruKey)
			}
		}
	}

	now := time.Now()

	c.items[key] = &Entry{
		value:     val,
		accessAt:  now,
		expiresAt: now.Add(effectiveTTL),
	}
	return c.items[key]
}

func (c *Cache) SnapshotStats() CacheStats {
	return CacheStats{
		getOps:       atomic.LoadInt64(&c.stats.getOps),
		setOps:       atomic.LoadInt64(&c.stats.setOps),
		hits:         atomic.LoadInt64(&c.stats.hits),
		misses:       atomic.LoadInt64(&c.stats.misses),
		evictions:    atomic.LoadInt64(&c.stats.evictions),
		expiredSwept: atomic.LoadInt64(&c.stats.expiredSwept),
		expiredLazy:  atomic.LoadInt64(&c.stats.expiredLazy),
		getLatencyNs: atomic.LoadInt64(&c.stats.getLatencyNs),
		setLatencyNs: atomic.LoadInt64(&c.stats.setLatencyNs),
	}
}

func avgLatency(totalNs, ops int64) time.Duration {
	if ops == 0 {
		return 0
	}
	return time.Duration(totalNs / ops)
}

func main() {
	cache := NewCache(100, 3*time.Second, 2*time.Second)
	defer cache.Close()

	for i := 0; i < 100; i++ {
		cache.Set(fmt.Sprintf("item-%02d", i), i)
	}

	readerCount := 100
	writerCount := 20
	runFor := 5 * time.Second

	fmt.Printf("[Benchmark] Spawning %d readers + %d writers...\n\n", readerCount, writerCount)

	var wg sync.WaitGroup
	stopAt := time.Now().Add(runFor)

	for i := 0; i < readerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for time.Now().Before(stopAt) {
				key := fmt.Sprintf("item-%02d", rand.Intn(150))
				cache.Get(key)
			}
		}()
	}

	for i := 0; i < writerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for time.Now().Before(stopAt) {
				key := fmt.Sprintf("item-%02d", rand.Intn(150))
				val := rand.Intn(10000)
				ttl := time.Duration(1+rand.Intn(4)) * time.Second
				cache.Set(key, val, ttl)
			}
		}()
	}

	wg.Wait()

	stats := cache.SnapshotStats()
	totalGets := stats.getOps
	hits := stats.hits
	misses := stats.misses
	hitRatio := 0.0
	missRatio := 0.0
	if totalGets > 0 {
		hitRatio = (float64(hits) / float64(totalGets)) * 100
		missRatio = (float64(misses) / float64(totalGets)) * 100
	}

	fmt.Printf("\n========== CACHE BENCHMARK (%s) ==========\n", runFor)
	fmt.Printf("Total Get ops:    %10d\n", stats.getOps)
	fmt.Printf("Total Set ops:    %10d\n", stats.setOps)
	fmt.Printf("Cache hits:       %10d (%0.1f%%)\n", hits, hitRatio)
	fmt.Printf("Cache misses:     %10d (%0.1f%%)\n", misses, missRatio)
	fmt.Printf("Evictions:        %10d\n", stats.evictions)
	fmt.Printf("Expired (swept):  %10d\n", stats.expiredSwept)
	fmt.Printf("Expired (lazy):   %10d\n", stats.expiredLazy)
	fmt.Printf("Avg Get latency:  ~%s\n", avgLatency(stats.getLatencyNs, stats.getOps))
	fmt.Printf("Avg Set latency:  ~%s\n", avgLatency(stats.setLatencyNs, stats.setOps))
	fmt.Println("===========================================")
}
