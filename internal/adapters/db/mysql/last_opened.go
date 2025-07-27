package mysql

import (
	"context"
	"database/sql"
	"errors"

	"gogs.utking.net/utking/spaces/internal/adapters/db"
	"gogs.utking.net/utking/spaces/internal/application/domain"
	"xorm.io/builder"
)

func (a *Adapter) GetLastOpened(
	ctx context.Context,
	itemType domain.LastOpenedType,
	userID string,
) (string, error) {
	var itemID string

	sqlBuilder := builder.Dialect(sqlDialect).
		Select("item_id").
		From(db.LastOpened{}.TableName()).
		Where(builder.And(
			builder.Eq{"item_type": string(itemType)},
			builder.Eq{"user_id": userID},
		))

	sqlStr, err := sqlBuilder.ToBoundSQL()
	if err != nil {
		return "", err
	}

	err = a.db.GetContext(ctx, &itemID, sqlStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// If no record is found, return an empty string
			return "", nil
		}

		return "", err
	}

	return itemID, nil
}

func (a *Adapter) SetLastOpened(
	ctx context.Context,
	itemType domain.LastOpenedType,
	userID, newItemID string,
) error {
	// delete if itemID is empty
	if newItemID == "" {
		sqlBuilder := builder.Dialect(sqlDialect).
			Delete().
			From(db.LastOpened{}.TableName()).
			Where(builder.And(
				builder.Eq{"item_type": string(itemType)},
				builder.Eq{"user_id": userID},
			))

		sqlStr, err := sqlBuilder.ToBoundSQL()
		if err != nil {
			return err
		}

		_, err = a.db.ExecContext(ctx, sqlStr)

		return err
	}

	// Check if the last opened item already exists
	existingItemID, err := a.GetLastOpened(ctx, itemType, userID)
	if err != nil {
		return err
	}

	if existingItemID != "" {
		// Do not update if the itemID is the same
		if existingItemID == newItemID {
			return nil
		}

		// If it exists, update the existing record
		sqlBuilder := builder.Dialect(sqlDialect).
			Update(builder.Eq{"item_id": newItemID}).
			From(db.LastOpened{}.TableName()).
			Where(builder.And(
				builder.Eq{"item_type": string(itemType)},
				builder.Eq{"user_id": userID},
			))

		sqlStr, bErr := sqlBuilder.ToBoundSQL()
		if bErr != nil {
			return bErr
		}

		_, err = a.db.ExecContext(ctx, sqlStr)

		return err
	}

	// If it does not exist, insert a new record
	sqlBuilder := builder.Dialect(sqlDialect).
		Insert(builder.Eq{
			"user_id":   userID,
			"item_id":   newItemID,
			"item_type": string(itemType),
		}).
		Into(db.LastOpened{}.TableName())

	sqlStr, err := sqlBuilder.ToBoundSQL()
	if err != nil {
		return err
	}

	_, err = a.db.ExecContext(ctx, sqlStr)

	return err
}
