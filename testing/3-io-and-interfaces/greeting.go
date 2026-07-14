// Package greeting shows how designing for testability makes tests trivial.
// The lesson: accept an io.Writer instead of printing straight to os.Stdout.
package greeting

import (
	"fmt"
	"io"
)

// Greet writes a greeting to out.
//
// Because it takes an io.Writer, a test can pass a *bytes.Buffer and read back
// exactly what was written — no global state, no stdout capture. In production
// you pass os.Stdout; in tests you pass a buffer.
func Greet(out io.Writer, prefix, name string) {
	fmt.Fprintf(out, "Hello, %s %s!\n", prefix, name)
}

// GreetStdout is the HARD-to-test version: it prints directly to stdout. It's
// here only to contrast with Greet — see the test for the awkward pipe dance
// you need to capture its output.
func GreetStdout(prefix, name string) {
	fmt.Printf("Hello, %s %s!\n", prefix, name)
}
