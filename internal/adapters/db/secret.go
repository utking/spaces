package db

import (
	"errors"
	"time"
)

// Secret represents a secret in the .
type Secret struct {
	Secret      []byte    `db:"secret"`
	Username    []byte    `db:"username"` // max len 128
	Name        string    `db:"name"`     // len 1-128
	URL         string    `db:"url"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	UserID      string    `db:"user_id"`
	ID          string    `db:"id"`   // primary key
	Tags        TagList   `db:"tags"` // JSON string, can be empty
}

// TableName returns the name of the table in the database.
func (Secret) TableName() string {
	return "password_record"
}

// Validate checks if the Secret is valid.
func (s *Secret) Validate() error {
	if len(s.Name) < 1 || len(s.Name) > 128 {
		return errors.New("secret name length must be between 1 and 128 characters")
	}

	if len(s.Username) > 128 {
		return errors.New("username length must not exceed 128 characters")
	}

	if s.UserID == "" {
		return errors.New("user ID must be set")
	}

	if len(s.Tags) == 0 {
		return errors.New("at least one tag must be set")
	}

	if len(s.Description) > 4096 {
		return errors.New("description length must not exceed 4096 characters")
	}

	if len(s.Secret) > 4096 {
		return errors.New("secret length must not exceed 4096 bytes")
	}

	return nil
}

type SecretExportItem struct {
	Tags        TagList `db:"tags"` // JSON string, can be empty
	Username    []byte  `db:"username"`
	Secret      []byte  `db:"secret"` // encrypted secret
	ID          string  `db:"id"`
	Name        string  `db:"name"`
	URL         string  `db:"url"`
	Description string  `db:"description"`
}
