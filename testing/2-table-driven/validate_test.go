package validate

import (
	"errors"
	"testing"
)

/*
TABLE-DRIVEN TESTS & SUBTESTS
=============================

When one function has many input/output cases, don't write one Test per case.
Instead put the cases in a slice ("the table") and loop over them. Wrapping the
loop body in t.Run(name, ...) creates a named SUBTEST, which:

  - shows each case separately in verbose output (TestX/valid, TestX/too_short)
  - keeps going after a failure (one bad case doesn't hide the others)
  - lets you run a single case: go test -run TestCheckUsername/contains_admin

This is the single most common Go testing pattern — learn it well.
*/

func TestCheckUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		want     bool
	}{
		{"valid", "greatusername", true},
		{"exactly six", "abcdef", true},
		{"too short", "test1", false},
		{"empty", "", false},
		{"contains admin", "adminuser", false},
		{"admin any case", "TheADMINhere", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := CheckUsername(tc.username); got != tc.want {
				t.Errorf("CheckUsername(%q) = %v, want %v", tc.username, got, tc.want)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantErr  error
	}{
		{"good", "validname", nil},
		{"too short", "abc", ErrInvalidUsername},
		{"admin blocked", "adminuser", ErrInvalidUsername},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := Login(tc.username)
			// errors.Is compares against the sentinel error value.
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Login(%q) error = %v, want %v", tc.username, err, tc.wantErr)
			}
		})
	}
}
