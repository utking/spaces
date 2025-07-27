// Package state provides the core application state management.
// It holds references to various services that the application uses.
package state

import (
	"gogs.utking.net/utking/spaces/internal/config"
	"gogs.utking.net/utking/spaces/internal/ports"
)

// State represents the core application state.
type State struct {
	Config      *config.Config
	Logger      ports.LoggingService
	Users       ports.UsersService
	SysStats    ports.SystemStatsService
	Notes       ports.NotesService
	Secrets     ports.SecretService
	Bookmarks   ports.BookmarkService
	Mailer      ports.NotificationService
	LastOpened  ports.LastOpenedService
	FileBrowser ports.FileBrowserService
}

// New creates a new instance of the State struct.
func New(
	config *config.Config,
	logger ports.LoggingService,
	users ports.UsersService,
	sysStats ports.SystemStatsService,
	notes ports.NotesService,
	secrets ports.SecretService,
	mailer ports.NotificationService,
	bookmarks ports.BookmarkService,
	lastOpened ports.LastOpenedService,
	fileBrowser ports.FileBrowserService,
) *State {
	return &State{
		Config:      config,
		Logger:      logger,
		Users:       users,
		SysStats:    sysStats,
		Notes:       notes,
		Secrets:     secrets,
		Bookmarks:   bookmarks,
		Mailer:      mailer,
		LastOpened:  lastOpened,
		FileBrowser: fileBrowser,
	}
}
