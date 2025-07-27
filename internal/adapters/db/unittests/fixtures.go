// Package unittests provides utilities for loading test fixtures into a database.
package unittests

import (
	"errors"
	"time"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/jmoiron/sqlx"
	"github.com/utking/spaces/internal/config"
)

var (
	fixtures *testfixtures.Loader
)

// initFixtures initializes the fixtures loader.
// It can use one of the following dialects:
//   - "postgresql"
//   - "timescaledb"
//   - "mysql"
//   - "mariadb"
//   - "sqlite"
//   - "sqlserver"
//
// Fixtures are loaded from the given directory.
func initFixtures(engine *sqlx.DB, fixturesDir string) error {
	var (
		dialect config.SQLDriver
		err     error
	)

	switch config.SQLDriver(engine.DriverName()) {
	case config.SQLDriverSQLite:
		dialect = config.SQLDriverSQLite
	case config.SQLDriverMySQL:
		dialect = config.SQLDriverMySQL
	default:
		return errors.New("unsupported SQL driver for fixtures: " + engine.DriverName())
	}

	fixtures, err = testfixtures.New(
		testfixtures.Database(engine.DB),
		testfixtures.Dialect(string(dialect)),
		testfixtures.DangerousSkipTestDatabaseCheck(),
		testfixtures.Directory(fixturesDir),
	)

	return err
}

// loadFixtures load fixtures for a test database.
func loadFixtures() error {
	var err error
	// (doubt) database transaction conflicts could occur and result in ROLLBACK? just try for a few times.
	for range 5 {
		if err = fixtures.Load(); err == nil {
			break
		}

		time.Sleep(time.Second)
	}

	return err
}
