package db

// UserStats represents the statistics of users in the system.
type UserStats struct {
	ActiveUsers   int64 `db:"active_users"`
	InactiveUsers int64 `db:"inactive_users"`
}

// TableName returns the name of the table in the database.
func (UserStats) TableName() string {
	return "user"
}

type NoteStats struct {
	NoteTags int64 `db:"note_tags"`
	Notes    int64 `db:"notes"`
}

// TableName returns the name of the table in the database.
func (NoteStats) TableName() string {
	return "note"
}

type SecretStats struct {
	SecretTags int64 `db:"secret_tags"`
	Secrets    int64 `db:"secrets"`
}

// TableName returns the name of the table in the database.
func (SecretStats) TableName() string {
	return "password_record"
}

type BookmarkStats struct {
	Bookmarks    int64 `db:"bookmarks"`
	BookmarkTags int64 `db:"bookmark_tags"`
}

// TableName returns the name of the table in the database.
func (BookmarkStats) TableName() string {
	return "bookmark"
}
