package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/utking/spaces/internal/infra/state"
)

var webMenu *WebMenu

const (
	labelTypeLink = "link"
	labelDivider  = "-"
)

// RegisterRoutes registers all the routes for the web application.
// Exposed to be used when creating the application.
func RegisterRoutes(e *echo.Echo, state *state.State) {
	webMenu = new(WebMenu)
	webMenu.UserItems = make([]map[string][]WebMenuItem, 0)

	// Home
	e.GET("/", getNotesWrapper(state.Notes, state.Users, state.LastOpened))
	// Ping
	e.GET("/ping", func(c echo.Context) error { return c.String(http.StatusOK, "pong") })

	setProfileRouting(e, state)
	setSelfRegisterRouting(e, state)
	setImportRouting(e, state)
	setNotesRouting(e, state)
	setSecretsRouting(e, state)
	setBookmarksRouting(e, state)
	setUsersRouting(e, state)
	setFilebrowserRouting(e, state)

	setUpMenu(webMenu)
}

func setBookmarksRouting(
	e *echo.Echo,
	state *state.State,
) {
	e.GET("/bookmarks", getBookmarksWrapper(state.Bookmarks, state.Users, state.LastOpened))
	e.POST("/bookmark/create", postBookmarkCreateWrapper(state.Bookmarks, state.Users))
	e.GET("/bookmark/:id/edit", getBookmarkEditWrapper(state.Bookmarks, state.Users))
	e.PUT("/bookmark/:id/edit", putBookmarkEditWrapper(state.Bookmarks, state.Users))
	e.DELETE("/bookmark/:id", deleteBookmarkWrapper(state.Bookmarks, state.Users))
	e.GET("/export/bookmarks", getExportBookmarksWrapper(state.Bookmarks, state.Users))
	e.GET("/search/bookmarks", getSearchBookmarksWrapper(state.Bookmarks, state.Users))
}

func setNotesRouting(
	e *echo.Echo,
	state *state.State,
) {
	e.GET("/notes", getNotesWrapper(state.Notes, state.Users, state.LastOpened))
	e.GET("/note/create", getNoteCreateWrapper(state.Notes, state.Users))
	e.POST("/note/create", postNoteCreateWrapper(state.Notes, state.Users))
	e.PUT("/notes", putNotesWrapper(state.Notes, state.Users))
	e.DELETE("/note/:id", deleteNoteWrapper(state.Notes, state.Users))
	e.GET("/export/notes", getExportNotesWrapper(state.Notes, state.Users))
	e.GET("/search/notes", getSearchNotesWrapper(state.Notes, state.Users))
}

func setSecretsRouting(
	e *echo.Echo,
	state *state.State,
) {
	e.GET("/secrets", getSecretsWrapper(state.Secrets, state.Users))
	e.GET("/secret/create", getSecretCreateWrapper(state.Secrets, state.Users))
	e.POST("/secret/create", postSecretCreateWrapper(state.Secrets, state.Users))
	e.PUT("/secrets", putSecretUpdateWrapper(state.Secrets, state.Users))
	e.DELETE("/secret/:id", deleteSecretWrapper(state.Secrets, state.Users))
	e.GET("/export/secrets", getExportSecretsWrapper())
	e.POST("/export/secrets", postExportSecretsWrapper(state.Secrets, state.Users, state.Secrets))
	e.GET("/search/secrets", getSearchSecretsWrapper(state.Secrets, state.Users))
	e.GET("/secrets/rotate-key", getSecretsRotateKeyWrapper())
	e.POST("/secrets/rotate-key", postSecretsRotateKeyWrapper(state.Secrets, state.Users, state.Logger))
}

func setUsersRouting(
	e *echo.Echo,
	state *state.State,
) {
	e.GET("/users", getUsersWrapper(state.Users))
	e.GET("/user/:id", getUserWrapper(state.Users))
	e.GET("/user/create", getUserCreateWrapper())
	e.POST("/user/create", postUserCreateWrapper(state.Users, state.Mailer, state.Logger, state.Config))
	e.GET("/user/:id/edit", getUserEditWrapper(state.Users))
	e.PUT("/user/:id/edit", postUserEditWrapper(state.Users))
	e.PUT("/users/settings", putUserSettingsWrapper(state.Users))
	e.DELETE("/user/:id", deleteUserWrapper(state.Users))
	e.GET("/verify-user", getUserVerifyWrapper(state.Users, state.Logger, state.Config))
}

func setImportRouting(
	e *echo.Echo,
	state *state.State,
) {
	e.GET("/import/notes", getImportNotesWrapper())
	e.POST("/import/notes", postImportNotesWrapper(state.Notes, state.Users))
	e.GET("/import/bookmarks", getImportBookmarksWrapper())
	e.POST("/import/bookmarks", postImportBookmarksWrapper(state.Bookmarks, state.Users))
	e.GET("/import/secrets", getImportSecretsWrapper())
	e.POST("/import/secrets", postImportSecretsWrapper(state.Secrets, state.Users))
}

func setFilebrowserRouting(
	e *echo.Echo,
	state *state.State,
) {
	e.GET("/filebrowser", getFileBrowserWrapper(state.FileBrowser, state.Users))
	e.POST("/filebrowser/upload", postFileBrowserUploadWrapper(state.FileBrowser, state.Users))
	e.POST("/filebrowser/folder", postFileBrowserNewFolderWrapper(state.FileBrowser, state.Users))
	e.GET("/filebrowser/view", getFileBrowserFileViewWrapper(state.FileBrowser, state.Users))
	e.GET("/filebrowser/download", getFileBrowserFileDownloadWrapper(state.FileBrowser, state.Users))
	e.POST("/filebrowser/rename", postFileBrowserFileRenameWrapper(state.FileBrowser, state.Users))
	e.DELETE("/filebrowser/delete", deleteFileBrowserFileWrapper(state.FileBrowser, state.Users))
	e.POST("/filebrowser/mode", postFileBrowserSetViewModeWrapper(state.Users))
}

func setSelfRegisterRouting(
	e *echo.Echo,
	state *state.State,
) {
	if state.Config.SelfRegistrationEnabled() {
		// Register
		e.Match([]string{"GET", "POST"}, "/register",
			getRegisterWrapper(state.Users, state.Mailer, state.Logger, state.Config))
		e.GET("/register-success", getRegisterSuccessWrapper(state.Config))
	}
}

func setUpMenu(webMenu *WebMenu) {
	webMenu.SimpleItems = append(
		webMenu.SimpleItems,
		WebMenuItem{Type: labelTypeLink, Title: "Notes", URIPath: "/notes"},
		WebMenuItem{Type: labelTypeLink, Title: "Bookmarks", URIPath: "/bookmarks"},
		WebMenuItem{Type: labelTypeLink, Title: "Secrets", URIPath: "/secrets"},
		WebMenuItem{Type: labelTypeLink, Title: "Files", URIPath: "/filebrowser"},
	)

	webMenu.AdminItems = append(webMenu.AdminItems, map[string][]WebMenuItem{
		"Admin": {
			{Type: labelTypeLink, Title: "Users", URIPath: "/users"},
			{Type: labelTypeLink, Title: labelDivider},
			{Type: labelTypeLink, Title: "System Stats", URIPath: "/system-stats"},
		},
	})

	webMenu.AccountItems = append(webMenu.AccountItems, map[string][]WebMenuItem{
		"Account": {
			{Type: labelTypeLink, Title: "Profile", URIPath: "/profile"},
			{Type: labelTypeLink, Title: "Change Password", URIPath: "/change-password"},
			{Type: labelTypeLink, Title: labelDivider},
			{Type: labelTypeLink, Title: "Import Notes", URIPath: "/import/notes"},
			{Type: labelTypeLink, Title: "Import Bookmarks", URIPath: "/import/bookmarks"},
			{Type: labelTypeLink, Title: "Import Secrets", URIPath: "/import/secrets"},
			{Type: labelTypeLink, Title: labelDivider},
			{Type: labelTypeLink, Title: "Rotate Encryption Key", URIPath: "/secrets/rotate-key"},
			{Type: labelTypeLink, Title: labelDivider},
			{Type: labelTypeLink, Title: "Logout", URIPath: "/logout"},
		},
	})
}

// setProfileRouting sets up the routing for user profile related actions.
func setProfileRouting(
	e *echo.Echo,
	state *state.State,
) {
	// Login and Logout
	e.Match([]string{"GET", "POST"}, "/login", getLoginWrapper(state.Users, state.Logger, state.Config))
	e.GET("/logout", getLogoutWrapper(state.Logger))

	// Profile
	e.GET("/profile", getProfileWrapper(state.Users, state.Notes, state.Secrets, state.Bookmarks))
	e.GET("/system-stats", getSystemStatsWrapper(state.SysStats, state.Users))
	e.GET("/secret-generator", getPasswordGeneratorWrapper())
	e.GET("/change-password", getChangePasswordWrapper())
	e.POST("/change-password", postChangePasswordWrapper(state.Users, state.Logger))
}
