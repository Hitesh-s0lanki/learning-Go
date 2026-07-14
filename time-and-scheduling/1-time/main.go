package main

import (
	"fmt"
	"time"
)

/*
TIMES AND DURATIONS (the time package)
======================================

Two core types:

  time.Time     -> an instant (a specific moment: 2009-11-10 23:00:00 UTC)
  time.Duration -> a length of time (2 hours, 100ms) stored as an int64 of
                   nanoseconds

Build durations by multiplying the unit constants:
  time.Second, time.Minute, time.Hour, time.Millisecond, ...
  5 * time.Minute, (1*time.Hour)+(30*time.Minute)

Move between instants with a duration:
  t.Add(d)   -> a new Time, d later (use a negative d to go back)
  t.Sub(u)   -> the Duration between two Times
  t.Before/After/Equal(u) -> compare instants
*/

func main() {
	// --- 1. The current instant and its parts ---
	now := time.Now()
	fmt.Printf("Current time: %s\n", now)
	fmt.Printf("Year: %d, Month: %s, Day: %d\n", now.Year(), now.Month(), now.Day())
	fmt.Printf("Hour: %d, Minute: %d, Second: %d\n", now.Hour(), now.Minute(), now.Second())
	fmt.Printf("Weekday: %s\n", now.Weekday())

	// --- 2. Build a specific instant with time.Date ---
	goLaunch := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	fmt.Printf("\nGo's launch (UTC): %s\n", goLaunch)

	// --- 3. Durations ---
	fiveMinutes := 5 * time.Minute
	oneHourThirty := time.Hour + 30*time.Minute
	fmt.Printf("\nFive minutes as string: %s\n", fiveMinutes)    // "5m0s"
	fmt.Printf("Five minutes in nanoseconds: %d\n", fiveMinutes) // an int64
	fmt.Printf("1h30m as string: %s\n", oneHourThirty)

	// --- 4. Arithmetic on instants ---
	fmt.Printf("\nIn 2 hours: %s\n", now.Add(2*time.Hour).Format(time.Kitchen))
	fmt.Printf("30 min ago: %s\n", now.Add(-30*time.Minute).Format(time.Kitchen))

	// --- 5. The duration between two instants ---
	since := now.Sub(goLaunch)
	fmt.Printf("\nSince Go launched: about %.0f days\n", since.Hours()/24)

	// --- 6. Comparing instants ---
	future := now.Add(time.Hour)
	fmt.Printf("\nfuture.After(now): %v\n", future.After(now))
	fmt.Printf("goLaunch.Before(now): %v\n", goLaunch.Before(now))

	// --- 7. Sleep pauses the current goroutine ---
	fmt.Println("\nSleeping 100ms...")
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Awake!")
}
