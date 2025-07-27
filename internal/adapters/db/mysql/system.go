package mysql

import (
	"context"
	"errors"

	"gogs.utking.net/utking/spaces/internal/adapters/db"
	"gogs.utking.net/utking/spaces/internal/application/domain"
	"xorm.io/builder"
)

// GetSystemStats retrieves system statistics from the database.
// It returns the number of active servers, access points, users, owners, and tarifs.
func (a *Adapter) GetSystemStats(ctx context.Context, uid string) (*domain.SystemStats, error) {
	var stats domain.SystemStats

	// Get user stats
	usersStats, uErr := a.getUsersStats(ctx)
	if uErr != nil {
		return nil, uErr
	}

	// Get note stats
	noteStats, nErr := a.getNoteStats(ctx, uid)
	if nErr != nil {
		return nil, nErr
	}

	// Get secret stats
	secretStats, sErr := a.getSecretStats(ctx, uid)
	if sErr != nil {
		return nil, sErr
	}

	// Get bookmark stats
	bookmarkStats, bErr := a.getBookmarksStats(ctx, uid)
	if bErr != nil {
		return nil, bErr
	}

	// combine all stats
	stats = domain.SystemStats{
		ActiveUsers:   usersStats.ActiveUsers,
		InactiveUsers: usersStats.InactiveUsers,
		NoteTags:      noteStats.NoteTags,
		Notes:         noteStats.Notes,
		SecretTags:    secretStats.SecretTags,
		Secrets:       secretStats.Secrets,
		Bookmarks:     bookmarkStats.Bookmarks,
		BookmarkTags:  bookmarkStats.BookmarkTags,
	}

	return &stats, nil
}

func (a *Adapter) getUsersStats(ctx context.Context) (*db.UserStats, error) {
	var (
		stats      db.UserStats
		statsItems []db.UserStats
	)

	sqlBuilder := builder.Dialect(sqlDialect).
		From(stats.TableName()).
		Select(
			"sum(status = 10) as active_users",
			"sum(status = 0) as inactive_users",
		).
		GroupBy("status")

	sql, err := sqlBuilder.ToBoundSQL()
	if err != nil {
		return nil, err
	}

	err = a.db.SelectContext(ctx, &statsItems, sql)
	if err != nil {
		return nil, err
	}

	if len(statsItems) > 0 {
		for _, item := range statsItems {
			stats.ActiveUsers += item.ActiveUsers
			stats.InactiveUsers += item.InactiveUsers
		}
	}

	return &stats, nil
}

// getNoteStats retrieves note statistics from the database.
func (a *Adapter) getNoteStats(ctx context.Context, uid string) (*db.NoteStats, error) {
	var (
		items []string
		stats db.NoteStats
		nErr  error
		tErr  error
	)

	stats.Notes, nErr = a.GetNotesCount(ctx, uid, nil)
	items, tErr = a.GetNoteTags(ctx, uid)
	stats.NoteTags = int64(len(items))

	err := errors.Join(nErr, tErr)

	return &stats, err
}

// getSecretStats retrieves secret statistics from the database.
func (a *Adapter) getSecretStats(ctx context.Context, uid string) (*db.SecretStats, error) {
	var (
		items []string
		stats db.SecretStats
		nErr  error
		tErr  error
	)

	stats.Secrets, nErr = a.GetSecretsCount(ctx, uid, nil)
	items, tErr = a.GetSecretTags(ctx, uid)
	stats.SecretTags = int64(len(items))

	err := errors.Join(nErr, tErr)

	return &stats, err
}

// getBookmarksStats retrieves bookmark statistics from the database.
func (a *Adapter) getBookmarksStats(ctx context.Context, uid string) (*db.BookmarkStats, error) {
	var (
		stats db.BookmarkStats
		nErr  error
	)

	stats.Bookmarks, nErr = a.GetBookmarksCount(ctx, uid, nil)
	items, tErr := a.GetBookmarkTags(ctx, uid)
	stats.BookmarkTags = int64(len(items))

	err := errors.Join(nErr, tErr)

	return &stats, err
}
