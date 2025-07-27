// Package middleware provides custom middleware for the Echo framework.
package middleware

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
)

// Custom AdminAccess middleware

// AdminAccess is a middleware to check if the user is allowed to access a resource.
func AdminAccess() echo.MiddlewareFunc {
	return AdminAccessWithConfig(AdminAccessConfig{
		Validator: DefaultCheck,
		Skipper:   func(_ echo.Context) bool { return false },
	})
}

type (
	// AdminAccessConfig defines the config for AdminAccess middleware.
	AdminAccessConfig struct {
		// Skipper defines a function to skip the middleware.
		Skipper echomw.Skipper

		// Validator is a function to validate user credentials.
		// Required.
		Validator AdminAccessValidator
	}

	// AdminAccessValidator defines a function to validate the user's access.
	// The user's role is stored in the session.
	AdminAccessValidator func(ctx echo.Context) error
)

var (
	// DefaultAdminAccessConfig is the default AdminAccess middleware config.
	DefaultAdminAccessConfig = AdminAccessConfig{
		Skipper: echomw.DefaultSkipper,
	}
)

// AdminAccessWithConfig returns the middleware built with the given config.
func AdminAccessWithConfig(config AdminAccessConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Validator == nil {
		panic("echo: AdminAccess middleware requires a validator function")
	}

	if config.Skipper == nil {
		config.Skipper = DefaultAdminAccessConfig.Skipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skip what is defined in the config
			if config.Skipper(c) {
				return next(c)
			}

			// Check the user's role
			if err := config.Validator(c); err == nil {
				return next(c)
			}

			// Return 403 Forbidden if the user is not an admin
			return echo.NewHTTPError(http.StatusForbidden, "You are not allowed to access this resource")
		}
	}
}

// DefaultCheck is a placeholder function.
func DefaultCheck(_ echo.Context) error {
	return errors.New("not implemented")
}
