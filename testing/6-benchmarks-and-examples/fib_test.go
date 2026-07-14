package fib

import (
	"fmt"
	"testing"
)

/*
BENCHMARKS AND EXAMPLES
=======================

Beyond TestXxx, the testing package supports two more function kinds:

  BenchmarkXxx(b *testing.B)  measures performance. Run the body b.N times;
                              the framework tunes b.N until timing is stable.
      go test -bench=. ./testing/6-benchmarks-and-examples
      go test -bench=. -benchmem ./...   # also report allocations

  ExampleXxx()                doubles as documentation AND a test. The text in
                              the "// Output:" comment is compared against what
                              the example prints; a mismatch fails `go test`.
      Examples also appear in `go doc` / pkg.go.dev, so they never go stale.
*/

// A normal correctness test first — benchmarks assume the code is correct.
func TestImplementationsAgree(t *testing.T) {
	for n := 0; n < 15; n++ {
		if it, rec := Iterative(n), Recursive(n); it != rec {
			t.Errorf("n=%d: Iterative=%d, Recursive=%d", n, it, rec)
		}
	}
}

func BenchmarkIterative(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Iterative(20)
	}
}

func BenchmarkRecursive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Recursive(20)
	}
}

// ExampleIterative is verified against its Output comment when tests run.
func ExampleIterative() {
	fmt.Println(Iterative(10))
	// Output: 55
}

// Example functions can also demonstrate a sequence of calls.
func ExampleIterative_sequence() {
	for n := 0; n < 8; n++ {
		fmt.Print(Iterative(n), " ")
	}
	// Output: 0 1 1 2 3 5 8 13
}
