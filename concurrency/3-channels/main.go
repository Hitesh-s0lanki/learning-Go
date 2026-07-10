package main

import (
	"fmt"
	"time"
)

/*
CHANNELS (unbuffered)
=====================

Channels are Go's way for goroutines to communicate and synchronize.
Slogan: "Don't communicate by sharing memory; share memory by communicating."

An unbuffered channel has NO capacity. This makes it a synchronization point:
  - A send  (ch <- v) blocks until another goroutine receives.
  - A receive (<-ch) blocks until another goroutine sends.

So an unbuffered channel is a "handoff": the sender and receiver must meet.

Syntax:
  ch := make(chan int)   // create a channel of int
  ch <- 42               // send 42 into the channel
  v := <-ch              // receive a value from the channel
*/

func main() {
	// --- Example 1: a single handoff between two goroutines ---
	ch := make(chan string)

	go func() {
		fmt.Println("goroutine: doing work...")
		time.Sleep(1 * time.Second)
		ch <- "result from goroutine" // blocks until main receives
	}()

	fmt.Println("main: waiting for a value on the channel")
	msg := <-ch // blocks until the goroutine sends
	fmt.Println("main: received ->", msg)

	// --- Example 2: using a channel to collect results from many workers ---
	results := make(chan int)
	numbers := []int{2, 4, 6, 8}

	for _, n := range numbers {
		go func(x int) {
			results <- x * x // send the square back
		}(n)
	}

	// We know exactly how many results to expect, so we receive that many times.
	sum := 0
	for range numbers {
		sum += <-results
	}
	fmt.Println("main: sum of squares =", sum)
}
