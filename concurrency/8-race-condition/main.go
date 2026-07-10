package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

/*
RACE CONDITIONS & THE RACE DETECTOR
===================================

A data race happens when two goroutines access the same variable at the same
time and at least one of them is writing. The result is undefined — you may
lose updates or read garbage.

This file shows THREE versions of "increment a counter 100000 times":

  1. buggy()   -> unsynchronized, loses updates (a real data race).
  2. mutexed() -> protected by a sync.Mutex (correct).
  3. atomicOp() -> uses sync/atomic (correct, lock-free, very fast for counters).

HOW TO DETECT RACES:
  Run with the built-in race detector:

      go run -race .

  It will print a "WARNING: DATA RACE" report pointing at buggy(). Always run
  your concurrent code with -race during development and in CI.
*/

func buggy() int {
	count := 0
	var wg sync.WaitGroup
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 1000 {
				count++ // DATA RACE: read-modify-write with no synchronization
			}
		}()
	}
	wg.Wait()
	return count
}

func mutexed() int {
	count := 0
	var mu sync.Mutex
	var wg sync.WaitGroup
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 1000 {
				mu.Lock()
				count++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	return count
}

func atomicOp() int64 {
	var count int64
	var wg sync.WaitGroup
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 1000 {
				atomic.AddInt64(&count, 1) // atomic: safe without a mutex
			}
		}()
	}
	wg.Wait()
	return atomic.LoadInt64(&count)
}

func main() {
	fmt.Println("buggy (unsynchronized):", buggy(), "<- usually < 100000, and varies each run")
	fmt.Println("mutexed:               ", mutexed(), "<- always 100000")
	fmt.Println("atomic:                ", atomicOp(), "<- always 100000")
	fmt.Println()
	fmt.Println("Run `go run -race .` in this folder to see the race detector flag buggy().")
}
