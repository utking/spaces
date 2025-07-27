package services

import (
	"context"
	"errors"

	"github.com/utking/spaces/internal/application/domain"
	"github.com/utking/spaces/internal/ports"
)

type BookmarkService struct {
	db ports.DBPort
}

// NewBookmarkService creates a new instance of BookmarkService.
func NewBookmarkService(db ports.DBPort) *BookmarkService {
	return &BookmarkService{
		db: db,
	}
}

// GetTags retrieves all unique tags for the user's bookmarks.
func (s *BookmarkService) GetTags(ctx context.Context, uid string) ([]string, error) {
	return s.db.GetBookmarkTags(ctx, uid)
}

func (s *BookmarkService) GetItems(
	ctx context.Context,
	uid string,
	req *domain.BookmarkSearchRequest,
) ([]domain.Bookmark, error) {
	return s.db.GetBookmarks(ctx, uid, req)
}

func (s *BookmarkService) GetCount(
	ctx context.Context,
	uid string,
	req *domain.BookmarkSearchRequest,
) (int64, error) {
	return s.db.GetBookmarksCount(ctx, uid, req)
}

func (s *BookmarkService) GetItem(ctx context.Context, uid, id string) (*domain.Bookmark, error) {
	// id must be given
	if id == "" {
		return nil, errors.New("bookmark ID must be provided")
	}

	return s.db.GetBookmark(ctx, uid, id)
}

func (s *BookmarkService) Create(ctx context.Context, uid string, req *domain.Bookmark) (string, error) {
	// tags must be given
	if len(req.Tags) == 0 {
		return "", errors.New("at least one tag must be provided")
	}

	if err := req.Validate(); err != nil {
		return "", err
	}

	return s.db.CreateBookmark(ctx, uid, req)
}

func (s *BookmarkService) Update(
	ctx context.Context,
	uid, id string,
	req *domain.Bookmark,
) (int64, error) {
	// id must be given
	if id == "" {
		return 0, errors.New("bookmark ID must be provided")
	}

	if len(req.Tags) == 0 {
		return 0, errors.New("at least one tag must be provided")
	}

	if err := req.Validate(); err != nil {
		return 0, err
	}

	return s.db.UpdateBookmark(ctx, uid, id, req)
}

func (s *BookmarkService) Delete(ctx context.Context, uid, id string) error {
	// id must be given
	if id == "" {
		return errors.New("bookmark ID must be provided")
	}

	return s.db.DeleteBookmark(ctx, uid, id)
}

func (s *BookmarkService) GetItemsMap(
	ctx context.Context,
	uid string,
	req *domain.BookmarkSearchRequest,
) ([]domain.Bookmark, error) {
	return s.db.GetBookmarksMap(ctx, uid, req)
}

func (s *BookmarkService) SearchItemsByTerm(
	ctx context.Context,
	uid string,
	req *domain.BookmarkSearchRequest,
) ([]domain.Bookmark, error) {
	return s.db.SearchBookmarksByTerm(ctx, uid, req)
}
