# Databases in Go (`database/sql`)

A hands-on tour of talking to a SQL database from Go — connecting, inserting,
querying, prepared statements, transactions, and the repository pattern. Each
numbered folder is a standalone program you can run on its own.

## Driver choice

These examples use **[`modernc.org/sqlite`](https://pkg.go.dev/modernc.org/sqlite)**,
a **pure-Go** SQLite driver — so everything runs with no C compiler and no CGO.
It registers itself under the driver name `"sqlite"`:

```go
import _ "modernc.org/sqlite"
db, _ := sql.Open("sqlite", "app.db")
```

The popular `github.com/mattn/go-sqlite3` driver works the same way but requires
CGO (a C toolchain). The `database/sql` code you write is identical either way —
only the import and driver name change.

## How to run

```bash
cd database/1-connecting
go run .
```

Every example creates its database in a **temp directory and deletes it on
exit**, so running them leaves your repo clean.

## Topics (learn them in order)

| # | Folder | Concept |
|---|--------|---------|
| 1 | [1-connecting](1-connecting/main.go) | `sql.Open`, `db.Ping`, `db.Exec` for DDL, the `*sql.DB` pool |
| 2 | [2-insert](2-insert/main.go) | `db.Exec` with `?` placeholders, `LastInsertId`, `RowsAffected`, bcrypt hashing |
| 3 | [3-select](3-select/main.go) | `QueryRow` vs `Query`, `rows.Scan`, `sql.ErrNoRows`, `rows.Err` |
| 4 | [4-prepare](4-prepare/main.go) | `Prepare`/`PrepareContext`, reusing a statement, context timeouts |
| 5 | [5-transactions](5-transactions/main.go) | `BeginTx`, `Commit`, `Rollback`, the `defer tx.Rollback()` idiom |
| 6 | [6-repository](6-repository/main.go) | repository pattern: `models/` + `repository/` behind an interface |

## Core mental model

- **`sql.Open` does not connect.** It builds a connection **pool** and validates
  args. Use `db.Ping()` to actually reach the database.
- **`*sql.DB` is a pool**, safe for concurrent use. Open it **once** and share
  it — never one-per-request.
- **`Exec` vs `Query`:** `Exec` for statements returning no rows (INSERT/UPDATE/
  DELETE/DDL); `Query`/`QueryRow` for SELECTs that return rows.
- **Always use placeholders** (`?` for SQLite/MySQL, `$1` for Postgres). Never
  concatenate values into SQL — that's SQL injection.
- **`Scan` is positional:** the SELECT column order must match the `&pointer`
  order you pass to `Scan`.

## Key habits

- `defer db.Close()` after opening; `defer rows.Close()` after a `Query`.
- After looping `rows.Next()`, always check `rows.Err()`.
- Handle `sql.ErrNoRows` explicitly for single-row lookups.
- Use the `*Context` methods (`ExecContext`, `QueryContext`, `BeginTx`) so slow
  queries can be cancelled or time out.
- Wrap multi-step writes in a transaction with `defer tx.Rollback()`, and
  `Commit` only once everything succeeds.
- Never store plaintext passwords — hash with `bcrypt` before insert.

## Key packages

- **`database/sql`** — the standard, driver-agnostic SQL API.
- **`modernc.org/sqlite`** — pure-Go SQLite driver (no CGO).
- **`golang.org/x/crypto/bcrypt`** — password hashing.
- **`context`** — timeouts and cancellation for queries and transactions.
