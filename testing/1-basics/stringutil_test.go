package basics

import "testing"

/*
THE ANATOMY OF A GO TEST
========================

Go's testing is built into the toolchain — no framework required.

Rules:
  - Test files end in _test.go and are excluded from normal builds.
  - Test functions are named TestXxx and take a single *testing.T.
  - A test FAILS when it calls t.Error/t.Errorf (continue) or t.Fatal/t.Fatalf
    (stop this test now). There is no "assert" keyword; you compare and report.

Run them:
  go test ./testing/1-basics        # run this package's tests
  go test -v ./testing/1-basics     # verbose: show each test
  go test -run Reverse ./...        # only tests whose name matches "Reverse"

The idiomatic failure message names the input, what you GOT, and what you WANT.
*/

func TestReverse(t *testing.T) {
	got := Reverse("hello")
	want := "olleh"
	if got != want {
		// Errorf marks the test failed but keeps running the rest of the func.
		t.Errorf("Reverse(%q) = %q, want %q", "hello", got, want)
	}
}

func TestReverseUnicode(t *testing.T) {
	// Reversing byte-by-byte would corrupt "é"; runes keep it intact.
	got := Reverse("héllo")
	want := "olléh"
	if got != want {
		// Fatalf stops THIS test immediately (use when continuing is pointless).
		t.Fatalf("Reverse(%q) = %q, want %q", "héllo", got, want)
	}
}

func TestReverseEmpty(t *testing.T) {
	if got := Reverse(""); got != "" {
		t.Errorf("Reverse(\"\") = %q, want empty string", got)
	}
}
