package main

import (
	"fmt"
	"sync"
)

/*
CLOSING CHANNELS & RANGE
========================

Closing a channel signals "no more values will ever be sent".

  close(ch)

Rules:
  1. Only the SENDER should close a channel, never the receiver.
  2. Sending on a closed channel PANICS.
  3. Receiving from a closed channel never blocks: it returns the zero value
     immediately once drained.
  4. The two-value receive tells you if the channel is still open:
        v, ok := <-ch   // ok == false means the channel is closed and drained
  5. `for v := range ch` loops over values until the channel is closed.

This is the idiomatic producer/consumer pattern.
*/

func producer(ch chan<- int, count int) {
	// chan<- int means "send-only" channel (a good habit for clarity).
	for i := 1; i <= count; i++ {
		ch <- i * i
	}
	close(ch) // tell the consumer we're done
}

func main() {
	// --- Example 1: range over a channel until it is closed ---
	ch := make(chan int)
	go producer(ch, 5)

	fmt.Println("ranging over channel:")
	for v := range ch { // stops automatically when producer calls close(ch)
		fmt.Println("  got", v)
	}

	// --- Example 2: the comma-ok form to detect a closed channel ---
	ch2 := make(chan string)
	go func() {
		ch2 <- "a"
		ch2 <- "b"
		close(ch2)
	}()

	for {
		v, ok := <-ch2
		if !ok {
			fmt.Println("ch2 is closed, stopping")
			break
		}
		fmt.Println("received:", v)
	}

	// --- Example 3: closing a channel to broadcast a signal to many goroutines ---
	// A closed channel unblocks ALL receivers at once, so an empty struct
	// channel is a common "done"/"quit" broadcast mechanism.
	done := make(chan struct{})
	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			<-done // every goroutine blocks here...
			fmt.Printf("goroutine %d released by close(done)\n", id)
		}(i)
	}

	fmt.Println("broadcasting stop signal by closing done channel")
	close(done) // ...and all are released at once
	wg.Wait()
}
