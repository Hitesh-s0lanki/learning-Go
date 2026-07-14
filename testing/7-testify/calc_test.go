package calc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
TESTING WITH TESTIFY (a popular third-party library)
====================================================

The standard library has no assert helpers on purpose. Many teams add
github.com/stretchr/testify for concise assertions and readable failure output.

Two flavors:
  assert.*   report the failure and CONTINUE the test.
  require.*  report the failure and STOP the test (like t.Fatal) — use when
             later lines would panic if this one failed (e.g. a nil result).

Standard library vs testify — same intent, less boilerplate:
  if got != want { t.Errorf(...) }   ->   assert.Equal(t, want, got)
  if err != nil { t.Fatalf(...) }    ->   require.NoError(t, err)

testify is optional. Everything in folders 1-6 uses only the standard library.
*/

func TestAdd(t *testing.T) {
	assert.Equal(t, 5, Add(2, 3))
	assert.Equal(t, 0, Add(-2, 2))
}

func TestDivOK(t *testing.T) {
	got, err := Div(10, 2)
	require.NoError(t, err) // stop here if it unexpectedly errored
	assert.Equal(t, 5, got)
}

func TestDivByZero(t *testing.T) {
	_, err := Div(1, 0)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrDivideByZero)
}

// testify also works cleanly with table-driven tests.
func TestDivTable(t *testing.T) {
	tests := []struct {
		name    string
		a, b    int
		want    int
		wantErr bool
	}{
		{"simple", 6, 3, 2, false},
		{"truncates", 7, 2, 3, false},
		{"by zero", 1, 0, 0, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Div(tc.a, tc.b)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}
