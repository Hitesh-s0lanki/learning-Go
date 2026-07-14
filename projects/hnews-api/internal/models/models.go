// Package models holds the domain types shared across the app. They are plain
// structs with JSON tags and know nothing about SQL, HTTP, or auth.
package models

import "time"

// User is an account in the system.
type User struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"` // never serialized to clients
	Avatar         string    `json:"avatar"`
	CreatedAt      time.Time `json:"created_at"`
}

// Post is a submitted link (Hacker News style).
type Post struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	UserID    int64     `json:"user_id"`
	Author    string    `json:"author"` // joined from users.name
	CreatedAt time.Time `json:"created_at"`

	CommentCount int `json:"comment_count"`
	VoteCount    int `json:"vote_count"`
}

// Comment belongs to a post and a user.
type Comment struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"post_id"`
	UserID    int64     `json:"user_id"`
	Author    string    `json:"author"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

// Filter carries list/pagination/search options coming from query params.
type Filter struct {
	Page     int
	PageSize int
	Sort     string // "new" (default) or "popular"
	Search   string
}

// Limit is the SQL LIMIT derived from the page size.
func (f Filter) Limit() int { return f.PageSize }

// Offset is the SQL OFFSET derived from page and page size.
func (f Filter) Offset() int { return (f.Page - 1) * f.PageSize }

// Metadata describes a paginated result set, returned alongside a list.
type Metadata struct {
	CurrentPage  int `json:"current_page"`
	PageSize     int `json:"page_size"`
	FirstPage    int `json:"first_page"`
	LastPage     int `json:"last_page"`
	TotalRecords int `json:"total_records"`
}

// CalculateMetadata builds pagination metadata from the total row count.
func CalculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{} // an empty result set has no pages
	}
	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     (totalRecords + pageSize - 1) / pageSize, // ceil division
		TotalRecords: totalRecords,
	}
}
