package services

import (
	"context"
	"runtime"

	"gogs.utking.net/utking/spaces/internal/application/domain"
	"gogs.utking.net/utking/spaces/internal/ports"
)

// SysStatService is a struct that implements the SysStatService interface.
type SysStatService struct {
	db ports.DBPort
}

// NewSysStatService creates a new instance of SysStatService.
func NewSysStatService(db ports.DBPort) *SysStatService {
	return &SysStatService{
		db: db,
	}
}

// GetStats retrieves the aggregated system stats from the database.
func (a *SysStatService) GetStats(ctx context.Context, uid string) (*domain.SystemStats, error) {
	stats, err := a.db.GetSystemStats(ctx, uid)
	if err != nil {
		return nil, err
	}

	var mem runtime.MemStats

	runtime.ReadMemStats(&mem)

	stats.MemoryAlloc = int64(mem.Alloc) //nolint:gosec // Just ignoting it here
	stats.Allocations = mem.Mallocs
	stats.TotalCPU = runtime.NumCPU()

	return stats, nil
}
