package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

/*
THE context PACKAGE
===================

context is the standard way to carry cancellation, deadlines, and timeouts
across goroutine boundaries. It answers the question: "how do I tell a
goroutine to STOP?"

Creating contexts:
  ctx := context.Background()                              // root, never cancels
  ctx, cancel := context.WithCancel(parent)               // cancel manually
  ctx, cancel := context.WithTimeout(parent, 2*time.Second) // auto-cancel after 2s
  ctx, cancel := context.WithDeadline(parent, someTime)   // auto-cancel at a time

Every goroutine that receives a ctx should:
  - watch  <-ctx.Done()  and return when it fires, and
  - check   ctx.Err()    to learn WHY (Canceled or DeadlineExceeded).

ALWAYS call cancel() (usually with defer) to release resources, even for
timeouts.
*/

// worker runs until it finishes its work OR the context is cancelled.
func worker(ctx context.Context, id int) {
	for i := 1; ; i++ {
		select {
		case <-ctx.Done():
			// ctx.Err() explains why we were told to stop.
			fmt.Printf("worker %d stopping: %v\n", id, ctx.Err())
			return
		default:
			fmt.Printf("worker %d working... step %d\n", id, i)
			time.Sleep(300 * time.Millisecond)
		}
	}
}

func main() {
	// --- Example 1: manual cancellation ---
	ctx, cancel := context.WithCancel(context.Background())
	go worker(ctx, 1)

	time.Sleep(1 * time.Second)
	fmt.Println("main: cancelling worker 1")
	cancel()                          // signal the worker to stop
	time.Sleep(200 * time.Millisecond) // give it a moment to print

	fmt.Println("---")

	// --- Example 2: timeout (context cancels itself after 800ms) ---
	ctx2, cancel2 := context.WithTimeout(context.Background(), 800*time.Millisecond)
	defer cancel2() // good hygiene even though it will time out on its own
	go worker(ctx2, 2)

	<-ctx2.Done() // block until the timeout fires
	time.Sleep(200 * time.Millisecond)

	fmt.Println("---")

	// --- Example 3: a function that respects a deadline ---
	result, err := doWork(ctx2) // ctx2 is already done, so this returns quickly
	if errors.Is(err, context.DeadlineExceeded) {
		fmt.Println("doWork: gave up because the deadline passed")
	} else {
		fmt.Println("doWork result:", result)
	}
}

// doWork simulates an operation that can be cancelled by its context.
func doWork(ctx context.Context) (string, error) {
	select {
	case <-time.After(2 * time.Second):
		return "completed", nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
