package main

import (
	"fmt"
	"time"
)

/*
TIMERS AND TICKERS
==================

  time.Timer  -> fires ONCE after a delay. Its channel .C receives one value.
  time.Ticker -> fires REPEATEDLY on an interval. Read .C in a loop.

Convenience wrappers:
  time.After(d)       -> a channel that receives once after d (great for timeouts)
  time.AfterFunc(d,f) -> runs f in its own goroutine after d, returns a Timer
  time.Tick(d)        -> a ticker channel (leaks; only for programs that run forever)

ALWAYS Stop a Ticker (and a Timer you may not drain) to release its resources —
`defer ticker.Stop()`.

This example uses short durations so it finishes in well under a second.
*/

func main() {
	// --- 1. Timer: fire once after a delay ---
	timer := time.NewTimer(150 * time.Millisecond)
	<-timer.C // block until it fires
	fmt.Println("timer fired after 150ms")

	// --- 2. time.After in a select: the classic timeout pattern ---
	work := make(chan string, 1)
	go func() {
		time.Sleep(80 * time.Millisecond)
		work <- "result"
	}()

	select {
	case r := <-work:
		fmt.Println("got work:", r)
	case <-time.After(200 * time.Millisecond):
		fmt.Println("timed out waiting for work")
	}

	// --- 3. AfterFunc: run a callback later without managing a channel ---
	done := make(chan struct{})
	time.AfterFunc(100*time.Millisecond, func() {
		fmt.Println("AfterFunc callback ran")
		close(done)
	})
	<-done

	// --- 4. Ticker: repeat on an interval, then stop ---
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	ticks := 0
	for range ticker.C {
		ticks++
		fmt.Printf("tick %d\n", ticks)
		if ticks == 5 {
			fmt.Println("stopping ticker")
			return
		}
	}
}
