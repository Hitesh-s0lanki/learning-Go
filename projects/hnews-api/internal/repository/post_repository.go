package repository

import (
	"context"
	"database/sql"
	"errors"

	"hnews-api/internal/models"
)

// PostRepository is the contract for posts, comments, and votes.
type PostRepository interface {
	Create(ctx context.Context, title, url string, userID int64) (*models.Post, error)
	GetByID(ctx context.Context, id int64) (*models.Post, error)
	List(ctx context.Context, f models.Filter) ([]models.Post, models.Metadata, error)

	AddVote(ctx context.Context, postID, userID int64) error
	RemoveVote(ctx context.Context, postID, userID int64) error

	AddComment(ctx context.Context, postID, userID int64, body string) (*models.Comment, error)
	ListComments(ctx context.Context, postID int64) ([]models.Comment, error)
}

type postRepo struct {
	db *sql.DB
}

// NewPostRepository returns a Postgres-backed PostRepository.
func NewPostRepository(db *sql.DB) PostRepository {
	return &postRepo{db: db}
}

func (r *postRepo) Create(ctx context.Context, title, url string, userID int64) (*models.Post, error) {
	const q = `
		INSERT INTO posts (title, url, user_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	var p = models.Post{Title: title, URL: url, UserID: userID}
	err := r.db.QueryRowContext(ctx, q, title, url, userID).Scan(&p.ID, &p.CreatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrDuplicate // duplicate title
		}
		return nil, err
	}
	// Re-read so the response includes the author name and zeroed counts.
	return r.GetByID(ctx, p.ID)
}

// selectPosts is the shared projection: a post plus its author and aggregated
// comment/vote counts. LEFT JOINs keep posts with zero comments/votes.
const selectPosts = `
	SELECT p.id, p.title, p.url, p.user_id, u.name AS author, p.created_at,
	       COUNT(DISTINCT c.id) AS comment_count,
	       COUNT(DISTINCT v.user_id) AS vote_count
	FROM posts p
	JOIN users u ON u.id = p.user_id
	LEFT JOIN comments c ON c.post_id = p.id
	LEFT JOIN votes v ON v.post_id = p.id`

func (r *postRepo) GetByID(ctx context.Context, id int64) (*models.Post, error) {
	q := selectPosts + `
		WHERE p.id = $1
		GROUP BY p.id, u.name`

	var p models.Post
	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&p.ID, &p.Title, &p.URL, &p.UserID, &p.Author, &p.CreatedAt,
		&p.CommentCount, &p.VoteCount,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *postRepo) List(ctx context.Context, f models.Filter) ([]models.Post, models.Metadata, error) {
	// ORDER BY is chosen from a fixed whitelist (never interpolate user input).
	orderBy := "p.created_at DESC"
	if f.Sort == "popular" {
		orderBy = "vote_count DESC, p.created_at DESC"
	}

	// count(*) OVER() returns the full match count alongside the page rows,
	// so we get total records and the page in a single round-trip.
	q := `
		SELECT COUNT(*) OVER() AS total, x.* FROM (
			` + selectPosts + `
			WHERE ($1 = '' OR p.title ILIKE '%' || $1 || '%')
			GROUP BY p.id, u.name
			ORDER BY ` + orderBy + `
			LIMIT $2 OFFSET $3
		) x`

	rows, err := r.db.QueryContext(ctx, q, f.Search, f.Limit(), f.Offset())
	if err != nil {
		return nil, models.Metadata{}, err
	}
	defer rows.Close()

	total := 0
	posts := []models.Post{}
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(
			&total,
			&p.ID, &p.Title, &p.URL, &p.UserID, &p.Author, &p.CreatedAt,
			&p.CommentCount, &p.VoteCount,
		); err != nil {
			return nil, models.Metadata{}, err
		}
		posts = append(posts, p)
	}
	if err := rows.Err(); err != nil {
		return nil, models.Metadata{}, err
	}

	meta := models.CalculateMetadata(total, f.Page, f.PageSize)
	return posts, meta, nil
}

func (r *postRepo) AddVote(ctx context.Context, postID, userID int64) error {
	const q = `INSERT INTO votes (post_id, user_id) VALUES ($1, $2)`
	_, err := r.db.ExecContext(ctx, q, postID, userID)
	if err != nil {
		if isUniqueViolation(err) {
			return ErrDuplicateVote
		}
		if isForeignKeyViolation(err) {
			return ErrNotFound // post doesn't exist
		}
		return err
	}
	return nil
}

func (r *postRepo) RemoveVote(ctx context.Context, postID, userID int64) error {
	const q = `DELETE FROM votes WHERE post_id = $1 AND user_id = $2`
	res, err := r.db.ExecContext(ctx, q, postID, userID)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound // no such vote to remove
	}
	return nil
}

func (r *postRepo) AddComment(ctx context.Context, postID, userID int64, body string) (*models.Comment, error) {
	const q = `
		INSERT INTO comments (post_id, user_id, body)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	c := models.Comment{PostID: postID, UserID: userID, Body: body}
	err := r.db.QueryRowContext(ctx, q, postID, userID, body).Scan(&c.ID, &c.CreatedAt)
	if err != nil {
		if isForeignKeyViolation(err) {
			return nil, ErrNotFound // post doesn't exist
		}
		return nil, err
	}

	// Fill in the author's display name for the response.
	_ = r.db.QueryRowContext(ctx, `SELECT name FROM users WHERE id = $1`, userID).Scan(&c.Author)
	return &c, nil
}

func (r *postRepo) ListComments(ctx context.Context, postID int64) ([]models.Comment, error) {
	const q = `
		SELECT c.id, c.post_id, c.user_id, u.name AS author, c.body, c.created_at
		FROM comments c
		JOIN users u ON u.id = c.user_id
		WHERE c.post_id = $1
		ORDER BY c.created_at ASC`

	rows, err := r.db.QueryContext(ctx, q, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []models.Comment{}
	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Author, &c.Body, &c.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, rows.Err()
}
