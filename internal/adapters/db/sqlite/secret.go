package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"slices"

	"github.com/utking/spaces/internal/adapters/db"
	"github.com/utking/spaces/internal/adapters/web/go_echo/helpers"
	"github.com/utking/spaces/internal/application/domain"
	"xorm.io/builder"
)

// GetSecretTags retrieves secret tags for a user based on the provided request.
func (a *Adapter) GetSecretTags(
	ctx context.Context,
	uid string,
) ([]string, error) {
	sqlBuilder := builder.Dialect(sqlDialect).
		Select("DISTINCT tags").
		From(db.Secret{}.TableName())

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
		return nil, errors.New("failed to select secrets tags")
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

// GetSecrets retrieves secrets for a user based on the provided request.
func (a *Adapter) GetSecrets(
	ctx context.Context,
	uid string,
	req *domain.SecretSearchRequest,
) ([]domain.Secret, error) {
	var dbItems []db.Secret

	sqlBuilder := builder.Dialect(sqlDialect).
		Select(
			"id",
			"name",
			"tags",
		).
		From(db.Secret{}.TableName()).
		Where(builder.Eq{"user_id": uid}).
		OrderBy("name")

	if req != nil && req.Name != "" {
		sqlBuilder = sqlBuilder.Where(builder.Like{"name", req.Name})
	}

	sqlStr, err := sqlBuilder.ToBoundSQL()
	if err != nil {
		return nil, errors.New("failed to build SQL query")
	}

	if err = a.db.SelectContext(ctx, &dbItems, sqlStr); err != nil {
		return nil, errors.New("failed to execute query")
	}

	items := make([]domain.Secret, 0, len(dbItems))
	for _, item := range dbItems {
		if req != nil && req.Tag != "" && !hasTag(item.Tags, req.Tag) {
			continue
		}

		items = append(items, domain.Secret{
			ID:   item.ID,
			Name: item.Name,
		})
	}

	return items, nil
}

func (a *Adapter) GetSecretsCount(
	ctx context.Context,
	uid string,
	req *domain.SecretSearchRequest,
) (int64, error) {
	var count int64
	var dbItems []db.Secret

	sqlBuilder := builder.Dialect(sqlDialect).
		Select("tags").
		From(db.Secret{}.TableName())

	if uid != "" {
		sqlBuilder = sqlBuilder.Where(builder.Eq{"user_id": uid})
	}

	if req != nil && req.Name != "" {
		sqlBuilder = sqlBuilder.Where(builder.Like{"name", req.Name})
	}

	sqlStr, err := sqlBuilder.ToBoundSQL()
	if err != nil {
		return 0, errors.New("failed to build SQL query")
	}

	if err = a.db.SelectContext(ctx, &dbItems, sqlStr); err != nil {
		return 0, errors.New("failed to execute query")
	}

	for _, item := range dbItems {
		if req != nil && req.Tag != "" && !hasTag(item.Tags, req.Tag) {
			continue
		}

		count++
	}

	return count, err
}

func (a *Adapter) GetSecret(ctx context.Context, uid, id string) (*domain.Secret, error) {
	var dbItem db.Secret

	sqlBuilder := builder.Dialect(sqlDialect).
		Select(
			"id",
			"name",
			"username",
			"description",
			"url",
			"tags",
			"secret",
		).
		From(db.Secret{}.TableName()).
		Where(builder.Eq{"user_id": uid}).
		Where(builder.Eq{"id": id})

	sqlStr, err := sqlBuilder.ToBoundSQL()
	if err != nil {
		return nil, errors.New("failed to build SQL query")
	}

	if err = a.db.GetContext(ctx, &dbItem, sqlStr); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("secret not found")
		}

		return nil, errors.New("failed to execute query")
	}

	item := &domain.Secret{
		ID:              dbItem.ID,
		Name:            dbItem.Name,
		EncodedUsername: dbItem.Username,
		Description:     dbItem.Description,
		URL:             dbItem.URL,
		Tags:            dbItem.Tags,
		EncodedSecret:   dbItem.Secret,
	}

	return item, nil
}

func (a *Adapter) CreateSecret(
	ctx context.Context,
	uid string,
	req *domain.Secret,
) (id string, err error) {
	if req == nil {
		return "", errors.New("request cannot be nil")
	}

	if len(req.Tags) == 0 {
		return "", errors.New("tags cannot be empty")
	}

	if req.EncodedSecret == nil {
		req.EncodedSecret = []byte{}
	}

	if req.EncodedUsername == nil {
		req.EncodedUsername = []byte{}
	}

	id = helpers.GenerateUUID()
	tags, _ := toJSONString(req.Tags)
	sqlBuilder := builder.Dialect(sqlDialect).
		Into(db.Secret{}.TableName()).
		Insert(
			builder.Eq{"id": id},
			builder.Eq{"user_id": uid},
			builder.Eq{"name": req.Name},
			builder.Eq{"username": req.EncodedUsername},
			builder.Eq{"url": req.URL},
			builder.Eq{"description": req.Description},
			builder.Eq{"tags": tags},
			builder.Eq{"secret": req.EncodedSecret},
		)

	sqlStr, args, sqlErr := sqlBuilder.ToSQL()
	if sqlErr != nil {
		return "", sqlErr
	}

	// execute the insert statement
	if _, err = a.db.ExecContext(ctx, sqlStr, args...); err != nil {
		// check if the error is a unique constraint violation
		if sqliteUniqViolation(err) {
			return "", errors.New("secret with this name already exists")
		}

		if sqliteConstraintViolation(err) {
			return "", errors.New("failed to create secret: constraint violation")
		}

		return "", err
	}

	return id, nil
}

func (a *Adapter) UpdateSecret(
	ctx context.Context,
	uid, id string,
	req *domain.Secret,
) (affected int64, err error) {
	if req == nil {
		return 0, errors.New("request cannot be nil")
	}

	if len(req.Tags) == 0 {
		return 0, errors.New("tags cannot be empty")
	}

	if req.EncodedSecret == nil {
		req.EncodedSecret = []byte{}
	}

	if req.EncodedUsername == nil {
		req.EncodedUsername = []byte{}
	}

	tags, _ := toJSONString(req.Tags)
	sqlBuilder := builder.Dialect(sqlDialect).
		From(db.Secret{}.TableName()).
		Update(
			builder.Eq{"name": req.Name},
			builder.Eq{"username": req.EncodedUsername},
			builder.Eq{"url": req.URL},
			builder.Eq{"description": req.Description},
			builder.Eq{"tags": tags},
			builder.Eq{"secret": req.EncodedSecret},
		).
		Where(
			builder.And(
				builder.Eq{"user_id": uid},
				builder.Eq{"id": id},
			),
		)

	sqlStr, args, sqlErr := sqlBuilder.ToSQL()
	if sqlErr != nil {
		return 0, sqlErr
	}

	// start transaction
	tx, txErr := a.db.BeginTx(ctx, nil)
	if txErr != nil {
		return 0, txErr
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
	}()

	belongs, _ := a.secretBelongsToUser(ctx, tx, uid, id)
	if !belongs {
		return 0, errors.New("the secret does not exist/belong to current user")
	}

	// execute the update statement
	if _, err = tx.ExecContext(ctx, sqlStr, args...); err != nil {
		return 0, err
	}

	return 1, tx.Commit()
}

func (a *Adapter) DeleteSecret(ctx context.Context, uid, id string) (err error) {
	// create a transaction to ensure atomicity
	tx, txErr := a.db.BeginTx(ctx, nil)
	if txErr != nil {
		return txErr
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
	}()

	// check if the secret exists and belongs to the user
	belongs, _ := a.secretBelongsToUser(ctx, tx, uid, id)
	if !belongs {
		err = errors.New("the secret does not exist or does not belong to the user")
		return nil
	}

	sqlBuilder := builder.Dialect(sqlDialect).
		Delete().
		From(db.Secret{}.TableName()).
		Where(builder.Eq{"user_id": uid}).
		Where(builder.Eq{"id": id})

	var sqlStr string

	if sqlStr, err = sqlBuilder.ToBoundSQL(); err != nil {
		return err
	}

	if _, err = tx.ExecContext(ctx, sqlStr); err != nil {
		return err
	}

	_ = tx.Commit()

	return nil
}

func (a *Adapter) secretBelongsToUser(
	ctx context.Context,
	tx *sql.Tx,
	uid, id string,
) (bool, error) {
	sqlBuilder := builder.Dialect(sqlDialect).
		Select("COUNT(1)").
		From(db.Secret{}.TableName()).
		Where(builder.Eq{"user_id": uid}).
		Where(builder.Eq{"id": id})

	sqlStr, err := sqlBuilder.ToBoundSQL()
	if err != nil {
		return false, err
	}

	row := tx.QueryRowContext(ctx, sqlStr)
	if row.Err() != nil {
		return false, row.Err()
	}

	var counter int
	if err = row.Scan(&counter); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil // No rows means the secret does not belong to the user
		}

		return false, err // Some other error occurred
	}

	return counter > 0, nil
}

func (a *Adapter) GetSecretsMap(
	ctx context.Context,
	uid string,
	_ *domain.SecretSearchRequest,
) ([]domain.SecretExportItem, error) {
	var dbItems []db.SecretExportItem

	sqlBuilder := builder.Dialect(sqlDialect).
		Select(
			"id",
			"tags",
			"name",
			"username",
			"url",
			"description",
			"secret",
		).
		From(db.Secret{}.TableName()).
		Where(builder.Eq{"user_id": uid})

	sqlStr, err := sqlBuilder.ToBoundSQL()
	if err != nil {
		return nil, err
	}

	if err = a.db.SelectContext(ctx, &dbItems, sqlStr); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// If no records are found, return an empty map
			return []domain.SecretExportItem{}, nil
		}

		return nil, err
	}

	// TODO: Filter by tag if given
	itemsMap := make([]domain.SecretExportItem, 0, len(dbItems))
	for _, item := range dbItems {
		itemsMap = append(itemsMap, domain.SecretExportItem{
			Name:        item.Name,
			Tags:        item.Tags,
			URL:         item.URL,
			Description: item.Description,
			// The following fields are optional and may be empty
			EncodedPassword: item.Secret,
			EncodedUsername: item.Username,
		})
	}

	return itemsMap, nil
}

// SearchSecretsByTerm retrieves secrets for a user based on a search term.
func (a *Adapter) SearchSecretsByTerm(
	ctx context.Context,
	uid string,
	req *domain.SecretRequest,
) ([]domain.Secret, error) {
	var dbItems []db.Secret

	sqlBuilder := builder.Dialect(sqlDialect).
		Select(
			"id",
			"tags",
			"name",
		).
		From(db.Secret{}.TableName()).
		Where(builder.Eq{"user_id": uid}).
		OrderBy("name")

	if req != nil {
		if req.Name != "" && req.Username != "" && req.URL != "" && req.Description != "" {
			sqlBuilder = sqlBuilder.Where(
				builder.Or(
					builder.Like{"name", req.Name},
					builder.Like{"username", req.Username},
					builder.Like{"url", req.URL},
					builder.Like{"description", req.Description},
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

	if err = a.db.SelectContext(ctx, &dbItems, sqlStr); err != nil {
		return nil, err
	}

	items := make([]domain.Secret, len(dbItems))
	for i, item := range dbItems {
		items[i] = domain.Secret{
			ID:   item.ID,
			Tags: item.Tags,
			Name: item.Name,
		}
	}

	return items, nil
}

// UpdateEncryptedSecrets updates encrypted secrets for a user.
// The updated is wrapped in a transaction to ensure atomicity.
func (a *Adapter) UpdateEncryptedSecrets(
	ctx context.Context,
	_ string, // uid is not used in this method, but kept for compatibility
	items map[string]domain.EncryptSecret,
) (err error) {
	if len(items) == 0 {
		return nil
	}

	for id, item := range items {
		if len(item.Password) == 0 {
			item.Password = nil
		}

		if len(item.Username) == 0 {
			item.Username = nil
		}

		sqlBuilder := builder.Dialect(sqlDialect).
			From(db.Secret{}.TableName()).
			Update(
				builder.Eq{"secret": item.Password},
				builder.Eq{"username": item.Username},
			).
			Where(builder.Eq{"id": id})

		sqlStr, args, rbErr := sqlBuilder.ToSQL()
		if rbErr != nil {
			return rbErr
		}

		// execute the update statement
		if _, err = a.db.ExecContext(ctx, sqlStr, args...); err != nil {
			return err
		}
	}

	return nil
}
