package main

import (
	"fmt"
	"sync"
	"time"
)

// createStage creates a pipeline stage with multiple workers
// It reads from inputChan, applies the function, and writes to outputChan
func createStage(stageId int, numWorkers int, inputChan <-chan int, outputChan chan<- int, fn func(int) int) {
	var wg sync.WaitGroup

	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for data := range inputChan {
				result := fn(data)
				fmt.Printf("Stage %d, Worker %d: input=%d, output=%d\n", stageId, workerID, data, result)
				if outputChan != nil {
					outputChan <- result
				}
			}
		}(w)
	}

	// Close output channel when all workers are done (in a separate goroutine)
	go func() {
		wg.Wait()
		if outputChan != nil {
			close(outputChan)
		}
		fmt.Printf("Stage %d completed\n", stageId)
	}()
}

func main() {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	numWorkers := 3

	// Create channels to connect stages
	// input -> stage0 -> chan1 -> stage1 -> chan2 -> stage2 -> results
	chan0 := make(chan int, len(input))   // Input channel
	chan1 := make(chan int, len(input))   // Stage 0 -> Stage 1
	chan2 := make(chan int, len(input))   // Stage 1 -> Stage 2
	results := make(chan int, len(input)) // Stage 2 -> Final results

	// Define stage functions
	stageFunctions := []func(int) int{
		// Stage 0: Pass through (identity)
		func(val int) int {
			return val
		},
		// Stage 1: Square the value
		func(val int) int {
			return val * val
		},
		// Stage 2: Filter - keep even, zero out odd
		func(val int) int {
			if val%2 == 0 {
				return val
			}
			return 0
		},
	}

	start := time.Now()

	// Start all stages CONCURRENTLY (this is the key!)
	createStage(0, numWorkers, chan0, chan1, stageFunctions[0])
	createStage(1, numWorkers, chan1, chan2, stageFunctions[1])
	createStage(2, numWorkers, chan2, results, stageFunctions[2])

	// Send input data to the first stage
	go func() {
		for _, v := range input {
			chan0 <- v
		}
		close(chan0) // Close input channel after sending all data
	}()

	// Collect final results
	var finalResults []int
	for r := range results {
		finalResults = append(finalResults, r)
	}

	elapsed := time.Since(start)

	fmt.Println("\n=== Pipeline Complete ===")
	fmt.Printf("Final results: %v\n", finalResults)
	fmt.Printf("Execution time: %s\n", elapsed)

	// // Run the other pipeline pattern example
	// fmt.Println("\n\n--- Running Alternative Pipeline Pattern ---\n")
	// RunPipeline()
}
