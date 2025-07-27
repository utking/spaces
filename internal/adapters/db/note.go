package db

import (
	"errors"
	"time"
)

// Note represents a note in the system.
type Note struct {
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Title     string    `db:"title"`   // len 1-128
	Content   string    `db:"content"` // TEXT
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"` // > 0
	Tags      TagList   `db:"tags"`    // JSON string, can be empty
}

// TableName returns the name of the table in the database.
func (Note) TableName() string {
	return "note"
}

// Validate checks the validity of the Note fields. Content can be empty.
func (o *Note) Validate() error {
	if len(o.Title) < 1 || len(o.Title) > 128 {
		return errors.New("title length must be between 1 and 128 characters")
	}

	if o.UserID == "" {
		return errors.New("user ID must exist")
	}

	return nil
}
