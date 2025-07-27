// Package goecho implements the error handler for the Echo framework.
package goecho

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	echo_log "github.com/labstack/gommon/log"
	"golang.org/x/time/rate"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gogs.utking.net/utking/spaces/internal/adapters/web/go_echo/handlers"
	"gogs.utking.net/utking/spaces/internal/adapters/web/go_echo/helpers"
	auth_ms "gogs.utking.net/utking/spaces/internal/adapters/web/go_echo/middleware"
	"gogs.utking.net/utking/spaces/internal/config"
	sess "gogs.utking.net/utking/spaces/internal/infra/session"
	"gogs.utking.net/utking/spaces/internal/infra/state"
	"gogs.utking.net/utking/spaces/static"
)

// Adapter is the struct implementing the web server using the Echo framework.
type Adapter struct {
	state *state.State
	port  uint
}

// NewAdapter creates a new instance of the Adapter struct.
func NewAdapter(port uint, state *state.State) *Adapter {
	return &Adapter{
		port:  port,
		state: state,
	}
}

// Run starts the web server and listens for incoming requests.
func (a *Adapter) Run() {
	e := echo.New()

	e.Pre(middleware.MethodOverrideWithConfig(middleware.MethodOverrideConfig{
		Getter: middleware.MethodFromForm("_method"),
	}))

	if err := InitTemplates(e, a.state.Logger, a.state.Config); err != nil {
		e.Logger.Fatalf("Failed to initialize templates: %w", err)
	}

	e.StaticFS("/assets", static.StaticFiles)
	e.Static("/uploads", "uploads")

	e.HideBanner = true
	e.HTTPErrorHandler = HTTPErrorHandler

	if a.state.Config.GetLogLevel() == "DEBUG" {
		e.Logger.SetLevel(echo_log.DEBUG)
	}

	setExtendedPanicHandler(e, a.state.Config)
	e.Use(middleware.Gzip())
	e.Use(middleware.CORS())
	setRateLimitConfig(e)

	// Access Logs Logger
	logFile, logFileErr := createLogFile(a.state.Config)
	if logFileErr != nil {
		e.Logger.Fatalf("Failed to open log file: %w", logFileErr)
	}

	defer logFile.Close()

	fileLogger := middleware.LoggerWithConfig(middleware.LoggerConfig{Output: logFile})
	e.Use(fileLogger)

	csrfConfig := middleware.CSRFConfig{
		TokenLookup:    "cookie:_csrf",
		CookiePath:     "/",
		CookieHTTPOnly: true,
		CookieSameSite: http.SameSiteDefaultMode,
	}

	// Set secure cookie if using TLS
	csrfConfig.CookieSecure = a.state.Config.GetWithTLS()
	if csrfConfig.CookieSecure {
		csrfConfig.CookieSameSite = http.SameSiteStrictMode
	}

	e.Use(middleware.CSRFWithConfig(csrfConfig))
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 30 * time.Second,
	}))

	e.Pre(middleware.RemoveTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		Skipper: func(c echo.Context) bool {
			return strings.HasPrefix(c.Request().URL.Path, "/uploads/")
		},
	}))

	// Session management
	var (
		store sessions.Store
		err   error
	)

	sessKey, sessSecret := a.state.Config.GetSessionSecretAndKey()
	if a.state.Config.GetSQLDriver() == config.SQLDriverSQLite {
		store, err = sess.NewSqliteStore(
			a.state.Config.GetDataSourceURL(),
			"sessions", "/",
			a.state.Config.GetSessionTTL(),
			[]byte(sessSecret), []byte(sessKey),
		)
	} else {
		store, err = sess.NewStore(a.state.Config)
	}
	if err != nil {
		e.Logger.Fatalf("Failed to create session store: %v", err)
	}

	e.Use(session.MiddlewareWithConfig(session.Config{
		Skipper: middleware.DefaultSkipper,
		Store:   store,
	}))

	setAuthConfig(e, a)
	setAdminAccessConfig(e)
	a.registerRoutes(e, a.state)

	// Start the server - either HTTP or HTTPS
	// depending on the configuration
	if a.state.Config.GetWithTLS() {
		setSecurityConfig(e)
		e.Logger.Fatal(e.StartTLS(
			fmt.Sprintf(":%d", a.port),
			a.state.Config.GetTLSCertFile(),
			a.state.Config.GetTLSKeyFile(),
		))
	} else {
		e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", a.port)))
	}
}

func (a *Adapter) registerRoutes(e *echo.Echo, state *state.State) {
	// Core API
	handlers.RegisterRoutes(e, state)
}

// setSecurityConfig sets the security configuration for the Echo framework.
func setSecurityConfig(e *echo.Echo) {
	e.Use(middleware.SecureWithConfig(
		middleware.SecureConfig{
			Skipper:               middleware.DefaultSkipper,
			XSSProtection:         middleware.DefaultSecureConfig.XSSProtection,
			ContentTypeNosniff:    middleware.DefaultSecureConfig.ContentTypeNosniff,
			XFrameOptions:         middleware.DefaultSecureConfig.XFrameOptions,
			HSTSPreloadEnabled:    false,
			ContentSecurityPolicy: `default-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:;`,
		},
	))
}

// setExtendedPanicHandler sets a custom panic handler for the Echo framework.
func setExtendedPanicHandler(e *echo.Echo, cfg *config.Config) {
	e.Use(middleware.RecoverWithConfig(
		middleware.RecoverConfig{
			StackSize: 4 << 10, // 1 KB
			LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
				c.Logger().Errorf("Recovered from panic: %v", err)
				c.Logger().Errorf("Stack trace: %s", stack)

				if cfg.IsDevMode() && err != nil {
					c.Logger().Debugf("Request details: %+v", c.Request())
				}

				var errStr string

				if cfg.IsDevMode() {
					errStr = helpers.ErrorMessage(err)
				} else {
					errStr = "An unexpected error occurred. Please try again later."
				}

				return c.Render(
					http.StatusInternalServerError,
					"errors/500.html",
					map[string]interface{}{
						"Title": "Internal Server Error",
						"Error": errStr,
					},
				)
			},
		},
	))
}

// setAuthConfig sets the authentication configuration for the Echo framework.
func setAuthConfig(e *echo.Echo, a *Adapter) {
	e.Use(auth_ms.AuthWithConfig(auth_ms.AuthConfig{
		LoginURL: "/login",
		Skipper: func(c echo.Context) bool {
			return c.Request().URL.Path == "/login" ||
				c.Request().URL.Path == "/register" ||
				c.Request().URL.Path == "/register-success" ||
				c.Request().URL.Path == "/ping" ||
				c.Request().URL.Path == "/verify-user" ||
				strings.HasPrefix(c.Request().URL.Path, "/assets")
		},
		Validator: func(username, password string, ctx echo.Context) (string, error) {
			return a.state.Users.ValidateUser(
				ctx.Request().Context(),
				username,
				password,
			)
		},
	}))
}

// setAdminAccessConfig sets the admin access configuration for the Echo framework.
func setAdminAccessConfig(e *echo.Echo) {
	e.Use(auth_ms.AdminAccessWithConfig(auth_ms.AdminAccessConfig{
		Validator: func(c echo.Context) error {
			isAdmin, aErr := sess.IsAdminSession(c)
			if aErr != nil || !isAdmin {
				return errors.New("access denied: user is not an admin")
			}

			return nil
		},
		Skipper: func(c echo.Context) bool {
			return !strings.HasPrefix(c.Request().URL.Path, "/user") &&
				!strings.HasPrefix(c.Request().URL.Path, "/system-stats")
		},
	}))
}

// setRateLimitConfig sets the rate limiting configuration for the Echo framework.
func setRateLimitConfig(e *echo.Echo) {
	e.Use(middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{
				Rate:      rate.Limit(10),
				Burst:     30,
				ExpiresIn: 3 * time.Minute,
			},
		),
		Skipper: func(c echo.Context) bool {
			// Skip rate limiting for specific paths
			return strings.HasPrefix(c.Request().URL.Path, "/uploads/") ||
				strings.HasPrefix(c.Request().URL.Path, "/assets/")
		},
	}))
}

// createLogFile creates a log file with the specified name and returns the file handle.
func createLogFile(cfg *config.Config) (*os.File, error) {
	logFile, logFileErr := os.OpenFile(
		cfg.GetAccessLogFilePath(),
		os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0o600,
	)

	if logFileErr != nil {
		return nil, logFileErr
	}

	return logFile, nil
}
