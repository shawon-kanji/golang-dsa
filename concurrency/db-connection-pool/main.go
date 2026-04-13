package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Pool struct {
	ID                string
	mu                sync.Mutex
	maxConnections    int
	activeConnections int
	cond              *sync.Cond
}

func (p *Pool) GetConnection() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	for {
		if p.activeConnections < p.maxConnections {
			p.activeConnections++
			return true
		}
		fmt.Println("No available connections, waiting...")
		p.cond.Wait()
	}
}

func (p *Pool) ReleaseConnection() {
	p.mu.Lock()
	defer p.mu.Unlock()
	fmt.Println("Releasing a connection..")
	if p.activeConnections > 0 {
		p.activeConnections--
		p.cond.Signal()
	}
}

func QueryDatabase() (interface{}, error) {
	// Simulate a database query
	randomNumber := rand.Intn(100)
	time.Sleep(time.Duration(randomNumber) * time.Millisecond)
	//simulate error at random
	if randomNumber < 20 {
		return nil, fmt.Errorf("database query failed")
	}
	uuid := uuid.New()

	return map[string]string{"id": uuid.String(), "name": "example", "status": "success"}, nil
}

func main() {
	pool := &Pool{
		ID:             "db-pool-1",
		maxConnections: 3,
	}
	pool.cond = sync.NewCond(&pool.mu)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if pool.GetConnection() {
				fmt.Printf("Goroutine %d got a connection\n", id)
				defer pool.ReleaseConnection()
				result, err := QueryDatabase()
				if err != nil {
					fmt.Printf("Goroutine %d: error querying database: %v\n", id, err)
				} else {
					fmt.Printf("Goroutine %d: got result: %v\n", id, result)
				}
			}
		}(i)
	}
	wg.Wait()
}
