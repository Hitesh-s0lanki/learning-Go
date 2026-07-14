# Testing in Go (`go test`)

A hands-on tour of Go's built-in testing toolkit — from your first `TestXxx`
function to table-driven tests, HTTP handler tests, benchmarks, examples, and
the popular testify library. Testing is part of the Go toolchain, so there is
nothing to install for folders 1–6.

Unlike the other sections, these folders aren't run with `go run` — they're run
with **`go test`**. Each folder is a small package with its code and its
`_test.go` file side by side.

## How to run

```bash
go test ./testing/...              # run every test in the section
go test -v ./testing/2-table-driven   # verbose: show each (sub)test
go test -run TestLogin ./...       # only tests whose name matches
go test -cover ./testing/...       # report statement coverage
go test -race ./testing/...        # detect data races (see the concurrency section)
go test -bench=. ./testing/6-benchmarks-and-examples   # run benchmarks
```

## Topics (learn them in order)

| # | Folder | Concept |
|---|--------|---------|
| 1 | [1-basics](1-basics/stringutil_test.go) | test anatomy: `TestXxx(t *testing.T)`, `t.Errorf` vs `t.Fatalf`, want/got |
| 2 | [2-table-driven](2-table-driven/validate_test.go) | table-driven tests + named subtests with `t.Run` |
| 3 | [3-io-and-interfaces](3-io-and-interfaces/greeting_test.go) | design for testability: inject an `io.Writer` instead of printing to stdout |
| 4 | [4-http-handlers](4-http-handlers/handlers_test.go) | `net/http/httptest`: `NewRequest` + `NewRecorder`, no server needed |
| 5 | [5-helpers-and-setup](5-helpers-and-setup/store_test.go) | `t.Helper`, `t.TempDir`, `t.Cleanup`, and `TestMain` for setup/teardown |
| 6 | [6-benchmarks-and-examples](6-benchmarks-and-examples/fib_test.go) | `BenchmarkXxx(b *testing.B)` and testable `ExampleXxx` functions |
| 7 | [7-testify](7-testify/calc_test.go) | concise assertions with the third-party `stretchr/testify` library |

## The rules of a Go test

- Test files must end in **`_test.go`** (they're excluded from normal builds).
- Test functions are named **`TestXxx`** and take a single **`*testing.T`**.
- There is **no assert keyword** in the standard library — you compare values
  yourself and call `t.Error`/`t.Fatal` to fail.
  - `t.Errorf` — mark failed, **keep going** (report multiple problems).
  - `t.Fatalf` — mark failed, **stop this test now** (when continuing is pointless).
- Write **table-driven tests with subtests** (`t.Run`) whenever one function has
  many cases — it's the dominant Go pattern.

## Function kinds the toolchain recognizes

| Signature | Purpose | Run with |
|---|---|---|
| `func TestXxx(t *testing.T)` | correctness | `go test` |
| `func BenchmarkXxx(b *testing.B)` | performance | `go test -bench=.` |
| `func ExampleXxx()` + `// Output:` | docs that are also tested | `go test` |
| `func TestMain(m *testing.M)` | package-wide setup/teardown | `go test` |

## Key habits

- Name failures with **input, got, and want** so a red test explains itself.
- Prefer **dependency injection** (`io.Writer`, interfaces) over global state —
  testable code is usually better-designed code (see folder 3).
- Use **`t.TempDir()`** and **`t.Cleanup()`** so tests never leave artifacts or
  leak resources.
- Mark assertion helpers with **`t.Helper()`** so failures point at the caller.
- Keep **`Example` functions** accurate — they're compiled and their output is
  checked, so they can't silently rot like comments do.

## Standard library vs testify

Folders 1–6 use only the standard library — the idiomatic baseline every Go
developer should know. Folder 7 shows [testify](https://github.com/stretchr/testify),
a widely-used library that trades a dependency for terser assertions
(`assert.Equal(t, want, got)`) and nicer diffs. Both are valid; learn the
standard library first so you understand what testify is doing for you.
