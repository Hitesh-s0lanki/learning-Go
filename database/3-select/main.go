package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

/*
QUERYING DATA (QueryRow vs Query + Scan)
========================================

Two read paths:

  db.QueryRow(sql, args...)  -> exactly ONE row. Scan straight into vars.
                                If nothing matches, Scan returns sql.ErrNoRows.
  db.Query(sql, args...)     -> MANY rows. Loop with rows.Next(), Scan each,
                                then ALWAYS check rows.Err() and rows.Close().

Scan copies column values into your Go variables IN ORDER — the order of the
columns in your SELECT must match the order of the &pointers you pass.

Golden rules for multi-row queries:
  - `defer rows.Close()` right after a successful Query.
  - check `rows.Err()` after the loop (a row may fail mid-iteration).
  - handle `sql.ErrNoRows` explicitly for single-row lookups.
*/

type User struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"` // never expose in JSON
	CreatedAt      time.Time `json:"created_at"`
}

var schema = `
CREATE TABLE IF NOT EXISTS users (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    name            TEXT NOT NULL,
    email           TEXT NOT NULL UNIQUE,
    hashed_password TEXT NOT NULL,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`

func main() {
	dir, err := os.MkdirTemp("", "go-db-select")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	db, err := sql.Open("sqlite", filepath.Join(dir, "app.db"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if _, err := db.Exec(schema); err != nil {
		log.Fatal(err)
	}
	seed(db)

	// --- 1. Single-row lookup with QueryRow ---
	u, err := getUserByEmail(db, "grace@example.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("found user: %s <%s>\n", u.Name, u.Email)

	// --- 2. Missing row returns sql.ErrNoRows ---
	if _, err := getUserByEmail(db, "nobody@example.com"); errors.Is(err, sql.ErrNoRows) {
		fmt.Println("lookup for missing email returned sql.ErrNoRows (as expected)")
	}

	// --- 3. Multi-row query, marshalled to JSON ---
	users, err := getUsers(db)
	if err != nil {
		log.Fatal(err)
	}
	bs, _ := json.MarshalIndent(users, "", "  ")
	fmt.Printf("\nall users (%d):\n%s\n", len(users), bs)
}

func getUserByEmail(db *sql.DB, email string) (*User, error) {
	const q = `SELECT id, name, email, hashed_password, created_at FROM users WHERE email = ?`
	var u User
	err := db.QueryRow(q, email).
		Scan(&u.ID, &u.Name, &u.Email, &u.HashedPassword, &u.CreatedAt)
	if err != nil {
		return nil, err // may be sql.ErrNoRows — let the caller decide
	}
	return &u, nil
}

func getUsers(db *sql.DB) ([]User, error) {
	const q = `SELECT id, name, email, hashed_password, created_at FROM users ORDER BY id`
	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.HashedPassword, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	// A row can fail mid-iteration — always check.
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func seed(db *sql.DB) {
	for _, s := range []struct{ name, email string }{
		{"Ada Lovelace", "ada@example.com"},
		{"Alan Turing", "alan@example.com"},
		{"Grace Hopper", "grace@example.com"},
	} {
		hp, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		if _, err := db.Exec(
			`INSERT INTO users (name, email, hashed_password) VALUES (?, ?, ?)`,
			s.name, s.email, string(hp),
		); err != nil {
			log.Fatal(err)
		}
	}
}
