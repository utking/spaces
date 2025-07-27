package ports

import (
	"context"

	"github.com/utking/spaces/internal/application/domain"
)

// NotesService is an interface that defines the methods for notes-related operations.
type NotesService interface {
	GetTags(ctx context.Context, uid string) ([]string, error)
	GetItems(ctx context.Context, uid string, req *domain.NoteSearchRequest) ([]domain.Note, error)
	SearchItemsByTerm(ctx context.Context, uid string, req *domain.NoteRequest) ([]domain.Note, error)
	GetCount(ctx context.Context, uid string, req *domain.NoteSearchRequest) (int64, error)
	GetItem(ctx context.Context, uid, id string) (*domain.Note, error)
	Create(ctx context.Context, uid string, req *domain.Note) (string, error)
	Update(ctx context.Context, uid, id string, req *domain.Note) (int64, error)
	Delete(ctx context.Context, uid, id string) error
	GetItemsMap(ctx context.Context, uid string, req *domain.NoteSearchRequest) ([]domain.Note, error)
}
