// Package models holds the plain data types shared across the app.
// Models know nothing about SQL or HTTP — they are just structs.
package models

import "time"

// User mirrors a row in the users table.
type User struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"` // never serialized
	CreatedAt      time.Time `json:"created_at"`
}
