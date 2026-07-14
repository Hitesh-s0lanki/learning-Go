// Package database opens the Postgres connection pool and applies migrations.
package database

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"sort"
	"time"

	// pgx's database/sql driver. Registered under the name "pgx".
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Connect opens a Postgres pool from a DSN, verifies it with a ping, and
// tunes the pool for a small service. The caller must Close the returned *sql.DB.
func Connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	// Reasonable pool defaults for a demo service.
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}
	return db, nil
}

// Migrate applies every *.sql migration in the given filesystem in filename
// order. Each file is expected to be idempotent (CREATE TABLE IF NOT EXISTS).
func Migrate(db *sql.DB, fsys fs.FS) error {
	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return fmt.Errorf("read migrations: %w", err)
	}

	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names) // 0001_..., 0002_..., applied in order

	for _, name := range names {
		sqlBytes, err := fs.ReadFile(fsys, name)
		if err != nil {
			return fmt.Errorf("read %s: %w", name, err)
		}
		if _, err := db.Exec(string(sqlBytes)); err != nil {
			return fmt.Errorf("apply %s: %w", name, err)
		}
	}
	return nil
}
