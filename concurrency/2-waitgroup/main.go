package main

import (
	"fmt"
	"sync"
	"time"
)

/*
sync.WaitGroup
==============

Sleeping to "wait" for goroutines (as in ../1-goroutines) is fragile — you
never really know how long they take. A sync.WaitGroup lets the main goroutine
block until a set of goroutines has finished.

Think of it as a counter:
  - wg.Add(n)  -> increase the counter by n (jobs still running)
  - wg.Done()  -> decrease the counter by 1 (one job finished)
  - wg.Wait()  -> block until the counter reaches 0

RULES (very important):
  1. Call wg.Add BEFORE starting the goroutine (not inside it) so there is no
     race between Add and Wait.
  2. Call wg.Done inside the goroutine, ideally with `defer` so it always runs.
  3. Always pass the WaitGroup by POINTER (*sync.WaitGroup). Copying it breaks
     the shared counter.
*/

func worker(id int, wg *sync.WaitGroup) {
	// defer guarantees Done runs even if the function returns early or panics.
	defer wg.Done()

	fmt.Printf("worker %d: starting\n", id)
	time.Sleep(time.Duration(id) * 200 * time.Millisecond)
	fmt.Printf("worker %d: finished\n", id)
}

func main() {
	var wg sync.WaitGroup

	totalJobs := 5

	for i := 1; i <= totalJobs; i++ {
		wg.Add(1) // register one job BEFORE launching it
		go worker(i, &wg)
	}

	fmt.Println("main: all workers launched, waiting for them to finish...")

	wg.Wait() // blocks here until every worker has called Done()

	fmt.Println("main: all workers done")
}
