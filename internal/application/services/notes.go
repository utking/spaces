package services

import (
	"context"
	"errors"

	"github.com/utking/spaces/internal/application/domain"
	"github.com/utking/spaces/internal/ports"
)

// NotesService is a struct that implements the NotesService interface.
type NotesService struct {
	db ports.DBPort
}

// NewNotesService creates a new instance of NotesService.
func NewNotesService(db ports.DBPort) *NotesService {
	return &NotesService{
		db: db,
	}
}

func (s *NotesService) GetTags(
	ctx context.Context,
	uid string,
) ([]string, error) {
	return s.db.GetNoteTags(ctx, uid)
}

func (s *NotesService) GetItems(
	ctx context.Context,
	uid string,
	req *domain.NoteSearchRequest,
) ([]domain.Note, error) {
	return s.db.GetNotes(ctx, uid, req)
}

func (s *NotesService) GetCount(
	ctx context.Context,
	uid string,
	req *domain.NoteSearchRequest,
) (int64, error) {
	return s.db.GetNotesCount(ctx, uid, req)
}

func (s *NotesService) GetItem(ctx context.Context, uid, id string) (*domain.Note, error) {
	return s.db.GetNote(ctx, uid, id)
}

func (s *NotesService) Create(ctx context.Context, uid string, req *domain.Note) (string, error) {
	if err := req.Validate(); err != nil {
		return "", err
	}

	return s.db.CreateNote(ctx, uid, req)
}

func (s *NotesService) Update(ctx context.Context, uid, id string, req *domain.Note) (int64, error) {
	// id must be given
	if id == "" {
		return 0, errors.New("note ID must be provided")
	}

	if err := req.Validate(); err != nil {
		return 0, err
	}

	return s.db.UpdateNote(ctx, uid, id, req)
}

func (s *NotesService) Delete(ctx context.Context, uid, id string) error {
	// id must be given
	if id == "" {
		return errors.New("note ID must be provided")
	}

	return s.db.DeleteNote(ctx, uid, id)
}

func (s *NotesService) GetItemsMap(
	ctx context.Context,
	uid string,
	req *domain.NoteSearchRequest,
) ([]domain.Note, error) {
	return s.db.GetNotesMap(ctx, uid, req)
}

// SearchItemsByTerm searches for notes by a search term.
func (s *NotesService) SearchItemsByTerm(
	ctx context.Context,
	uid string,
	req *domain.NoteRequest,
) ([]domain.Note, error) {
	return s.db.SearchNotesByTerm(ctx, uid, req)
}
