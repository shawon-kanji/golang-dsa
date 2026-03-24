package main

import (
	"context"
	"fmt"
	"time"
)

// DemoAfterFunc demonstrates context.AfterFunc (Go 1.21+)
// Use this to register a cleanup function that runs when context is done
func DemoAfterFunc() {
	fmt.Println("\n=== context.AfterFunc Demo ===")

	ctx, cancel := context.WithCancel(context.Background())

	// Register a function to be called when context is cancelled
	// Returns a stop function to unregister if needed
	stop := context.AfterFunc(ctx, func() {
		fmt.Println("AfterFunc: Cleanup function called!")
		fmt.Println("AfterFunc: Closing resources, flushing buffers...")
	})

	// The stop function can be used to prevent the callback from running
	_ = stop // We won't use it here, but you can call stop() to unregister

	fmt.Println("Main: Doing some work...")
	time.Sleep(1 * time.Second)

	fmt.Println("Main: Cancelling context...")
	cancel()

	// Give AfterFunc time to execute
	time.Sleep(100 * time.Millisecond)
}

// DemoAfterFuncWithStop shows how to stop the registered function
func DemoAfterFuncWithStop() {
	fmt.Println("\n=== context.AfterFunc with Stop Demo ===")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Register cleanup
	stop := context.AfterFunc(ctx, func() {
		fmt.Println("AfterFunc: This should NOT print!")
	})

	fmt.Println("Main: Unregistering the AfterFunc...")
	wasRegistered := stop() // Returns true if the function was registered and is now stopped
	fmt.Printf("Main: Was registered: %v\n", wasRegistered)

	cancel() // Even after cancellation, the function won't run
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Main: Context cancelled, but AfterFunc didn't run because we stopped it")
}

// DemoAfterFuncMultiple shows multiple AfterFunc registrations
func DemoAfterFuncMultiple() {
	fmt.Println("\n=== Multiple AfterFunc Demo ===")

	ctx, cancel := context.WithCancel(context.Background())

	// Register multiple cleanup functions
	context.AfterFunc(ctx, func() {
		fmt.Println("Cleanup 1: Closing database connection")
	})
	context.AfterFunc(ctx, func() {
		fmt.Println("Cleanup 2: Flushing cache")
	})
	context.AfterFunc(ctx, func() {
		fmt.Println("Cleanup 3: Saving state")
	})

	time.Sleep(500 * time.Millisecond)
	cancel()
	time.Sleep(100 * time.Millisecond)
}
