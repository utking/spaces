package db

import (
	"errors"
	"time"

	"github.com/utking/spaces/internal/application/domain"
)

// User represents a user in the database.
type User struct {
	RoleName        *string   `db:"role_name"` // Additional fields - calculated or filled by JOINs
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
	Username        string    `db:"username"`                 // len 4-16, required
	PasswordHash    string    `db:"password_hash"`            // max len 255, required
	Email           string    `db:"email"`                    // max len 255, required
	AuthKey         string    `db:"auth_key"`                 // max len 32
	ActivationToken string    `db:"account_activation_token"` // max len 255, optional
	ID              string    `db:"id"`
	Status          int64     `db:"status"`
}

// TableName returns the name of the table in the database.
func (User) TableName() string {
	return "user"
}

// Validate checks the validity of the UserItem fields.
func (o *User) Validate() error {
	if len(o.Username) < 4 || len(o.Username) > 16 {
		return errors.New("username length must be between 4 and 16 characters")
	}

	if len(o.Email) > 255 {
		return errors.New("email length must be less than 255 characters")
	}

	if o.Email == "" {
		return errors.New("email is required")
	}

	return nil
}

// ValidateUpdate checks the validity of the UserItem fields for update.
func (o *User) ValidateUpdate() error {
	if len(o.Email) > 255 {
		return errors.New("email length must be less than 255 characters")
	}

	if o.Email == "" {
		return errors.New("email is required")
	}

	return nil
}

type AuthAssignment struct {
	UserID   string `db:"user_id"`   // User ID
	RoleName string `db:"item_name"` // Role name
}

// TableName returns the name of the table in the database for AuthAssignment.
func (AuthAssignment) TableName() string {
	return "auth_assignment"
}

// Validate checks the validity of the AuthAssignment fields.
func (a *AuthAssignment) Validate() error {
	if a.UserID == "" {
		return errors.New("user_id must be set")
	}

	if a.RoleName == "" {
		return errors.New("role name is required")
	}

	return nil
}

// UserSettings represents the settings for a user in the database.
type UserSettings struct {
	UserID string `db:"user_id"`
	Value  string `db:"value"` // JSON encoded settings
}

// TableName returns the name of the table in the database for UserSettings.
func (UserSettings) TableName() string {
	return "user_settings"
}

// ToStruct converts the UserSettings to a struct from domain.UserSettings.
func (s *UserSettings) ToStruct() (*domain.UserSettings, error) {
	settings, err := domain.UserSettings{}.FromJSON(s.Value)
	if err != nil {
		return &domain.UserSettings{}, nil
	}

	return settings, nil
}
