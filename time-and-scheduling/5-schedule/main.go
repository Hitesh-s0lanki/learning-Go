package main

import (
	"fmt"
	"sync"
	"time"
)

/*
A MINI SCHEDULER (composing timers, tickers, and goroutines)
============================================================

This ties the section together: a small in-memory scheduler that can run a task
ONCE after a delay (time.AfterFunc) or REPEATEDLY on an interval (a goroutine
driving a time.Ticker, cancellable via a stop channel).

Key concurrency ideas shown here:
  - each recurring task owns a stopChan; closing it signals the goroutine to exit
  - a select waits on EITHER the ticker OR the stop signal
  - a sync.WaitGroup lets StopAll block until every task goroutine has finished

Unlike a naive scheduler, this one is self-terminating: main schedules some work,
lets it run briefly, then calls StopAll and returns cleanly.
*/

// task is one recurring job the scheduler manages.
type task struct {
	name     string
	stopChan chan struct{}
	wg       sync.WaitGroup
}

// Scheduler runs one-off and recurring tasks.
type Scheduler struct {
	mu    sync.Mutex
	tasks map[string]*task
	wg    sync.WaitGroup // tracks every goroutine (one-off and recurring)
}

func NewScheduler() *Scheduler {
	return &Scheduler{tasks: make(map[string]*task)}
}

func stamp() string { return time.Now().Format("15:04:05.000") }

// Once runs action a single time after delay.
func (s *Scheduler) Once(name string, delay time.Duration, action func()) {
	fmt.Printf("[%s] scheduled one-off %q in %s\n", stamp(), name, delay)
	s.wg.Add(1)
	time.AfterFunc(delay, func() {
		defer s.wg.Done()
		fmt.Printf("[%s] running one-off %q\n", stamp(), name)
		action()
	})
}

// Every runs action every interval (after an initial delay) until stopped.
func (s *Scheduler) Every(name string, initialDelay, interval time.Duration, action func()) {
	t := &task{name: name, stopChan: make(chan struct{})}
	s.mu.Lock()
	s.tasks[name] = t
	s.mu.Unlock()

	fmt.Printf("[%s] scheduled recurring %q every %s\n", stamp(), name, interval)
	s.wg.Add(1)
	t.wg.Add(1)
	go func() {
		defer s.wg.Done()
		defer t.wg.Done()

		// Wait out the initial delay, but bail early if stopped.
		select {
		case <-time.After(initialDelay):
		case <-t.stopChan:
			return
		}

		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			action()
			select {
			case <-ticker.C:
				// loop and run again
			case <-t.stopChan:
				fmt.Printf("[%s] recurring %q stopped\n", stamp(), name)
				return
			}
		}
	}()
}

// StopAll signals every recurring task to stop and waits for all goroutines.
func (s *Scheduler) StopAll() {
	fmt.Printf("[%s] stopping all tasks...\n", stamp())
	s.mu.Lock()
	for _, t := range s.tasks {
		close(t.stopChan)
	}
	s.tasks = make(map[string]*task)
	s.mu.Unlock()

	s.wg.Wait() // block until every goroutine (recurring + one-off) has returned
	fmt.Printf("[%s] all tasks stopped\n", stamp())
}

func main() {
	s := NewScheduler()

	s.Once("send-welcome-email", 200*time.Millisecond, func() {
		fmt.Println("    -> welcome email sent")
	})

	s.Every("db-backup", 100*time.Millisecond, 300*time.Millisecond, func() {
		fmt.Println("    -> database backup written")
	})

	// Let the tasks run for a bit, then shut everything down cleanly.
	time.Sleep(1 * time.Second)
	s.StopAll()
	fmt.Println("scheduler finished")
}
