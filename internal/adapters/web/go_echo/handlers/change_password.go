package handlers

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/utking/spaces/internal/adapters/web/go_echo/helpers"
	"github.com/utking/spaces/internal/ports"
)

// getChangePasswordWrapper returns a handler function for the change password page.
func getChangePasswordWrapper() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(
			http.StatusOK,
			"users/change-password.html",
			map[string]interface{}{
				"Title": "Change Password",
			},
		)
	}
}

// postChangePasswordWrapper returns a handler function for processing the change password form.
func postChangePasswordWrapper(
	api ports.UsersService,
	logger ports.LoggingService,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		var err error

		userID := GetUserID(c, api)
		oldPassword := strings.TrimSpace(c.FormValue("current_password"))
		newPassword := strings.TrimSpace(c.FormValue("new_password"))

		if oldPassword == newPassword {
			return c.Render(
				http.StatusBadRequest,
				"users/change-password.html",
				map[string]interface{}{
					"Title": "Change Password",
					"Error": "new password cannot be the same as the old one",
				})
		}

		user, err := api.GetItem(c.Request().Context(), userID)
		if err != nil || user == nil {
			logger.Error(
				c.Request().Context(),
				"Failed to get user item for password change",
				ports.NewLoggerBag("error", err),
				ports.NewLoggerBag("user_id", userID),
			)

			return c.Render(
				http.StatusInternalServerError,
				"users/change-password.html",
				map[string]interface{}{
					"Title": "Change Password",
					"Error": "Failed to retrieve user information",
				})
		}

		if _, err = api.ValidateUser(c.Request().Context(), user.Username, oldPassword); err != nil {
			return c.Render(
				http.StatusBadRequest,
				"users/change-password.html",
				map[string]interface{}{
					"Title": "Change Password",
					"Error": helpers.ErrorMessage(err),
				})
		}

		if err = api.ChangePassword(c.Request().Context(), userID, newPassword); err != nil {
			return c.Render(
				http.StatusInternalServerError,
				"users/change-password.html",
				map[string]interface{}{
					"Title": "Change Password",
					"Error": helpers.ErrorMessage(err),
				})
		}

		return c.Render(
			http.StatusOK,
			"users/change-password.html",
			map[string]interface{}{
				"Title": "Change Password",
				"Ok":    "Password changed successfully",
			})
	}
}
