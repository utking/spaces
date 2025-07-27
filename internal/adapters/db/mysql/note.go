package mysql

import (
	"context"
	"errors"
	"slices"

	"github.com/utking/spaces/internal/adapters/db"
	"github.com/utking/spaces/internal/adapters/web/go_echo/helpers"
	"github.com/utking/spaces/internal/application/domain"
	"xorm.io/builder"
)

// GetNoteTags retrieves note tags for a specific user based on the provided request parameters.
func (a *Adapter) GetNoteTags(
	ctx context.Context,
	uid string,
) ([]string, error) {
	sqlBuilder := builder.Dialect(sqlDialect).
		Select("DISTINCT tags").
		From(db.Note{}.TableName())

	if uid != "" {
		sqlBuilder = sqlBuilder.Where(builder.Eq{"user_id": uid})
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

// GetNotes retrieves notes for a specific user based on the provided request parameters.
func (a *Adapter) GetNotes(
	ctx context.Context,
	uid string,
	req *domain.NoteSearchRequest,
) ([]domain.Note, error) {
	if req == nil {
		// If no request is provided, return an empty slice
		return nil, nil
	}

	var dbItems []db.Note

	sqlBuilder := builder.Dialect(sqlDialect).
		Select(
			"id",
			"title",
		).
		From(db.Note{}.TableName()).
		Where(builder.Eq{"user_id": uid}).
		OrderBy("title")

	// req != nil that is checked above
	if req.Content != "" {
		sqlBuilder = sqlBuilder.Where(builder.Like{"content", req.Content})
	}

	if req.Title != "" {
		sqlBuilder = sqlBuilder.Where(builder.Like{"title", req.Title})
	}

	sqlBuilder = sqlBuilder.Where(
		builder.Expr("? MEMBER OF(tags)", req.Tag),
	)

	sqlStr, err := sqlBuilder.ToBoundSQL()
	if err != nil {
		return nil, errors.New("failed to build SQL query")
	}

	err = a.db.SelectContext(ctx, &dbItems, sqlStr)
	if err != nil {
		return nil, errors.New("failed to execute query")
	}

	items := make([]domain.Note, len(dbItems))
	for i, item := range dbItems {
		items[i] = domain.Note{
			ID:    item.ID,
			Title: item.Title,
		}
	}

	return items, nil
}

func (a *Adapter) GetNotesCount(
	ctx context.Context,
	uid string,
	req *domain.NoteSearchRequest,
) (int64, error) {
	var count int64

	sqlBuilder := builder.Dialect(sqlDialect).
		Select("COUNT(1) as `count`").
		From(db.Note{}.TableName())

	if uid != "" {
		sqlBuilder = sqlBuilder.Where(builder.Eq{"user_id": uid})
	}

	if req != nil {
		if req.Title != "" {
			sqlBuilder = sqlBuilder.Where(builder.Like{"title", req.Title})
		}

		if req.Tag != "" {
			sqlBuilder = sqlBuilder.Where(
				builder.Expr("? MEMBER OF(tags)", req.Tag),
			)
		}
	}

	sqlStr, err := sqlBuilder.ToBoundSQL()
	if err != nil {
		return 0, errors.New("failed to build SQL query")
	}

	err = a.db.GetContext(ctx, &count, sqlStr)

	return count, err
}

func (a *Adapter) GetNote(ctx context.Context, uid, id string) (*domain.Note, error) {
	var dbItem db.Note

	sqlBuilder := builder.Dialect(sqlDialect).
		Select(
			"id",
			"tags",
			"title",
			"content",
		).
		From(db.Note{}.TableName()).
		Where(builder.And(
			builder.Eq{"user_id": uid},
			builder.Eq{"id": id},
		))

	sqlStr, err := sqlBuilder.ToBoundSQL()
	if err != nil {
		return nil, err
	}

	err = a.db.GetContext(ctx, &dbItem, sqlStr)
	if err != nil {
		return nil, err
	}

	item := &domain.Note{
		ID:      dbItem.ID,
		Tags:    dbItem.Tags,
		Title:   dbItem.Title,
		Content: dbItem.Content,
	}

	return item, nil
}

func (a *Adapter) CreateNote(ctx context.Context, uid string, req *domain.Note) (string, error) {
	if req == nil {
		return "", errors.New("note request cannot be nil")
	}

	if valErr := req.Validate(); valErr != nil {
		return "", valErr
	}

	id := helpers.GenerateUUID()
	tags, _ := toJSONString(req.Tags)
	sqlBuilder := builder.Dialect(sqlDialect).
		Into(db.Note{}.TableName()).
		Insert(
			builder.Eq{"id": id},
			builder.Eq{"tags": tags},
			builder.Eq{"user_id": uid},
			builder.Eq{"title": req.Title},
			builder.Eq{"content": req.Content},
		)

	// INFO: Cannot use ToBoundSQL here because it will ruin \n in the content field
	sqlStr, args, err := sqlBuilder.ToSQL()
	if err != nil {
		return "", err
	}

	_, err = a.db.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		// Check for duplicate entry error (MySQL error code 1062)
		if mySQLDuplicatePKError(err) {
			return "", errors.New("note with this title already exists")
		}

		return "", err
	}

	return id, nil
}

func (a *Adapter) UpdateNote(ctx context.Context, uid, id string, req *domain.Note) (int64, error) {
	if req == nil {
		return 0, errors.New("note request cannot be nil")
	}

	if valErr := req.Validate(); valErr != nil {
		return 0, valErr
	}

	tags, _ := toJSONString(req.Tags)
	sqlBuilder := builder.Dialect(sqlDialect).
		From(db.Note{}.TableName()).
		Update(
			builder.Eq{"content": req.Content},
			builder.Eq{"title": req.Title},
			builder.Eq{"tags": tags},
		).
		Where(
			builder.And(
				builder.Eq{"user_id": uid},
				builder.Eq{"id": id},
			),
		)

	// INFO: Cannot use ToBoundSQL here because it will ruin \n in the content field
	sqlStr, args, err := sqlBuilder.ToSQL()
	if err != nil {
		return 0, err
	}

	if _, err = a.db.ExecContext(ctx, sqlStr, args...); err != nil {
		// Check for duplicate entry error (MySQL error code 1062)
		if mySQLDuplicatePKError(err) {
			return 0, errors.New("note with this title already exists")
		}

		return 0, err
	}

	return 1, err
}

func (a *Adapter) DeleteNote(ctx context.Context, uid, id string) error {
	sqlBuilder := builder.Dialect(sqlDialect).
		Delete().
		From(db.Note{}.TableName()).
		Where(
			builder.And(
				builder.Eq{"user_id": uid},
				builder.Eq{"id": id},
			),
		)

	sqlStr, err := sqlBuilder.ToBoundSQL()
	if err != nil {
		return err
	}

	_, err = a.db.ExecContext(ctx, sqlStr)

	return err
}

func (a *Adapter) GetNotesMap(
	ctx context.Context,
	uid string, req *domain.NoteSearchRequest,
) ([]domain.Note, error) {
	var (
		dbItems []db.Note
		items   = make([]domain.Note, 0)
	)

	sqlBuilder := builder.Dialect(sqlDialect).
		Select("title", "content", "tags").
		From(db.Note{}.TableName()).
		Where(builder.Eq{"user_id": uid})

	if req != nil {
		if req.Title != "" {
			sqlBuilder = sqlBuilder.Where(builder.Like{"title", req.Title})
		}

		if req.Tag != "" {
			sqlBuilder = sqlBuilder.Where(
				builder.Expr("? MEMBER OF(tags)", req.Tag),
			)
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

	for _, item := range dbItems {
		items = append(items, domain.Note{
			Title:   item.Title,
			Content: item.Content,
			Tags:    item.Tags,
		})
	}

	return items, nil
}

// SearchNotesByTerm retrieves notes for a specific user based on a search term.
func (a *Adapter) SearchNotesByTerm(
	ctx context.Context,
	uid string,
	req *domain.NoteRequest,
) ([]domain.Note, error) {
	var dbItems []db.Note

	sqlBuilder := builder.Dialect(sqlDialect).
		Select(
			"id",
			"tags",
			"title",
		).
		From(db.Note{}.TableName()).
		Where(builder.Eq{"user_id": uid}).
		OrderBy("title")

	if req != nil {
		if req.Content != "" && req.Title != "" {
			sqlBuilder = sqlBuilder.Where(
				builder.Or(
					builder.Like{"content", req.Content},
					builder.Like{"title", req.Title},
				),
			)
		}

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

	items := make([]domain.Note, len(dbItems))
	for i, item := range dbItems {
		items[i] = domain.Note{
			ID:    item.ID,
			Tags:  item.Tags,
			Title: item.Title,
		}
	}

	return items, nil
}
