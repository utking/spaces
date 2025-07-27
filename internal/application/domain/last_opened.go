package domain

const (
	LastOpenedTypeNote     = LastOpenedType("note_id")
	LastOpenedTypeBookmark = LastOpenedType("bookmark_tag")
)

type LastOpenedType string

type LastOpened struct {
	UserID string         `json:"user_id"`
	ItemID string         `json:"item_id"`
	Type   LastOpenedType `json:"item_type"`
}
