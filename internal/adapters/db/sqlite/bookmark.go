package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"

	"github.com/utking/spaces/internal/adapters/db"
	"github.com/utking/spaces/internal/adapters/web/go_echo/helpers"
	"github.com/utking/spaces/internal/application/domain"
	"xorm.io/builder"
)

func (a *Adapter) GetBookmarkTags(
	ctx context.Context,
	userID string,
) ([]string, error) {
	sqlBuilder := builder.Dialect(sqlDialect).
		Select("DISTINCT tags").
		From(db.Bookmark{}.TableName())

	if userID != "" {
		sqlBuilder = sqlBuilder.Where(builder.Eq{"user_id": userID})
	}

	sqlStr, err := sqlBuilder.ToBoundSQL()
	if err != nil {
		return nil, errors.New("failed to build SQL query")
	}

	var tags []db.TagList

	err = a.db.SelectContext(ctx, &tags, sqlStr)
	if err != nil {
		return nil, errors.New("failed to execute query")
	}

	// Convert db.TagList to []string
	var resultMap = make(map[string]struct{}, 0)

	for _, tagList := range tags {
		for _, tag := range tagList {
			if tag != "" {
				resultMap[tag] = struct{}{}
			}
		}
	}

	result := make([]string, 0, len(resultMap))
	for tag := range resultMap {
		result = append(result, tag)
	}

	slices.Sort(result)

	return result, nil
}

func (a *Adapter) GetBookmarks(
	ctx context.Context,
	userID string,
	req *domain.BookmarkSearchRequest,
) ([]domain.Bookmark, error) {
	var items []db.Bookmark

	sqlBuilder := builder.Dialect(sqlDialect).
		Select(
			"id", "user_id", "title", "url", "tags",
		).
		From(db.Bookmark{}.TableName()).
		Where(builder.Eq{"user_id": userID}).
		OrderBy("title ASC")

	if req != nil {
		if req.URL != "" {
			sqlBuilder = sqlBuilder.Where(builder.Like{"url", req.URL})
		}

		if req.Title != "" {
			sqlBuilder = sqlBuilder.Where(builder.Like{"title", req.Title})
		}

		if req.RequestPageMeta.Limit > 0 {
			sqlBuilder = sqlBuilder.Limit(int(req.Limit))
		}
	}

	sqlStr, err := sqlBuilder.ToBoundSQL()
	if err != nil {
		return nil, errors.New("failed to build SQL query")
	}

	err = a.db.SelectContext(ctx, &items, sqlStr)
	if err != nil {
		return nil, errors.New("failed to execute query")
	}

	bookmarks := make([]domain.Bookmark, 0, len(items))
	for _, item := range items {
		if req != nil && req.Tag != "" && !hasTag(item.Tags, req.Tag) {
			continue
		}

		bookmarks = append(bookmarks, domain.Bookmark{
			ID:     item.ID,
			UserID: item.UserID,
			Title:  item.Title,
			URL:    item.URL,
			Tags:   item.Tags,
		})
	}

	return bookmarks, nil
}

func (a *Adapter) GetBookmarksCount(
	ctx context.Context,
	userID string,
	req *domain.BookmarkSearchRequest,
) (int64, error) {
	var items []db.Bookmark

	sqlBuilder := builder.Dialect(sqlDialect).
		Select("tags").
		From(db.Bookmark{}.TableName())

	if userID != "" {
		sqlBuilder = sqlBuilder.Where(builder.Eq{"user_id": userID})
	}

	if req != nil {
		if req.Title != "" {
			sqlBuilder = sqlBuilder.Where(builder.Like{"title", req.Title})
		}

		if req.URL != "" {
			sqlBuilder = sqlBuilder.Where(builder.Like{"url", req.URL})
		}
	}

	sqlStr, err := sqlBuilder.ToBoundSQL()
	if err != nil {
		return 0, errors.New("failed to build SQL query")
	}

	err = a.db.SelectContext(ctx, &items, sqlStr)
	if err != nil {
		return 0, errors.New("failed to execute query")
	}

	var count int64

	for _, item := range items {
		if req != nil && req.Tag != "" && !hasTag(item.Tags, req.Tag) {
			continue
		}

		count++
	}

	return count, nil
}

func (a *Adapter) GetBookmark(ctx context.Context, userID, id string) (*domain.Bookmark, error) {
	sqlBuilder := builder.Dialect(sqlDialect).
		Select(
			"id", "title", "url", "tags",
		).
		From(db.Bookmark{}.TableName()).
		Where(builder.Eq{"user_id": userID, "id": id})

	sqlStr, err := sqlBuilder.ToBoundSQL()
	if err != nil {
		return nil, errors.New("failed to build SQL query")
	}

	var item db.Bookmark

	err = a.db.GetContext(ctx, &item, sqlStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("bookmark not found")
		}

		return nil, errors.New("failed to execute query")
	}

	return &domain.Bookmark{
		ID:    item.ID,
		Title: item.Title,
		URL:   item.URL,
		Tags:  item.Tags,
	}, nil
}

func (a *Adapter) CreateBookmark(
	ctx context.Context,
	userID string,
	req *domain.Bookmark,
) (string, error) {
	req.Trim()

	if err := req.Validate(); err != nil {
		return "", err
	}

	req.ID = helpers.GenerateUUID()
	tags, _ := toJSONString(req.Tags)

	insertMap := builder.Eq{
		"id":      req.ID,
		"user_id": userID,
		"title":   req.Title,
		"url":     req.URL,
		"tags":    tags,
	}

	sqlBuilder := builder.Dialect(sqlDialect).
		Insert(insertMap).
		Into(db.Bookmark{}.TableName())

	sqlStr, args, err := sqlBuilder.ToSQL()
	if err != nil {
		return "", errors.New("failed to build SQL query")
	}

	_, err = a.db.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		return "", errors.New("failed to create bookmark")
	}

	return req.ID, nil
}

func (a *Adapter) UpdateBookmark(
	ctx context.Context,
	userID, id string,
	req *domain.Bookmark,
) (int64, error) {
	req.Trim()

	if err := req.Validate(); err != nil {
		return 0, err
	}

	tags, _ := toJSONString(req.Tags)

	updateMap := builder.Eq{
		"title": req.Title,
		"url":   req.URL,
		"tags":  tags,
	}

	sqlBuilder := builder.Dialect(sqlDialect).
		Update(updateMap).
		From(db.Bookmark{}.TableName()).
		Where(builder.Eq{"user_id": userID, "id": id})

	sqlStr, args, err := sqlBuilder.ToSQL()
	if err != nil {
		return 0, errors.New("failed to build SQL query")
	}

	result, err := a.db.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		return 0, errors.New("failed to update bookmark")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}

func (a *Adapter) DeleteBookmark(ctx context.Context, userID, id string) error {
	sqlBuilder := builder.Dialect(sqlDialect).
		Delete().
		From(db.Bookmark{}.TableName()).
		Where(builder.Eq{"user_id": userID, "id": id})

	sqlStr, err := sqlBuilder.ToBoundSQL()
	if err != nil {
		return errors.New("failed to build SQL query")
	}

	result, err := a.db.ExecContext(ctx, sqlStr)
	if err != nil {
		return errors.New("failed to delete bookmark")
	}

	_, err = result.RowsAffected()

	return err
}

func (a *Adapter) GetBookmarksMap(
	ctx context.Context,
	uid string,
	req *domain.BookmarkSearchRequest,
) ([]domain.Bookmark, error) {
	var (
		dbItems []db.Bookmark
		items   = make([]domain.Bookmark, 0)
	)

	sqlBuilder := builder.Dialect(sqlDialect).
		Select("title", "url", "tags").
		From(db.Bookmark{}.TableName()).
		Where(builder.Eq{"user_id": uid})

	if req != nil {
		if req.Title != "" {
			sqlBuilder = sqlBuilder.Where(builder.Like{"title", req.Title})
		}
	}

	sqlStr, err := sqlBuilder.ToBoundSQL()
	if err != nil {
		return nil, err
	}

	err = a.db.SelectContext(ctx, &dbItems, sqlStr)
	if err != nil {
		return nil, err
	}

	// TODO: filter by tag is given
	for _, item := range dbItems {
		items = append(items, domain.Bookmark{
			Title: item.Title,
			URL:   item.URL,
			Tags:  item.Tags,
		})
	}

	return items, nil
}

// SearchBookmarksByTerm retrieves bookmarks for a specific user based on a search term.
func (a *Adapter) SearchBookmarksByTerm(
	ctx context.Context,
	uid string,
	req *domain.BookmarkSearchRequest,
) ([]domain.Bookmark, error) {
	var dbItems []db.Bookmark

	sqlBuilder := builder.Dialect(sqlDialect).
		Select(
			"id",
			"tags",
			"url",
			"title",
		).
		From(db.Bookmark{}.TableName()).
		Where(builder.Eq{"user_id": uid}).
		OrderBy("title")

	if req != nil {
		sqlBuilder = sqlBuilder.Where(
			builder.Or(
				builder.Like{"url", req.URL},
				builder.Like{"title", req.Title},
				builder.Like{"tags", req.Title},
			),
		)

		if req.Limit > 0 {
			sqlBuilder = sqlBuilder.Limit(int(req.Limit))
		}
	}

	sqlStr, err := sqlBuilder.ToBoundSQL()
	if err != nil {
		return nil, err
	}

	err = a.db.SelectContext(ctx, &dbItems, sqlStr)
	if err != nil {
		return nil, err
	}

	items := make([]domain.Bookmark, len(dbItems))
	for i, item := range dbItems {
		items[i] = domain.Bookmark{
			ID:    item.ID,
			URL:   item.URL,
			Tags:  item.Tags,
			Title: item.Title,
		}
	}

	return items, nil
}
