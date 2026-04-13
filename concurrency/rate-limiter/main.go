package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Request struct {
	ID     string
	body   map[string]interface{}
	params map[string]string
	query  map[string]string
}

type Ratelimiter struct {
	mu              sync.Mutex
	strategy        string
	limit           int           // Max requests for both strategies
	availableTokens int           // For token_bucket
	window          time.Duration // For sliding_window
	frequency       time.Duration // For token_bucket
	timestamps      []time.Time   // For sliding_window log
}

func (r *Ratelimiter) MakeRequest(req Request) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if the request can be processed based on the rate limiting strategy
	if r.strategy == "token_bucket" {
		if r.availableTokens > 0 {
			r.availableTokens--
			// Process the request
			return nil
		}

		// No available tokens, reject the request
		return fmt.Errorf("rate limit exceeded")
	} else if r.strategy == "sliding_window" {
		now := time.Now()
		threshold := now.Add(-r.window)

		// Prune timestamps older than the window
		var valid []time.Time
		for _, t := range r.timestamps {
			if t.After(threshold) {
				valid = append(valid, t)
			}
		}
		r.timestamps = valid

		// Check if we are within the limit
		if len(r.timestamps) < r.limit {
			r.timestamps = append(r.timestamps, now)
			return nil
		}

		return fmt.Errorf("rate limit exceeded")
	}

	// Other rate limiting strategies can be implemented here
	return nil
}

func (r *Ratelimiter) Init(ctx context.Context) {
	// Initialize the rate limiter based on the strategy
	if r.strategy == "token_bucket" {
		r.mu.Lock()
		r.availableTokens = r.limit
		r.mu.Unlock()

		go func() {
			ticker := time.NewTicker(r.frequency)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					r.RefillTokens()
				case <-ctx.Done():
					return
				}
			}
		}()
	} else if r.strategy == "sliding_window" {
		r.mu.Lock()
		r.timestamps = make([]time.Time, 0, r.limit)
		r.mu.Unlock()
	}
}

func (r *Ratelimiter) RefillTokens() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.strategy == "token_bucket" {
		r.availableTokens = r.limit
	}
}

func main() {
	rl := Ratelimiter{
		strategy:  "sliding_window",
		limit:     5,
		window:    time.Second * 2, // Allow 5 requests per 2 seconds
		frequency: time.Second * 5, // Keep this if using token_bucket
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rl.Init(ctx)

	for i := 0; i < 50; i++ {
		req := Request{
			ID: fmt.Sprintf("req-%d", i+1),
		}
		err := rl.MakeRequest(req)
		if err != nil {
			fmt.Printf("Request %s failed: %v\n", req.ID, err)
		} else {
			fmt.Printf("Request %s succeeded\n", req.ID)
		}
		time.Sleep(100 * time.Millisecond)
	}
}
