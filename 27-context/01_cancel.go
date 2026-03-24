package main

import (
	"context"
	"fmt"
	"time"
)

// DemoCancel demonstrates context.WithCancel
// Use this when you want to manually cancel an operation
func DemoCancel() {
	fmt.Println("\n=== context.WithCancel Demo ===")

	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	// Start a goroutine that does work
	go func() {
		for i := 1; ; i++ {
			select {
			case <-ctx.Done():
				fmt.Println("Worker: Received cancel signal, stopping...")
				return
			default:
				fmt.Printf("Worker: Doing iteration %d\n", i)
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	// Let the worker run for 2 seconds, then cancel
	time.Sleep(2 * time.Second)
	fmt.Println("Main: Calling cancel()...")
	cancel()

	// Give the goroutine time to print its exit message
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("Main: Context error: %v\n", ctx.Err())
}
