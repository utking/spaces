package domain

import (
	"errors"
	"strings"
	"time"
)

// Note represents a note in the system.
type Note struct {
	ID        string    `form:"id"      json:"id"`
	Title     string    `form:"title"   json:"title"`
	Content   string    `form:"content" json:"content"`
	Tags      []string  `form:"tags"    json:"tags"` // JSON string, can be empty
	CreatedAt time.Time `               json:"-"`
	UpdatedAt time.Time `               json:"-"`
}

// Trim trims the strings in the Note struct.
func (n *Note) Trim() {
	n.Title = strings.TrimSpace(n.Title)
	n.Content = strings.TrimSpace(n.Content)
	n.ID = strings.TrimSpace(n.ID)
}

// Validate checks the validity of the User struct fields.
func (n *Note) Validate() error {
	if len(n.Title) < 1 || len(n.Title) > 128 {
		return errors.New("title length must be between 1 and 128 characters")
	}

	if len(n.Tags) == 0 {
		return errors.New("at least one tag must be provided")
	}

	return nil
}

// NoteSearchRequest represents a request for searching notes.
type NoteSearchRequest struct {
	NoteID  string `query:"note_id"`
	Tag     string `query:"tag"`
	Title   string `query:"title"`   // Search by title
	Content string `query:"content"` // Search by content
	RequestPageMeta
}

// NoteRequest represents a request for creating/updating notes.
type NoteRequest struct {
	Title   string   `form:"title"   json:"title"`
	Content string   `form:"content" json:"content"`
	NoteID  string   `form:"note_id" json:"note_id"`
	Tags    []string `form:"tags"    json:"tags"` // JSON string, can be empty
	RequestPageMeta
}

// Trim trims the strings in the NoteRequest.
func (req *NoteRequest) Trim() {
	req.Title = strings.TrimSpace(req.Title)
	req.Content = strings.TrimSpace(req.Content)
	req.NoteID = strings.TrimSpace(req.NoteID)
}

// Validate checks the validity of the NoteRequest struct fields.
func (req *NoteRequest) Validate() error {
	if len(req.Title) < 1 || len(req.Title) > 128 {
		return errors.New("title length must be between 1 and 128 characters")
	}

	if len(req.Tags) == 0 {
		return errors.New("at least one tag must be provided")
	}

	return nil
}
