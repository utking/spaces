package ports

import (
	"context"

	"github.com/utking/spaces/internal/application/domain"
)

type BookmarkService interface {
	GetTags(ctx context.Context, uid string) ([]string, error)
	GetItems(ctx context.Context, uid string, req *domain.BookmarkSearchRequest) ([]domain.Bookmark, error)
	SearchItemsByTerm(ctx context.Context, uid string, req *domain.BookmarkSearchRequest) ([]domain.Bookmark, error)
	GetCount(ctx context.Context, uid string, req *domain.BookmarkSearchRequest) (int64, error)
	GetItem(ctx context.Context, uid, id string) (*domain.Bookmark, error)
	Create(ctx context.Context, uid string, req *domain.Bookmark) (string, error)
	Update(ctx context.Context, uid, id string, req *domain.Bookmark) (int64, error)
	Delete(ctx context.Context, uid, id string) error
	GetItemsMap(ctx context.Context, uid string, req *domain.BookmarkSearchRequest) ([]domain.Bookmark, error)
}
