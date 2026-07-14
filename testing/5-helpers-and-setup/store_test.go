package store

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

/*
HELPERS, TEMP DIRS, CLEANUP, AND TestMain
=========================================

Four tools that keep bigger test suites clean:

  t.Helper()   marks a function as a helper, so a failure it reports points at
               the CALLER's line, not inside the helper.
  t.TempDir()  a fresh temp directory, automatically deleted after the test —
               perfect for file-based tests, no manual cleanup.
  t.Cleanup()  registers a function to run when the test finishes (LIFO),
               a cleaner alternative to defer that also works inside helpers.
  TestMain     one function that runs BEFORE/AFTER all tests in the package —
               use it for shared setup/teardown (fixtures, containers, etc.).
*/

// TestMain wraps the whole package's tests with setup and teardown. It MUST
// call m.Run() and exit with its result.
func TestMain(m *testing.M) {
	fmt.Println("== setup: starting store tests ==")
	code := m.Run() // runs every Test* in this package
	fmt.Println("== teardown: store tests complete ==")
	os.Exit(code)
}

// newTestStore is a HELPER: it builds a Store in a temp file. Because it calls
// t.Helper(), any t.Fatal it triggers is blamed on the calling test's line.
func newTestStore(t *testing.T) *Store {
	t.Helper()
	dir := t.TempDir() // auto-removed when the test ends
	return New(filepath.Join(dir, "store.json"))
}

// assertGet is a HELPER for the repeated "get and compare" assertion.
func assertGet(t *testing.T, s *Store, key, want string) {
	t.Helper()
	got, ok := s.Get(key)
	if !ok {
		t.Fatalf("key %q not found", key)
	}
	if got != want {
		t.Fatalf("Get(%q) = %q, want %q", key, got, want)
	}
}

func TestSetAndGet(t *testing.T) {
	s := newTestStore(t)

	if err := s.Set("lang", "go"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	assertGet(t, s, "lang", "go")
}

func TestOverwrite(t *testing.T) {
	s := newTestStore(t)

	_ = s.Set("k", "first")
	_ = s.Set("k", "second") // overwrite
	assertGet(t, s, "k", "second")
}

func TestMissingKey(t *testing.T) {
	s := newTestStore(t)

	if _, ok := s.Get("nope"); ok {
		t.Error("expected missing key to report ok=false")
	}
}

func TestCleanupRuns(t *testing.T) {
	// t.Cleanup runs after the test — useful for closing things opened in helpers.
	t.Cleanup(func() {
		t.Log("cleanup ran after TestCleanupRuns")
	})
	newTestStore(t)
}
