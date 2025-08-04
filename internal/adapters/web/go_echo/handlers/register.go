package handlers

import (
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

// getRegisterWrapper returns a handler function for the registration page.
func getRegisterWrapper(
	api ports.UsersService,
	mailerAPI ports.NotificationService,
	logger ports.LoggingService,
	cfg *config.Config,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			err      error
			username string
			code     = http.StatusOK
		)

		// Check if the user is already logged in
		if username, _ = session.GetSessionUsername(c); username != "" {
			return c.Redirect(http.StatusSeeOther, "/")
		}

		if !cfg.SelfRegistrationEnabled() {
			// If self-registration is disabled, redirect to the login page
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		// Process the registration form submission
		if c.Request().Method == http.MethodPost {
			if err = processRegisterPost(c, api, mailerAPI, logger, cfg); err != nil {
				code = http.StatusBadRequest
			} else {
				return c.Redirect(http.StatusSeeOther, "register-success")
			}
		}

		// Render the registration page
		return c.Render(
			code,
			"register.html",
			map[string]interface{}{
				"Title": "Sign Up",
				"Error": helpers.ErrorMessage(err),
			},
		)
	}
}

// processRegisterPost processes the registration form submission.
func processRegisterPost(
	c echo.Context,
	api ports.UsersService,
	mailerAPI ports.NotificationService,
	logger ports.LoggingService,
	cfg *config.Config,
) error {
	var (
		query = new(domain.User)
		token string
	)

	_ = c.Bind(query)
	// trim all strings in the query
	query.Normalize()
	query.SendNotification = true // self-registration always sends a notification

	err := query.Validate()
	if err == nil {
		_, token, err = api.Create(c.Request().Context(), query)
	}

	if err != nil {
		return err
	}

	if query.SendNotification {
		// Send a welcome email
		messageBody, renderErr := mailer.RenderTemplate(
			c.Request().Context(),
			"welcome.html",
			map[string]interface{}{
				"Username": query.Username,
				"AppName":  cfg.GetAppName(),
				"VerificationLink": fmt.Sprintf(
					"%s%s",
					cfg.GetEmailVerificationLink(),
					token,
				),
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
		} else if mailerAPI != nil {
			message := &domain.Notification{
				To:      query.Email,
				Title:   fmt.Sprintf("Welcome to %s", cfg.GetAppName()),
				Message: messageBody,
			}

			if err = mailerAPI.Send(c.Request().Context(), message); err != nil {
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

	return err
}

// getRegisterSuccessWrapper returns a handler function for the registration success page.
func getRegisterSuccessWrapper(
	cfg *config.Config,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !cfg.SelfRegistrationEnabled() {
			// If self-registration is disabled, redirect to the index page
			return c.Redirect(http.StatusSeeOther, "/")
		}

		username, _ := session.GetSessionUsername(c)
		if username != "" {
			// If the user is already logged in, redirect to the index page
			return c.Redirect(http.StatusSeeOther, "/")
		}

		// Render the registration success page
		return c.Render(
			http.StatusOK,
			"register_success.html",
			map[string]interface{}{})
	}
}
