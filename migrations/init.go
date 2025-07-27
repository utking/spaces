// Package migrations provides access to SQL migration files used in the application.
package migrations

import "embed"

var (
	//go:embed mysql/*.sql
	MySQLFiles embed.FS
	//go:embed sqlite3/*.sql
	SQLiteFiles embed.FS
)
