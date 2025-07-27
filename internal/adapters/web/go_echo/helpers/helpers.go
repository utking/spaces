// Package helpers provides utility functions for the application.
package helpers

import (
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gogs.utking.net/utking/spaces/internal"
)

// ErrorMessage returns a string representation of the error.
//   - If the error is nil, an empty string is returned.
//   - If the error is not nil, the error message is returned.
func ErrorMessage(err error) string {
	if err == nil {
		return ""
	}

	return err.Error()
}

// GetIDParam returns the ID from the URL `id` param.
func GetIDParam(c echo.Context) (id string) {
	return c.Param("id")
}

// GenerateUUID generates a new UUID v4 string.
func GenerateUUID() string {
	return uuid.New().String()
}

// GetReleaseVersion returns the release version of the application.
func GetReleaseVersion() string {
	return internal.Version
}

// FileIsImage checks if the given file name has an image file extension.
// A templates helper function that checks if a file is an image based on its extension.
func FileIsImage(fileName string) bool {
	fileName = strings.ToLower(fileName)

	return strings.HasSuffix(fileName, ".png") ||
		strings.HasSuffix(fileName, ".jpg") ||
		strings.HasSuffix(fileName, ".jpeg") ||
		strings.HasSuffix(fileName, ".gif") ||
		strings.HasSuffix(fileName, ".webp") ||
		strings.HasSuffix(fileName, ".svg")
}

// FileIsViewable checks if the given file name is viewable based on its extension.
// A templates helper function that checks if a file is viewable based on its extension.
func FileIsViewable(fileName string) bool {
	fileName = strings.ToLower(fileName)

	return strings.HasSuffix(fileName, ".md") ||
		strings.HasSuffix(fileName, ".yaml") ||
		strings.HasSuffix(fileName, ".yml") ||
		strings.HasSuffix(fileName, ".json") ||
		strings.HasSuffix(fileName, ".txt") ||
		strings.HasSuffix(fileName, ".pdf") ||
		strings.HasSuffix(fileName, ".txt")
}

// FileIconNameFromExt returns the icon name based on the file extension.
// A templates helper function that returns the icon name based on the file extension.
func FileIconNameFromExt(fileName string) string {
	fileName = strings.ToLower(fileName)

	switch {
	case strings.HasSuffix(fileName, ".zip"),
		strings.HasSuffix(fileName, ".tar"),
		strings.HasSuffix(fileName, ".tar.gz"):
		return "zip"
	case strings.HasSuffix(fileName, ".yaml"),
		strings.HasSuffix(fileName, ".yml"):
		return "yaml"
	case strings.HasSuffix(fileName, ".md"):
		return "md"
	case strings.HasSuffix(fileName, ".png"),
		strings.HasSuffix(fileName, ".jpg"),
		strings.HasSuffix(fileName, ".jpeg"),
		strings.HasSuffix(fileName, ".gif"),
		strings.HasSuffix(fileName, ".webp"),
		strings.HasSuffix(fileName, ".svg"):

		return "image"
	case strings.HasSuffix(fileName, ".json"):
		return "json"
	case strings.HasSuffix(fileName, ".pdf"):
		return "pdf"
	default:
		return "file"
	}
}
