package services

import (
	"context"
	"errors"

	"gogs.utking.net/utking/spaces/internal/application/domain"
	"gogs.utking.net/utking/spaces/internal/ports"
)

type LastOpenedService struct {
	db ports.DBPort
}

// NewLastOpenedService creates a new instance of LastOpenedService.
func NewLastOpenedService(db ports.DBPort) *LastOpenedService {
	return &LastOpenedService{
		db: db,
	}
}

// GetLastOpened retrieves the last opened item ID for a given type and user ID.
func (s *LastOpenedService) GetLastOpened(
	ctx context.Context,
	itemType domain.LastOpenedType,
	userID string,
) (string, error) {
	if string(itemType) == "" || userID == "" {
		return "", errors.New("item type and user ID must be provided")
	}

	return s.db.GetLastOpened(ctx, itemType, userID)
}

// SetLastOpened sets the last opened item ID for a given type and user ID.
func (s *LastOpenedService) SetLastOpened(
	ctx context.Context,
	itemType domain.LastOpenedType,
	userID, itemID string,
) error {
	if string(itemType) == "" || userID == "" || itemID == "" {
		return errors.New("item type, user ID, and item ID must be provided")
	}

	return s.db.SetLastOpened(ctx, itemType, userID, itemID)
}
