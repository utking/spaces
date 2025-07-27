package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/labstack/echo/v4"
	"gogs.utking.net/utking/spaces/internal/adapters/web/go_echo/helpers"
	"gogs.utking.net/utking/spaces/internal/application/domain"
	"gogs.utking.net/utking/spaces/internal/ports"
)

// getImportNotesWrapper is a wrapper for the import notes handler.
func getImportNotesWrapper() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(
			http.StatusOK,
			"import/notes.html",
			map[string]interface{}{
				"Title": "Import Notes",
			},
		)
	}
}

// postImportNotesWrapper is a wrapper for the import notes handler.
// If accepts a POST request to import notes from a JSON file.
func postImportNotesWrapper(
	api ports.NotesService,
	userAPI ports.UsersService,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			userID = GetUserID(c, userAPI)
			code   = http.StatusOK
		)

		// Multipart form
		form, fErr := c.MultipartForm()
		if fErr != nil {
			return fErr
		}

		var (
			errList        []error
			totalNotes     int
			processedNotes int
		)

		files := form.File["files"]
		if len(files) == 0 {
			return c.Render(
				http.StatusBadRequest,
				"import/notes.html",
				map[string]interface{}{
					"Title": "Import Notes",
					"Error": "No files provided for import",
				},
			)
		}

		for _, file := range files {
			// read the file into notes map
			notes, rErr := readNotesFromFile(file)

			if rErr != nil || notes == nil {
				return c.Render(
					http.StatusBadRequest,
					"import/notes.html",
					map[string]interface{}{
						"Title": "Import Notes",
						"Error": helpers.ErrorMessage(rErr),
					},
				)
			}

			totalNotes = len(notes)

			for _, n := range notes {
				// Create note
				if _, cErr := api.Create(c.Request().Context(), userID, &n); cErr != nil {
					errList = append(
						errList,
						fmt.Errorf("failed to create note %q: %w; ", n.Title, cErr),
					)
					code = http.StatusInternalServerError
				} else {
					processedNotes++
				}
			}
		}

		return c.Render(
			code,
			"import/notes.html",
			map[string]interface{}{
				"Title":     "Import Notes",
				"Total":     totalNotes,
				"Processed": processedNotes,
				"Errors":    errList,
			},
		)
	}
}

// readNotesFromFile reads notes from a JSON file and returns them as a map.
func readNotesFromFile(file *multipart.FileHeader) ([]domain.Note, error) {
	// Open the file
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file %q: %w", file.Filename, err)
	}

	defer src.Close()

	// Decode the JSON file
	var notes = make([]domain.Note, 0)
	if err = json.NewDecoder(src).Decode(&notes); err != nil {
		return nil, fmt.Errorf("failed to decode JSON file: %w", err)
	}

	return notes, nil
}

// getImportBookmarksWrapper is a wrapper for the import bookmarks handler.
func getImportBookmarksWrapper() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(
			http.StatusOK,
			"import/bookmarks.html",
			map[string]interface{}{
				"Title": "Import Bookmarks",
			},
		)
	}
}

// postImportBookmarksWrapper is a wrapper for the import bookmarks handler.
func postImportBookmarksWrapper(
	api ports.BookmarkService,
	userAPI ports.UsersService,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			userID = GetUserID(c, userAPI)
			code   = http.StatusOK
		)

		// Multipart form
		form, err := c.MultipartForm()
		if err != nil {
			return err
		}

		var errList []error

		files := form.File["files"]
		if len(files) == 0 {
			return c.Render(
				http.StatusBadRequest,
				"import/bookmarks.html",
				map[string]interface{}{
					"Title": "Import Bookmarks",
					"Error": "No files provided for import",
				},
			)
		}

		for _, file := range files {
			bookmarks, rErr := readBookmarksFromFile(file)

			if rErr != nil || bookmarks == nil {
				return c.Render(
					http.StatusBadRequest,
					"import/bookmarks.html",
					map[string]interface{}{
						"Title": "Import Bookmarks",
						"Error": helpers.ErrorMessage(rErr),
					},
				)
			}

			for _, b := range bookmarks {
				if _, err = api.Create(c.Request().Context(), userID, &b); err != nil {
					errList = append(
						errList,
						fmt.Errorf("failed to create bookmark %q: %w; ", b.Title, err),
					)
					code = http.StatusInternalServerError
				}
			}
		}

		return c.Render(
			code,
			"import/bookmarks.html",
			map[string]interface{}{
				"Title":    "Import Bookmarks",
				"Imported": len(files),
				"Errors":   errList,
			},
		)
	}
}

// readBookmarksFromFile reads bookmarks from a JSON file and returns them as a slice.
func readBookmarksFromFile(file *multipart.FileHeader) ([]domain.Bookmark, error) {
	// Open the file
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file %q: %w", file.Filename, err)
	}

	defer src.Close()

	var bookmarks []domain.Bookmark
	if err = json.NewDecoder(src).Decode(&bookmarks); err != nil {
		return nil, fmt.Errorf("failed to decode JSON file: %w", err)
	}

	return bookmarks, nil
}

// getImportSecretsWrapper is a wrapper for the import secrets handler.
func getImportSecretsWrapper() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(
			http.StatusOK,
			"import/secrets.html",
			map[string]interface{}{
				"Title": "Import Secrets",
			},
		)
	}
}

// postImportSecretsWrapper is a wrapper for the import secrets handler.
func postImportSecretsWrapper(
	api ports.SecretService,
	userAPI ports.UsersService,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			userID = GetUserID(c, userAPI)
			code   = http.StatusOK
		)

		authKey, akErr := userAPI.GetAuthKey(c.Request().Context(), userID)
		if akErr != nil {
			return c.Render(
				http.StatusInternalServerError,
				"import/secrets.html",
				map[string]interface{}{
					"Title": "Import Secrets",
					"Error": "failed to get encryption key, cannot proceed with import",
				},
			)
		}

		// Multipart form
		form, formErr := c.MultipartForm()
		if formErr != nil {
			return formErr
		}

		var errList []error

		files := form.File["files"]
		if len(files) == 0 {
			return c.Render(
				http.StatusBadRequest,
				"import/secrets.html",
				map[string]interface{}{
					"Title": "Import Secrets",
					"Error": "No files provided for import",
				},
			)
		}

		for _, file := range files {
			secrets, err := readSecretsFromFile(
				c.Request().Context(),
				file,
				api,
				[]byte(c.FormValue("password")),
			)

			if err != nil {
				return c.Render(
					http.StatusBadRequest,
					"import/secrets.html",
					map[string]interface{}{
						"Title": "Import Secrets",
						"Error": helpers.ErrorMessage(err),
					},
				)
			}

			for _, item := range secrets {
				var encErr error

				s := domain.Secret{
					Tags:        item.Tags,
					Name:        item.Name,
					URL:         item.URL,
					Description: item.Description,
				}

				// encrypt the secret
				if s.EncodedSecret, encErr = encryptString(
					c.Request().Context(),
					api,
					authKey,
					item.Password,
				); encErr != nil {
					errList = append(
						errList,
						fmt.Errorf("failed to encrypt secret %q: %w", s.Name, encErr),
					)
					code = http.StatusInternalServerError
					continue
				}

				// encrypt the username
				if s.EncodedUsername, encErr = encryptString(
					c.Request().Context(),
					api,
					authKey,
					item.Username,
				); encErr != nil {
					errList = append(
						errList,
						fmt.Errorf("failed to encrypt username %q: %w", s.Name, encErr),
					)
					code = http.StatusInternalServerError
					continue
				}

				if _, cErr := api.Create(c.Request().Context(), userID, &s); cErr != nil {
					errList = append(
						errList,
						fmt.Errorf("failed to create secret %q: %w; ", s.Name, cErr),
					)
					code = http.StatusInternalServerError
				}
			}
		}

		return c.Render(
			code,
			"import/secrets.html",
			map[string]interface{}{
				"Title":    "Import Secrets",
				"Imported": len(files),
				"Errors":   errList,
			},
		)
	}
}

// readSecretsFromFile reads secrets from a JSON file and returns them as a slice.
func readSecretsFromFile(
	ctx context.Context,
	file *multipart.FileHeader,
	secretsAPI ports.SecretService,
	authKey []byte,
) ([]domain.SecretExportItem, error) {
	// Open the file
	src, fErr := file.Open()
	if fErr != nil {
		return nil, errors.New("failed to open file")
	}

	defer src.Close()

	// Decode the JSON file
	// If the file is empty, return an empty slice
	if file.Size == 0 {
		return []domain.SecretExportItem{}, nil
	}

	// read []bytes from the file into encData
	encData := bytes.NewBuffer(nil)

	if bRead, err := io.Copy(encData, src); err != nil || bRead == 0 {
		return nil, errors.New("failed to read the file or file is empty")
	}

	secrets := make([]domain.SecretExportItem, 0)

	nonce := encData.Bytes()[:12] // Assuming the first 12 bytes are the nonce
	fData := encData.Bytes()[12:] // The rest is the encrypted data

	// Decrypt the data
	if data, decErr := secretsAPI.Decrypt(
		ctx,     /* context */
		nonce,   /* nonce */
		fData,   /* encrypted data */
		authKey, /* auth key for decryption */
	); decErr != nil {
		return nil, errors.New("failed to decrypt file")
	} else if uErr := json.Unmarshal(data, &secrets); uErr != nil {
		return nil, errors.New("failed to decode JSON file")
	}

	return secrets, nil
}
