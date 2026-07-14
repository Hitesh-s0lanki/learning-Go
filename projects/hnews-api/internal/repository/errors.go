// Package repository contains all SQL data access, hidden behind interfaces so
// handlers depend on behavior, not on database/sql. It translates driver
// errors (sql.ErrNoRows, unique-violation) into domain errors defined here.
package repository

import "errors"

var (
	// ErrNotFound means a requested row does not exist.
	ErrNotFound = errors.New("record not found")
	// ErrDuplicate means a unique constraint was violated (e.g. email/title).
	ErrDuplicate = errors.New("record already exists")
	// ErrDuplicateVote means the user already voted on this post.
	ErrDuplicateVote = errors.New("already voted")
)
