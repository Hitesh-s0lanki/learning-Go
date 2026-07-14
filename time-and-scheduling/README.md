# Time & Scheduling in Go (`time`, `math/rand/v2`)

A hands-on tour of working with time: instants and durations, formatting and
parsing, timers and tickers, random numbers, a small scheduler, and time zones.
Each folder is a standalone `package main` you can run on its own.

## How to run

```bash
cd time-and-scheduling/1-time
go run .
```

Every example finishes in about a second and self-terminates — the timer,
ticker, and scheduler demos use short durations (milliseconds) instead of
blocking for minutes or hours.

## Topics (learn them in order)

| # | Folder | Concept |
|---|--------|---------|
| 1 | [1-time](1-time/main.go) | `time.Time` & `time.Duration`, `time.Date`, `Add`/`Sub`, `Before`/`After`, `Sleep` |
| 2 | [2-formatting](2-formatting/main.go) | the reference layout `2006-01-02 15:04:05`, `Format`, `Parse`, RFC3339 |
| 3 | [3-timer-and-ticker](3-timer-and-ticker/main.go) | `Timer`, `Ticker`, `time.After` (timeouts), `AfterFunc` |
| 4 | [4-random](4-random/main.go) | `math/rand/v2`: `IntN`, `Shuffle`, `Perm`, seeded PCG for reproducibility |
| 5 | [5-schedule](5-schedule/main.go) | a mini scheduler composing timers, tickers, goroutines, and a `WaitGroup` |
| 6 | [6-timezone](6-timezone/main.go) | `LoadLocation`, `In`, converting one instant across zones, `time/tzdata` |

## The reference layout (the #1 gotcha)

Go doesn't use `%Y-%m-%d`. You spell the layout using this one magic reference
time — the numbers `1 2 3 4 5 6 -7`:

```
Mon Jan 2 15:04:05 MST 2006
     │   │  │  │  │      └── 2006 = year
     │   │  │  │  └───────── 05   = second
     │   │  │  └──────────── 04   = minute
     │   │  └─────────────── 15   = hour (24h)  |  03 = hour (12h) with PM
     │   └────────────────── 2    = day
     └────────────────────── Jan  = month (01)
```

So `now.Format("2006-01-02")` gives `2025-07-15`, and `time.Parse("15:04", s)`
reads a 24-hour clock. Handy constants: `time.RFC3339`, `time.DateOnly`,
`time.DateTime`, `time.Kitchen`.

## Key habits

- Build durations from the unit constants: `5 * time.Minute`, not raw numbers.
- Move between instants with `Add` (use a negative duration to go back) and
  measure gaps with `Sub`.
- **Always `Stop()` a `Ticker`** (and any `Timer` you might not drain) —
  `defer ticker.Stop()` — or it leaks.
- Use `select { case <-work: ... case <-time.After(d): ... }` for timeouts.
- Store and transmit time in **UTC** (or RFC3339 with an offset); convert to a
  local zone only when displaying to a human.
- Seed your own `rand.New(rand.NewPCG(...))` when you need **reproducible**
  randomness (tests, simulations).

## Key packages

- **`time`** — instants, durations, formatting/parsing, timers, tickers, zones.
- **`math/rand/v2`** — fast, auto-seeded pseudo-random numbers (Go 1.22+).
  **Not** cryptographically secure — use **`crypto/rand`** for tokens/keys.
- **`time/tzdata`** — a blank import that embeds the IANA zone database so
  `LoadLocation` works even on systems without OS timezone files.
