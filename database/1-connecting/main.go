package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	// Pure-Go SQLite driver — no C compiler / CGO required.
	// The blank import registers the driver under the name "sqlite".
	_ "modernc.org/sqlite"
)

/*
CONNECTING TO A DATABASE (database/sql)
=======================================

Go talks to SQL databases through the standard `database/sql` package plus a
DRIVER for your specific database. You import the driver for its side effects
(the blank `_` import) so it registers itself, then open by driver name:

  db, err := sql.Open("sqlite", dataSourceName)

Important facts:
  - sql.Open does NOT actually connect. It just validates arguments and
    prepares a connection POOL. Call db.Ping() to force a real connection and
    confirm the database is reachable.
  - *sql.DB is a POOL of connections, safe for concurrent use. Open it once
    and share it — do NOT open a new one per request.
  - Always `defer db.Close()`.

We use modernc.org/sqlite (pure Go) so this runs anywhere without a C toolchain.
The classic github.com/mattn/go-sqlite3 driver works too but needs CGO.
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
	// Keep the repo clean: put the database file in a temp dir and remove it.
	dir, err := os.MkdirTemp("", "go-db-connect")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	dbPath := filepath.Join(dir, "app.db")

	// 1. Open the pool (no connection happens yet).
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		fmt.Println("closing database connection")
		if err := db.Close(); err != nil {
			log.Printf("error closing db: %v", err)
		}
	}()

	// 2. Ping forces a real connection and verifies the DB is reachable.
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("database connection established")

	// 3. Run DDL to create our table. Exec is for statements with no rows back.
	if _, err := db.Exec(schema); err != nil {
		log.Fatal(err)
	}
	fmt.Println("users table created")

	// 4. Sanity-check the table exists by querying sqlite's catalog.
	var name string
	row := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users'`)
	if err := row.Scan(&name); err != nil {
		log.Fatal(err)
	}
	fmt.Println("verified table exists:", name)
}
