package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	"hnews-api/internal/models"
)

// UserRepository is the contract for user persistence.
type UserRepository interface {
	Create(ctx context.Context, name, email, hashedPassword, avatar string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id int64) (*models.User, error)
}

type userRepo struct {
	db *sql.DB
}

// NewUserRepository returns a Postgres-backed UserRepository.
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{db: db}
}

const userColumns = `id, name, email, hashed_password, avatar, created_at`

func (r *userRepo) Create(ctx context.Context, name, email, hashedPassword, avatar string) (*models.User, error) {
	const q = `
		INSERT INTO users (name, email, hashed_password, avatar)
		VALUES ($1, $2, $3, $4)
		RETURNING ` + userColumns

	var u models.User
	err := r.db.QueryRowContext(ctx, q, name, email, hashedPassword, avatar).
		Scan(&u.ID, &u.Name, &u.Email, &u.HashedPassword, &u.Avatar, &u.CreatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrDuplicate
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	const q = `SELECT ` + userColumns + ` FROM users WHERE email = $1`
	return r.scanOne(r.db.QueryRowContext(ctx, q, email))
}

func (r *userRepo) GetByID(ctx context.Context, id int64) (*models.User, error) {
	const q = `SELECT ` + userColumns + ` FROM users WHERE id = $1`
	return r.scanOne(r.db.QueryRowContext(ctx, q, id))
}

func (r *userRepo) scanOne(row *sql.Row) (*models.User, error) {
	var u models.User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.HashedPassword, &u.Avatar, &u.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// isUniqueViolation reports whether err is a Postgres unique-constraint error
// (SQLSTATE 23505). This lets us return a friendly domain error.
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

// isForeignKeyViolation reports whether err is a Postgres foreign-key error
// (SQLSTATE 23503) — e.g. voting/commenting on a post that doesn't exist.
func isForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23503"
}
