package main

import "fmt"

/*
BUFFERED CHANNELS
=================

A buffered channel has a capacity. Sends only block when the buffer is FULL;
receives only block when the buffer is EMPTY.

  ch := make(chan int, 3)  // buffer can hold 3 values

Compared to unbuffered channels:
  - Unbuffered: every send waits for a receiver (tight synchronization).
  - Buffered:   the sender can get ahead by up to `cap` values before blocking.

Use buffered channels when you want to decouple producer and consumer speeds,
or as a simple semaphore to limit concurrency.

  len(ch) -> number of values currently queued
  cap(ch) -> total capacity of the buffer
*/

func main() {
	ch := make(chan int, 3)

	// These three sends do NOT block, because the buffer has room for 3.
	ch <- 10
	ch <- 20
	ch <- 30
	fmt.Printf("after 3 sends: len=%d cap=%d\n", len(ch), cap(ch))

	// A 4th send here (without a receiver) would block forever and deadlock,
	// because the buffer is full. Try uncommenting to see the runtime panic:
	// ch <- 40

	// Drain the buffer.
	fmt.Println("received:", <-ch)
	fmt.Println("received:", <-ch)
	fmt.Println("received:", <-ch)
	fmt.Printf("after draining: len=%d cap=%d\n", len(ch), cap(ch))

	// --- Buffered channel as a "semaphore" to limit concurrency ---
	// Only 2 workers may hold a slot at once.
	const maxConcurrent = 2
	sem := make(chan struct{}, maxConcurrent)
	done := make(chan int)

	for i := 1; i <= 5; i++ {
		go func(id int) {
			sem <- struct{}{}        // acquire a slot (blocks if 2 already held)
			fmt.Printf("worker %d running\n", id)
			<-sem                    // release the slot
			done <- id
		}(i)
	}

	for range 5 {
		<-done
	}
	fmt.Println("all workers finished (max 2 ran at once)")
}
