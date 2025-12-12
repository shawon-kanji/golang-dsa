package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	jobs := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	numWorkers := 3
	var receiver = make(chan int, 2)
	var wg sync.WaitGroup

	start := time.Now()
	for i := 0; i <= numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for jobData := range receiver {
				randomNumber := rand.Intn(5)
				time.Sleep(time.Second * time.Duration(randomNumber))
				fmt.Printf("Worker : %d value: %d\n", workerID, jobData*jobData)
			}
		}(i)
	}

	for _, i := range jobs {
		receiver <- i
	}

	close(receiver)
	wg.Wait()
	elapsed := time.Since(start)

	fmt.Printf("Time of execution :: %s", elapsed)

	RunWorkerPool()
}
