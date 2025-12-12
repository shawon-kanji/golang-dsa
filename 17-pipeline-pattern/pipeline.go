package main

import (
	"fmt"
	"time"
)

// Stage 1: Generator - produces numbers and sends them to the output channel
func generator(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			fmt.Printf("Stage 1 (Generator): Sending %d\n", n)
			out <- n
		}
		close(out)
	}()
	return out
}

// Stage 2: Square - receives numbers, squares them, sends to output channel
func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			result := n * n
			fmt.Printf("Stage 2 (Square): %d -> %d\n", n, result)
			out <- result
		}
		close(out)
	}()
	return out
}

// Stage 3: Filter - receives numbers, filters even ones, sends to output channel
func filterEven(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			if n%2 == 0 {
				fmt.Printf("Stage 3 (Filter): %d is even, keeping\n", n)
				out <- n
			} else {
				fmt.Printf("Stage 3 (Filter): %d is odd, discarding\n", n)
			}
		}
		close(out)
	}()
	return out
}

// Stage 4: Consumer - final stage that collects results
func consumer(in <-chan int) []int {
	var results []int
	for n := range in {
		fmt.Printf("Stage 4 (Consumer): Received %d\n", n)
		results = append(results, n)
	}
	return results
}

// RunPipeline demonstrates the pipeline pattern
func RunPipeline() {
	fmt.Println("=== Pipeline Pattern Demo ===\n")

	start := time.Now()

	// Build the pipeline: generator -> square -> filterEven -> consumer
	// All stages run concurrently!

	// Stage 1: Generate numbers
	numbers := generator(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	// Stage 2: Square each number (runs concurrently with Stage 1)
	squared := square(numbers)

	// Stage 3: Filter even numbers (runs concurrently with Stage 1 & 2)
	filtered := filterEven(squared)

	// Stage 4: Consume results (runs concurrently with all previous stages)
	results := consumer(filtered)

	elapsed := time.Since(start)

	fmt.Println("\n=== Final Results ===")
	fmt.Printf("Even squares: %v\n", results)
	fmt.Printf("Execution time: %s\n", elapsed)
}
