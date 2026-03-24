package main

import (
	"context"
	"fmt"
	"time"
)

// DemoDeadline demonstrates context.WithDeadline
// Use this when you have an absolute time by which an operation must complete
func DemoDeadline() {
	fmt.Println("\n=== context.WithDeadline Demo ===")

	// Set a deadline 2 seconds from now
	deadline := time.Now().Add(2 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	// Print the deadline
	if d, ok := ctx.Deadline(); ok {
		fmt.Printf("Deadline set to: %v\n", d.Format("15:04:05.000"))
	}

	// Simulate work that takes 3 seconds
	select {
	case <-time.After(3 * time.Second):
		fmt.Println("Work completed!")
	case <-ctx.Done():
		fmt.Printf("Deadline exceeded: %v\n", ctx.Err())
	}
}

// DemoDeadlineSuccess shows a deadline that doesn't expire
func DemoDeadlineSuccess() {
	fmt.Println("\n=== context.WithDeadline Success Demo ===")

	deadline := time.Now().Add(3 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	// Simulate work that completes in 1 second
	select {
	case <-time.After(1 * time.Second):
		fmt.Println("Work completed before deadline!")
	case <-ctx.Done():
		fmt.Printf("Deadline exceeded: %v\n", ctx.Err())
	}
}
