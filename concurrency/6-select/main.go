package main

import (
	"fmt"
	"time"
)

/*
THE select STATEMENT
====================

`select` lets a goroutine wait on MULTIPLE channel operations at once.
It blocks until one of its cases can proceed; if several are ready, it picks
one at random.

  select {
  case v := <-ch1:      // ready when ch1 has a value
  case ch2 <- x:        // ready when ch2 can accept a value
  case <-time.After(d): // a timeout (time.After returns a channel)
  default:              // runs immediately if no other case is ready
  }

Common uses:
  - Combining several channels (e.g. results + a quit signal).
  - Timeouts with time.After.
  - Non-blocking send/receive with `default`.
*/

func main() {
	// --- Example 1: receive from whichever channel is ready first ---
	c1 := make(chan string)
	c2 := make(chan string)

	go func() { time.Sleep(300 * time.Millisecond); c1 <- "from c1" }()
	go func() { time.Sleep(600 * time.Millisecond); c2 <- "from c2" }()

	for range 2 {
		select {
		case msg := <-c1:
			fmt.Println("select got:", msg)
		case msg := <-c2:
			fmt.Println("select got:", msg)
		}
	}

	// --- Example 2: timeout ---
	slow := make(chan string)
	go func() { time.Sleep(2 * time.Second); slow <- "slow result" }()

	select {
	case r := <-slow:
		fmt.Println("received:", r)
	case <-time.After(500 * time.Millisecond):
		fmt.Println("timed out waiting for slow result")
	}

	// --- Example 3: non-blocking receive with default ---
	empty := make(chan int)
	select {
	case v := <-empty:
		fmt.Println("got", v)
	default:
		fmt.Println("nothing ready, not blocking (default case)")
	}

	// --- Example 4: a ticker loop that stops on a quit signal ---
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()
	quit := make(chan struct{})

	go func() {
		time.Sleep(1 * time.Second)
		close(quit)
	}()

	for {
		select {
		case t := <-ticker.C:
			fmt.Println("tick at", t.Format("15:04:05.000"))
		case <-quit:
			fmt.Println("quit signal received, stopping ticker loop")
			return
		}
	}
}
