package domain

import (
	"errors"
	"strings"
	"time"
)

type Bookmark struct {
	ID        string    `json:"-"     form:"id"`      // len 36
	UserID    string    `json:"-"     form:"user_id"` // len 36
	Title     string    `json:"title" form:"title"`   // len 1-255
	URL       string    `json:"url"   form:"url"`     // len 1-4096
	CreatedAt time.Time `json:"-"`
	Tags      []string  `json:"tags"  form:"tags"` // JSON string, can be empty
}

// Trim trims the strings in the Bookmark struct.
func (b *Bookmark) Trim() {
	b.Title = strings.TrimSpace(b.Title)
	b.URL = strings.TrimSpace(b.URL)
}

// Validate checks the validity of the Bookmark struct fields.
func (b *Bookmark) Validate() error {
	var err error

	if b.Title == "" {
		err = errors.Join(err, errors.New("title cannot be empty;"))
	}

	if b.URL == "" {
		err = errors.Join(err, errors.New("URL cannot be empty;"))
	}

	if len(b.Title) > 255 {
		err = errors.Join(err, errors.New("title cannot be longer that 255 characters;"))
	}

	if len(b.URL) > 4096 {
		err = errors.Join(err, errors.New("URL cannot be longer that 4096 characters;"))
	}

	return err
}

// BookmarkSearchRequest represents a request for searching bookmarks.
type BookmarkSearchRequest struct {
	UserID string `query:"user_id"`
	Title  string `query:"title"`
	URL    string `query:"url"`
	Tag    string `query:"tag"`
	RequestPageMeta
}
