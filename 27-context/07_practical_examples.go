package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// DemoHTTPStyleContext simulates how context is used in HTTP handlers
func DemoHTTPStyleContext() {
	fmt.Println("\n=== HTTP-Style Context Usage Demo ===")

	// Simulate incoming request with context
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// Add request metadata
	ctx = context.WithValue(ctx, contextKey("requestID"), "req-"+fmt.Sprint(rand.Intn(10000)))
	ctx = context.WithValue(ctx, contextKey("userAgent"), "Mozilla/5.0")

	// Simulate handler -> service -> repository chain
	result, err := httpHandler(ctx)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
		return
	}
	fmt.Printf("Request succeeded: %s\n", result)
}

func httpHandler(ctx context.Context) (string, error) {
	reqID := ctx.Value(contextKey("requestID"))
	fmt.Printf("[Handler] Processing request %v\n", reqID)
	return userService(ctx)
}

func userService(ctx context.Context) (string, error) {
	fmt.Println("[Service] Fetching user data...")
	return userRepository(ctx)
}

func userRepository(ctx context.Context) (string, error) {
	fmt.Println("[Repository] Querying database...")

	// Simulate DB query
	select {
	case <-time.After(1 * time.Second):
		return "User: John Doe", nil
	case <-ctx.Done():
		return "", fmt.Errorf("database query cancelled: %w", ctx.Err())
	}
}

// DemoParallelWithContext shows cancelling parallel operations
func DemoParallelWithContext() {
	fmt.Println("\n=== Parallel Operations with Context Demo ===")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	results := make(chan string, 3)
	errors := make(chan error, 3)

	services := []string{"ServiceA", "ServiceB", "ServiceC"}

	for _, svc := range services {
		wg.Add(1)
		go func(serviceName string) {
			defer wg.Done()
			result, err := callService(ctx, serviceName)
			if err != nil {
				errors <- err
				return
			}
			results <- result
		}(svc)
	}

	// Wait in a separate goroutine
	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	// Collect results
	for result := range results {
		fmt.Printf("Got result: %s\n", result)
	}
	for err := range errors {
		fmt.Printf("Got error: %v\n", err)
	}
}

func callService(ctx context.Context, name string) (string, error) {
	// Random delay between 500ms and 2500ms
	delay := time.Duration(500+rand.Intn(2000)) * time.Millisecond

	select {
	case <-time.After(delay):
		return fmt.Sprintf("%s responded in %v", name, delay), nil
	case <-ctx.Done():
		return "", fmt.Errorf("%s cancelled: %w", name, ctx.Err())
	}
}

// DemoPropagation shows how cancellation propagates through context hierarchy
func DemoPropagation() {
	fmt.Println("\n=== Context Propagation Demo ===")

	// Root context
	rootCtx, rootCancel := context.WithCancel(context.Background())

	// Child context (inherits cancellation from root)
	childCtx, childCancel := context.WithCancel(rootCtx)
	defer childCancel()

	// Grandchild context
	grandchildCtx, grandchildCancel := context.WithCancel(childCtx)
	defer grandchildCancel()

	// Monitor all contexts
	go monitor("Root", rootCtx)
	go monitor("Child", childCtx)
	go monitor("Grandchild", grandchildCtx)

	time.Sleep(500 * time.Millisecond)
	fmt.Println("\nMain: Cancelling ROOT context...")
	rootCancel()

	time.Sleep(500 * time.Millisecond)
	fmt.Println("\nNote: All child contexts were cancelled when root was cancelled!")
}

func monitor(name string, ctx context.Context) {
	<-ctx.Done()
	fmt.Printf("  [%s] Context done! Error: %v\n", name, ctx.Err())
}

// DemoGracefulShutdown simulates graceful shutdown pattern
func DemoGracefulShutdown() {
	fmt.Println("\n=== Graceful Shutdown Demo ===")

	// Main context for all operations
	ctx, cancel := context.WithCancel(context.Background())

	// Start some workers
	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go worker(ctx, &wg, i)
	}

	// Simulate running for a while
	time.Sleep(1500 * time.Millisecond)

	// Initiate shutdown
	fmt.Println("\nMain: Initiating graceful shutdown...")
	cancel()

	// Wait for all workers to finish
	wg.Wait()
	fmt.Println("Main: All workers stopped gracefully!")
}

func worker(ctx context.Context, wg *sync.WaitGroup, id int) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d: Cleaning up and stopping...\n", id)
			time.Sleep(100 * time.Millisecond) // Simulate cleanup
			return
		default:
			fmt.Printf("Worker %d: Processing task...\n", id)
			time.Sleep(400 * time.Millisecond)
		}
	}
}
