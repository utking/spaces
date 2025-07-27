package unittests

import (
	"errors"
	"fmt"
	"log"
	"path"
	"runtime"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	"github.com/utking/spaces/internal/config"
	"github.com/utking/spaces/migrations"
	"xorm.io/builder"
)

// CreateMySQLTestEngine creates a new test engine. Uses an in-memory sqlite database.
func CreateMySQLTestEngine() (*sqlx.DB, error) {
	x, err := sqlx.Connect(
		"mysql",
		"root:secret@tcp(localhost:13306)/test?charset=utf8mb4&parseTime=True&multiStatements=true",
	)
	if err != nil {
		if strings.Contains(err.Error(), "unknown driver") {
			return nil, fmt.Errorf(`mysql requires: import _ "github.com/go-sql-driver/mysql"
					or -tags mysql\n%w`, err)
		}

		return nil, err
	}

	return x, err
}

// CreateTestEngine creates a new test engine. Uses an in-memory sqlite database.
func CreateTestEngine() (*sqlx.DB, error) {
	x, err := sqlx.Connect("sqlite", "file::memory:?_txlock=immediate")
	if err != nil {
		if strings.Contains(err.Error(), "unknown driver") {
			return nil, fmt.Errorf(`sqlite requires: import _ "modernc.org/sqlite"
					or -tags sqlite,sqlite_unlock_notify\n%w`, err)
		}

		return nil, err
	}

	return x, err
}

// CreateTestDatabase creates a new test database.
func CreateTestDatabase(x *sqlx.DB) error {
	if err := createSchema(x); err != nil {
		return fmt.Errorf("error preparing test schemas, %w", err)
	}

	_, thisFilePath, _, _ := runtime.Caller(0) //nolint:dogsled // This is a common pattern
	fixturesDir := path.Join(thisFilePath, "..", "..", "_fixtures")

	if err := initFixtures(x, fixturesDir); err != nil {
		return fmt.Errorf("could not load fixtures, %w", err)
	}

	return loadFixtures()
}

// createSchema creates the schema for the test database.
func createSchema(engine *sqlx.DB) error {
	// Create tables that will be used in the tests
	switch engine.DriverName() {
	case string(config.SQLDriverSQLite):
		dbDriver, err := sqlite.WithInstance(engine.DB, &sqlite.Config{})
		if err != nil {
			return fmt.Errorf("failed to create SQLite instance: %w", err)
		}

		dataDriver, err := iofs.New(migrations.SQLiteFiles, builder.SQLITE)
		if err != nil {
			return fmt.Errorf("failed to get migration files for SQLite: %w", err)
		}

		m, err := migrate.NewWithInstance(
			"iofs",         // sourceName
			dataDriver,     // sourceInstance
			builder.SQLITE, // databaseName
			dbDriver,       // database instance
		)
		if err != nil {
			log.Fatalf("failed to create migration instance: %+v", err)
		}

		if err = m.Up(); err != nil {
			return fmt.Errorf("failed to apply migration: %w", err)
		}
	case string(config.SQLDriverMySQL):
		dbDriver, err := mysql.WithInstance(engine.DB, &mysql.Config{})
		if err != nil {
			return fmt.Errorf("failed to create MySQL instance: %w", err)
		}

		dataDriver, err := iofs.New(migrations.MySQLFiles, builder.MYSQL)
		if err != nil {
			return fmt.Errorf("failed to get migration files for SQLite: %w", err)
		}

		m, err := migrate.NewWithInstance(
			"iofs",        // sourceName
			dataDriver,    // sourceInstance
			builder.MYSQL, // databaseName
			dbDriver,      // database instance
		)
		if err != nil {
			log.Fatalf("failed to create migration instance: %+v", err)
		}

		if err = m.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				// No changes to apply, this is not an error
				return nil
			}

			return fmt.Errorf("failed to apply migration: %w", err)
		}
	default:
		return fmt.Errorf("unsupported database driver: %s", engine.DriverName())
	}

	return nil
}
