package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"gogs.utking.net/utking/spaces/internal/adapters/web/go_echo/helpers"
	"gogs.utking.net/utking/spaces/internal/application/domain"
	"gogs.utking.net/utking/spaces/internal/ports"
)

// getBookmarksWrapper is a wrapper for the bookmarks handler.
func getBookmarksWrapper(
	api ports.BookmarkService,
	userAPI ports.UsersService,
	lastOpened ports.LastOpenedService,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			userID = GetUserID(c, userAPI)
			req    = new(domain.BookmarkSearchRequest)
			items  = make([]domain.Bookmark, 0)
			err    error
		)

		_ = c.Bind(req)

		if req.Tag == "" {
			// check existing last opened tag
			lastOpenedTag, _ := lastOpened.GetLastOpened(
				c.Request().Context(), domain.LastOpenedTypeBookmark, userID)

			if lastOpenedTag != "" {
				return c.Redirect(
					http.StatusSeeOther,
					fmt.Sprintf("/bookmarks?tag=%s", lastOpenedTag),
				)
			}
		}

		tags, _ := api.GetTags(c.Request().Context(), userID)
		// if no tag is specified, do not load bookmark items
		if req.Tag != "" {
			items, err = api.GetItems(c.Request().Context(), userID, req)
			// save last opened tag
			_ = lastOpened.SetLastOpened(
				c.Request().Context(), domain.LastOpenedTypeBookmark, userID, req.Tag)
		}

		return c.Render(
			http.StatusOK,
			"bookmarks/index.html",
			map[string]interface{}{
				"Title":      "Bookmarks",
				"Items":      items,
				"Tags":       tags,
				"Query":      req,
				"ItemsCount": len(items),
				"TagsCount":  len(tags),
				"Error":      helpers.ErrorMessage(err),
			},
		)
	}
}

// postBookmarkCreateWrapper is a wrapper for the bookmark creation handler.
// Returns a JSON response with the created bookmark ID or an error message.
func postBookmarkCreateWrapper(
	api ports.BookmarkService,
	userAPI ports.UsersService,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			userID = GetUserID(c, userAPI)
			req    = new(domain.Bookmark)
		)

		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"Error": "Invalid request"})
		}

		id, err := api.Create(c.Request().Context(), userID, req)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"Error": helpers.ErrorMessage(err)})
		}

		return c.JSON(http.StatusOK, map[string]string{"ID": id})
	}
}

// deleteBookmarkWrapper is a wrapper for the bookmark deletion handler.
func deleteBookmarkWrapper(
	api ports.BookmarkService,
	userAPI ports.UsersService,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := GetUserID(c, userAPI)
		id := c.Param("id")

		if id == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"Error": "Bookmark ID is required"})
		}

		if err := api.Delete(c.Request().Context(), userID, id); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"Error": helpers.ErrorMessage(err)})
		}

		return c.NoContent(http.StatusNoContent)
	}
}

// getBookmarkEditWrapper is a wrapper for the bookmark edit handler.
func getBookmarkEditWrapper(
	api ports.BookmarkService,
	userAPI ports.UsersService,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			userID = GetUserID(c, userAPI)
			id     = c.Param("id")
			code   = http.StatusOK
		)

		tags, _ := api.GetTags(c.Request().Context(), userID)
		item, err := api.GetItem(c.Request().Context(), userID, id)

		if err != nil {
			if item == nil {
				code = http.StatusNotFound
			} else {
				code = http.StatusInternalServerError
			}
		}

		return c.Render(
			code,
			"bookmarks/edit.html",
			map[string]interface{}{
				"Title": "Edit Bookmark",
				"Item":  item,
				"Tags":  tags,
				"Error": helpers.ErrorMessage(err),
			},
		)
	}
}

// putBookmarkEditWrapper is a wrapper for the bookmark edit handler.
func putBookmarkEditWrapper(
	api ports.BookmarkService,
	userAPI ports.UsersService,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			userID = GetUserID(c, userAPI)
			id     = c.Param("id")
			req    = new(domain.Bookmark)
		)

		_ = c.Bind(req)

		if _, err := api.Update(c.Request().Context(), userID, id, req); err != nil {
			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"Error": helpers.ErrorMessage(err),
				},
			)
		}

		return c.NoContent(http.StatusNoContent)
	}
}

// getExportBookmarksWrapper is a wrapper for the export bookmarks handler.
// it compiles a map of the user's bookmarks with their tags, exporting them to a JSON file,
// and returns a downloadable file to the user.
func getExportBookmarksWrapper(
	api ports.BookmarkService,
	userAPI ports.UsersService,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		const fileName = "bookmarks_export.json"
		var eFile *os.File

		items, err := api.GetItemsMap(c.Request().Context(), GetUserID(c, userAPI), nil)
		if err == nil {
			eFile, err = saveBookmarksToFile(items)
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

// saveBookmarksToFile is a helper function that saves the bookmarks data to a file.
func saveBookmarksToFile(
	items []domain.Bookmark,
) (*os.File, error) {
	data, err := json.Marshal(items)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal bookmarks data: %w", err)
	}

	eFile, err := os.CreateTemp("", "bookmarks_export*.json")
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

// getSearchBookmarksWrapper is a wrapper for the search bookmarks handler.
func getSearchBookmarksWrapper(
	api ports.BookmarkService,
	userAPI ports.UsersService,
) echo.HandlerFunc {
	type BookmarkItem struct {
		ID   string `json:"id"`
		Text string `json:"text"`
	}

	return func(c echo.Context) error {
		term := c.QueryParam("term")
		if term == "" {
			return c.JSON(
				http.StatusOK,
				map[string]interface{}{
					"items": []BookmarkItem{},
				},
			)
		}

		userID := GetUserID(c, userAPI)

		items, _ := api.SearchItemsByTerm(c.Request().Context(), userID, &domain.BookmarkSearchRequest{
			Title: term,
			URL:   term,
			RequestPageMeta: domain.RequestPageMeta{
				Limit: 10,
			},
		})

		// filter items based on the search term
		filteredItems := make([]BookmarkItem, len(items))
		for idx, item := range items {
			filteredItems[idx] = BookmarkItem{
				ID:   item.URL,
				Text: item.Title,
			}
		}

		return c.JSON(
			http.StatusOK,
			map[string]interface{}{
				"items": filteredItems,
			},
		)
	}
}
