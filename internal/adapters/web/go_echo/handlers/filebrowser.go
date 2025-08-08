package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"path"
	"slices"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/utking/spaces/internal/adapters/web/go_echo/helpers"
	"github.com/utking/spaces/internal/application/domain"
	"github.com/utking/spaces/internal/infra/session"
	"github.com/utking/spaces/internal/ports"
)

const (
	maxFileSize = 100 * 1024 * 1024 // 100 MB
)

// getFileBrowserWrapper returns a handler function that serves the file browser page.
func getFileBrowserWrapper(
	fileBrowser ports.FileBrowserService,
	usersService ports.UsersService,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		tplFile := "filebrowser/index.html"
		viewMode := "list"
		basePath := fileBrowser.CleanPath(c.QueryParam("path"))
		parentPath := fileBrowser.CleanPath(basePath + "/..")
		viewModeTiles, _ := session.GetStrVar(c, "file_browser_tiles")
		if viewModeTiles == "true" {
			tplFile = "filebrowser/index-tiles.html"
			viewMode = "tile"
		}

		// Get the user ID from the session
		userID := GetUserID(c, usersService)
		if userID == "" {
			return c.Render(
				http.StatusInternalServerError,
				tplFile,
				map[string]interface{}{
					"Title":      "Files",
					"BasePath":   basePath,
					"ParentPath": parentPath,
					"Mode":       viewMode,
					"Error":      "User ID not found in session",
					"Files":      nil,
				},
			)
		}

		filesList, err := fileBrowser.ListFiles(
			c.Request().Context(),
			userID,
			basePath,
		)

		// Render the file browser page with the user ID
		return c.Render(
			http.StatusOK,
			tplFile,
			map[string]interface{}{
				"Title":      "Files",
				"BasePath":   basePath,
				"ParentPath": parentPath,
				"Count":      len(filesList),
				"Mode":       viewMode,
				"UserID":     userID,
				"Files":      filesList,
				"Error":      helpers.ErrorMessage(err),
			})
	}
}

// postFileBrowserUploadWrapper returns a handler function that handles file uploads in the file browser.
func postFileBrowserUploadWrapper(
	fileBrowser ports.FileBrowserService,
	usersService ports.UsersService,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			userID = GetUserID(c, usersService)
			code   = http.StatusOK
		)

		if userID == "" {
			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"Error": "User ID not found in session",
				})
		}

		// Multipart form
		form, fErr := c.MultipartForm()
		if fErr != nil {
			return fErr
		}

		// path from the form
		if len(form.Value["path"]) == 0 {
			return c.JSON(
				http.StatusBadRequest,
				map[string]string{
					"Error": "Path not provided in the form",
				})
		}

		basePath := fileBrowser.CleanPath(form.Value["path"][0])

		var (
			errList        []error
			processedItems int
		)

		files := form.File["files"]
		if len(files) == 0 {
			return c.JSON(
				http.StatusBadRequest,
				map[string]string{
					"Error": "No files provided for upload",
				})
		}

		for _, file := range files {
			// save the file to the user's directory and with basePath
			if file.Size < 0 {
				errList = append(errList, fmt.Errorf("file size is invalid: %s", file.Filename))
				continue
			}

			if file.Size > maxFileSize { // Limit file size to 100MB
				errList = append(errList, fmt.Errorf("file size exceeds limit (100MB): %s", file.Filename))
				continue
			}

			// get file data as []byte
			fileReader, err := file.Open()
			if err != nil {
				errList = append(errList, err)
				continue
			}

			fileData := make([]byte, file.Size)
			if _, err = fileReader.Read(fileData); err != nil {
				errList = append(errList, err)
				code = http.StatusInternalServerError
				_ = fileReader.Close()
				continue
			}

			if err = fileBrowser.UploadFile(
				c.Request().Context(),
				userID,
				path.Join(basePath, file.Filename),
				fileData,
			); err != nil {
				errList = append(errList, fmt.Errorf("failed to upload file %s", file.Filename))
				code = http.StatusInternalServerError
				_ = fileReader.Close()
				continue
			}

			processedItems++
		}

		return c.JSON(
			code,
			map[string]interface{}{
				"TotalFiles":     len(files),
				"ProcessedFiles": processedItems,
				"Error":          helpers.ErrorMessage(errors.Join(errList...)),
			})
	}
}

// postFileBrowserNewFolderWrapper returns a handler function that handles the
// creation of new folders in the file browser.
func postFileBrowserNewFolderWrapper(
	fileBrowser ports.FileBrowserService,
	usersService ports.UsersService,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			userID = GetUserID(c, usersService)
			code   = http.StatusOK
		)

		if userID == "" {
			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"Error": "User ID not found in session",
				})
		}

		folderName := c.FormValue("name")
		if folderName == "" {
			return c.JSON(
				http.StatusBadRequest,
				map[string]string{
					"Error": "Folder name is required",
				})
		}

		basePath := fileBrowser.CleanPath(c.FormValue("path"))

		if err := fileBrowser.CreateFolder(
			c.Request().Context(),
			userID,
			path.Join(basePath, folderName),
		); err != nil {
			code = http.StatusInternalServerError
			return c.JSON(
				code,
				map[string]string{
					"Error": helpers.ErrorMessage(err),
				})
		}

		return c.JSON(code, map[string]string{})
	}
}

// getFileBrowserFileViewWrapper returns a handler function that serves a file
// if it's an image, PDF, txt, markdown, or json.
func getFileBrowserFileViewWrapper(
	fileBrowser ports.FileBrowserService,
	usersService ports.UsersService,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := GetUserID(c, usersService)
		if userID == "" {
			return c.Blob(http.StatusBadRequest, "text/plain", []byte("User ID not found in session"))
		}

		filePath := fileBrowser.CleanPath(c.QueryParam("path"))
		if filePath == "" {
			return c.Blob(http.StatusBadRequest, "text/plain", []byte("File path is required"))
		}

		content, contentType, err := fileBrowser.GetFileContent(
			c.Request().Context(),
			userID,
			filePath,
		)
		if err != nil {
			return c.Blob(http.StatusInternalServerError, "text/plain", []byte("Failed to get file content"))
		}

		return c.Blob(http.StatusOK, contentType, content)
	}
}

// postFileBrowserFileRenameWrapper returns a handler function that handles file renaming in the file browser.
func postFileBrowserFileRenameWrapper(
	fileBrowser ports.FileBrowserService,
	usersService ports.UsersService,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			userID = GetUserID(c, usersService)
			code   = http.StatusOK
		)

		if userID == "" {
			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"Error": "User ID not found in session",
				})
		}

		folderPath := c.FormValue("path")
		oldName := c.FormValue("old_name")
		newName := c.FormValue("new_name")

		oldPath := fileBrowser.CleanPath(path.Join(folderPath, oldName))
		newPath := fileBrowser.CleanPath(path.Join(folderPath, newName))

		if err := fileBrowser.RenameFile(c.Request().Context(), userID, oldPath, newPath); err != nil {
			code = http.StatusInternalServerError
			return c.JSON(code, map[string]string{"Error": err.Error()})
		}

		return c.NoContent(code)
	}
}

// deleteFileBrowserFileWrapper returns a handler function that handles file deletion in the file browser.
func deleteFileBrowserFileWrapper(
	fileBrowser ports.FileBrowserService,
	usersService ports.UsersService,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := GetUserID(c, usersService)
		if userID == "" {
			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"Error": "User ID not found in session",
				})
		}

		filePath := c.FormValue("path")
		fileName := c.FormValue("name")

		fullPath := fileBrowser.CleanPath(path.Join(filePath, fileName))
		if fullPath == "" || fullPath == "/" {
			return c.NoContent(http.StatusOK)
		}

		if err := fileBrowser.DeleteFile(c.Request().Context(), userID, fullPath); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"Error": err.Error()})
		}

		return c.NoContent(http.StatusOK)
	}
}

// getFileBrowserFileDownloadWrapper returns a handler function that handles file downloads in the file browser.
func getFileBrowserFileDownloadWrapper(
	fileBrowser ports.FileBrowserService,
	usersService ports.UsersService,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := GetUserID(c, usersService)
		if userID == "" {
			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"Error": "User ID not found in session",
				})
		}

		folderPath := c.QueryParam("path")
		fileName := c.QueryParam("name")

		filePath := fileBrowser.CleanPath(path.Join(folderPath, fileName))
		if filePath == "" {
			return c.JSON(
				http.StatusBadRequest,
				map[string]string{
					"Error": "File path is required",
				})
		}

		_, err := fileBrowser.FileExists(c.Request().Context(), userID, filePath)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"Error": err.Error()})
		}

		internalName, err := fileBrowser.FileInternalName(c.Request().Context(), userID, filePath)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"Error": err.Error()})
		}

		return c.Attachment(
			internalName,
			path.Base(filePath), // Use the base name of the file for download
		)
	}
}

// postFileBrowserSetViewModeWrapper returns a handler function that sets the view mode for the file browser.
// Sets the selected mode in the session for the user.
func postFileBrowserSetViewModeWrapper(
	usersService ports.UsersService,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			userID = GetUserID(c, usersService)
			code   = http.StatusOK
		)

		curSettings, usErr := usersService.GetUserSettings(c.Request().Context(), userID)
		if usErr != nil {
			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"Error": "Could not read user settings",
				})
		}

		viewMode := c.FormValue("mode")
		if viewMode == "" || !slices.Contains(
			[]string{domain.FileBrowserViewModeList, domain.FileBrowserViewModeTiles},
			viewMode,
		) {
			viewMode = domain.FileBrowserViewModeList // Default view mode
		}

		switch viewMode {
		case domain.FileBrowserViewModeList:
			curSettings.FileBrowserTiles = false
		case domain.FileBrowserViewModeTiles:
			curSettings.FileBrowserTiles = true
		}

		usErr = usersService.UpdateUserSettings(c.Request().Context(), userID, curSettings)
		if usErr != nil {
			code = http.StatusInternalServerError
			return c.JSON(
				code,
				map[string]string{
					"Error": fmt.Sprintf("Failed to set view mode: %v", usErr),
				})
		}

		// Set the view mode in the session
		_ = session.SetStrVar(c, "file_browser_tiles", strconv.FormatBool(curSettings.FileBrowserTiles))

		return c.NoContent(code)
	}
}
