package ports

import (
	"context"

	"github.com/utking/spaces/internal/application/domain"
)

// DBPort is an interface that defines the methods for database operations.
type DBPort interface {
	// System stats
	GetSystemStats(context.Context, string) (*domain.SystemStats, error)

	// Notes
	GetNoteTags(ctx context.Context, uid string) ([]string, error)
	GetNotes(ctx context.Context, uid string, req *domain.NoteSearchRequest) ([]domain.Note, error)
	SearchNotesByTerm(ctx context.Context, uid string, req *domain.NoteRequest) ([]domain.Note, error)
	GetNotesCount(ctx context.Context, uid string, req *domain.NoteSearchRequest) (int64, error)
	GetNote(ctx context.Context, uid, id string) (*domain.Note, error)
	CreateNote(ctx context.Context, uid string, req *domain.Note) (string, error)
	UpdateNote(ctx context.Context, uid, id string, req *domain.Note) (int64, error)
	DeleteNote(ctx context.Context, uid, id string) error
	GetNotesMap(ctx context.Context, uid string, req *domain.NoteSearchRequest) ([]domain.Note, error)

	// Users
	GetUsers(ctx context.Context, req *domain.UserRequest) ([]domain.User, error)
	GetUsersCount(ctx context.Context, req *domain.UserRequest) (int64, error)
	GetUser(ctx context.Context, id string) (*domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
	CreateUser(ctx context.Context, req *domain.User) (string, string, error)
	UpdateUser(ctx context.Context, id string, req *domain.UserUpdate) (int64, error)
	DeleteUser(ctx context.Context, id string) error
	SetUserVerified(ctx context.Context, token string) (*domain.User, error)
	ChangePassword(ctx context.Context, id string, newPassword string) error
	GetUserAuthKey(ctx context.Context, id string) ([]byte, error)
	UpdateUserAuthKey(ctx context.Context, uid string, newEncKey []byte) error
	// User Settings
	GetUserSettings(ctx context.Context, id string) (*domain.UserSettings, error)
	UpdateUserSettings(ctx context.Context, id string, settings *domain.UserSettings) error

	// Secrets
	GetSecretTags(ctx context.Context, uid string) ([]string, error)
	GetSecrets(ctx context.Context, uid string, req *domain.SecretSearchRequest) ([]domain.Secret, error)
	SearchSecretsByTerm(ctx context.Context, uid string, req *domain.SecretRequest) ([]domain.Secret, error)
	GetSecretsCount(ctx context.Context, uid string, req *domain.SecretSearchRequest) (int64, error)
	GetSecret(ctx context.Context, uid, id string) (*domain.Secret, error)
	CreateSecret(ctx context.Context, uid string, req *domain.Secret) (string, error)
	UpdateSecret(ctx context.Context, uid, id string, req *domain.Secret) (int64, error)
	DeleteSecret(ctx context.Context, uid, id string) error
	UpdateEncryptedSecrets(ctx context.Context, uid string, items map[string]domain.EncryptSecret) error

	GetSecretsMap(
		ctx context.Context,
		uid string,
		req *domain.SecretSearchRequest,
	) ([]domain.SecretExportItem, error)

	// Bookmarks
	GetBookmarkTags(ctx context.Context, uid string) ([]string, error)
	GetBookmarks(ctx context.Context, uid string, req *domain.BookmarkSearchRequest) ([]domain.Bookmark, error)
	SearchBookmarksByTerm(ctx context.Context, uid string, req *domain.BookmarkSearchRequest) ([]domain.Bookmark, error)
	GetBookmarksCount(ctx context.Context, uid string, req *domain.BookmarkSearchRequest) (int64, error)
	GetBookmark(ctx context.Context, uid, id string) (*domain.Bookmark, error)
	CreateBookmark(ctx context.Context, uid string, req *domain.Bookmark) (string, error)
	UpdateBookmark(ctx context.Context, uid, id string, req *domain.Bookmark) (int64, error)
	DeleteBookmark(ctx context.Context, uid, id string) error
	GetBookmarksMap(ctx context.Context, uid string, req *domain.BookmarkSearchRequest) ([]domain.Bookmark, error)

	// Last Opened
	GetLastOpened(ctx context.Context, itemType domain.LastOpenedType, uid string) (string, error)
	SetLastOpened(ctx context.Context, itemType domain.LastOpenedType, uid string, itemID string) error
}
