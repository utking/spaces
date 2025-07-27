// Package sqlite implements the SQLite database adapter for the application.
package sqlite

import (
	"encoding/json"
	"errors"

	"github.com/jmoiron/sqlx"
	"modernc.org/sqlite"
	"xorm.io/builder"
)

const (
	sqlDialect = builder.SQLITE
	driverName = "sqlite"
)

// Adapter is a struct that holds the database connection.
// It is used to interact with the SQLite database.
type Adapter struct {
	db *sqlx.DB
}

// NewAdapter creates a new Adapter instance with the given data source URL.
// It connects to the SQLite database and returns an error if the connection fails.
func NewAdapter(dataSourceURL string) (*Adapter, error) {
	db, err := sqlx.Connect(driverName, dataSourceURL)
	if err != nil {
		return nil, err
	}

	return &Adapter{db: db}, nil
}

// NewAdapterWithDB creates a new Adapter instance with an existing sqlx.DB connection.
func NewAdapterWithDB(db *sqlx.DB) *Adapter {
	return &Adapter{db: db}
}

// toJSONString converts []string to a JSON string.
func toJSONString(anyVal []string) (string, error) {
	// Ensure the output is not nul for empty input
	if len(anyVal) == 0 {
		return "[]", nil // Return empty JSON array if input is empty
	}

	jsonBytes, err := json.Marshal(anyVal)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

func hasTag(tags []string, tag string) bool {
	for _, t := range tags {
		if t == tag {
			return true
		}
	}

	return false
}

func sqliteHasErrorCode(err error, code int) bool {
	var mysqlErr *sqlite.Error
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Code() == code
	}

	return false
}

// sqliteUniqViolation checks if the error is a SQLite UNIQUE constraint violation.
func sqliteUniqViolation(err error) bool {
	return sqliteHasErrorCode(err, 2067) // SQLITE_CONSTRAINT_UNIQUE
}

// sqliteNotNullViolation checks if the error is a SQLite NOT NULL constraint violation.
func sqliteConstraintViolation(err error) bool {
	return sqliteHasErrorCode(err, 1299) // not null constraint failed
}
