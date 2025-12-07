package main

import (
	"fmt"
	"log"
)

func riskyOperation() {
	panic("something went terribly wrong!")
}

func safeCall(fn func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic: %v", r)
			log.Printf("Stack trace would go here")
		}
	}()
	fmt.Println("--executing safecall ---")
	fn()
	return nil
}

func main() {
	// Without recover - program crashes
	// riskyOperation()

	// With recover - error is handled gracefully
	defer func() {
		fmt.Println("-- main defer func ---")
	}()
	err := safeCall(riskyOperation)
	if err != nil {
		fmt.Println("Caught:", err)
	}

	fmt.Println("Program continues...")
}
