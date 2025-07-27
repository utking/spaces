// Package filesystem provides an adapter for file system operations.
package filesystem

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"gogs.utking.net/utking/spaces/internal/application/domain"
)

type FSAdapter struct {
	BasePath string
}

// NewAdapter creates a new instance of FSAdapter with the specified base path.
func NewAdapter(basePath string) *FSAdapter {
	return &FSAdapter{
		BasePath: basePath,
	}
}

// CreateUserDataDirectory creates a data directory for a user based on their user ID.
// Full path to the user's data directory is ${BasePath}/${uid}.
func (b *FSAdapter) CreateUserDataDirectory(_ context.Context, uid string) error {
	// Validate uid
	if uid == "" {
		return fmt.Errorf("invalid user id: %s", uid)
	}

	// Construct the full path for the user's data directory
	userPath := path.Join(b.BasePath, "users", uid)
	// Create the directory if it does not exist
	if err := createDirectory(userPath); err != nil {
		return fmt.Errorf("failed to create user data directory: %w", err)
	}

	return nil
}

// GetDiskUsage calculates the disk usage for a user based on their user ID.
// Full path to the user's data directory is ${BasePath}/${uid}.
// If uid is 0, it returns the total disk usage for all users.
func (b *FSAdapter) GetDiskUsage(ctx context.Context, uid string) (int64, error) {
	// Implement the logic to calculate disk usage for a user
	// based on the uid and the BasePath.
	if uid == "" {
		// Calculate total disk usage for all users
		return b.calculateTotalDiskUsage(ctx, b.BasePath)
	}

	// Calculate disk usage for a specific user
	userPath := path.Join(b.BasePath, "users", uid)

	// Run create user data directory to ensure the path exists
	if err := b.CreateUserDataDirectory(ctx, uid); err != nil {
		return 0, fmt.Errorf("failed to create user data directory: %w", err)
	}

	return b.calculateTotalDiskUsage(ctx, userPath)
}

// calculateTotalDiskUsage walks through the directory and calculates the total size of files.
func (b *FSAdapter) calculateTotalDiskUsage(_ context.Context, dirPath string) (int64, error) {
	var size int64

	err := filepath.Walk(dirPath, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			size += info.Size()
		}

		return err
	})

	return size, err
}

// createDirectory creates a directory if it does not exist.
func createDirectory(dirName string) error {
	err := os.MkdirAll(dirName, 0o700)
	if err != nil {
		return errors.New("failed to create profile directory")
	}

	return nil
}

// ListFiles lists files in the specified directory for a given user.
func (b *FSAdapter) ListFiles(_ context.Context, uid, filePath string) ([]domain.FileInfo, error) {
	// Validate uid
	if uid == "" {
		return nil, errors.New("invalid user id")
	}

	filePath = b.CleanPath(filePath)

	// Construct the full path for the user's data directory
	userPath := path.Join(b.BasePath, "users", uid, filePath)
	// Ensure the directory exists
	if _, err := os.Stat(userPath); os.IsNotExist(err) {
		return nil, errors.New("directory does not exist")
	}

	// Read the directory and list files
	files, err := os.ReadDir(userPath)
	if err != nil {
		return nil, errors.New("failed to read directory")
	}

	fileInfos := make([]domain.FileInfo, 0, len(files))
	for _, file := range files {
		fileInfo := domain.FileInfo{
			Name:  file.Name(),
			Path:  path.Join(filePath, file.Name()),
			IsDir: file.IsDir(),
		}

		// If it's a file, get its size
		if !file.IsDir() {
			info, fiErr := file.Info()
			if fiErr != nil {
				return nil, errors.New("failed to get file info")
			}

			fileInfo.Size = info.Size()
			fileInfo.Modified = info.ModTime().Format(time.RFC822)

			// get mime type
			mType, mErr := mimetype.DetectFile(path.Join(userPath, file.Name()))
			if mErr == nil {
				fileInfo.MimeType = mType.String()
			} else {
				fileInfo.MimeType = "application/octet-stream" // Default MIME type if detection fails
			}
		}

		fileInfos = append(fileInfos, fileInfo)
	}

	return fileInfos, nil
}

func (b *FSAdapter) CleanPath(filePath string) string {
	// Clean the path to remove any relative components
	return filepath.Clean("/" + filePath)
}

// UploadFile uploads a file to the specified path for a given user.
func (b *FSAdapter) UploadFile(_ context.Context, uid, filePath string, content []byte) error {
	// Validate uid
	if uid == "" {
		return errors.New("invalid user id")
	}

	filePath = b.CleanPath(filePath)

	// Construct the full path for the user's data directory
	userPath := path.Join(b.BasePath, "users", uid, filePath)

	// Ensure the directory exists
	if err := createDirectory(filepath.Dir(userPath)); err != nil {
		return errors.New("failed to create directory for file upload")
	}

	// Write the content to the file
	if err := os.WriteFile(userPath, content, 0o600); err != nil {
		return errors.New("failed to write file")
	}

	return nil
}

// CreateFolder creates a folder at the specified path for a given user.
func (b *FSAdapter) CreateFolder(_ context.Context, uid, folderPath string) error {
	// Validate uid
	if uid == "" {
		return errors.New("invalid user id")
	}

	folderPath = b.CleanPath(folderPath)

	// Construct the full path for the user's data directory
	userPath := path.Join(b.BasePath, "users", uid, folderPath)

	// Ensure the directory does not already exist
	if _, err := os.Stat(userPath); !os.IsNotExist(err) {
		return errors.New("path already exists")
	}

	// Create the directory
	if err := createDirectory(userPath); err != nil {
		return errors.New("failed to create folder")
	}

	return nil
}

// GetFileContent retrieves the content of a file at the specified path for a given user.
// Returns the content as a byte slice and the file type (mime).
func (b *FSAdapter) GetFileContent(_ context.Context, uid, filePath string) ([]byte, string, error) {
	// Validate uid
	if uid == "" {
		return nil, "", errors.New("invalid user id")
	}

	filePath = b.CleanPath(filePath)

	// Construct the full path for the user's data directory
	userPath := path.Join(b.BasePath, "users", uid, filePath)

	// Read the file content
	content, err := os.ReadFile(userPath)
	if err != nil {
		return nil, "", errors.New("failed to read file")
	}

	// Get the file MIME type
	if len(content) == 0 {
		return content, "text/plain", nil // Return empty content with a default MIME type
	}

	ext := mimetype.Detect(content).String()

	return content, ext, nil
}

// RenameFile renames a file from oldPath to newPath for a given user.
func (b *FSAdapter) RenameFile(_ context.Context, uid, oldPath, newPath string) error {
	// Validate uid
	if uid == "" {
		return errors.New("invalid user id")
	}

	oldPath = b.CleanPath(oldPath)
	newPath = b.CleanPath(newPath)

	// Construct the full paths for the user's data directory
	userDir := path.Join(b.BasePath, "users", uid)
	oldFilePath := path.Join(userDir, oldPath)
	newFilePath := path.Join(userDir, newPath)

	if oldFilePath == newFilePath {
		return nil // No need to rename if paths are the same
	}

	// Ensure the old file exists
	if _, err := os.Stat(oldFilePath); os.IsNotExist(err) {
		return errors.New("file does not exist")
	}

	// Rename the file
	if err := os.Rename(oldFilePath, newFilePath); err != nil {
		return errors.New("failed to rename file")
	}

	return nil
}

// DeleteFile deletes a file at the specified path for a given user.
func (b *FSAdapter) DeleteFile(_ context.Context, uid, filePath string) error {
	// Validate uid
	if uid == "" {
		return errors.New("invalid user id")
	}

	filePath = b.CleanPath(filePath)

	// Construct the full path for the user's data directory
	userPath := path.Join(b.BasePath, "users", uid, filePath)

	if filePath == "/" {
		return nil // Do not allow deletion of the root directory
	}

	// Ensure the file exists
	if _, err := os.Stat(userPath); os.IsNotExist(err) {
		return errors.New("file does not exist")
	}

	// Delete the file
	if err := os.RemoveAll(userPath); err != nil {
		return errors.New("failed to delete file")
	}

	return nil
}

// FileExists checks if a file exists at the specified path for a given user.
func (b *FSAdapter) FileExists(_ context.Context, uid, filePath string) (bool, error) {
	// Validate uid
	if uid == "" {
		return false, errors.New("invalid user id")
	}

	filePath = b.CleanPath(filePath)

	// Construct the full path for the user's data directory
	userPath := path.Join(b.BasePath, "users", uid, filePath)

	// Check if the file exists
	if _, err := os.Stat(userPath); os.IsNotExist(err) {
		return false, nil // File does not exist
	}

	return true, nil // File exists
}

// FileInternalName returns the internal name of a file for a given user.
func (b *FSAdapter) FileInternalName(_ context.Context, uid, filePath string) (string, error) {
	// Validate uid
	if uid == "" {
		return "", errors.New("invalid user id")
	}

	filePath = b.CleanPath(filePath)

	// Construct the full path for the user's data directory
	userPath := path.Join(b.BasePath, "users", uid, filePath)

	// Check if the file exists
	if _, err := os.Stat(userPath); os.IsNotExist(err) {
		return "", errors.New("file does not exist")
	}

	return userPath, nil
}
