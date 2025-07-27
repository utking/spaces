// Package session provides session management for the application
// using the gorilla/sessions package.
package session

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/srinathgs/mysqlstore"
	"github.com/utking/spaces/internal/config"
)

func NewStore(cfg *config.Config) (sessions.Store, error) {
	secret, key := cfg.GetSessionSecretAndKey()

	return mysqlstore.NewMySQLStore(
		cfg.GetDataSourceURL(), "sessions",
		"/", cfg.GetSessionTTL(),
		[]byte(secret), []byte(key),
	)
}

// StartSession initializes a new session for the user.
func StartSession(
	c echo.Context,
	user string,
	role string,
	ttlSec int,
	useSecure bool,
) error {
	session, err := session.Get("session", c)
	if err != nil && session == nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	mode := http.SameSiteLaxMode
	if useSecure {
		mode = http.SameSiteStrictMode
	}

	session.Values["username"] = user
	session.Values["role"] = role
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   ttlSec,
		HttpOnly: true,
		Secure:   useSecure,
		SameSite: mode,
	}

	if err = session.Save(c.Request(), c.Response()); err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	return nil
}

// TerminateSession clears the session for the user
// and removes the username from the session values
// and sets the MaxAge to -1 to delete the session cookie
// and remove the session from the server
// and the client.
func TerminateSession(
	c echo.Context,
) error {
	session, err := session.Get("session", c)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	session.Values["username"] = ""
	session.Values["role"] = ""
	session.Options.MaxAge = -1

	if err = session.Save(c.Request(), c.Response()); err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	return nil
}

// GetSessionUsername retrieves the username from the session
// and returns it as a string.
func GetSessionUsername(
	c echo.Context,
) (string, error) {
	session, err := session.Get("session", c)
	if err != nil {
		return "", fmt.Errorf("failed to get session: %w", err)
	}

	username, ok := session.Values["username"].(string)
	if !ok {
		return "", errors.New("failed to get username from session")
	}

	return username, nil
}

// IsAdminSession checks if the user's role is admin
// and returns true if it is, otherwise false.
func IsAdminSession(
	c echo.Context,
) (bool, error) {
	session, err := session.Get("session", c)
	if err != nil {
		return false, fmt.Errorf("failed to get session: %w", err)
	}

	role, ok := session.Values["role"].(string)
	if !ok || role != "admin" {
		return false, nil
	}

	return true, nil
}

// SetStrVar sets a variable in the session.
func SetStrVar(
	c echo.Context,
	key string,
	value string,
) error {
	session, err := session.Get("session", c)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	session.Values[key] = value

	if err = session.Save(c.Request(), c.Response()); err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	return nil
}

// GetStrVar retrieves a variable from the session.
func GetStrVar(
	c echo.Context,
	key string,
) (string, error) {
	session, err := session.Get("session", c)
	if err != nil {
		return "", fmt.Errorf("failed to get session: %w", err)
	}

	value, ok := session.Values[key].(string)
	if !ok {
		return "", fmt.Errorf("failed to get variable %s from session", key)
	}

	return value, nil
}
