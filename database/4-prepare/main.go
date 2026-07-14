package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

/*
PREPARED STATEMENTS (db.Prepare + Context)
==========================================

A prepared statement is a query the database parses and plans ONCE, then you
execute MANY times with different values. Benefits:

  - speed: the DB reuses the query plan instead of re-parsing every time
  - safety: values are always sent as parameters (injection-safe by design)

  stmt, err := db.Prepare(`INSERT ... VALUES (?, ?)`)
  defer stmt.Close()          // statements hold resources — always close
  stmt.Exec(a, b)             // run it as many times as you like

Context variants let you attach timeouts / cancellation to a call:

  stmt.ExecContext(ctx, ...)  // cancel if ctx is done
  db.QueryContext(ctx, ...)

Prefer the *Context methods in real services so a slow query can be cancelled.
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
	dir, err := os.MkdirTemp("", "go-db-prepare")
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

	// A context with a timeout guards the whole batch of inserts.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// --- Prepare ONCE, execute MANY times ---
	stmt, err := db.PrepareContext(ctx, `INSERT INTO users (name, email, hashed_password) VALUES (?, ?, ?)`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	people := []struct{ name, email string }{
		{"User One", "one@example.com"},
		{"User Two", "two@example.com"},
		{"User Three", "three@example.com"},
		{"User Four", "four@example.com"},
	}

	for _, p := range people {
		hp, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal(err)
		}
		res, err := stmt.ExecContext(ctx, p.name, p.email, string(hp))
		if err != nil {
			log.Fatal(err)
		}
		id, _ := res.LastInsertId()
		fmt.Printf("inserted %-10s -> id %d (reused prepared statement)\n", p.name, id)
	}

	var count int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users`).Scan(&count); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\ntotal users: %d\n", count)
}
