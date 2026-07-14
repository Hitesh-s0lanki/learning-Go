package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"go-learning/database/6-repository/repository"

	_ "modernc.org/sqlite"
)

/*
THE REPOSITORY PATTERN
======================

As an app grows, sprinkling raw SQL through your handlers becomes a mess: it's
hard to test, hard to change databases, and business logic gets tangled with
query strings.

The repository pattern fixes this by putting ALL data access behind an
interface:

  models/      -> plain structs (the data)
  repository/  -> the UserRepository interface + a SQL implementation (the how)
  main.go      -> depends only on the interface (the what)

Benefits:
  - swap the backend (SQLite -> Postgres) without touching callers
  - unit-test business logic against a fake repository, no DB required
  - domain errors (repository.ErrNotFound) instead of leaking sql.ErrNoRows

This file just wires everything together and exercises the repository.
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
	dir, err := os.MkdirTemp("", "go-db-repo")
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

	ctx := context.Background()

	// main() knows only the interface, never the SQL behind it.
	users := repository.NewUserRepository(db)

	// --- 1. Create through the repository ---
	created, err := users.Create(ctx, "Ada Lovelace", "ada@example.com", "s3cret")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("created user id=%d at %s\n", created.ID, created.CreatedAt.Format("2006-01-02 15:04:05"))

	if _, err := users.Create(ctx, "Alan Turing", "alan@example.com", "s3cret"); err != nil {
		log.Fatal(err)
	}

	// --- 2. Look up by email ---
	found, err := users.GetByEmail(ctx, "ada@example.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("fetched by email: %s <%s>\n", found.Name, found.Email)

	// --- 3. Domain error instead of sql.ErrNoRows ---
	if _, err := users.GetByEmail(ctx, "ghost@example.com"); errors.Is(err, repository.ErrNotFound) {
		fmt.Println("missing user surfaced as repository.ErrNotFound")
	}

	// --- 4. List everyone ---
	all, err := users.List(ctx)
	if err != nil {
		log.Fatal(err)
	}
	bs, _ := json.MarshalIndent(all, "", "  ")
	fmt.Printf("\nall users (%d):\n%s\n", len(all), bs)
}
