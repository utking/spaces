// Package handlers provides HTTP handlers for the Echo web framework.
package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/utking/spaces/internal/infra/session"
	"github.com/utking/spaces/internal/ports"
)

// GetMenu returns the web menu tree.
// It is used in the layout template to render the menu.
func GetMenu() WebMenu {
	if webMenu == nil {
		return WebMenu{}
	}

	return *webMenu
}

// GetUserID returns the user ID for the username from the session.
//   - If the username is not found, 0 is returned.
func GetUserID(c echo.Context, userAPI ports.UsersService) (userID string) {
	username, err := session.GetSessionUsername(c)
	if err != nil {
		return ""
	}

	user, err := userAPI.GetByUsername(c.Request().Context(), username)
	if err != nil {
		return ""
	}

	return user.ID
}
