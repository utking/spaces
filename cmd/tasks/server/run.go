// Package server provides the command to start the HTTP server.
package server

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/utking/spaces/internal/adapters/cryptor"
	db_mysql "github.com/utking/spaces/internal/adapters/db/mysql"
	db_sqlite "github.com/utking/spaces/internal/adapters/db/sqlite"
	"github.com/utking/spaces/internal/adapters/filesystem"
	"github.com/utking/spaces/internal/adapters/logger"
	"github.com/utking/spaces/internal/adapters/notification/mailer"
	web "github.com/utking/spaces/internal/adapters/web/go_echo"
	"github.com/utking/spaces/internal/application/services"
	"github.com/utking/spaces/internal/config"
	"github.com/utking/spaces/internal/infra/state"
	"github.com/utking/spaces/internal/ports"
	"xorm.io/builder"
)

// Init initializes the server command and adds it to the root command.
func Init(rootCmd *cobra.Command) {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start HTTP server",
	Run: func(_ *cobra.Command, _ []string) {
		var (
			dbAdapter ports.DBPort
			err       error
		)

		cfg := config.New()

		aesCryptor := cryptor.New()
		fsAdapter := filesystem.NewAdapter(cfg.GetDataBasePath())

		if cfg.GetSQLDriver() == builder.MYSQL {
			log.Println("Using MySQL database")
			dbAdapter, err = db_mysql.NewAdapter(cfg.GetDataSourceURL())
		} else {
			log.Println("Using SQLite database")
			dbAdapter, err = db_sqlite.NewAdapter(cfg.GetDataSourceURL())
		}

		if err != nil {
			log.Fatalf("failed to connect to DB: %+v", err)
		}

		mailerAdapter := mailer.New(
			cfg.GetSMTPHost(),
			cfg.GetSMTPPort(),
			cfg.GetSMTPUsername(),
			cfg.GetSMTPPassword(),
			cfg.GetSMTPFrom(),
			cfg.GetSMTPUseTLS(),
		)

		notesService := services.NewNotesService(dbAdapter)
		usersService := services.NewUsersService(dbAdapter, fsAdapter)
		sysStatsService := services.NewSysStatService(dbAdapter)
		secretsService := services.NewSecretService(dbAdapter, aesCryptor)
		bookmarkService := services.NewBookmarkService(dbAdapter)
		lastOpenedService := services.NewLastOpenedService(dbAdapter)
		fileBrowser := filesystem.NewFileBrowserAdapter(cfg.GetDataBasePath())

		// App Logs Logger
		logFile, logFileErr := os.OpenFile(
			cfg.GetAppLogFilePath(),
			os.O_RDWR|os.O_CREATE|os.O_APPEND,
			0o600,
		)

		if logFileErr != nil {
			log.Fatalf("Failed to open log file: %v", logFileErr)
		}
		defer logFile.Close()

		logAdapter := logger.NewAdapter(
			logFile,
			cfg.GetLogLevel(),
		)

		// state with all services
		state := state.New(
			cfg,               /* Config */
			logAdapter,        /* LoggingService */
			usersService,      /* UsersService */
			sysStatsService,   /* SysStatService */
			notesService,      /* NotesService */
			secretsService,    /* SecretService */
			mailerAdapter,     /* NotificationService */
			bookmarkService,   /* BookmarkService */
			lastOpenedService, /* LastOpenedService */
			fileBrowser,       /* FileBrowserService */
		)

		httpAdapter := web.NewAdapter(uint(cfg.GetApplicationPort()), state)

		httpAdapter.Run()
	},
}
