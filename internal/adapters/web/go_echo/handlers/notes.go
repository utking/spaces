package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"gogs.utking.net/utking/spaces/internal/adapters/web/go_echo/helpers"
	"gogs.utking.net/utking/spaces/internal/application/domain"
	"gogs.utking.net/utking/spaces/internal/ports"
)

func getNotesWrapper(
	api ports.NotesService,
	userAPI ports.UsersService,
	lastOpened ports.LastOpenedService,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			query = new(domain.NoteSearchRequest)
			note  = new(domain.Note)
			code  = http.StatusOK
			ctx   = c.Request().Context()
		)

		userID := GetUserID(c, userAPI)
		_ = c.Bind(query)

		// If no specific note or tag is requested, check the last opened note
		if query.Tag == "" && query.NoteID == "" {
			if redirectURL := lastOpenRedirectURL(
				ctx, userID, api, lastOpened,
			); redirectURL != "" {
				// redirect is found
				return c.Redirect(http.StatusSeeOther, redirectURL)
			}

			// If no last opened note, delete the last opened note ID from the DB
			_ = lastOpened.SetLastOpened(ctx, domain.LastOpenedTypeNote, userID, "")
			// Reset the query to avoid confusion
			query.NoteID = ""
		} else if query.NoteID != "" {
			// If a specific note is requested, set the last opened note ID
			_ = lastOpened.SetLastOpened(ctx, domain.LastOpenedTypeNote, userID, query.NoteID)
		}

		noteReq := &domain.NoteSearchRequest{Tag: query.Tag, NoteID: query.NoteID}

		items, iRrr := api.GetItems(ctx, userID, noteReq)
		tags, tErr := api.GetTags(ctx, userID)

		if len(items) > 0 && query.NoteID != "" {
			note, _ = api.GetItem(ctx, userID, query.NoteID)
		}

		err := errors.Join(tErr, iRrr)
		if err != nil {
			code = http.StatusInternalServerError
		}

		return c.Render(
			code,
			"notes/index.html",
			map[string]interface{}{
				"Title":      "Notes",
				"Items":      items,
				"Item":       note,
				"ItemsCount": len(items),
				"Tags":       tags,
				"Error":      helpers.ErrorMessage(err),
				"Query":      query,
				"TagsCount":  len(tags),
			},
		)
	}
}

// putNotesWrapper is a wrapper for the notes post handler.
// It handles the request and response for updating a note.
func putNotesWrapper(
	api ports.NotesService,
	userAPI ports.UsersService,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		note := new(domain.NoteRequest)

		_ = c.Bind(note)

		userID := GetUserID(c, userAPI)
		_, err := api.Update(c.Request().Context(), userID, note.NoteID, &domain.Note{
			Tags:    note.Tags,
			Title:   note.Title,
			Content: note.Content,
		})

		if err != nil {
			return c.JSON(
				http.StatusInternalServerError,
				map[string]interface{}{
					"Error": helpers.ErrorMessage(err),
				},
			)
		}

		return c.JSON(
			http.StatusOK,
			map[string]interface{}{
				"ID": note.NoteID,
			},
		)
	}
}

// getNoteCreateWrapper is a wrapper for the note create handler.
// Renders a form for creating a new note.
func getNoteCreateWrapper(
	api ports.NotesService,
	userAPI ports.UsersService,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			code  = http.StatusOK
			query = new(domain.NoteSearchRequest)
		)

		_ = c.Bind(query)

		userID := GetUserID(c, userAPI)
		items, _ := api.GetItems(c.Request().Context(), userID, query)
		tags, err := api.GetTags(c.Request().Context(), userID)

		if err != nil {
			code = http.StatusInternalServerError
		}

		return c.Render(
			code,
			"notes/create.html",
			map[string]interface{}{
				"Title":      "Create Note",
				"Tags":       tags,
				"Query":      query,
				"ItemsCount": len(items),
				"Error":      helpers.ErrorMessage(err),
			},
		)
	}
}

// postNoteCreateWrapper is a wrapper for the note create handler.
// It handles the request and response for creating a note.
// JSON response contains the error message if any and the created note ID.
func postNoteCreateWrapper(
	api ports.NotesService,
	userAPI ports.UsersService,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			note = new(domain.NoteRequest)
			code = http.StatusOK
		)

		err := c.Bind(note)
		if err != nil {
			return c.JSON(
				http.StatusBadRequest,
				map[string]interface{}{
					"Error": helpers.ErrorMessage(err),
				},
			)
		}

		userID := GetUserID(c, userAPI)
		noteID, err := api.Create(c.Request().Context(), userID, &domain.Note{
			Title:   note.Title,
			Content: note.Content,
			Tags:    note.Tags,
		})

		if err != nil {
			return c.JSON(
				http.StatusInternalServerError,
				map[string]interface{}{
					"Error": helpers.ErrorMessage(err),
				})
		}

		return c.JSON(
			code,
			map[string]interface{}{
				"ID":  noteID,
				"Tag": note.Tags[0], // Assuming at least one tag is provided
			},
		)
	}
}

// deleteNoteWrapper is a wrapper for the note delete handler.
// It handles the request and response for deleting a note.
// JSON response contains the error message if any.
func deleteNoteWrapper(
	api ports.NotesService,
	userAPI ports.UsersService,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			noteID = helpers.GetIDParam(c)
			code   = http.StatusOK
		)

		userID := GetUserID(c, userAPI)
		err := api.Delete(c.Request().Context(), userID, noteID)

		if err != nil {
			code = http.StatusInternalServerError
		}

		return c.JSON(
			code,
			map[string]interface{}{
				"Error": helpers.ErrorMessage(err),
			},
		)
	}
}

// getExportNotesWrapper is a wrapper for the notes export handler.
// it compiles a map of the user's notes and tags, exporting them to a JSON file,
// and returns a downloadable file to the user.
func getExportNotesWrapper(
	api ports.NotesService,
	userAPI ports.UsersService,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		const fileName = "notes_export.json"
		var eFile *os.File

		items, err := api.GetItemsMap(c.Request().Context(), GetUserID(c, userAPI), nil)
		if err == nil {
			eFile, err = saveNotesToFile(items)
			if err == nil {
				defer func() {
					_ = os.Remove(eFile.Name()) // Clean up the temporary file after sending
				}()
			}
		}

		if err != nil {
			return c.JSON(
				http.StatusInternalServerError,
				map[string]interface{}{
					"Error": helpers.ErrorMessage(err),
				},
			)
		}

		return c.Attachment(eFile.Name(), fileName)
	}
}

// saveNotesToFile is a helper function that saves the notes data to a file.
func saveNotesToFile(
	items []domain.Note,
) (*os.File, error) {
	if len(items) == 0 {
		return nil, errors.New("no notes to export")
	}

	data, err := json.Marshal(items)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal notes data: %w", err)
	}

	eFile, err := os.CreateTemp("", "notes_export*.json")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}

	_, err = eFile.Write(data)

	if err != nil {
		return nil, fmt.Errorf("failed to write data to file: %w", err)
	}

	if err = eFile.Close(); err != nil {
		return nil, fmt.Errorf("failed to close file: %w", err)
	}

	return eFile, nil
}

// getSearchNotesWrapper is a wrapper for the notes search handler.
// Returns a JSON with "items" key containing the search results slice.
func getSearchNotesWrapper(
	api ports.NotesService,
	userAPI ports.UsersService,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		term := c.QueryParam("term")
		if term == "" {
			return c.JSON(
				http.StatusOK,
				map[string]interface{}{
					"items": []domain.Note{},
				},
			)
		}

		userID := GetUserID(c, userAPI)
		items, _ := api.SearchItemsByTerm(
			c.Request().Context(),
			userID,
			&domain.NoteRequest{
				Title:   term,
				Content: term,
				RequestPageMeta: domain.RequestPageMeta{
					Limit: 10,
				},
			})

		return c.JSON(
			http.StatusOK,
			map[string]interface{}{
				"items": items,
			},
		)
	}
}

func lastOpenRedirectURL(
	ctx context.Context,
	userID string,
	api ports.NotesService,
	lastOpened ports.LastOpenedService,
) string {
	var tagName string

	lastVisitedID, lErr := lastOpened.GetLastOpened(ctx, domain.LastOpenedTypeNote, userID)
	if lErr == nil && lastVisitedID != "" {
		// get the last opened note to check it exists and to have its ID and tag ID
		existing, nErr := api.GetItem(ctx, userID, lastVisitedID)
		if nErr == nil && existing != nil {
			if len(existing.Tags) > 0 {
				tagName = existing.Tags[0] // Use the first tag if available
			}

			// redirect is found
			return fmt.Sprintf("/notes?note_id=%s&tag=%s", existing.ID, tagName)
		}
	}

	return ""
}
