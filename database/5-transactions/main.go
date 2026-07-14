package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

/*
TRANSACTIONS (Begin / Commit / Rollback)
========================================

A transaction groups several statements so they ALL succeed or ALL fail — never
half. Classic case: create a user AND their profile row together. If the second
insert fails, the first must be undone.

  tx, err := db.BeginTx(ctx, nil)
  defer tx.Rollback()          // safe no-op if we already committed
  ... tx.ExecContext(...) ...
  return tx.Commit()           // commit only if everything worked

The `defer tx.Rollback()` idiom is the key trick: if any step returns early with
an error, the deferred Rollback undoes the partial work. If we reach Commit, the
later Rollback becomes a harmless no-op.

This example shows BOTH outcomes: one transaction that commits, and one that is
rolled back because the profile insert violates a constraint.
*/

var schema = `
CREATE TABLE IF NOT EXISTS users (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    name            TEXT NOT NULL,
    email           TEXT NOT NULL UNIQUE,
    hashed_password TEXT NOT NULL,
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS profiles (
    user_id    INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    avatar     TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`

func main() {
	dir, err := os.MkdirTemp("", "go-db-tx")
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

	// --- 1. Happy path: user + profile committed together ---
	id, err := createUserWithProfile(ctx, db, "Grace Hopper", "grace@example.com", "password", "grace.png")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("committed user+profile, id=%d\n", id)

	// --- 2. Failure path: avatar too long -> whole tx rolls back ---
	// (empty avatar violates NOT NULL via our own guard to force a rollback)
	_, err = createUserWithProfile(ctx, db, "Broken User", "broken@example.com", "password", "")
	if err != nil {
		fmt.Println("second insert failed, transaction rolled back:", err)
	}

	// --- 3. Prove the rollback: "broken@example.com" must NOT exist ---
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM users WHERE email = ?`, "broken@example.com").Scan(&count); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nusers with rolled-back email present in DB: %d (expected 0)\n", count)

	// Total users should be exactly 1 (only Grace).
	if err := db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("total users: %d\n", count)
}

func createUserWithProfile(ctx context.Context, db *sql.DB, name, email, password, avatar string) (int64, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	// Undo everything if we return before Commit. No-op after a successful Commit.
	defer tx.Rollback()

	hp, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	res, err := tx.ExecContext(ctx,
		`INSERT INTO users (name, email, hashed_password) VALUES (?, ?, ?)`,
		name, email, string(hp))
	if err != nil {
		return 0, err
	}
	userID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Guard: a profile must have an avatar. Returning here triggers the
	// deferred Rollback, so the user row above is undone too.
	if avatar == "" {
		return 0, errors.New("avatar is required")
	}

	if _, err := tx.ExecContext(ctx,
		`INSERT INTO profiles (user_id, avatar) VALUES (?, ?)`,
		userID, avatar); err != nil {
		return 0, err
	}

	// Everything worked — make it permanent.
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return userID, nil
}
