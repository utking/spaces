package db

import "github.com/utking/spaces/internal/application/domain"

type LastOpened struct {
	UserID string                `db:"user_id"`
	ItemID string                `db:"item_id"`
	Type   domain.LastOpenedType `db:"item_type"`
}

// TableName returns the name of the table for LastOpened.
func (LastOpened) TableName() string {
	return "last_opened"
}
