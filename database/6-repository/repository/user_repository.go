// Package repository isolates all SQL for the users table behind a small,
// testable interface. The rest of the app depends on the interface, not on
// database/sql — so you could swap SQLite for Postgres, or a fake for tests,
// without touching business logic.
package repository

import (
	"context"
	"database/sql"
	"errors"

	"go-learning/database/6-repository/models"

	"golang.org/x/crypto/bcrypt"
)

// ErrNotFound is returned when a lookup matches no rows. Callers check this
// instead of the driver-specific sql.ErrNoRows, keeping them decoupled.
var ErrNotFound = errors.New("user not found")

// UserRepository is the contract the rest of the app codes against.
type UserRepository interface {
	Create(ctx context.Context, name, email, password string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	List(ctx context.Context) ([]models.User, error)
}

// sqlUserRepository is the concrete SQL-backed implementation.
type sqlUserRepository struct {
	db *sql.DB
}

// NewUserRepository wires a *sql.DB into the repository. Returning the
// interface type keeps callers unaware of the concrete struct.
func NewUserRepository(db *sql.DB) UserRepository {
	return &sqlUserRepository{db: db}
}

const columns = `id, name, email, hashed_password, created_at`

func (r *sqlUserRepository) Create(ctx context.Context, name, email, password string) (*models.User, error) {
	hp, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	res, err := r.db.ExecContext(ctx,
		`INSERT INTO users (name, email, hashed_password) VALUES (?, ?, ?)`,
		name, email, string(hp))
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Read the row back so the caller gets DB-populated fields (created_at).
	return r.getByID(ctx, id)
}

func (r *sqlUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	err := r.db.QueryRowContext(ctx,
		`SELECT `+columns+` FROM users WHERE email = ?`, email).
		Scan(&u.ID, &u.Name, &u.Email, &u.HashedPassword, &u.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound // translate driver error to our domain error
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *sqlUserRepository) List(ctx context.Context) ([]models.User, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT `+columns+` FROM users ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.HashedPassword, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *sqlUserRepository) getByID(ctx context.Context, id int64) (*models.User, error) {
	var u models.User
	err := r.db.QueryRowContext(ctx,
		`SELECT `+columns+` FROM users WHERE id = ?`, id).
		Scan(&u.ID, &u.Name, &u.Email, &u.HashedPassword, &u.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}
