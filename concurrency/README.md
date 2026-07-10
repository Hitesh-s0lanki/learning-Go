# Concurrency in Go

A hands-on tour of Go's concurrency model, from a single goroutine up to
real-world patterns. Each folder is a standalone `package main` you can run on
its own.

## How to run

```bash
cd concurrency/1-goroutines
go run .

# to detect data races (essential for concurrent code):
go run -race .
```

## Topics (learn them in order)

| # | Folder | Concept |
|---|--------|---------|
| 1 | [1-goroutines](1-goroutines/main.go) | Starting goroutines with `go`; why `main` exiting kills them |
| 2 | [2-waitgroup](2-waitgroup/main.go) | `sync.WaitGroup` to wait for goroutines to finish |
| 3 | [3-channels](3-channels/main.go) | Unbuffered channels: send/receive as a synchronized handoff |
| 4 | [4-buffered-channels](4-buffered-channels/main.go) | Buffered channels, `len`/`cap`, semaphore for bounded concurrency |
| 5 | [5-channel-closing](5-channel-closing/main.go) | `close`, `range` over a channel, comma-ok, broadcast via close |
| 6 | [6-select](6-select/main.go) | `select` for multiplexing channels, timeouts, non-blocking ops |
| 7 | [7-mutex](7-mutex/main.go) | `sync.Mutex` / `sync.RWMutex` for shared state |
| 8 | [8-race-condition](8-race-condition/main.go) | Data races, the `-race` detector, `sync/atomic` |
| 9 | [9-worker-pool](9-worker-pool/main.go) | Worker pool pattern (bounded parallelism) |
| 10 | [10-context](10-context/main.go) | `context` for cancellation, timeouts, and deadlines |
| 11 | [11-sync-once](11-sync-once/main.go) | `sync.Once` for exactly-once initialization |
| 12 | [12-patterns](12-patterns/main.go) | Pipeline and fan-out / fan-in patterns |

## Key mental models

- **"Don't communicate by sharing memory; share memory by communicating."**
  Prefer passing data over channels to locking shared variables — but use a
  `Mutex` when you genuinely have shared state (a counter, cache, map).
- A **goroutine** is not an OS thread — it's cheap; thousands are fine.
- An **unbuffered** channel synchronizes sender and receiver (a handoff);
  a **buffered** channel decouples them up to its capacity.
- Only the **sender** closes a channel. Sending on a closed channel panics.
- Always develop concurrent code with **`go run -race .`** enabled.
- Use **`context`** to tell goroutines when to stop.

## Why Go's concurrency is different (and more useful)

This is *why* Go became the go-to language for backends, cloud infra, and
networking tools (Docker, Kubernetes, and most modern infra are written in Go).

### 1. Goroutines vs. OS threads

In most languages (Java, C++, C#, Python), concurrency means **OS threads**.
Each thread costs ~1–8 MB of stack and switching between them needs the kernel,
so you can run *thousands*, but not *millions* — hence thread pools to ration
them.

Go's **goroutines** start at ~2 KB and grow/shrink on demand. The Go runtime
multiplexes millions of them onto a small number of OS threads (the **M:N
scheduler**), and switching happens in user space, not the kernel.

```
Java/C++:   1 task = 1 OS thread   → thousands max, heavy
Go:         1 task = 1 goroutine   → millions, cheap
```

This is why a Go web server can handle 100k concurrent connections with one
goroutine each, while other languages need async/event-loop tricks to match it.

### 2. Concurrency is built into the language, not a library

In Java/Python/C++, concurrency lives in **libraries bolted on later**
(`Thread`, `ExecutorService`, `asyncio`, `std::thread`). In Go it's **keywords
and built-in types**:

- `go f()` — start concurrency (one word)
- `chan` — a first-class channel type
- `select` — a language statement for multiplexing
- `<-` — the channel operator

```go
go doWork()                       // Go
```
```java
executor.submit(() -> doWork());  // Java
```

### 3. Philosophy: "share memory by communicating"

Most languages default to **shared memory + locks** — error-prone (deadlocks,
races, forgotten unlocks). Go's slogan is **"Don't communicate by sharing
memory; share memory by communicating."** Instead of many threads fighting over
one variable, you **pass the data through a channel** so only one goroutine owns
it at a time. Go still has mutexes, but channels are the idiomatic first choice
and make data ownership explicit and safe by design.

### 4. vs. async/await (JS, Python, Rust, C#)

`async`/`await` solved "threads are expensive" but introduced **function
coloring**: once a function is `async`, every caller must be `async` too, and
blocking code silently poisons the event loop.

```python
async def fetch():        # "colored" — callers must await, must be async too
    await something()
```

Go has **no `async` keyword and no coloring**. You write plain,
*synchronous-looking* code, and when a goroutine blocks on I/O the scheduler
parks it and runs another. It reads like blocking code but scales like async.

```go
data := fetchURL(url)     // looks blocking, but the goroutine yields under the hood
```

### 5. First-class safety tooling

Go ships a **race detector** built into the toolchain (`go run -race .`) —
few languages give you a production-grade data-race detector for free. Combined
with `context` for clean cancellation, Go gives you a full concurrency *system*,
not just primitives.

### Summary

| | Other languages | Go |
|---|---|---|
| Unit of concurrency | OS thread (heavy) | Goroutine (~2KB, cheap) |
| Scale | Thousands | Millions |
| Scheduling | Kernel | Go runtime (user space, M:N) |
| Syntax | Library APIs | Built-in keywords (`go`, `chan`, `select`) |
| Default model | Shared memory + locks | Communicate over channels |
| Async style | `async`/`await` (function coloring) | Plain synchronous-looking code |
| Race detection | External/rare | Built into the toolchain (`-race`) |

**Bottom line:** Go makes concurrency **cheap, simple, and safe enough to use
everywhere** — so instead of avoiding it, you design *with* it. That's exactly
why it dominates cloud/backend/networking.
