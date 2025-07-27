// Package middleware provides custom middleware for the Echo framework.
package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	"gogs.utking.net/utking/spaces/internal/infra/session"
)

// Custom Authentication middleware

// AuthCheck is a middleware to check if the user is authenticated.
func AuthCheck() echo.MiddlewareFunc {
	return AuthWithConfig(AuthConfig{
		Validator: DefaultValidateUser,
		LoginURL:  "/login",
		Skipper:   func(c echo.Context) bool { return c.Request().URL.Path == "/login" },
	})
}

type (
	// AuthConfig defines the config for Auth middleware.
	AuthConfig struct {
		// Skipper defines a function to skip the middleware.
		Skipper echomw.Skipper

		// Validator is a function to validate user credentials.
		// Required.
		Validator AuthValidator

		// LoginURL must define a login page
		// Required
		LoginURL string
	}

	// AuthValidator defines a function to validate Auth credentials.
	AuthValidator func(string, string, echo.Context) (string, error)
)

var (
	// DefaultAuthConfig is the default Auth middleware config.
	DefaultAuthConfig = AuthConfig{
		Skipper: echomw.DefaultSkipper,
	}
)

// AuthWithConfig returns an Auth middleware with config.
// See `Auth()`.
func AuthWithConfig(config AuthConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Validator == nil {
		panic("echo: auth middleware requires a validator function")
	}

	if config.LoginURL == "" {
		panic("echo: auth middleware requires a login URL")
	}

	if config.Skipper == nil {
		config.Skipper = DefaultAuthConfig.Skipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skip what is defined in the config
			if config.Skipper(c) {
				return next(c)
			}

			// Check if the user is already authenticated
			if username, sErr := session.GetSessionUsername(c); sErr == nil && strings.TrimSpace(username) != "" {
				return next(c)
			}

			// if application/json, return 401
			if c.Request().Header.Get(echo.HeaderContentType) == echo.MIMEApplicationJSON {
				return c.JSON(http.StatusUnauthorized, map[string]string{})
			}

			// Redirect to the login page
			return c.Redirect(http.StatusSeeOther, config.LoginURL)
		}
	}
}

// DefaultValidateUser is a placeholder function for validating user credentials.
func DefaultValidateUser(_, _ string, _ echo.Context) (string, error) {
	return "", errors.New("not implemented")
}
