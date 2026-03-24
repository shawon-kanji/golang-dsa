package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// DemoCancelCause demonstrates context.WithCancelCause (Go 1.20+)
// Use this when you want to provide a specific reason for cancellation
func DemoCancelCause() {
	fmt.Println("\n=== context.WithCancelCause Demo ===")

	ctx, cancel := context.WithCancelCause(context.Background())

	// Start a worker
	go func() {
		select {
		case <-ctx.Done():
			// Get the cause of cancellation
			cause := context.Cause(ctx)
			fmt.Printf("Worker: Stopped! Cause: %v\n", cause)
		}
	}()

	// Simulate some condition that requires cancellation
	time.Sleep(1 * time.Second)

	// Cancel with a specific cause
	customError := errors.New("database connection lost")
	cancel(customError)

	time.Sleep(100 * time.Millisecond)

	// ctx.Err() still returns context.Canceled
	fmt.Printf("ctx.Err(): %v\n", ctx.Err())
	// But context.Cause() returns our custom error
	fmt.Printf("context.Cause(): %v\n", context.Cause(ctx))
}

// DemoTimeoutWithCause shows how timeout contexts also have causes
func DemoTimeoutWithCause() {
	fmt.Println("\n=== Timeout Context Cause Demo ===")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Wait for timeout
	<-ctx.Done()

	fmt.Printf("ctx.Err(): %v\n", ctx.Err())
	fmt.Printf("context.Cause(): %v\n", context.Cause(ctx))
	// For timeout, cause equals the error (DeadlineExceeded)
}
