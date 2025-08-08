package ports

import (
	"context"

	"github.com/utking/spaces/internal/application/domain"
)

// UsersService is an interface that defines the methods for user-related operations.
type UsersService interface {
	GetItems(ctx context.Context, req *domain.UserRequest) ([]domain.User, error)
	GetCount(ctx context.Context, req *domain.UserRequest) (int64, error)
	GetItem(ctx context.Context, id string) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	Create(ctx context.Context, req *domain.User) (string, string, error)
	Update(ctx context.Context, id string, req *domain.UserUpdate) (int64, error)
	Delete(ctx context.Context, id string) error
	ValidateUser(ctx context.Context, username, password string) (string, error)
	VerifyUser(ctx context.Context, token string) (*domain.User, error)
	ChangePassword(ctx context.Context, id string, newPassword string) error
	GetDiskUsage(ctx context.Context, id string) (int64, error)
	CreateDataDirectory(ctx context.Context, uid string) error
	GetAuthKey(ctx context.Context, id string) ([]byte, error)
	UpdateAuthKey(ctx context.Context, uid string, newEncKey []byte) error
	GetUserSettings(ctx context.Context, id string) (*domain.UserSettings, error)
	UpdateUserSettings(ctx context.Context, id string, settings *domain.UserSettings) error
}
