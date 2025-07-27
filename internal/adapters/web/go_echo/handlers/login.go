package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/utking/spaces/internal/adapters/web/go_echo/helpers"
	"github.com/utking/spaces/internal/config"
	"github.com/utking/spaces/internal/infra/session"
	log_port "github.com/utking/spaces/internal/ports"
)

const loginTemplate = "login.html"

func getLoginWrapper(
	api log_port.UsersService,
	logger log_port.LoggingService,
	cfg *config.Config,
) echo.HandlerFunc {
	sessTTLSec := cfg.GetSessionTTL()

	return func(c echo.Context) error {
		var (
			err      error
			username string
			code     = http.StatusOK
		)

		// Check if the user is already logged in
		if username, err = session.GetSessionUsername(c); err == nil && username != "" {
			return c.Redirect(http.StatusSeeOther, "/")
		}

		err = nil // Reset error to nil for the login form

		// Process the login form submission
		if c.Request().Method == http.MethodPost {
			err = processLoginPost(c, api, logger, sessTTLSec, cfg)
			if err != nil {
				code = http.StatusUnauthorized
			} else {
				return nil
			}
		}

		// Render the login page
		return c.Render(
			code,
			loginTemplate,
			map[string]interface{}{
				"Title":            "Login",
				"Error":            helpers.ErrorMessage(err),
				"SelfRegistration": cfg.SelfRegistrationEnabled(),
			},
		)
	}
}

func processLoginPost(
	c echo.Context,
	api log_port.UsersService,
	logger log_port.LoggingService,
	sessTTLSec int,
	cfg *config.Config,
) error {
	params, errParams := c.FormParams()
	if errParams != nil {
		return errors.New("invalid login parameters")
	}

	username, password :=
		strings.TrimSpace(params.Get("username")),
		strings.TrimSpace(params.Get("password"))

	if username == "" || password == "" {
		return errors.New("username and password are required")
	}

	if userRole, err := api.ValidateUser(c.Request().Context(), username, password); err == nil {
		if err = session.StartSession(c, username, userRole, sessTTLSec, cfg.GetWithTLS()); err != nil {
			logger.Error(
				c.Request().Context(),
				fmt.Sprintf("failed to start session: %v", err),
			)
		}

		// check if user profile directory exists and create it if not
		if err = api.CreateDataDirectory(c.Request().Context(), GetUserID(c, api)); err != nil {
			logger.Error(
				c.Request().Context(),
				fmt.Sprintf("failed to create user profile directory: %v", err),
			)
		}

		return c.Redirect(http.StatusSeeOther, "/")
	}

	logger.Warn(c.Request().Context(),
		"failed login attempt",
		log_port.NewLoggerBag("username", username),
		log_port.NewLoggerBag("ip", c.RealIP()),
		log_port.NewLoggerBag("user-agent", c.Request().UserAgent()),
	)

	return errors.New("invalid username or password")
}

func getLogoutWrapper(
	logger log_port.LoggingService,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := session.TerminateSession(c); err != nil {
			logger.Error(
				c.Request().Context(),
				fmt.Sprintf("failed to terminate session: %v", err),
			)
		}

		return c.Redirect(http.StatusSeeOther, "/login")
	}
}
