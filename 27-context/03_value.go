package main

import (
	"context"
	"fmt"
)

// Custom types for context keys to avoid collisions
type contextKey string

const (
	userIDKey    contextKey = "userID"
	requestIDKey contextKey = "requestID"
	roleKey      contextKey = "role"
)

// DemoValue demonstrates context.WithValue
// Use this to pass request-scoped data through the call chain
func DemoValue() {
	fmt.Println("\n=== context.WithValue Demo ===")

	// Create a context with values (chained)
	ctx := context.Background()
	ctx = context.WithValue(ctx, userIDKey, "user-123")
	ctx = context.WithValue(ctx, requestIDKey, "req-abc-456")
	ctx = context.WithValue(ctx, roleKey, "admin")

	// Simulate passing context through multiple layers
	handleRequest(ctx)
}

func handleRequest(ctx context.Context) {
	fmt.Println("Handler: Processing request...")

	// Extract values from context
	userID := ctx.Value(userIDKey)
	requestID := ctx.Value(requestIDKey)

	fmt.Printf("Handler: UserID=%v, RequestID=%v\n", userID, requestID)

	// Pass to next layer
	processData(ctx)
}

func processData(ctx context.Context) {
	fmt.Println("Processor: Processing data...")

	// Access values without knowing they were set by handler
	role := ctx.Value(roleKey)
	userID := ctx.Value(userIDKey)

	if role == "admin" {
		fmt.Printf("Processor: User %v has admin privileges\n", userID)
	}

	// Trying to access a non-existent key returns nil
	nonExistent := ctx.Value("nonexistent")
	fmt.Printf("Processor: Non-existent key returns: %v\n", nonExistent)
}

// DemoValueBestPractices shows best practices for context values
func DemoValueBestPractices() {
	fmt.Println("\n=== context.WithValue Best Practices ===")

	// BAD: Using string keys (can cause collisions)
	// ctx = context.WithValue(ctx, "userID", "123") // DON'T DO THIS

	// GOOD: Using typed keys
	type myKey struct{}
	ctx := context.WithValue(context.Background(), myKey{}, "safe-value")

	// Values are immutable - creating a new value doesn't affect parent
	ctxChild := context.WithValue(ctx, userIDKey, "child-user")

	fmt.Printf("Parent context userID: %v\n", ctx.Value(userIDKey))
	fmt.Printf("Child context userID: %v\n", ctxChild.Value(userIDKey))
	fmt.Printf("Child can access parent's value: %v\n", ctxChild.Value(myKey{}))
}
