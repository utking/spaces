package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"gogs.utking.net/utking/spaces/internal/adapters/web/go_echo/helpers"
	"gogs.utking.net/utking/spaces/internal/infra/session"
	"gogs.utking.net/utking/spaces/internal/ports"
)

// getProfileWrapper returns a handler function for the profile page.
func getProfileWrapper(
	userAPI ports.UsersService,
	notes ports.NotesService,
	secrets ports.SecretService,
	bookmarks ports.BookmarkService,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			code        = http.StatusOK
			userID      = GetUserID(c, userAPI)
			username, _ = session.GetSessionUsername(c)
		)

		diskUse, err := userAPI.GetDiskUsage(c.Request().Context(), userID)
		if err != nil {
			code = http.StatusInternalServerError
			err = errors.New("failed to get disk usage for user")
		}

		notesCount, ncErr := notes.GetCount(c.Request().Context(), userID, nil)
		noteTags, ntcErr := notes.GetTags(c.Request().Context(), userID)
		secretsCount, scErr := secrets.GetCount(c.Request().Context(), userID, nil)
		secretTags, stcErr := secrets.GetTags(c.Request().Context(), userID)
		bookmarksCount, bcErr := bookmarks.GetCount(c.Request().Context(), userID, nil)
		bookmarkTags, _ := bookmarks.GetTags(c.Request().Context(), userID)

		err = errors.Join(err, ncErr, ntcErr, scErr, stcErr, bcErr)

		return c.Render(
			code,
			"users/profile.html",
			map[string]interface{}{
				"Title":             "User Profile",
				"Error":             helpers.ErrorMessage(err),
				"Username":          username,
				"DiskUse":           diskUse,
				"BookmarksCount":    bookmarksCount,
				"BookmarkTagsCount": len(bookmarkTags),
				"NotesCount":        notesCount,
				"NoteTagsCount":     len(noteTags),
				"SecretsCount":      secretsCount,
				"SecretTagsCount":   len(secretTags),
			},
		)
	}
}
