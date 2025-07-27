package domain_test

import (
	"testing"

	"github.com/utking/spaces/internal/application/domain"
)

func TestNoteValidateOk(t *testing.T) {
	n := &domain.Note{
		ID:      "123e4567-e89b-12d3-a456-426614174000",
		Title:   "Example Note",
		Content: "This is an example note content.",
		Tags:    []string{"tag1", "tag2"},
	}

	if err := n.Validate(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestNoteValidateErr(t *testing.T) {
	tests := []struct {
		name string
		n    *domain.Note
	}{
		{"EmptyTitle", &domain.Note{Title: "", Content: "Content"}},
		{"LongTitle", &domain.Note{Title: "a" + string(make([]byte, 129)), Content: "Content"}},
		{"EmptyTags", &domain.Note{Title: "Title", Content: "Content", Tags: []string{}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.Validate(); err == nil {
				t.Errorf("expected error for %s, got nil", tt.name)
			}
		})
	}
}

func TestNoteRequestValidateOk(t *testing.T) {
	req := &domain.NoteRequest{
		Title:   "Valid Title",
		Content: "Content",
		Tags:    []string{"tag1", "tag2"},
	}

	if err := req.Validate(); err != nil {
		t.Errorf("expected no error for valid request, got %v", err)
	}
}

func TestNoteRequestValidateFail(t *testing.T) {
	req := &domain.NoteRequest{
		Title:   "",
		Content: "Content",
		Tags:    []string{},
	}

	if err := req.Validate(); err == nil {
		t.Error("expected error for empty title and tags, got nil")
	}
}
