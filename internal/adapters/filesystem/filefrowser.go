package filesystem

import (
	"context"

	"github.com/utking/spaces/internal/application/domain"
)

type FileBrowserAdapter struct {
	fsService *FSAdapter
}

// NewFileBrowserAdapter creates a new instance of FileBrowserAdapter with the specified base path.
func NewFileBrowserAdapter(basePath string) *FileBrowserAdapter {
	return &FileBrowserAdapter{
		fsService: NewAdapter(basePath),
	}
}

// ListFiles lists files in the specified directory for a given user.
func (f *FileBrowserAdapter) ListFiles(ctx context.Context, uid, path string) ([]domain.FileInfo, error) {
	return f.fsService.ListFiles(ctx, uid, path)
}

// CleanPath cleans the provided path by removing any leading or trailing slashes.
func (f *FileBrowserAdapter) CleanPath(path string) string {
	return f.fsService.CleanPath(path)
}

// UploadFile uploads a file to the specified path for a given user.
func (f *FileBrowserAdapter) UploadFile(ctx context.Context, uid, path string, content []byte) error {
	return f.fsService.UploadFile(ctx, uid, path, content)
}

// CreateFolder creates a new folder at the specified path for a given user.
func (f *FileBrowserAdapter) CreateFolder(ctx context.Context, uid, path string) error {
	return f.fsService.CreateFolder(ctx, uid, path)
}

// GetFileContent retrieves the content of a file at the specified path for a given user.
func (f *FileBrowserAdapter) GetFileContent(ctx context.Context, uid, path string) ([]byte, string, error) {
	return f.fsService.GetFileContent(ctx, uid, path)
}

// RenameFile renames a file from oldPath to newPath for a given user.
func (f *FileBrowserAdapter) RenameFile(ctx context.Context, uid, oldPath, newPath string) error {
	return f.fsService.RenameFile(ctx, uid, oldPath, newPath)
}

// DeleteFile deletes a file at the specified path for a given user.
func (f *FileBrowserAdapter) DeleteFile(ctx context.Context, uid, path string) error {
	return f.fsService.DeleteFile(ctx, uid, path)
}

// FileExists checks if a file exists at the specified path for a given user.
func (f *FileBrowserAdapter) FileExists(ctx context.Context, uid, path string) (bool, error) {
	return f.fsService.FileExists(ctx, uid, path)
}

// FileInternalName returns the internal name of a file based on its path.
func (f *FileBrowserAdapter) FileInternalName(ctx context.Context, uid, path string) (string, error) {
	return f.fsService.FileInternalName(ctx, uid, path)
}
