package domain

const (
	FileBrowserViewModeTiles = "tile" // View mode for displaying files as tiles
	FileBrowserViewModeList  = "list" // View mode for displaying files as a list
)

// FileInfo represents information about a file or directory in the file browser.
type FileInfo struct {
	Name     string `json:"name"`      // Name of the file or directory
	Modified string `json:"modified"`  // Last modified time in RFC3339 format
	Path     string `json:"path"`      // Full path of the file or directory
	MimeType string `json:"mime_type"` // MIME type of the file
	Size     int64  `json:"size"`      // Size of the file in bytes
	IsDir    bool   `json:"is_dir"`    // true if it's a directory, false if it's a file
}
