package sqlite_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gogs.utking.net/utking/spaces/internal/adapters/db/sqlite"
	"gogs.utking.net/utking/spaces/internal/adapters/db/unittests"
	"gogs.utking.net/utking/spaces/internal/application/domain"
)

func TestGetTags(t *testing.T) {
	db, dbErr := unittests.CreateTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := sqlite.NewAdapterWithDB(db)
	// get tags for all users
	tags, err := dbAdapter.GetBookmarkTags(t.Context(), "")
	if assert.NoError(t, err, "GetBookmarkTags error") {
		assert.Len(t, tags, 5, "Wrong number of tags returned")
	}

	// get tags for a specific user
	userID := "uuid-u-3456-7890-1234"
	tags, err = dbAdapter.GetBookmarkTags(t.Context(), userID)
	if assert.NoError(t, err, "GetBookmarkTags error") {
		assert.Len(t, tags, 3, "Wrong number of tags returned for user")
	}
}

func TestGetBookmarks(t *testing.T) {
	db, dbErr := unittests.CreateTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := sqlite.NewAdapterWithDB(db)
	// get bookmarks for a specific user, with no filters
	userID := "uuid-u-3456-7890-1234"
	bookmarks, err := dbAdapter.GetBookmarks(t.Context(), userID, nil)
	if assert.NoError(t, err, "GetBookmarks error") {
		assert.Len(t, bookmarks, 2, "Wrong number of bookmarks returned for user")
		for _, bookmark := range bookmarks {
			assert.Equal(t, userID, bookmark.UserID, "Bookmark title mismatch")
			assert.NotEmpty(t, bookmark.Title, "First bookmark title should not be empty")
			assert.NotEmpty(t, bookmark.URL, "First bookmark URL should not be empty")
		}
	}

	// get bookmarks for a wrong user ID
	userID = "wrong-uuid"
	bookmarks, err = dbAdapter.GetBookmarks(t.Context(), userID, nil)
	if assert.NoError(t, err, "GetBookmarks error") {
		assert.Empty(t, bookmarks)
	}

	// get bookmarks for a specific user, with filters
	userID = "uuid-u-3456-7890-1234"
	req := &domain.BookmarkSearchRequest{
		Title: "Fourth Bookmark",
	}
	bookmarks, err = dbAdapter.GetBookmarks(t.Context(), userID, req)
	if assert.NoError(t, err, "GetBookmarks error") {
		assert.Len(t, bookmarks, 1, "Wrong number of bookmarks returned for user with filters")
	}
}

func TestGetBookmarksCount(t *testing.T) {
	db, dbErr := unittests.CreateTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := sqlite.NewAdapterWithDB(db)
	// get bookmarks count for a specific user, with no filters
	userID := "uuid-u-3456-7890-1234"
	count, err := dbAdapter.GetBookmarksCount(t.Context(), userID, nil)
	if assert.NoError(t, err, "GetBookmarksCount error") {
		assert.Equal(t, int64(2), count, "Wrong number of bookmarks returned for user")
	}

	// get for a wrong user ID
	userID = "wrong-uuid"
	count, err = dbAdapter.GetBookmarksCount(t.Context(), userID, nil)
	if assert.NoError(t, err, "GetBookmarksCount error") {
		assert.Equal(t, int64(0), count, "Wrong number of bookmarks returned for user")
	}
}

func TestGetOneBookmark(t *testing.T) {
	db, dbErr := unittests.CreateTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := sqlite.NewAdapterWithDB(db)
	// get one bookmark for a specific user
	userID := "uuid-u-3456-7890-1234"
	bookmarkID := "uuid-4567-8901-2345"
	bookmark, err := dbAdapter.GetBookmark(t.Context(), userID, bookmarkID)
	if assert.NoError(t, err, "GetOneBookmark error") {
		assert.Equal(t, bookmarkID, bookmark.ID, "Bookmark ID mismatch")
		assert.NotEmpty(t, bookmark.Title, "Bookmark title should not be empty")
		assert.NotEmpty(t, bookmark.URL, "Bookmark URL should not be empty")
		assert.NotEmpty(t, bookmark.Tags, "Bookmark tags should not be empty")
	}

	// get one bookmark for a wrong user ID
	bookmarkID = "uuid-b-1234-5678-xxxx"
	bookmark, err = dbAdapter.GetBookmark(t.Context(), userID, bookmarkID)
	if assert.Error(t, err) {
		assert.Nil(t, bookmark)
	}
}

func TestDeleteBookmark(t *testing.T) {
	db, dbErr := unittests.CreateTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := sqlite.NewAdapterWithDB(db)
	// delete a bookmark for a specific user
	userID := "uuid-u-3456-7890-1234"
	bookmarkID := "uuid-4567-8901-2345"
	err := dbAdapter.DeleteBookmark(t.Context(), userID, bookmarkID)
	if assert.NoError(t, err, "DeleteBookmark error") {
		// verify that the bookmark is deleted
		bookmark, getErr := dbAdapter.GetBookmark(t.Context(), userID, bookmarkID)
		assert.Error(t, getErr, "Expected error when getting deleted bookmark")
		assert.Nil(t, bookmark, "Bookmark should be nil after deletion")
	}

	// try deleting it again - should not cause an error
	err = dbAdapter.DeleteBookmark(t.Context(), userID, bookmarkID)
	if assert.NoError(t, err, "DeleteBookmark error on non-existing bookmark") {
		// verify that the bookmark is still deleted
		bookmark, getErr := dbAdapter.GetBookmark(t.Context(), userID, bookmarkID)
		if assert.Error(t, getErr, "Expected error when getting deleted bookmark") {
			assert.Nil(t, bookmark, "Bookmark should be nil after deletion")
		}
	}
}

func TestUpdateBookmark(t *testing.T) {
	db, dbErr := unittests.CreateTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := sqlite.NewAdapterWithDB(db)
	// update a bookmark for a specific user
	userID := "uuid-u-3456-7890-1234"
	bookmarkID := "uuid-4567-8901-2345"
	req := &domain.Bookmark{
		ID:    bookmarkID,
		Title: "Updated Bookmark Title",
		URL:   "https://updated.example.com",
		Tags:  []string{"updated", "example"},
	}

	rowsAffected, err := dbAdapter.UpdateBookmark(t.Context(), userID, bookmarkID, req)
	if assert.NoError(t, err, "UpdateBookmark error") {
		assert.Equal(t, int64(1), rowsAffected, "Wrong number of rows affected")
	}

	// verify that the bookmark is updated
	bookmark, err := dbAdapter.GetBookmark(t.Context(), userID, bookmarkID)
	if assert.NoError(t, err, "GetBookmark after update error") {
		assert.Equal(t, req.Title, bookmark.Title, "Bookmark title mismatch after update")
		assert.Equal(t, req.URL, bookmark.URL, "Bookmark URL mismatch after update")
		assert.ElementsMatch(t, req.Tags, bookmark.Tags, "Bookmark tags mismatch after update")
	}

	// try updating a bookmark with invalid data
	reqInvalid := &domain.Bookmark{
		ID:    bookmarkID,
		Title: "", // empty title
		URL:   "https://updated.example.com",
		Tags:  []string{"updated", "example"},
	}

	_, err = dbAdapter.UpdateBookmark(t.Context(), userID, bookmarkID, reqInvalid)
	if assert.Error(t, err, "Expected error when updating bookmark with invalid data") {
		assert.Equal(t, "title cannot be empty;", err.Error(), "Error message mismatch for invalid update")
	}
}

func TestCreateBookmark(t *testing.T) {
	db, dbErr := unittests.CreateTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := sqlite.NewAdapterWithDB(db)
	// create a new bookmark for a specific user
	userID := "uuid-u-3456-7890-1234"
	req := &domain.Bookmark{
		UserID: userID,
		Title:  "New Bookmark",
		URL:    "https://new.example.com",
		Tags:   []string{"new", "example"},
	}

	id, err := dbAdapter.CreateBookmark(t.Context(), userID, req)
	if assert.NoError(t, err, "CreateBookmark error") {
		assert.NotEmpty(t, id, "Bookmark ID should not be empty after creation")
	}

	// verify that the bookmark is created
	bookmark, err := dbAdapter.GetBookmark(t.Context(), userID, id)
	if assert.NoError(t, err, "GetBookmark after create error") {
		assert.Equal(t, req.Title, bookmark.Title, "Bookmark title mismatch after create")
		assert.Equal(t, req.URL, bookmark.URL, "Bookmark URL mismatch after create")
		assert.ElementsMatch(t, req.Tags, bookmark.Tags, "Bookmark tags mismatch after create")
	}

	// try creating a bookmark with invalid data
	reqInvalid := &domain.Bookmark{
		UserID: userID,
		Title:  "", // empty title
		URL:    "https://new.example.com",
		Tags:   []string{"new", "example"},
	}

	idInvalid, err := dbAdapter.CreateBookmark(t.Context(), userID, reqInvalid)
	if assert.Error(t, err, "Expected error when creating bookmark with invalid data") {
		assert.Empty(t, idInvalid, "Bookmark ID should be empty for invalid creation")
		assert.Equal(t, "title cannot be empty;", err.Error(), "Error message mismatch for invalid creation")
	}
}

func TestGetBookmarksMap(t *testing.T) {
	db, dbErr := unittests.CreateTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := sqlite.NewAdapterWithDB(db)
	// get bookmarks map for a specific user
	userID := "uuid-u-3456-7890-1234"
	bookmarksMap, err := dbAdapter.GetBookmarksMap(t.Context(), userID, nil)
	if assert.NoError(t, err, "GetBookmarksMap error") {
		assert.Len(t, bookmarksMap, 2, "Wrong number of bookmarks returned for user")
		for _, bookmark := range bookmarksMap {
			assert.Empty(t, bookmark.UserID, "Bookmark UserID should be empty in map")
			assert.Empty(t, bookmark.ID, "Bookmark ID should be empty")
			assert.NotEmpty(t, bookmark.Title, "Bookmark title should not be empty")
			assert.NotEmpty(t, bookmark.URL, "Bookmark URL should not be empty")
		}
	}

	// get bookmarks map for a wrong user ID
	userID = "wrong-uuid"
	bookmarksMap, err = dbAdapter.GetBookmarksMap(t.Context(), userID, nil)
	if assert.NoError(t, err, "GetBookmarksMap error") {
		assert.Empty(t, bookmarksMap)
	}

	// empty userID should return no bookmarks
	bookmarksMap, err = dbAdapter.GetBookmarksMap(t.Context(), "", nil)
	if assert.NoError(t, err, "GetBookmarksMap error with empty userID") {
		assert.Empty(t, bookmarksMap, "Expected no bookmarks for empty userID")
	}
}

func TestSearchBookmarksByTerm(t *testing.T) {
	db, dbErr := unittests.CreateTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := sqlite.NewAdapterWithDB(db)
	// search bookmarks by term for a specific user
	userID := "uuid-u-3456-7890-1234"
	term := "Fourth"
	req := &domain.BookmarkSearchRequest{
		Title: term,
		URL:   term,
		RequestPageMeta: domain.RequestPageMeta{
			Limit: 2,
		},
	}
	bookmarks, err := dbAdapter.SearchBookmarksByTerm(t.Context(), userID, req)
	if assert.NoError(t, err, "SearchBookmarksByTerm error") {
		assert.Len(t, bookmarks, 1, "Wrong number of bookmarks returned for search term")
		assert.Equal(t, "Fourth Bookmark", bookmarks[0].Title, "Bookmark title mismatch for search term")
	}

	// search with a term that does not match any bookmarks
	term = "NonExistingTerm"
	req = &domain.BookmarkSearchRequest{
		Title: term,
		URL:   term,
		RequestPageMeta: domain.RequestPageMeta{
			Limit: 2,
		},
	}
	bookmarks, err = dbAdapter.SearchBookmarksByTerm(t.Context(), userID, req)
	if assert.NoError(t, err, "SearchBookmarksByTerm error") {
		assert.Empty(t, bookmarks, "Expected no bookmarks for non-existing search term")
	}
}
