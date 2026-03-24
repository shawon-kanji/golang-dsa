package main

import (
	"context"
	"fmt"
	"time"
)

// DemoWithoutCancel demonstrates context.WithoutCancel (Go 1.21+)
// Use this when you need to continue work even after parent context is cancelled
// Common use case: logging, metrics, cleanup that must complete
func DemoWithoutCancel() {
	fmt.Println("\n=== context.WithoutCancel Demo ===")

	// Create a parent context with timeout
	parentCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Create a child that won't be cancelled
	detachedCtx := context.WithoutCancel(parentCtx)

	// Values are still inherited
	parentCtx = context.WithValue(parentCtx, contextKey("traceID"), "trace-123")
	detachedCtx = context.WithoutCancel(parentCtx)

	fmt.Println("Main: Starting operations...")

	// Start a goroutine with regular context
	go func() {
		select {
		case <-parentCtx.Done():
			fmt.Println("Regular Worker: Parent cancelled, I must stop!")
		case <-time.After(3 * time.Second):
			fmt.Println("Regular Worker: Completed work")
		}
	}()

	// Start a goroutine with detached context
	go func() {
		// This goroutine continues even after parent is cancelled
		select {
		case <-detachedCtx.Done():
			fmt.Println("Detached Worker: Context done (this won't happen!)")
		case <-time.After(2 * time.Second):
			fmt.Println("Detached Worker: Completed critical cleanup work!")
			// Can still access parent's values
			fmt.Printf("Detached Worker: TraceID = %v\n", detachedCtx.Value(contextKey("traceID")))
		}
	}()

	// Wait for both to complete
	time.Sleep(3 * time.Second)
}

// DemoWithoutCancelUseCase shows a practical use case
func DemoWithoutCancelUseCase() {
	fmt.Println("\n=== WithoutCancel Use Case: Audit Logging ===")

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	ctx = context.WithValue(ctx, contextKey("userID"), "user-456")
	defer cancel()

	// Simulate an operation that might timeout
	err := performOperation(ctx)
	if err != nil {
		fmt.Printf("Operation failed: %v\n", err)

		// Even though operation failed, we MUST log it
		// Use WithoutCancel to ensure logging completes
		logCtx := context.WithoutCancel(ctx)
		auditLog(logCtx, "operation_failed", err)
	}

	time.Sleep(1 * time.Second)
}

func performOperation(ctx context.Context) error {
	select {
	case <-time.After(1 * time.Second):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func auditLog(ctx context.Context, event string, err error) {
	// Simulate slow audit logging (e.g., to external service)
	time.Sleep(500 * time.Millisecond)
	userID := ctx.Value(contextKey("userID"))
	fmt.Printf("AUDIT: Event=%s, User=%v, Error=%v\n", event, userID, err)
}
