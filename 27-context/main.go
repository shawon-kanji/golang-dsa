package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	// 1. Create a root context
	ctx := context.Background()

	// 2. Create a child context with a timeout of 2 seconds.
	// The `cancel` function is returned to release resources if we finish early.
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel() // Always defer cancel to avoid memory leaks!

	fmt.Println("Main: Starting operation...")

	// 3. Pass the context into the function doing the work
	result, err := slowOperation(ctx)
	if err != nil {
		fmt.Printf("Main: Operation failed: %v\n", err)
	} else {
		fmt.Printf("Main: Operation success: %s\n", result)
	}
}

// slowOperation simulates a task that takes time
func slowOperation(ctx context.Context) (string, error) {
	// A channel to simulate the work finishing
	done := make(chan string)

	go func() {
		// Simulate work that takes 4 seconds (longer than our 2s timeout)
		time.Sleep(4 * time.Second)
		done <- "Work Complete!"
		close(done)
	}()

	// for i := range done {
	// 	fmt.Print(i)
	// }

	// The select statement blocks until one of the cases is ready.
	select {
	case res := <-done:
		// The work finished successfully
		return res, nil
	case <-ctx.Done():
		// The context timed out or was cancelled before work finished
		return "", ctx.Err()
	}
	// return "", nil
}
