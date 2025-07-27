package ports

import (
	"context"

	"github.com/utking/spaces/internal/application/domain"
)

// FileSystem is an interface that defines the methods for file system-related operations.
type FileSystem interface {
	GetDiskUsage(ctx context.Context, uid string) (int64, error)
	CreateUserDataDirectory(ctx context.Context, uid string) error
	CreateFolder(ctx context.Context, uid, path string) error
	ListFiles(ctx context.Context, uid, path string) ([]domain.FileInfo, error)
	UploadFile(ctx context.Context, uid, path string, fileData []byte) error
	GetFileContent(ctx context.Context, uid, path string) ([]byte, string, error)
	FileExists(ctx context.Context, uid, path string) (bool, error)
	FileInternalName(ctx context.Context, uid, path string) (string, error)
	DeleteFile(ctx context.Context, uid, path string) error
	RenameFile(ctx context.Context, uid, oldPath string, newPath string) error
	CleanPath(path string) string
}

// FileBrowserService is an interface that defines the methods for file browser-related operations.
type FileBrowserService interface {
	CreateFolder(ctx context.Context, uid, path string) error
	ListFiles(ctx context.Context, uid, path string) ([]domain.FileInfo, error)
	UploadFile(ctx context.Context, uid, path string, fileData []byte) error
	GetFileContent(ctx context.Context, uid, path string) ([]byte, string, error)
	FileExists(ctx context.Context, uid, path string) (bool, error)
	FileInternalName(ctx context.Context, uid, path string) (string, error)
	DeleteFile(ctx context.Context, uid, path string) error
	RenameFile(ctx context.Context, uid, oldPath string, newPath string) error
	CleanPath(path string) string
}
