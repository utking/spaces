package domain

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	// UserInactive represents the deleted status of the user.
	UserInactive = 0
	// UserActive represents the active status of the user.
	UserActive     = 10
	hashCost       = 13
	UserRoleUser   = "user"
	UserRoleAdmin  = "admin"
	charset        = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_-"
	activeStatus   = "Active"
	inactiveStatus = "Inactive"
)

var userRoles = []string{UserRoleAdmin, UserRoleUser}

func UserRolesList() []string {
	return userRoles
}

// GetPasswordHash generates a bcrypt hash of the given password.
func GetPasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
	return string(bytes), err
}

// PasswordVerify checks if the provided password matches the hashed password.
func PasswordVerify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

var (
	statusMap = map[int64]Status{
		UserInactive: {
			Title: inactiveStatus,
			ID:    UserInactive,
		},
		UserActive: {
			Title: activeStatus,
			ID:    UserActive,
		},
	}
)

// User represents a user in the system.
type User struct {
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	ID               string    `json:"id"            form:"id"`
	Username         string    `json:"username"      form:"username"` // len 4-16
	Password         string    `                     form:"password"` // len 8-32
	PasswordHash     string    `json:"password_hash"`                 // len 8-255
	Email            string    `json:"email"         form:"email"`
	PasswordConfirm  string    `                     form:"password_confirm"` // must match Password, not stored
	AuthKey          string    `json:"auth_key"`
	RoleName         string    // Filled by JOINs
	Status           int64     `json:"status"        form:"status"`
	SendNotification bool      `                     form:"send_notification"`
}

// Normalize normalizes the User struct fields. Username and Email are trimmed
// and converted to lowercase.
func (u *User) Normalize() {
	u.Username = strings.TrimSpace(u.Username)
	u.Email = strings.TrimSpace(u.Email)

	u.Password = strings.TrimSpace(u.Password)
	u.PasswordConfirm = strings.TrimSpace(u.PasswordConfirm)

	// lowercase the username and email
	u.Username = strings.ToLower(u.Username)
	u.Email = strings.ToLower(u.Email)
}

// Validate checks the validity of the User struct fields.
func (u *User) Validate() error {
	if len(u.Username) < 4 || len(u.Username) > 16 {
		return errors.New("username length must be between 4 and 16 characters")
	}

	if len(u.Password) < 8 || len(u.Password) > 32 {
		return errors.New("password length must be between 8 and 32 characters")
	}

	if len(u.Email) > 255 {
		return errors.New("email length must be less than 255 characters")
	}

	if u.Email == "" {
		return errors.New("email cannot be empty")
	}

	if u.Password != "" && u.Password != u.PasswordConfirm {
		return errors.New("password and password confirmation do not match")
	}

	return nil
}

func Int63() int64 {
	var b [8]byte
	_, _ = rand.Read(b[:]) // Fill the byte array with secure random bytes

	return int64(binary.LittleEndian.Uint64(b[:])) & (1<<63 - 1) //nolint:gosec // Won't overflow
}

// GenerateRandomString generates a random string of specified length.
// The string generated matches [A-Za-z0-9_-]+ and is transparent to URL-encoding.
func GenerateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[Int63()%int64(len(charset))]
	}

	return string(b)
}

// UserRequest represents a request for searching users.
type UserRequest struct {
	Status   *int64 `query:"status"`
	Username string `query:"username"`
	Email    string `query:"email"`
	RequestPageMeta
}

// Trim trims the strings in the UserRequest.
func (req *UserRequest) Trim() {
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
}

// UserStatus represents the status of a user.
type UserStatus struct{}

// GetStatusMap returns a map of all user statuses.
func (UserStatus) GetStatusMap() map[int64]Status {
	return statusMap
}

// UserUpdate represents a request to update user information.
type UserUpdate struct {
	ID              string `form:"id"`
	Username        string `form:"username"` // len 4-16
	Password        string `form:"password"` // len 8-32
	Email           string `form:"email"`
	PasswordConfirm string `form:"password_confirm"` // must match Password, not stored
	RoleName        string `form:"role_name"`        // len 1+
	Status          int64  `form:"status"`
}

// Normalize normalizes the User struct fields. Username and Email are trimmed
// and converted to lowercase.
func (u *UserUpdate) Normalize() {
	u.Username = strings.TrimSpace(u.Username)
	u.Email = strings.TrimSpace(u.Email)
	u.RoleName = strings.TrimSpace(u.RoleName)

	u.Password = strings.TrimSpace(u.Password)
	u.PasswordConfirm = strings.TrimSpace(u.PasswordConfirm)

	// lowercase the username and email
	u.Username = strings.ToLower(u.Username)
	u.Email = strings.ToLower(u.Email)
}

// Validate checks the validity of the UserUpdate struct fields.
func (u *UserUpdate) Validate() error {
	if len(u.Email) > 255 {
		return errors.New("email length must be less than 255 characters")
	}

	if u.Password != "" && u.Password != u.PasswordConfirm {
		return errors.New("password and password confirmation do not match")
	}

	if strings.TrimSpace(u.RoleName) == "" {
		return errors.New("role name cannot be empty")
	}

	return nil
}
