// Package migrate provides the command to apply database migrations.
package migrate

import (
	"errors"
	"log"

	"gogs.utking.net/utking/spaces/internal/config"
	"gogs.utking.net/utking/spaces/migrations"
	"xorm.io/builder"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/spf13/cobra"
)

var down bool

// Init initializes the server command and adds it to the root command.
func Init(rootCmd *cobra.Command) {
	migrateCmd.PersistentFlags().BoolVarP(&down, "down", "d", false, "run a migration down")
	rootCmd.AddCommand(migrateCmd)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Apply database migrations",
	Run: func(cmd *cobra.Command, _ []string) {
		cfg := config.New()

		var (
			data   source.Driver
			dsnStr string
			err    error
		)

		if cfg.GetSQLDriver() == config.SQLDriverSQLite {
			cmd.Println("Using SQLite database for migrations")

			data, err = iofs.New(migrations.SQLiteFiles, builder.SQLITE)
			if err != nil {
				log.Fatalf("get migration files: %+v", err)
			}

			dsnStr = string(config.SQLDriverSQLite) + "://" + cfg.GetDataSourceURL()
		} else {
			cmd.Println("Using MySQL database for migrations")

			data, err = iofs.New(migrations.MySQLFiles, builder.MYSQL)
			if err != nil {
				log.Fatalf("get migration files: %+v", err)
			}

			dsnStr = builder.MYSQL + "://" + cfg.GetDataSourceURL()
		}

		m, err := migrate.NewWithSourceInstance(
			"iofs", // sourceName
			data,   // sourceInstance
			dsnStr, // database URL
		)
		if err != nil {
			log.Fatalf("failed to create migration instance: %+v", err)
		}

		version, dirty, _ := m.Version()
		cmd.Printf("Previous version: %d. With errors?: %+v\n", version, dirty)

		if down {
			if err = m.Steps(-1); err != nil {
				if errors.Is(err, migrate.ErrNoChange) {
					cmd.Println("No new migrations to apply")
				} else if errors.Is(err, migrate.ErrLocked) {
					cmd.Println("Migration is locked")
				} else {
					log.Fatalf("failed to apply migration: %+v", err)
				}
			} else {
				cmd.Println("Migrations applied successfully")
			}
		} else {
			if err = m.Up(); err != nil {
				if errors.Is(err, migrate.ErrNoChange) {
					cmd.Println("No new migrations to apply")
				} else if errors.Is(err, migrate.ErrLocked) {
					cmd.Println("Migration is locked")
				} else {
					log.Fatalf("failed to apply migration: %+v", err)
				}
			} else {
				cmd.Println("Migrations applied successfully")
			}
		}
	},
}
