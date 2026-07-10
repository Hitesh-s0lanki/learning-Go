package main

import (
	"fmt"
	"sync"
)

/*
CONCURRENCY PATTERNS: PIPELINE and FAN-OUT / FAN-IN
===================================================

Once you know goroutines and channels, you compose them into patterns.

PIPELINE
  A series of stages connected by channels. Each stage is a goroutine that
  receives from an inbound channel, does work, and sends to an outbound one.
      generate -> square -> print

FAN-OUT / FAN-IN
  Fan-out:  start multiple goroutines reading from the SAME input channel to
            parallelize a slow stage.
  Fan-in:   merge multiple result channels back into ONE channel.

These build on the "sender closes the channel" rule from ../5-channel-closing.
*/

// generate is a pipeline SOURCE: it emits the numbers, then closes its output.
func generate(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

// square is a pipeline STAGE: read a number, send back its square.
func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n * n
		}
	}()
	return out
}

// merge implements FAN-IN: combine several channels into one.
func merge(channels ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	for _, c := range channels {
		wg.Add(1)
		go func(ch <-chan int) {
			defer wg.Done()
			for v := range ch {
				out <- v
			}
		}(c)
	}

	// Close the merged channel once every input has been fully drained.
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	// --- Simple pipeline: generate -> square -> print ---
	fmt.Println("pipeline (generate -> square):")
	for result := range square(generate(1, 2, 3, 4, 5)) {
		fmt.Println("  ", result)
	}

	// --- Fan-out / Fan-in ---
	// One source, then TWO square stages reading from it in parallel (fan-out),
	// then merged back into one channel (fan-in).
	fmt.Println("fan-out (2 workers) + fan-in:")
	source := generate(1, 2, 3, 4, 5, 6, 7, 8)

	worker1 := square(source) // both read from the same source channel
	worker2 := square(source)

	sum := 0
	for result := range merge(worker1, worker2) {
		fmt.Println("   got", result)
		sum += result
	}
	fmt.Println("sum of squares:", sum)
}
