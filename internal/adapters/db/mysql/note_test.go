//go:build mysql
// +build mysql

package mysql_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/utking/spaces/internal/adapters/db/mysql"
	"github.com/utking/spaces/internal/adapters/db/unittests"
	"github.com/utking/spaces/internal/application/domain"
)

func TestGetNoteTags(t *testing.T) {
	db, dbErr := unittests.CreateMySQLTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := mysql.NewAdapterWithDB(db)

	// expected tags for all users. sorted alphabetically
	expectedTags := []string{
		"example2", "example3", "example4", "test", "test2", "test3", "test4",
	}
	// get tags for all users
	tags, err := dbAdapter.GetNoteTags(t.Context(), "")
	if assert.NoError(t, err) {
		assert.Len(t, tags, 7, "Wrong number of tags returned for all users")
		assert.Equal(t, expectedTags, tags, "Wrong tags returned for all users")
	}

	// get tags for a specific user
	userID := "uuid-user-12345"

	// sorted alphabetically expected tags for user
	expectedUserTags := []string{
		"test", "test2", "test3", "test4",
	}

	tags, err = dbAdapter.GetNoteTags(t.Context(), userID)
	if assert.NoError(t, err) {
		assert.Len(t, tags, 4, "Wrong number of tags returned for user")
		assert.Equal(t, expectedUserTags, tags, "Wrong tags returned for user")
	}

	// get tags for a non-existing user
	tags, err = dbAdapter.GetNoteTags(t.Context(), "non-existing-user")
	if assert.NoError(t, err) {
		assert.Empty(t, tags, "Expected no tags for non-existing user")
	}
}

func TestGetNodes(t *testing.T) {
	db, dbErr := unittests.CreateMySQLTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := mysql.NewAdapterWithDB(db)
	userID := "uuid-user-12345"
	filter := &domain.NoteSearchRequest{Tag: "test"}

	// get all notes for a specific user
	items, err := dbAdapter.GetNotes(t.Context(), userID, filter)
	if assert.NoError(t, err) {
		assert.Len(t, items, 1, "Wrong number of notes returned for user")
		assert.Equal(t, "Sample Note Title", items[0].Title)
		assert.Empty(t, items[0].Content, "Content should be empty while listing notes")
	}

	// no notes without a tag filter
	items, err = dbAdapter.GetNotes(t.Context(), userID, nil)
	if assert.NoError(t, err) {
		assert.Empty(t, items, "Expected no notes without a tag filter")
	}

	// search by tag and title substring
	items, err = dbAdapter.GetNotes(t.Context(), userID, &domain.NoteSearchRequest{
		Tag:   "test",
		Title: "Sample",
	})
	if assert.NoError(t, err) {
		assert.Len(t, items, 1, "Wrong number of notes returned for tag and title search")
		assert.Equal(t, "Sample Note Title", items[0].Title)
		assert.Empty(t, items[0].Content, "Content should be empty while listing notes")
	}

	// search by tag and content substring
	items, err = dbAdapter.GetNotes(t.Context(), userID, &domain.NoteSearchRequest{
		Tag:     "test",
		Content: "Sample",
	})
	if assert.NoError(t, err) {
		assert.Len(t, items, 1, "Wrong number of notes returned for tag and content search")
		assert.Equal(t, "Sample Note Title", items[0].Title)
		assert.Empty(t, items[0].Content, "Content should be empty while listing notes")
	}
}

func TestGetNotesCount(t *testing.T) {
	db, dbErr := unittests.CreateMySQLTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := mysql.NewAdapterWithDB(db)
	userID := "uuid-user-12345"

	// get notes count for a specific user. no filter should count all notes
	count, err := dbAdapter.GetNotesCount(t.Context(), userID, nil)
	if assert.NoError(t, err) {
		assert.EqualValues(t, 3, count, "Wrong notes count for user")
	}

	// get notes count for all users
	count, err = dbAdapter.GetNotesCount(t.Context(), "", nil)
	if assert.NoError(t, err) {
		assert.EqualValues(t, 6, count, "Wrong notes count for all users")
	}
}

func TestGetNote(t *testing.T) {
	db, dbErr := unittests.CreateMySQLTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := mysql.NewAdapterWithDB(db)
	userID := "uuid-user-12345"
	noteID := "uuid-note-12345"

	// get a specific note by ID for a user
	note, err := dbAdapter.GetNote(t.Context(), userID, noteID)
	if assert.NoError(t, err) {
		assert.Equal(t, "Sample Note Title", note.Title)
		assert.NotEmpty(t, note.Content, "Content should not be empty for existing note")
	}

	// try to get a non-existing note
	note, err = dbAdapter.GetNote(t.Context(), userID, "non-existing-note")
	if assert.Error(t, err) {
		assert.Nil(t, note, "Expected no note for non-existing ID")
	}

	// try a wrong user ID for a note that exists
	note, err = dbAdapter.GetNote(t.Context(), "wrong-user-id", noteID)
	if assert.Error(t, err) {
		assert.Nil(t, note, "Expected no note for wrong user ID")
	}
}

func TestDeleteNote(t *testing.T) {
	db, dbErr := unittests.CreateMySQLTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := mysql.NewAdapterWithDB(db)
	userID := "uuid-user-12345"
	noteID := "uuid-note-12345"

	// delete a specific note by ID for a user
	err := dbAdapter.DeleteNote(t.Context(), userID, noteID)
	if assert.NoError(t, err) {
		// verify the note is deleted
		note, noteErr := dbAdapter.GetNote(t.Context(), userID, noteID)
		assert.Error(t, noteErr, "Expected error when getting deleted note")
		assert.Nil(t, note, "Expected no note after deletion")
	}

	// try to delete a non-existing note
	assert.NoError(t, dbAdapter.DeleteNote(t.Context(), userID, "non-existing-note"))

	// try to delete a note with a wrong user ID
	err = dbAdapter.DeleteNote(t.Context(), "wrong-user-id", noteID)
	assert.NoError(t, err)
}

func TestCreateNote(t *testing.T) {
	db, dbErr := unittests.CreateMySQLTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := mysql.NewAdapterWithDB(db)
	userID := "uuid-user-12345"

	// create a new note
	note := &domain.Note{
		Title:   "New Note Title",
		Content: "This is the content of the new note.",
		Tags:    []string{"test", "example"},
	}

	id, err := dbAdapter.CreateNote(t.Context(), userID, note)
	if assert.NoError(t, err) {
		assert.NotEmpty(t, id, "Expected non-empty ID for created note")

		// verify the note is created
		createdNote, getErr := dbAdapter.GetNote(t.Context(), userID, id)
		if assert.NoError(t, getErr) {
			assert.Equal(t, note.Title, createdNote.Title)
			assert.Equal(t, note.Content, createdNote.Content)
			assert.ElementsMatch(t, note.Tags, createdNote.Tags)
		}
	}

	// try to create a note with an empty title
	note.Title = ""
	id, err = dbAdapter.CreateNote(t.Context(), userID, note)
	if assert.Error(t, err, "Expected error when creating a note with empty title") {
		assert.Empty(t, id, "Expected empty ID when creation fails")
	}

	// try creating a note without tags. validation must fail
	note.Title = "Valid Note Title"
	note.Tags = nil
	id, err = dbAdapter.CreateNote(t.Context(), userID, note)
	if assert.Error(t, err, "Expected error when creating a note without tags") {
		assert.Empty(t, id, "Expected empty ID when creation fails without tags")
	}

	// test duplicate title for the same user
	note.Title = "Duplicate Note Title"
	note.Tags = []string{"test", "duplicate"}
	id, err = dbAdapter.CreateNote(t.Context(), userID, note)
	if assert.NoError(t, err) {
		assert.NotEmpty(t, id, "Expected non-empty ID for created note with duplicate title")

		// verify the note is created
		createdNoteID, getErr := dbAdapter.CreateNote(t.Context(), userID, note)
		if assert.Error(t, getErr, "Expected error when creating a note with duplicate title") {
			assert.Empty(t, createdNoteID, "Expected empty ID when creating a note with duplicate title")
			assert.Equal(t, "note with this title already exists", getErr.Error(),
				"Expected specific error message for duplicate title")
		}
	}
}

func TestUpdateNote(t *testing.T) {
	db, dbErr := unittests.CreateMySQLTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := mysql.NewAdapterWithDB(db)
	userID := "uuid-user-12345"
	noteID := "uuid-note-12345"

	// update an existing note
	updateReq := &domain.Note{
		Title:   "Updated Note Title",
		Content: "Updated content for the note.",
		Tags:    []string{"updated", "example"},
	}

	// create a note first to update it
	_, createErr := dbAdapter.CreateNote(t.Context(), userID, updateReq)
	if assert.NoError(t, createErr) {
		rowsAffected, err := dbAdapter.UpdateNote(t.Context(), userID, noteID, updateReq)
		if assert.Error(t, err) {
			assert.EqualValues(t, 0, rowsAffected,
				"Expected no rows to be affected when updating a note that does not exist")

			// verify the note is not updated
			updatedNote, getErr := dbAdapter.GetNote(t.Context(), userID, noteID)
			if assert.NoError(t, getErr) {
				assert.NotEqual(t, updateReq.Title, updatedNote.Title)
				assert.NotEqual(t, updateReq.Content, updatedNote.Content)
				assert.NotElementsMatch(t, updateReq.Tags, updatedNote.Tags)
			}
		}
	}

	// try to update a non-existing note. it should fail silently
	rowsAffected, err := dbAdapter.UpdateNote(t.Context(), userID, "non-existing-note", updateReq)
	if assert.NoError(t, err) {
		assert.EqualValues(t, 1, rowsAffected, "Expected one row to be fake-updated for non-existing note")
	}

	// try to update a note with an empty title
	updateReq.Title = ""
	rowsAffected, err = dbAdapter.UpdateNote(t.Context(), userID, noteID, updateReq)
	if assert.Error(t, err) {
		assert.EqualValues(t, 0, rowsAffected,
			"Expected no rows to be affected when updating a note with empty title")
	}

	// try updating a note without tags. validation must fail
	updateReq.Title = "Valid Update Title"
	updateReq.Tags = nil
	rowsAffected, err = dbAdapter.UpdateNote(t.Context(), userID, noteID, updateReq)
	if assert.Error(t, err) {
		assert.EqualValues(t, 0,
			rowsAffected,
			"Expected no rows to be affected when updating a note without tags")
	}

	// try updating a note with a wrong user ID
	rowsAffected, err = dbAdapter.UpdateNote(t.Context(), "wrong-user-id", noteID, updateReq)
	if assert.Error(t, err) {
		assert.EqualValues(t, 0, rowsAffected, "Expected no rows to be affected when updating with wrong user ID")
	}

	// verify the note is not updated
	updatedNote, getErr := dbAdapter.GetNote(t.Context(), userID, noteID)
	if assert.NoError(t, getErr) {
		assert.NotEqual(t, updateReq.Title, updatedNote.Title,
			"Expected note title to remain unchanged after failed update with wrong user ID")
		assert.NotEqual(t, updateReq.Content, updatedNote.Content,
			"Expected note content to remain unchanged after failed update with wrong user ID")
		assert.NotEqual(t, updateReq.Tags, updatedNote.Tags,
			"Expected note tags to remain unchanged after failed update with wrong user ID")
	}
}

func TestGetNotesMap(t *testing.T) {
	db, dbErr := unittests.CreateMySQLTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := mysql.NewAdapterWithDB(db)
	userID := "uuid-user-12345"

	// get notes map for a specific user
	notesMap, err := dbAdapter.GetNotesMap(t.Context(), userID, nil)
	if assert.NoError(t, err) {
		assert.Len(t, notesMap, 3, "Wrong number of notes returned for user")
		// must have empty ID. non-empty Title and Content
		for _, note := range notesMap {
			assert.Empty(t, note.ID, "Expected non-empty note ID")
			assert.NotEmpty(t, note.Title, "Expected non-empty note title")
			assert.NotEmpty(t, note.Content, "Expected non-empty note content")
		}
	}

	// get notes map for all users must return no notes
	allNotesMap, err := dbAdapter.GetNotesMap(t.Context(), "", nil)
	if assert.NoError(t, err) {
		assert.Empty(t, allNotesMap, "Expected no notes for all users without a user ID")
	}
}

func TestSearchNotesByTerm(t *testing.T) {
	db, dbErr := unittests.CreateMySQLTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := mysql.NewAdapterWithDB(db)
	userID := "uuid-user-12345"
	searchReq := &domain.NoteRequest{
		Title:   "Sample",
		Content: "Sample",
	}

	// search notes by term for a specific user
	items, err := dbAdapter.SearchNotesByTerm(t.Context(), userID, searchReq)
	if assert.NoError(t, err) {
		assert.Len(t, items, 2, "Wrong number of notes returned for search term")
	}

	searchReq.Title = "NonExistingTitle"
	searchReq.Content = "NonExistingContent"

	// search with a term that does not match any notes
	items, err = dbAdapter.SearchNotesByTerm(t.Context(), userID, searchReq)
	if assert.NoError(t, err) {
		assert.Empty(t, items, "Expected no notes for non-existing search term")
	}
}
