// Package fib provides Fibonacci implementations used to demonstrate
// benchmarks and testable examples.
package fib

// Iterative returns the nth Fibonacci number in O(n) time.
func Iterative(n int) int {
	if n < 2 {
		return n
	}
	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

// Recursive returns the nth Fibonacci number the naive (slow) way, O(2^n).
// It exists so a benchmark can show just how much slower it is.
func Recursive(n int) int {
	if n < 2 {
		return n
	}
	return Recursive(n-1) + Recursive(n-2)
}
