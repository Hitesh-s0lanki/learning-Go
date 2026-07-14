// Package calc is a tiny calculator used to demonstrate testing with the
// popular third-party assertion library, testify.
package calc

import "errors"

// ErrDivideByZero is returned by Div when the divisor is zero.
var ErrDivideByZero = errors.New("divide by zero")

// Add returns a + b.
func Add(a, b int) int { return a + b }

// Div returns a / b, or ErrDivideByZero when b is zero.
func Div(a, b int) (int, error) {
	if b == 0 {
		return 0, ErrDivideByZero
	}
	return a / b, nil
}
