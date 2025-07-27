package domain_test

import (
	"testing"

	"github.com/utking/spaces/internal/application/domain"
)

func TestBookmarkValidateNoErr(t *testing.T) {
	b := &domain.Bookmark{
		ID:     "123e4567-e89b-12d3-a456-426614174000",
		UserID: "123e4567-e89b-12d3-a456-426614174001",
		Title:  "Example Bookmark",
		URL:    "https://example.com",
		Tags:   []string{"tag1", "tag2"},
	}

	if err := b.Validate(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestBookmarkValidateErr(t *testing.T) {
	tests := []struct {
		name string
		b    *domain.Bookmark
	}{
		{"EmptyTitle", &domain.Bookmark{Title: "", URL: "https://example.com"}},
		{"EmptyURL", &domain.Bookmark{Title: "Example Bookmark", URL: ""}},
		{"LongTitle", &domain.Bookmark{Title: "a" + string(make([]byte, 256)), URL: "https://example.com"}},
		{"LongURL", &domain.Bookmark{Title: "Example Bookmark", URL: "https://" + string(make([]byte, 4097))}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.Validate(); err == nil {
				t.Errorf("expected error for %s, got nil", tt.name)
			}
		})
	}
}
