// Package mysql implements the MySQL database adapter for the application.
package mysql

import (
	"encoding/json"
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"xorm.io/builder"
)

const (
	sqlDialect = builder.MYSQL
)

// Adapter is a struct that holds the database connection.
// It is used to interact with the MySQL database.
type Adapter struct {
	db *sqlx.DB
}

// NewAdapter creates a new Adapter instance with the given data source URL.
// It connects to the MySQL database and returns an error if the connection fails.
func NewAdapter(dataSourceURL string) (*Adapter, error) {
	db, err := sqlx.Connect(sqlDialect, dataSourceURL)
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

func mySQLHasErrorCode(err error, code uint16) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == code
	}

	return false
}

func mySQLDuplicatePKError(err error) bool {
	return mySQLHasErrorCode(err, 1062) // MySQL error code for duplicate entry
}

func mySQLParentKeyViolationError(err error) bool {
	return mySQLHasErrorCode(err, 1451) // MySQL error code for foreign key constraint fails
}

func mySQLForeignKeyViolationError(err error) bool {
	return mySQLHasErrorCode(err, 1452) // MySQL error code for foreign key constraint fails
}

func mySQLConstraintViolationError(err error) bool {
	return mySQLHasErrorCode(err, 1048) // MySQL error code for not null constraint fails
}
