package main

import (
	"fmt"
	"sync"
)

type Singleflight struct {
	mu       sync.Mutex
	inflight map[string]*call
}

type call struct {
	wg  sync.WaitGroup
	val any
	err error
}

func (g *Singleflight) Do(key string, fn func() (any, error)) (any, error, bool) {
	g.mu.Lock()
	if c, ok := g.inflight[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err, true
	}

	c := &call{}
	c.wg.Add(1)
	g.inflight[key] = c
	g.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.inflight, key)
	g.mu.Unlock()

	return c.val, c.err, false
}

func main() {
	workers := 49
	result := make(chan int, workers)

	sf := Singleflight{
		inflight: make(map[string]*call),
	}

	backend := map[string]int{
		"key1": 42,
		"key2": 55,
	}

	var cacheMu sync.RWMutex
	cache := make(map[string]int)

	var wg sync.WaitGroup
	for i := 1; i <= workers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			key := "key1"

			cacheMu.RLock()
			v, ok := cache[key]
			cacheMu.RUnlock()
			if ok {
				fmt.Println("Worker", i, "reading", key, "from cache")
				result <- v
				return
			}

			val, err, shared := sf.Do(key, func() (any, error) {
				fmt.Println("Worker", i, "reading", key, "from datasource")
				return backend[key], nil
			})
			if err != nil {
				return
			}

			intVal, ok := val.(int)
			if !ok {
				fmt.Println("Worker", i, "type assertion failed")
				return
			}

			cacheMu.Lock()
			cache[key] = intVal
			cacheMu.Unlock()

			if shared {
				fmt.Println("Worker", i, "received", key, "from in-flight request")
			}

			result <- intVal
		}(i)
	}

	wg.Wait()
	close(result)

	for v := range result {
		fmt.Println("Result:", v)
	}

}
