package db

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type Bookmark struct {
	ID        string    `db:"id"`      // UUID
	UserID    string    `db:"user_id"` // UUID
	Title     string    `db:"title"`   // 1-255 characters
	URL       string    `db:"url"`     // 1-4096 characters
	CreatedAt time.Time `db:"created_at"`
	Tags      TagList   `db:"tags"` // JSON string, can be empty
}

// TableName returns the name of the database table for bookmarks.
func (Bookmark) TableName() string {
	return "bookmark"
}

type TagList []string

// Scan implements the sqlx Scanner interface for TagList.
func (pc *TagList) Scan(val interface{}) error {
	switch v := val.(type) {
	case []uint8:
		_ = json.Unmarshal(v, &pc)
		return nil
	case string:
		if v == "" {
			*pc = TagList{}
			return nil
		}

		return json.Unmarshal([]byte(v), &pc)
	default:
		pc = new(TagList)

		return nil
	}
}

func (pc *TagList) Value() (driver.Value, error) {
	return json.Marshal(pc)
}
