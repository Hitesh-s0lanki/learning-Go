package main

import (
	"fmt"
	"time"
)

/*
GOROUTINES
==========

A goroutine is a lightweight thread managed by the Go runtime (not the OS).
You start one by putting the keyword `go` before a function call.

Key ideas:
  1. `go f()` runs f() concurrently and returns immediately — it does NOT wait.
  2. The `main` function itself runs in the "main goroutine".
  3. When main() returns, the program exits — even if other goroutines
     are still running. That is why we often need to wait (here with Sleep,
     later with sync.WaitGroup — see ../2-waitgroup).
  4. Goroutines are cheap: you can spawn thousands. They start with a tiny
     stack (~2KB) that grows and shrinks as needed.
*/

func sayHello(message string, delay time.Duration) {
	time.Sleep(delay)
	fmt.Println("goroutine says:", message)
}

func main() {
	fmt.Println("main goroutine: start")

	// Each of these runs concurrently. Order of output is NOT guaranteed —
	// it depends on the delays and how the scheduler runs them.
	go sayHello("hello after 1s (A)", 1*time.Second)
	go sayHello("hello after 1s (B)", 1*time.Second)
	go sayHello("hello after 2s", 2*time.Second)
	go sayHello("hello after 3s", 3*time.Second)

	// This prints immediately, before the goroutines above finish, because
	// `go` does not block.
	fmt.Println("main goroutine: launched all goroutines (did not wait)")

	// If we removed this Sleep, main() would return right away and the
	// program would exit before any goroutine printed anything. This is a
	// crude way to wait — real code uses sync.WaitGroup or channels.
	time.Sleep(4 * time.Second)

	fmt.Println("main goroutine: done")
}
