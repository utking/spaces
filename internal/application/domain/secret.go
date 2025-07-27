package domain

import (
	"errors"
	"time"
)

// Secret represents a secret in the system.
type Secret struct {
	CreatedAt       time.Time `json:"-"`
	UpdatedAt       time.Time `json:"-"`
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	Name            string    `json:"name"` // len 1-128
	URL             string    `json:"url"`
	Description     string    `json:"description"`
	Tags            []string  `json:"tags"`     // JSON string, can be empty
	EncodedUsername []byte    `json:"username"` // len 0-1024, encrypted username, can be empty
	EncodedSecret   []byte    `json:"secret"`   // len 0-4096, encrypted secret, can be empty
	Password        string    `json:"-"`        // filled by a separate call on read. no write
	Username        string    `json:"-"`        // filled by a separate call on read. no write
}

// Validate checks if the Secret is valid.
func (s *Secret) Validate() error {
	if s.Name == "" || len(s.Name) > 128 {
		return errors.New("name must be between 1 and 128 characters")
	}
	if len(s.URL) > 256 {
		return errors.New("url must not exceed 256 characters")
	}
	if len(s.Description) > 512 {
		return errors.New("description must not exceed 512 characters")
	}
	if len(s.EncodedUsername) > 1024 {
		return errors.New("username must not exceed 1024 characters")
	}
	if len(s.EncodedSecret) > 4096 {
		return errors.New("secret must not exceed 4096 bytes")
	}

	if len(s.Tags) == 0 {
		return errors.New("tags cannot be empty")
	}

	return nil
}

// SecretSearchRequest represents a request for searching notes.
type SecretSearchRequest struct {
	Name     string `query:"name"`
	Tag      string `query:"tag"`
	SecretID string `query:"secret_id"`
	RequestPageMeta
}

// SecretRequest represents a request for creating/updating secrets.
type SecretRequest struct {
	Tags                []string `json:"tags"           form:"tags"`
	Name                string   `json:"name"           form:"name"`
	Username            string   `json:"username"       form:"username"`
	URL                 string   `json:"url"            form:"url"`
	Description         string   `json:"description"    form:"description"`
	PasswordSecretValue string   `json:"secret_value"   form:"secret_value"`
	UsernameSecretValue string   `json:"username_value" form:"username_value"`
	SecretID            string   `json:"secret_id"      form:"secret_id"`
	RequestPageMeta
}

// EncryptSecret - struct to update the secret field only.
type EncryptSecret struct {
	ID       string `json:"id"`
	Password []byte `json:"secret"`
	Username []byte `json:"username"`
}

// SecretEncodeRequest represents a request for decoding a secret.
type SecretEncodeRequest struct {
	PlainText []byte
}

type SecretExportRequest struct {
	Password string `form:"password"`
}

// Validate checks if the SecretExportRequest is valid.
func (r SecretExportRequest) Validate() error {
	if r.Password == "" {
		return errors.New("password cannot be empty")
	}

	if len(r.Password) != 32 {
		return errors.New("the password must be exactly 32 characters long")
	}

	return nil
}

type SecretExportItem struct {
	EncodedPassword []byte   `json:"-"`
	EncodedUsername []byte   `json:"-"`
	Tags            []string `json:"tags"`
	Username        string   `json:"username"`
	Name            string   `json:"name"`
	Password        string   `json:"password"`
	URL             string   `json:"url"`
	Description     string   `json:"description"`
}
