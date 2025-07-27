package ports

import (
	"context"

	"github.com/utking/spaces/internal/application/domain"
)

// SystemStatsService is an interface that defines the methods for system statistics-related operations.
type SystemStatsService interface {
	GetStats(ctx context.Context, uid string) (*domain.SystemStats, error)
}
