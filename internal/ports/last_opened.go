package ports

import (
	"context"

	"github.com/utking/spaces/internal/application/domain"
)

type LastOpenedService interface {
	GetLastOpened(ctx context.Context, itemType domain.LastOpenedType, uid string) (string, error)
	SetLastOpened(ctx context.Context, itemType domain.LastOpenedType, uid string, itemID string) error
}
