package sqlite_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/utking/spaces/internal/adapters/db/sqlite"
	"github.com/utking/spaces/internal/adapters/db/unittests"
	"github.com/utking/spaces/internal/application/domain"
)

func TestGetLastOpened(t *testing.T) {
	db, dbErr := unittests.CreateTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := sqlite.NewAdapterWithDB(db)
	userID := "uuid-user-12345"

	// test last note opened
	lastNoteID, err := dbAdapter.GetLastOpened(t.Context(), domain.LastOpenedTypeNote, userID)
	if assert.NoError(t, err) {
		assert.Equal(t, "uuid-note-12345", lastNoteID)
	}

	// test last bookmark opened
	lastBookmarkID, err := dbAdapter.GetLastOpened(t.Context(), domain.LastOpenedTypeBookmark, userID)
	if assert.NoError(t, err) {
		assert.Equal(t, "uuid-bookmark-54321", lastBookmarkID)
	}
}

func TestSetLastOpened(t *testing.T) {
	db, dbErr := unittests.CreateTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := sqlite.NewAdapterWithDB(db)
	userID := "uuid-user-12345"

	// test set last note opened
	err := dbAdapter.SetLastOpened(t.Context(), domain.LastOpenedTypeNote, userID, "uuid-note-67890")
	if assert.NoError(t, err) {
		// Verify that the last opened note was set correctly
		lastNoteID, getErr := dbAdapter.GetLastOpened(t.Context(), domain.LastOpenedTypeNote, userID)
		assert.NoError(t, getErr)
		assert.Equal(t, "uuid-note-67890", lastNoteID)
	}

	// test set last bookmark opened
	err = dbAdapter.SetLastOpened(t.Context(), domain.LastOpenedTypeBookmark, userID, "uuid-bookmark-09876")
	if assert.NoError(t, err) {
		// Verify that the last opened bookmark was set correctly
		lastBookmarkID, getErr := dbAdapter.GetLastOpened(t.Context(), domain.LastOpenedTypeBookmark, userID)
		assert.NoError(t, getErr)
		assert.Equal(t, "uuid-bookmark-09876", lastBookmarkID)
	}
}
