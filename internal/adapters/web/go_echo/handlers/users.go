package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/utking/spaces/internal/adapters/notification/mailer"
	"github.com/utking/spaces/internal/adapters/web/go_echo/helpers"
	"github.com/utking/spaces/internal/application/domain"
	"github.com/utking/spaces/internal/config"
	"github.com/utking/spaces/internal/infra/session"
	"github.com/utking/spaces/internal/ports"
)

func getUsersWrapper(
	api ports.UsersService,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			statuses = domain.UserStatus{}.GetStatusMap()
			query    = new(domain.UserRequest)
			code     = http.StatusOK
		)

		_ = c.Bind(query)
		// trim all strings in the query
		query.Trim()

		withStatus := c.QueryParam("status") != ""
		if !withStatus {
			query.Status = nil
		}

		items, err := api.GetItems(c.Request().Context(), query)
		if err != nil {
			code = http.StatusInternalServerError
		}

		return c.Render(
			code,
			"users/index.html",
			map[string]interface{}{
				"Title":    "Users",
				"Items":    items,
				"Statuses": statuses,
				"Error":    helpers.ErrorMessage(err),
				"Query":    query,
				"Count":    len(items),
			},
		)
	}
}

func getUserWrapper(
	api ports.UsersService,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			code     = http.StatusOK
			statuses = domain.UserStatus{}.GetStatusMap()
			userID   = helpers.GetIDParam(c)
			diskUse  = int64(0)
		)

		item, err := api.GetItem(c.Request().Context(), userID)
		if err != nil {
			code = http.StatusInternalServerError
		} else {
			// Get the disk usage for the user
			diskUse, _ = api.GetDiskUsage(c.Request().Context(), userID)
		}

		return c.Render(
			code,
			"users/view.html",
			map[string]interface{}{
				"Title":    "User Details",
				"Item":     item,
				"Statuses": statuses,
				"Error":    helpers.ErrorMessage(err),
				"DiskUse":  diskUse,
			})
	}
}

func getUserCreateWrapper() echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			code     = http.StatusOK
			statuses = domain.UserStatus{}.GetStatusMap()
		)

		return c.Render(
			code,
			"users/create.html",
			map[string]interface{}{
				"Title":    "Create User",
				"Statuses": statuses,
			})
	}
}

func postUserCreateWrapper(
	api ports.UsersService,
	mailerService ports.NotificationService,
	logger ports.LoggingService,
	cfg *config.Config,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			id    string
			query = new(domain.User)
			token string
		)

		_ = c.Bind(query)
		// trim all strings in the query
		query.Trim()

		err := query.Validate()
		if err == nil {
			id, token, err = api.Create(c.Request().Context(), query)
		}

		if err != nil {
			statuses := domain.UserStatus{}.GetStatusMap()

			return c.Render(
				http.StatusInternalServerError,
				"users/create.html",
				map[string]interface{}{
					"Title":    "New User",
					"Query":    query,
					"Statuses": statuses,
					"Error":    helpers.ErrorMessage(err),
				})
		}

		if query.SendNotification {
			// Send a welcome email
			messageBody, renderErr := mailer.RenderTemplate(
				c.Request().Context(),
				"welcome.html",
				map[string]interface{}{
					"Username":         query.Username,
					"AppName":          cfg.GetAppName(),
					"VerificationLink": cfg.GetEmailVerificationLink() + token,
				},
			)

			if renderErr != nil {
				logger.Error(
					c.Request().Context(),
					"Failed to render welcome email template",
					ports.NewLoggerBag(
						"error", helpers.ErrorMessage(renderErr),
					),
					ports.NewLoggerBag(
						"username", query.Username,
					),
				)
			} else if mailerService != nil {
				message := &domain.Notification{
					To:      query.Email,
					Title:   fmt.Sprintf("Welcome to %s", cfg.GetAppName()),
					Message: messageBody,
				}

				if err = mailerService.Send(c.Request().Context(), message); err != nil {
					logger.Error(
						c.Request().Context(),
						"Failed to send welcome email",
						ports.NewLoggerBag(
							"error", helpers.ErrorMessage(err),
						),
						ports.NewLoggerBag(
							"username", query.Username,
						),
					)
				}
			}
		}

		return c.Redirect(
			http.StatusSeeOther,
			fmt.Sprintf("/user/%s", id),
		)
	}
}

func getUserEditWrapper(
	api ports.UsersService,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			code     = http.StatusOK
			statuses = domain.UserStatus{}.GetStatusMap()
		)

		item, err := api.GetItem(c.Request().Context(), helpers.GetIDParam(c))
		if err != nil {
			code = http.StatusInternalServerError
		}

		return c.Render(
			code,
			"users/edit.html",
			map[string]interface{}{
				"Title":    "Edit User",
				"Item":     item,
				"Statuses": statuses,
				"Roles":    domain.UserRolesList(),
				"Error":    helpers.ErrorMessage(err),
			})
	}
}

func postUserEditWrapper(
	api ports.UsersService,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			code   = http.StatusOK
			id     = helpers.GetIDParam(c)
			query  = new(domain.UserUpdate)
			userID = GetUserID(c, api)
		)

		if userID == id {
			return c.JSON(
				http.StatusBadRequest,
				map[string]interface{}{
					"ID":    id,
					"Error": "You cannot edit your own user profile",
				})
		}

		err := c.Bind(query)
		if err != nil {
			code = http.StatusBadRequest

			return c.JSON(
				code,
				map[string]interface{}{
					"ID":    id,
					"Error": helpers.ErrorMessage(err),
				})
		}

		// trim all strings in the query and validate
		query.Trim()

		if err = query.Validate(); err != nil {
			code = http.StatusBadRequest

			return c.JSON(
				code,
				map[string]interface{}{
					"ID":    id,
					"Error": helpers.ErrorMessage(err),
				})
		}

		_, err = api.Update(c.Request().Context(), id, query)
		if err != nil {
			code = http.StatusInternalServerError
		}

		return c.JSON(
			code,
			map[string]interface{}{
				"ID":    id,
				"Error": helpers.ErrorMessage(err),
			})
	}
}

func deleteUserWrapper(
	api ports.UsersService,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			code = http.StatusOK
			id   = helpers.GetIDParam(c)
		)

		err := api.Delete(c.Request().Context(), id)
		if err != nil {
			code = http.StatusInternalServerError
		}

		return c.JSON(
			code,
			map[string]interface{}{
				"Error": helpers.ErrorMessage(err),
			})
	}
}

// getUserVerifyWrapper returns a handler function that verifies a user by a given verification token.
func getUserVerifyWrapper(
	api ports.UsersService,
	logger ports.LoggingService,
	cfg *config.Config,
) echo.HandlerFunc {
	sessTTLSec := cfg.GetSessionTTL()

	return func(c echo.Context) error {
		var (
			err  error
			user *domain.User
		)

		token := c.QueryParam("token")
		if token == "" {
			err = errors.New("verification token is required")
		} else {
			user, err = api.VerifyUser(c.Request().Context(), token)
		}

		if err != nil {
			logger.Error(
				c.Request().Context(),
				"Failed to verify user",
				ports.NewLoggerBag(
					"error", helpers.ErrorMessage(err),
				),
				ports.NewLoggerBag(
					"token", token,
				),
			)

			return c.Render(
				http.StatusInternalServerError,
				"users/verify-error.html",
				map[string]interface{}{
					"Title": "Verification Error",
					"Error": helpers.ErrorMessage(err),
				},
			)
		}

		_ = session.StartSession(c, user.Username, "user", sessTTLSec, cfg.GetWithTLS())

		return c.Redirect(http.StatusSeeOther, "/") // Redirect to the home page after verification
	}
}
