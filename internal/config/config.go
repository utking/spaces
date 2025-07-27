// Package config provides configuration management for the application.
package config

import (
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"xorm.io/builder"
)

// EnvValue is a type alias for string to represent environment variable values.
type EnvValue string

// SQLDriver is a type alias for string to represent SQL drivers.
type SQLDriver string

const (
	// UsersPageSize is the page size.
	UsersPageSize = 20
	// NotesPageSize is the page size.
	NotesPageSize = 100
	// SecretsPageSize is the page size for secrets.
	SecretsPageSize = 100
	// SQLDriverMySQL is the MySQL driver.
	SQLDriverMySQL SQLDriver = builder.MYSQL
	// SQLDriverSQLite is the SQLite driver.
	SQLDriverSQLite SQLDriver = "sqlite"
	// trueStr is a string representation of true.
	trueStr = "true"
)

type Config struct{}

// New creates a new instance of Config.
func New() *Config {
	if err := godotenv.Load(".env.dev"); err != nil {
		if err = godotenv.Load(".env"); err != nil {
			log.Println(err)
		}
	}

	return &Config{}
}

// getEnvValue retrieves the value of an environment variable by its name.
// Will log a fatal error and terminate the program if the variable is not set
// or is empty and there are no default values provided.
func getEnvValue(name string, defaultVal ...string) string {
	value := strings.Trim(os.Getenv(name), " ")
	if value == "" {
		if len(defaultVal) > 0 {
			value = defaultVal[0]
		} else {
			log.Fatalf("environment variable %s is not set", name)
		}
	}

	return value
}

// GetDataSourceURL returns the data source URL for the database connection
// based on the environment variable DATA_SOURCE_URL.
// It logs a fatal error and terminates the program if the URL is not set
// or is empty.
func (c *Config) GetDataSourceURL() string {
	return getEnvValue("DATA_SOURCE_URL")
}

// GetApplicationPort returns the port on which the application will listen
// for incoming requests.
func (c *Config) GetApplicationPort() int {
	portVal := getEnvValue("APPLICATION_PORT", "8080")
	port, err := strconv.Atoi(portVal)

	if err != nil {
		log.Fatalf("port %s is invalid", portVal)
	}

	return port
}

// GetSessionSecretAndKey returns session secret and encryption key
// for encrypting and decrypting session data.
func (c *Config) GetSessionSecretAndKey() (string, string) {
	secret := getEnvValue("SESSION_SECRET")
	key := getEnvValue("SESSION_KEY")

	return secret, key
}

// GetWithTLS returns whether to use TLS for the application.
func (c *Config) GetWithTLS() bool {
	withTLS := getEnvValue("USE_TLS", "false")
	if withTLS == trueStr || withTLS == "1" {
		return true
	}

	return false
}

// GetTLSCertFile returns the path to the TLS certificate file.
func (c *Config) GetTLSCertFile() string {
	return getEnvValue("TLS_CERT_FILE")
}

// GetTLSKeyFile returns the path to the TLS key file.
func (c *Config) GetTLSKeyFile() string {
	return getEnvValue("TLS_KEY_FILE")
}

// getLogsDir returns the directory where logs are stored.
func getLogsDir() string {
	return getEnvValue("LOGS_DIR", ".")
}

// GetAccessLogFilePath returns the path to the log file.
func (c *Config) GetAccessLogFilePath() string {
	return path.Join(getLogsDir(), "web.log.json")
}

// GetAppLogFilePath returns the path to the error log file.
func (c *Config) GetAppLogFilePath() string {
	return path.Join(getLogsDir(), "app.log.json")
}

// GetLogLevel returns the log level for the application.
func (c *Config) GetLogLevel() string {
	return strings.ToUpper(getEnvValue("LOG_LEVEL", "DEBUG"))
}

// IsDevMode returns whether the application is running in development mode.
func (c *Config) IsDevMode() bool {
	devMode := strings.ToUpper(getEnvValue("APP_ENV", "PROD"))
	return devMode == "DEV" || devMode == "DEVELOPMENT"
}

// GetDataBasePath returns the path to the user profiles directory.
func (c *Config) GetDataBasePath() string {
	return getEnvValue("DATA_DIR_PATH", "./data")
}

// GetSMTPHost returns the SMTP host for sending emails.
func (c *Config) GetSMTPHost() string {
	return getEnvValue("SMTP_HOST", "localhost")
}

// GetSMTPPort returns the SMTP port for sending emails.
func (c *Config) GetSMTPPort() int32 {
	portVal := getEnvValue("SMTP_PORT", "25")
	port, err := strconv.ParseUint(portVal, 10, 32)

	if err != nil {
		log.Fatalf("port %s is invalid", portVal)
	}

	return int32(port) //nolint:gosec //ParseUint is safe here
}

// GetSMTPUsername returns the SMTP username for sending emails.
func (c *Config) GetSMTPUsername() string {
	return getEnvValue("SMTP_USERNAME", "")
}

// GetSMTPPassword returns the SMTP password for sending emails.
func (c *Config) GetSMTPPassword() string {
	return getEnvValue("SMTP_PASSWORD", "")
}

// GetSMTPFrom returns the email address from which emails will be sent.
func (c *Config) GetSMTPFrom() string {
	return getEnvValue("SMTP_FROM", "")
}

// GetSMTPUseTLS returns whether to use TLS for the SMTP connection.
func (c *Config) GetSMTPUseTLS() bool {
	useTLS := getEnvValue("SMTP_USE_TLS", "false")
	return useTLS == trueStr || useTLS == "1"
}

// GetAppName returns the name of the application.
func (c *Config) GetAppName() string {
	return getEnvValue("APP_NAME", "Spaces")
}

// GetEmailVerificationLink returns the email verification link template.
func (c *Config) GetEmailVerificationLink() string {
	return getEnvValue("EMAIL_VERIFICATION_LINK", "http://localhost/user/verify?token=")
}

// GetSessionTTL returns the session TTL in seconds.
func (c *Config) GetSessionTTL() int {
	ttlVal := getEnvValue("SESSION_TTL", "3600") // Default to 1 hour
	ttlSec, err := strconv.Atoi(ttlVal)

	if err != nil || ttlSec < 3600 {
		ttlSec = int(time.Hour) // Ensure TTL is at least 1 hour
	}

	return ttlSec
}

// SelfRegistrationEnabled returns whether registration is enabled.
func (c *Config) SelfRegistrationEnabled() bool {
	registrationEnabled := getEnvValue("SELF_REGISTRATION", "false")
	return registrationEnabled == trueStr || registrationEnabled == "1"
}

// GetSQLDriver returns the SQL driver to be used for database connections.
func (c *Config) GetSQLDriver() SQLDriver {
	driver := getEnvValue("SQL_DRIVER", builder.MYSQL)
	switch strings.ToLower(driver) {
	case builder.MYSQL:
		return SQLDriverMySQL
	case "sqlite", "sqlite3":
		return SQLDriverSQLite
	default:
		log.Fatalf("unsupported SQL driver: %s", driver)
		return ""
	}
}
