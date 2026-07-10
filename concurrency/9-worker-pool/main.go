package main

import (
	"fmt"
	"sync"
	"time"
)

/*
WORKER POOL PATTERN (mini project)
==================================

A worker pool limits how many goroutines process work at once. Instead of
spawning one goroutine per job (which could be millions), you start a FIXED
number of workers that all pull jobs from a shared `jobs` channel and push
results to a shared `results` channel.

Flow:
  main ---> [ jobs channel ] ---> workers (N of them) ---> [ results channel ] ---> main

This is one of the most useful real-world concurrency patterns: rate-limiting,
bounded parallelism, and clean shutdown all fall out of it naturally.
*/

type Job struct {
	ID     int
	Number int
}

type Result struct {
	JobID  int
	Square int
}

// worker pulls jobs until the jobs channel is closed, then returns.
func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs { // exits automatically when jobs is closed
		fmt.Printf("worker %d processing job %d (n=%d)\n", id, job.ID, job.Number)
		time.Sleep(200 * time.Millisecond) // simulate work
		results <- Result{JobID: job.ID, Square: job.Number * job.Number}
	}
}

func main() {
	const numWorkers = 3
	const numJobs = 9

	jobs := make(chan Job, numJobs)
	results := make(chan Result, numJobs)

	// 1. Start a fixed pool of workers.
	var wg sync.WaitGroup
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	// 2. Send all the jobs, then close the jobs channel so workers know to stop.
	for j := 1; j <= numJobs; j++ {
		jobs <- Job{ID: j, Number: j}
	}
	close(jobs)

	// 3. Close the results channel once ALL workers are done.
	//    We do this in a separate goroutine so the range below can start
	//    receiving immediately instead of deadlocking.
	go func() {
		wg.Wait()
		close(results)
	}()

	// 4. Collect the results (order is not guaranteed).
	total := 0
	for r := range results {
		fmt.Printf("  result: job %d -> %d\n", r.JobID, r.Square)
		total += r.Square
	}
	fmt.Println("sum of all squares:", total)
}
