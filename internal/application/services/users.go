package services

import (
	"context"
	"errors"

	"github.com/utking/spaces/internal/application/domain"
	"github.com/utking/spaces/internal/ports"
)

const (
	fakeHash = "$2a$13$Xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx.xxx"
)

// UsersService is a struct that implements the UsersService interface.
type UsersService struct {
	db ports.DBPort
	fs ports.FileSystem
}

// NewUsersService creates a new instance of UsersService.
func NewUsersService(db ports.DBPort, fs ports.FileSystem) *UsersService {
	return &UsersService{
		db: db,
		fs: fs,
	}
}

// GetItems retrieves a list of users from the database based on the provided
// request parameters.
func (a *UsersService) GetItems(ctx context.Context, req *domain.UserRequest) ([]domain.User, error) {
	return a.db.GetUsers(ctx, req)
}

// GetCount retrieves the total number of users in the database based on
// the provided request parameters.
func (a *UsersService) GetCount(ctx context.Context, req *domain.UserRequest) (int64, error) {
	return a.db.GetUsersCount(ctx, req)
}

// GetItem retrieves a single user from the database based on the provided ID.
func (a *UsersService) GetItem(ctx context.Context, id string) (*domain.User, error) {
	return a.db.GetUser(ctx, id)
}

// Create creates a new user in the database based on the provided
// request parameters. It returns the ID of the newly created user.
// It also validates the request parameters before creating.
func (a *UsersService) Create(ctx context.Context, req *domain.User) (id, token string, err error) {
	if err = req.Validate(); err != nil {
		return id, token, err
	}

	return a.db.CreateUser(ctx, req)
}

// Update updates an existing user in the database based on the provided
// ID and request parameters. It returns the number of rows affected.
// It also validates the request parameters before updating.
func (a *UsersService) Update(ctx context.Context, id string, req *domain.UserUpdate) (int64, error) {
	if err := req.Validate(); err != nil {
		return 0, err
	}

	return a.db.UpdateUser(ctx, id, req)
}

// Delete deletes a user from the database based on the provided ID.
func (a *UsersService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("user ID must be provided")
	}

	return a.db.DeleteUser(ctx, id)
}

// ValidateUser validates the provided username and password against the
// database. It returns an error if the validation fails.
func (a *UsersService) ValidateUser(ctx context.Context, username, password string) (string, error) {
	user, err := a.db.GetUserByUsername(ctx, username)
	if err != nil || user == nil {
		// should not return an error here, to prevent information leakage
		// and time-based attacks
		_ = domain.PasswordVerify("no-password", fakeHash)

		return "", errors.New("invalid username or password")
	}

	if !(domain.PasswordVerify(password, user.PasswordHash)) {
		return "", errors.New("invalid username or password")
	}

	return user.RoleName, nil
}

// GetByUsername retrieves a user from the database based on the provided
// username. It returns the user object if found, or an error if not.
// INFO: This method finds only active users, not all users.
func (a *UsersService) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	return a.db.GetUserByUsername(ctx, username)
}

// VerifyUser verifies a user by a given verification token.
func (a *UsersService) VerifyUser(ctx context.Context, token string) (*domain.User, error) {
	if token == "" {
		return nil, errors.New("verification token is required")
	}

	return a.db.SetUserVerified(ctx, token)
}

// ChangePassword changes the password of a user by the provided ID.
func (a *UsersService) ChangePassword(ctx context.Context, id, newPassword string) error {
	if newPassword == "" {
		return errors.New("new password is required")
	}

	return a.db.ChangePassword(ctx, id, newPassword)
}

// GetDiskUsage retrieves the disk usage for a user by their ID.
func (a *UsersService) GetDiskUsage(ctx context.Context, id string) (int64, error) {
	// do not check for empty ID here, as it is valid to get disk usage for all users
	return a.fs.GetDiskUsage(ctx, id)
}

// CreateDataDirectory creates a data directory for a user by their ID.
func (a *UsersService) CreateDataDirectory(ctx context.Context, userID string) error {
	if userID == "" {
		return errors.New("invalid user ID")
	}

	return a.fs.CreateUserDataDirectory(ctx, userID)
}

// GetAuthKey retrieves the authentication key for a user by their ID.
func (a *UsersService) GetAuthKey(ctx context.Context, id string) ([]byte, error) {
	if id == "" {
		return nil, errors.New("user ID must be provided")
	}

	authKey, err := a.db.GetUserAuthKey(ctx, id)
	if err != nil {
		return nil, err
	}

	if len(authKey) == 0 {
		return nil, errors.New("auth key not found")
	}

	return authKey, nil
}

// UpdateAuthKey updates the authentication key for a user by their ID.
func (a *UsersService) UpdateAuthKey(ctx context.Context, uid string, newEncKey []byte) error {
	if uid == "" {
		return errors.New("user ID must be provided")
	}

	if len(newEncKey) == 0 {
		return errors.New("new encryption key must be provided")
	}

	return a.db.UpdateUserAuthKey(ctx, uid, newEncKey)
}

// GetUserSettings retrieves the user settings for a user by their ID.
func (a *UsersService) GetUserSettings(
	ctx context.Context,
	id string,
) (*domain.UserSettings, error) {
	if id == "" {
		return nil, errors.New("user ID must be provided")
	}

	settings, err := a.db.GetUserSettings(ctx, id)
	if err != nil {
		return nil, err
	}

	if settings == nil {
		settings = &domain.UserSettings{}
	}

	return settings, nil
}

// UpdateUserSettings updates the user settings for a user by their ID.
func (a *UsersService) UpdateUserSettings(
	ctx context.Context,
	id string,
	settings *domain.UserSettings,
) error {
	if id == "" {
		return errors.New("user ID must be provided")
	}

	if settings == nil {
		return errors.New("settings must not be nil")
	}

	return a.db.UpdateUserSettings(ctx, id, settings)
}
