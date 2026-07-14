# hnews-api — a Hacker News–style REST API in Go

A complete, production-shaped REST API for a Hacker News clone: register/login
with JWT auth, submit link posts, upvote them, and comment. Backed by
**PostgreSQL**, built with the **standard library** `net/http` router (Go 1.22+),
`database/sql` + [pgx](https://github.com/jackc/pgx), and JWT bearer tokens.

This is a REST/JSON reimagining of the server-rendered
[learning-go-hnews](https://github.com/joefazee/learning-go-hnews) project,
rebuilt on Postgres with a clean, layered architecture.

## Features

- **JWT authentication** — register & login return a signed bearer token.
- **Posts** — submit links, list with pagination, search, and `new`/`popular`
  sorting (vote & comment counts aggregated via SQL).
- **Votes** — one upvote per user per post (enforced by a composite primary key).
- **Comments** — add and list comments per post.
- **Real-world touches** — bcrypt passwords, request validation with per-field
  errors, panic recovery, structured request logging, graceful shutdown,
  embedded SQL migrations, and connection-pool tuning.

## Architecture

```
projects/hnews-api/
├── cmd/api/main.go              # entry point: config, DB, migrate, serve, graceful shutdown
├── migrations/                  # embedded *.sql migrations (applied on startup)
│   ├── migrations.go            # //go:embed *.sql
│   └── 0001_init.sql
└── internal/
    ├── config/                  # env-driven configuration
    ├── database/                # Postgres pool + migration runner
    ├── models/                  # domain structs (User, Post, Comment, Filter, Metadata)
    ├── auth/                    # bcrypt passwords + JWT token manager (+ tests)
    ├── repository/              # all SQL, behind UserRepository / PostRepository interfaces
    └── api/                     # HTTP: server, routes, middleware, handlers, validation, JSON
```

The dependency flow is one-directional: `api → repository → database`, with
`models` shared and `auth`/`config` as leaf packages. Handlers depend on
repository **interfaces**, so the storage layer is swappable and testable.

## Requirements

- Go 1.22+ (uses method-based `http.ServeMux` routing)
- PostgreSQL 13+

## Getting started

```bash
cd projects/hnews-api

# 1. Create a database
createdb hnews

# 2. Configure (copy and edit)
cp .env.example .env
export $(grep -v '^#' .env | xargs)

# 3. Run — migrations are applied automatically on startup
go run ./cmd/api
# or: make run
```

The server listens on `PORT` (default `8080`). Configuration comes from the
environment (see `.env.example`): `DATABASE_URL` and `JWT_SECRET` are required.

## API reference

Base path: `/api/v1`. All request/response bodies are JSON. Protected endpoints
require an `Authorization: Bearer <token>` header.

| Method & path | Auth | Description |
|---|:---:|---|
| `GET /health` | — | Liveness check |
| `POST /api/v1/auth/register` | — | Create an account, returns a token |
| `POST /api/v1/auth/login` | — | Log in, returns a token |
| `GET /api/v1/users/me` | ✓ | The current user |
| `GET /api/v1/posts` | — | List posts (paginated/searchable/sortable) |
| `POST /api/v1/posts` | ✓ | Submit a link |
| `GET /api/v1/posts/{id}` | — | A single post with counts |
| `POST /api/v1/posts/{id}/vote` | ✓ | Upvote a post |
| `DELETE /api/v1/posts/{id}/vote` | ✓ | Remove your upvote |
| `GET /api/v1/posts/{id}/comments` | — | List a post's comments |
| `POST /api/v1/posts/{id}/comments` | ✓ | Comment on a post |

**List query params:** `page` (default 1), `page_size` (1–100, default 20),
`sort` (`new` | `popular`), `q` (title search, case-insensitive).

### Example session

```bash
# Register (returns {"user": {...}, "token": "..."})
curl -s -X POST localhost:8080/api/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"name":"Ada","email":"ada@example.com","password":"password123"}'

TOKEN="<paste token>"

# Submit a post
curl -s -X POST localhost:8080/api/v1/posts \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"title":"The Go Programming Language","url":"https://go.dev"}'

# Upvote it
curl -s -X POST localhost:8080/api/v1/posts/1/vote -H "Authorization: Bearer $TOKEN"

# Comment
curl -s -X POST localhost:8080/api/v1/posts/1/comments \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"body":"Great language."}'

# List popular posts
curl -s "localhost:8080/api/v1/posts?sort=popular&page_size=10"
```

## Responses & status codes

- Success bodies are enveloped: `{"post": {...}}`, `{"posts": [...], "metadata": {...}}`.
- Errors are uniform: `{"error": "message"}`; validation failures return
  `422` with `{"errors": {"field": "reason"}}`.
- Codes used: `200`, `201`, `400` (bad JSON), `401` (auth), `404`, `409`
  (duplicate email/title/vote), `422` (validation), `500`.

## Development

```bash
make test    # run tests (auth unit tests need no database)
make vet     # go vet
make build   # build ./bin/api
```

## Notes

- Passwords are hashed with bcrypt; the hash is never returned (`json:"-"`).
- SQL uses parameterized queries only (`$1, $2, ...`) — no string interpolation.
- Postgres error codes are translated to domain errors (`23505` → duplicate,
  `23503` → not found) so handlers stay driver-agnostic.
- Migrations are embedded in the binary and applied idempotently at startup.
