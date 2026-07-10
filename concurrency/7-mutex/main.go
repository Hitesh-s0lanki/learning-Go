package main

import (
	"fmt"
	"sync"
)

/*
MUTEX & SYNCHRONIZATION
=======================

Channels are great for passing ownership of data. But sometimes you just need
several goroutines to safely update SHARED state (a counter, a map, a cache).
For that, use a mutex ("mutual exclusion") from the sync package.

  var mu sync.Mutex
  mu.Lock()   // only one goroutine can hold the lock at a time
  // ... critical section: touch shared state here ...
  mu.Unlock()

sync.RWMutex is an optimization: many readers OR one writer.
  mu.RLock()/mu.RUnlock() for reads, mu.Lock()/mu.Unlock() for writes.

A clean pattern is to bundle the mutex with the data it protects.
*/

// SafeCounter can be incremented from many goroutines safely.
type SafeCounter struct {
	mu    sync.Mutex
	count int
}

func (c *SafeCounter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock() // defer guarantees the lock is released
	c.count++
}

func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

func main() {
	counter := &SafeCounter{}
	var wg sync.WaitGroup

	// 100 goroutines each incrementing 1000 times = 100000 total.
	// Without the mutex this would lose updates (see ../8-race-condition).
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 1000 {
				counter.Inc()
			}
		}()
	}

	wg.Wait()
	fmt.Println("final count:", counter.Value(), "(expected 100000)")

	// --- RWMutex example: a concurrent-safe map with many readers ---
	cache := struct {
		sync.RWMutex
		data map[string]int
	}{data: make(map[string]int)}

	cache.Lock()
	cache.data["answer"] = 42
	cache.Unlock()

	var rwg sync.WaitGroup
	for i := range 5 {
		rwg.Add(1)
		go func(id int) {
			defer rwg.Done()
			cache.RLock() // multiple readers can hold RLock at once
			v := cache.data["answer"]
			cache.RUnlock()
			fmt.Printf("reader %d saw answer=%d\n", id, v)
		}(i)
	}
	rwg.Wait()
}
