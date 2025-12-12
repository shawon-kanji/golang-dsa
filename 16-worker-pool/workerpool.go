package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Job represents a unit of work
type Job struct {
	ID   int
	Data int
}

// Result represents the output of a job
type Result struct {
	JobID    int
	WorkerID int
	Output   int
}

// worker processes jobs from the jobs channel and sends results to the results channel
func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		// Simulate work with random delay
		randomDelay := rand.Intn(3) + 1
		time.Sleep(time.Second * time.Duration(randomDelay))

		result := Result{
			JobID:    job.ID,
			WorkerID: id,
			Output:   job.Data * job.Data,
		}
		results <- result
	}
}

func RunWorkerPool() {
	const numJobs = 10
	const numWorkers = 3

	jobs := make(chan Job, numJobs)
	results := make(chan Result, numJobs)
	var wg sync.WaitGroup

	start := time.Now()

	// Start workers
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	// Send jobs
	for j := 1; j <= numJobs; j++ {
		jobs <- Job{ID: j, Data: j}
	}
	close(jobs)

	// Wait for all workers in a separate goroutine and close results
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	for result := range results {
		fmt.Printf("Job %d processed by Worker %d: %d^2 = %d\n",
			result.JobID, result.WorkerID, result.JobID, result.Output)
	}

	elapsed := time.Since(start)
	fmt.Printf("\nTotal execution time: %s\n", elapsed)
}
