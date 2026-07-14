package main

import (
	"fmt"
	"time"
)

/*
FORMATTING AND PARSING TIME (the reference layout)
==================================================

Go does NOT use "%Y-%m-%d" style codes. Instead you write an EXAMPLE of the
layout using this one specific reference time:

  Mon Jan 2 15:04:05 MST 2006     (i.e. 01/02 03:04:05PM '06 -0700)

Read it as the numbers 1 2 3 4 5 6 7:
  month=1  day=2  hour=3(PM) or 15(24h)  minute=4  second=5  year=6  zone=7

So "2006-01-02" means YYYY-MM-DD, and "15:04:05" means 24-hour clock.

  t.Format(layout)         -> string
  time.Parse(layout, str)  -> Time  (parses in UTC unless the layout has a zone)

The time package ships common layouts as constants: time.RFC3339, time.Kitchen,
time.ANSIC, time.DateOnly, time.DateTime, ...
*/

func main() {
	now := time.Now()

	// --- 1. Formatting with custom layouts ---
	fmt.Printf("Default:        %s\n", now)
	fmt.Printf("YYYY-MM-DD:     %s\n", now.Format("2006-01-02"))
	fmt.Printf("US 12-hour:     %s\n", now.Format("01/02/2006 03:04:05 PM"))
	fmt.Printf("Readable:       %s\n", now.Format("Mon, Jan 2, 2006"))

	// --- 2. Formatting with the built-in layout constants ---
	fmt.Printf("RFC3339:        %s\n", now.Format(time.RFC3339))
	fmt.Printf("Kitchen:        %s\n", now.Format(time.Kitchen))
	fmt.Printf("DateTime:       %s\n", now.Format(time.DateTime))

	// --- 3. Parsing a date-only string (result is at midnight UTC) ---
	const layout = "2006-01-02"
	d, err := time.Parse(layout, "2025-07-15")
	if err != nil {
		fmt.Println("parse error:", err)
		return
	}
	fmt.Printf("\nParsed %q -> %d-%s-%d (%s)\n", "2025-07-15", d.Year(), d.Month(), d.Day(), d.Weekday())

	// --- 4. Parsing an RFC3339 timestamp WITH a zone offset ---
	ts, err := time.Parse(time.RFC3339, "2025-12-25T10:00:00+01:00")
	if err != nil {
		fmt.Println("parse error:", err)
		return
	}
	_, offset := ts.Zone()
	fmt.Printf("Parsed RFC3339: %s (offset %d seconds)\n", ts, offset)
	fmt.Printf("Same moment UTC: %s\n", ts.UTC())

	// --- 5. A parse failure reports exactly what didn't match ---
	if _, err := time.Parse("2006-01-02", "15/07/2025"); err != nil {
		fmt.Printf("\nexpected parse failure: %v\n", err)
	}
}
