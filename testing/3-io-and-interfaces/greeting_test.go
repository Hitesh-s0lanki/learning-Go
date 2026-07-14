package greeting

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

/*
TESTING OUTPUT: INJECT AN io.Writer
====================================

The easy, idiomatic way to test something that "produces output" is to have it
write to an io.Writer you provide. In the test you hand it a *bytes.Buffer and
assert on its contents. No globals, no cleanup.

The alternative — capturing os.Stdout — is fiddly and shown at the bottom as a
cautionary contrast. Prefer dependency injection.
*/

// The good way: pass a buffer, read it back.
func TestGreet(t *testing.T) {
	var buf bytes.Buffer
	Greet(&buf, "Mr.", "Joseph")

	want := "Hello, Mr. Joseph!\n"
	if got := buf.String(); got != want {
		t.Errorf("Greet wrote %q, want %q", got, want)
	}
}

// The awkward way: temporarily replace os.Stdout with a pipe to capture output.
// This is why functions should accept an io.Writer instead.
func TestGreetStdout(t *testing.T) {
	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	os.Stdout = w
	t.Cleanup(func() { os.Stdout = orig }) // always restore, even on failure

	GreetStdout("Mr.", "Joseph")
	w.Close()

	var buf strings.Builder
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("io.Copy: %v", err)
	}

	want := "Hello, Mr. Joseph!\n"
	if got := buf.String(); got != want {
		t.Errorf("GreetStdout printed %q, want %q", got, want)
	}
}
