package main

import (
	"container/heap"
	"context"
	"fmt"
	"hash/fnv"
	"sync"
	"time"
)

type Priority string

var HIGH Priority = "HIGH"
var MEDIUM Priority = "MEDIUM"
var LOW Priority = "LOW"

type Task struct {
	ID       string
	Priority Priority
	Deadline time.Duration
}

type Result struct {
	TaskID              string
	Output              string
	Priority            Priority
	processedByWorkerId int
	Error               error
}

func priorityRank(p Priority) int {
	switch p {
	case HIGH:
		return 3
	case MEDIUM:
		return 2
	case LOW:
		return 1
	default:
		return 0
	}
}

type TaskPriorityQueue []Task

func (pq TaskPriorityQueue) Len() int { return len(pq) }

func (pq TaskPriorityQueue) Less(i, j int) bool {
	ri := priorityRank(pq[i].Priority)
	rj := priorityRank(pq[j].Priority)

	// Higher priority first
	if ri != rj {
		return ri > rj
	}

	// Tie-breaker: earlier deadline first
	return pq[i].Deadline < pq[j].Deadline
}

func (pq TaskPriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }

func (pq *TaskPriorityQueue) Push(x any) {
	*pq = append(*pq, x.(Task))
}

func (pq *TaskPriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

func taskContextWithDeadline(deadline time.Duration) (context.Context, context.CancelFunc) {
	return context.WithDeadline(context.Background(), time.Now().Add(deadline))
}

func simulatedProcessingDuration(task Task) time.Duration {
	var base time.Duration
	switch task.Priority {
	case HIGH:
		base = 300 * time.Millisecond
	case MEDIUM:
		base = 800 * time.Millisecond
	case LOW:
		base = 1300 * time.Millisecond
	default:
		base = 700 * time.Millisecond
	}

	h := fnv.New32a()
	_, _ = h.Write([]byte(task.ID))
	jitter := time.Duration(h.Sum32()%900) * time.Millisecond

	return base + jitter
}

func main() {
	WORKERS := 5
	taskQueue := make(chan Task)
	result := make(chan Result, WORKERS)
	start := time.Now()
	taskArray := []Task{
		{ID: "task1", Priority: HIGH, Deadline: 2 * time.Second},
		{ID: "task2", Priority: MEDIUM, Deadline: 1 * time.Second},
		{ID: "task3", Priority: HIGH, Deadline: 3 * time.Second},
		{ID: "task4", Priority: LOW, Deadline: 500 * time.Millisecond},
		{ID: "task5", Priority: HIGH, Deadline: 4 * time.Second},
		{ID: "task6", Priority: LOW, Deadline: 700 * time.Millisecond},
		{ID: "task7", Priority: MEDIUM, Deadline: 1500 * time.Millisecond},
		{ID: "task8", Priority: HIGH, Deadline: 2500 * time.Millisecond},
		{ID: "task9", Priority: MEDIUM, Deadline: 2200 * time.Millisecond},
		{ID: "task10", Priority: LOW, Deadline: 900 * time.Millisecond},
		{ID: "task11", Priority: HIGH, Deadline: 5 * time.Second},
		{ID: "task12", Priority: MEDIUM, Deadline: 2800 * time.Millisecond},
		{ID: "task13", Priority: LOW, Deadline: 1200 * time.Millisecond},
		{ID: "task14", Priority: HIGH, Deadline: 3500 * time.Millisecond},
		{ID: "task15", Priority: MEDIUM, Deadline: 1800 * time.Millisecond},
		{ID: "task16", Priority: LOW, Deadline: 600 * time.Millisecond},
		{ID: "task17", Priority: HIGH, Deadline: 4200 * time.Millisecond},
		{ID: "task18", Priority: MEDIUM, Deadline: 2 * time.Second},
		{ID: "task19", Priority: LOW, Deadline: 1300 * time.Millisecond},
		{ID: "task20", Priority: HIGH, Deadline: 2700 * time.Millisecond},
	}

	wg := sync.WaitGroup{}
	fmt.Printf("[Scheduler] Started with %d workers and %d queued tasks\n", WORKERS, len(taskArray))
	for i := range WORKERS {
		wg.Add(1)
		go func(i int) {
			fmt.Printf("[Worker-%d] Ready\n", i+1)
			for task := range taskQueue {
				ctx, cancel := taskContextWithDeadline(task.Deadline)
				processingTime := simulatedProcessingDuration(task)
				fmt.Printf("[Worker-%d] Executing %s (priority=%s, work=%s, deadline=%s)\n", i+1, task.ID, task.Priority, processingTime, task.Deadline)
				workDone := time.After(processingTime)
				select {
				case <-ctx.Done():
					fmt.Printf("[Worker-%d] %s TIMEOUT (%v)\n", i+1, task.ID, ctx.Err())
					result <- Result{
						TaskID:              task.ID,
						Output:              fmt.Sprintf("task %s failed due to deadline", task.ID),
						Priority:            task.Priority,
						processedByWorkerId: i,
						Error:               fmt.Errorf("task %s failed due to deadline", task.ID),
					}
				case <-workDone:
					fmt.Printf("[Worker-%d] %s completed\n", i+1, task.ID)
					result <- Result{
						TaskID:              task.ID,
						Output:              fmt.Sprintf("result of task %s", task.ID),
						Priority:            task.Priority,
						processedByWorkerId: i,
					}
				}
				cancel()
			}
			wg.Done()
		}(i)
	}

	// go func() {
	// 	for i, task := range taskArray {
	// 		fmt.Printf("Submitting task %d with priority %s and deadline %s\n", i, task.Priority, task.Deadline)
	// 		taskQueue <- task
	// 	}
	// 	close(taskQueue)
	// }()

	go func() {
		pq := &TaskPriorityQueue{}
		heap.Init(pq)
		for _, task := range taskArray {
			heap.Push(pq, task)
		}
		idx := 0
		for pq.Len() > 0 {
			task := heap.Pop(pq).(Task)
			fmt.Printf("[Scheduler] Submitted #%02d -> %s (priority=%s, deadline=%s)\n", idx+1, task.ID, task.Priority, task.Deadline)
			taskQueue <- task
			idx++
		}
		fmt.Println("[Scheduler] Submission complete, closing task queue")
		close(taskQueue)
	}()

	go func() {
		wg.Wait()
		fmt.Println("[Scheduler] All workers drained, closing results channel")
		close(result)
	}()
	success := 0
	failed := 0
	for res := range result {
		status := "OK"
		if res.Error != nil {
			status = "FAILED"
		}
		fmt.Printf("[Result] %s | task=%s priority=%s worker=%d\n", status, res.TaskID, res.Priority, res.processedByWorkerId+1)
		if res.Error != nil {
			failed++
		} else {
			success++
		}
	}

	fmt.Println("\n========== SCHEDULER REPORT ==========")
	fmt.Printf("Tasks submitted:  %d\n", len(taskArray))
	fmt.Printf("Tasks completed:  %d\n", success)
	fmt.Printf("Tasks failed:     %d\n", failed)
	fmt.Printf("Total duration:   %s\n", time.Since(start).Round(time.Millisecond))
	fmt.Println("======================================")

}
