package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

/*
INSERTING DATA (db.Exec + placeholders)
=======================================

Use db.Exec for statements that DON'T return rows (INSERT, UPDATE, DELETE, DDL):

  result, err := db.Exec("INSERT INTO users (...) VALUES (?, ?)", a, b)

ALWAYS use placeholders (`?` for SQLite/MySQL, `$1` for Postgres) for values.
Never build SQL by string concatenation — that's how SQL injection happens.
The driver sends your values separately from the query text.

The returned sql.Result gives you:
  result.LastInsertId() -> the auto-generated primary key (driver dependent)
  result.RowsAffected() -> how many rows changed

Never store plaintext passwords. We hash with bcrypt before inserting.
*/

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
	dir, err := os.MkdirTemp("", "go-db-insert")
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

	// --- 1. Insert a few users and print their new IDs ---
	seed := []struct{ name, email, password string }{
		{"Ada Lovelace", "ada@example.com", "s3cret1"},
		{"Alan Turing", "alan@example.com", "s3cret2"},
		{"Grace Hopper", "grace@example.com", "s3cret3"},
	}
	for _, s := range seed {
		id, err := createUser(db, s.name, s.email, s.password)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("inserted %-14s -> id %d\n", s.name, id)
	}

	// --- 2. UNIQUE constraint kicks in on a duplicate email ---
	if _, err := createUser(db, "Impostor", "ada@example.com", "x"); err != nil {
		fmt.Println("\nduplicate email correctly rejected:", err)
	}

	// --- 3. UPDATE shows RowsAffected ---
	res, err := db.Exec(`UPDATE users SET name = ? WHERE email = ?`, "Augusta Ada King", "ada@example.com")
	if err != nil {
		log.Fatal(err)
	}
	n, _ := res.RowsAffected()
	fmt.Printf("\nupdate affected %d row(s)\n", n)
}

func createUser(db *sql.DB, name, email, password string) (int64, error) {
	// Hash the password; store the hash, never the plaintext.
	hp, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	stmt := `INSERT INTO users (name, email, hashed_password) VALUES (?, ?, ?)`
	res, err := db.Exec(stmt, name, email, string(hp))
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
