package main

import (
	"fmt"
	"time"

	// Embed the IANA timezone database into the binary so LoadLocation works
	// even on systems that don't ship tzdata (minimal containers, some Windows).
	_ "time/tzdata"
)

/*
TIME ZONES (time.Location)
==========================

A time.Time carries a *time.Location. The SAME instant displays differently in
different zones — the underlying moment never changes, only its presentation.

  time.LoadLocation("America/New_York") -> a *Location (IANA name)
  t.In(loc)   -> the same instant, rendered in loc
  t.UTC()     -> the same instant in UTC
  t.Local()   -> the same instant in the machine's local zone

Build a zoned instant directly by passing a Location to time.Date. Always store
and transmit times in UTC (or RFC3339 with an offset); convert to a local zone
only for display.

Importing time/tzdata (blank import above) makes these names resolve anywhere.
*/

func main() {
	newYork, _ := time.LoadLocation("America/New_York")
	tokyo, _ := time.LoadLocation("Asia/Tokyo")
	london, _ := time.LoadLocation("Europe/London")

	// --- 1. One instant, shown in several zones ---
	now := time.Now()
	fmt.Println("the same moment, everywhere:")
	fmt.Printf("  UTC:      %s\n", now.UTC().Format(time.RFC3339))
	fmt.Printf("  New York: %s\n", now.In(newYork).Format(time.RFC3339))
	fmt.Printf("  London:   %s\n", now.In(london).Format(time.RFC3339))
	fmt.Printf("  Tokyo:    %s\n", now.In(tokyo).Format(time.RFC3339))

	// --- 2. A meeting defined in one zone, read in others ---
	meeting := time.Date(2025, time.December, 25, 10, 0, 0, 0, newYork)
	fmt.Printf("\nmeeting (New York): %s\n", meeting.Format("Mon Jan 2 3:04 PM"))
	fmt.Printf("  in London: %s\n", meeting.In(london).Format("Mon Jan 2 3:04 PM"))
	fmt.Printf("  in Tokyo:  %s\n", meeting.In(tokyo).Format("Mon Jan 2 3:04 PM"))
	fmt.Printf("  in UTC:    %s\n", meeting.UTC().Format("Mon Jan 2 3:04 PM"))

	// --- 3. Parse an offset timestamp, then convert ---
	ts, err := time.Parse(time.RFC3339, "2025-06-15T14:30:00-07:00")
	if err != nil {
		fmt.Println("parse error:", err)
		return
	}
	fmt.Printf("\nparsed %s\n", ts.Format(time.RFC3339))
	fmt.Printf("  in New York: %s\n", ts.In(newYork).Format(time.RFC3339))
	fmt.Printf("  in Tokyo:    %s\n", ts.In(tokyo).Format(time.RFC3339))
}
