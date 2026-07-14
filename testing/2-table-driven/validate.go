// Package validate holds simple username rules — a good fit for table-driven
// tests because the logic is "many inputs -> expected outputs".
package validate

import (
	"errors"
	"strings"
)

// ErrInvalidUsername is returned by Login for a username that fails the rules.
var ErrInvalidUsername = errors.New("invalid username")

// CheckUsername reports whether a username is acceptable: at least 6 characters
// and not containing the word "admin".
func CheckUsername(username string) bool {
	if len(username) < 6 {
		return false
	}
	if strings.Contains(strings.ToLower(username), "admin") {
		return false
	}
	return true
}

// Login validates the username and returns an error if it is unacceptable.
func Login(username string) error {
	if !CheckUsername(username) {
		return ErrInvalidUsername
	}
	return nil
}
