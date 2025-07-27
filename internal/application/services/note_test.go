package services_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/utking/spaces/internal/application/domain"
	"github.com/utking/spaces/internal/application/services"
	"github.com/utking/spaces/internal/ports"
)

func TestGetNoteTags(t *testing.T) {
	tagsInDB := []string{"tag1", "tag2"}

	dbPort := ports.NewMockDBPort(t)
	dbPort.On("GetNoteTags", mock.Anything, "some-user-id").Return(tagsInDB, nil)

	svc := services.NewNotesService(dbPort)

	tags, err := svc.GetTags(t.Context(), "some-user-id")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(tags) != len(tagsInDB) {
		t.Fatalf("expected some tags, got none")
	}

	dbPort.AssertExpectations(t)
}

func TestGetNoteTagsError(t *testing.T) {
	dbPort := ports.NewMockDBPort(t)
	dbPort.On("GetNoteTags", mock.Anything, "some-user-id").Return([]string{}, errors.New("some error"))

	svc := services.NewNotesService(dbPort)

	_items, err := svc.GetTags(t.Context(), "some-user-id")
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if len(_items) != 0 {
		t.Fatalf("expected no tags, got %d", len(_items))
	}

	dbPort.AssertExpectations(t)
}

func TestGetItems(t *testing.T) {
	itemsInDB := []domain.Note{
		{ID: "1", Title: "Note 1", Content: "Content 1", Tags: []string{"tag1"}},
		{ID: "2", Title: "Note 2", Content: "Content 2", Tags: []string{"tag2"}},
	}

	dbPort := ports.NewMockDBPort(t)
	dbPort.On("GetNotes", mock.Anything, "some-user-id", mock.Anything).Return(itemsInDB, nil)

	svc := services.NewNotesService(dbPort)

	items, err := svc.GetItems(t.Context(), "some-user-id", &domain.NoteSearchRequest{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(items) != len(itemsInDB) {
		t.Fatalf("expected %d items, got %d", len(itemsInDB), len(items))
	}

	dbPort.AssertExpectations(t)
}

func TestGetItemsError(t *testing.T) {
	itemsInDB := []domain.Note{}

	dbPort := ports.NewMockDBPort(t)
	dbPort.On("GetNotes", mock.Anything, "some-user-id", mock.Anything).
		Return(itemsInDB, errors.New("some error"))

	svc := services.NewNotesService(dbPort)

	_items, err := svc.GetItems(t.Context(), "some-user-id", &domain.NoteSearchRequest{})
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if len(_items) != 0 {
		t.Fatalf("expected no items, got %d", len(_items))
	}

	dbPort.AssertExpectations(t)
}

func TestGetCount(t *testing.T) {
	countInDB := int64(5)

	dbPort := ports.NewMockDBPort(t)
	dbPort.On("GetNotesCount", mock.Anything, "some-user-id", mock.Anything).Return(countInDB, nil)

	svc := services.NewNotesService(dbPort)

	count, err := svc.GetCount(t.Context(), "some-user-id", nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if count != countInDB {
		t.Fatalf("expected count %d, got %d", countInDB, count)
	}

	dbPort.AssertExpectations(t)
}

func TestGetCountError(t *testing.T) {
	dbPort := ports.NewMockDBPort(t)
	dbPort.On("GetNotesCount", mock.Anything, "some-user-id", mock.Anything).
		Return(int64(0), errors.New("some error"))

	svc := services.NewNotesService(dbPort)

	count, err := svc.GetCount(t.Context(), "some-user-id", nil)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if count != 0 {
		t.Fatalf("expected count 0, got %d", count)
	}

	dbPort.AssertExpectations(t)
}

func TestGetItem(t *testing.T) {
	itemInDB := &domain.Note{
		ID:      "1",
		Title:   "Note 1",
		Content: "Content 1",
		Tags:    []string{"tag1"},
	}

	dbPort := ports.NewMockDBPort(t)
	dbPort.On("GetNote", mock.Anything, "some-user-id", "1").Return(itemInDB, nil)

	svc := services.NewNotesService(dbPort)

	item, err := svc.GetItem(t.Context(), "some-user-id", "1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if item.ID != itemInDB.ID {
		t.Fatalf("expected item ID %s, got %s", itemInDB.ID, item.ID)
	}

	dbPort.AssertExpectations(t)
}

func TestGetItemError(t *testing.T) {
	dbPort := ports.NewMockDBPort(t)
	dbPort.On("GetNote", mock.Anything, "some-user-id", "1").Return(nil, errors.New("some error"))

	svc := services.NewNotesService(dbPort)

	item, err := svc.GetItem(t.Context(), "some-user-id", "1")
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if item != nil {
		t.Fatalf("expected no item, got %v", item)
	}

	dbPort.AssertExpectations(t)
}

func TestCreateNote(t *testing.T) {
	itemToCreate := &domain.Note{
		Title:   "New Note",
		Content: "New Content",
		Tags:    []string{"tag1"},
	}

	expectedID := "new-note-id"

	dbPort := ports.NewMockDBPort(t)
	dbPort.On("CreateNote", mock.Anything, "some-user-id", itemToCreate).Return(expectedID, nil)

	svc := services.NewNotesService(dbPort)

	id, err := svc.Create(t.Context(), "some-user-id", itemToCreate)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if id != expectedID {
		t.Fatalf("expected ID %s, got %s", expectedID, id)
	}

	dbPort.AssertExpectations(t)
}

// Must fail on note.Validate and not call CreateNote in the DB port.
func TestCreateNoteErrorNoTags(t *testing.T) {
	itemToCreate := &domain.Note{
		Title:   "New Note",
		Content: "New Content",
		Tags:    []string{},
	}

	dbPort := ports.NewMockDBPort(t)
	svc := services.NewNotesService(dbPort)

	id, err := svc.Create(t.Context(), "some-user-id", itemToCreate)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if id != "" {
		t.Fatalf("expected empty ID, got %s", id)
	}

	dbPort.AssertNotCalled(t, "CreateNote", mock.Anything, "some-user-id", itemToCreate)
}

// Must fail on note.Validate and not call CreateNote in the DB port.
func TestCreateNoteErrorNoTitle(t *testing.T) {
	itemToCreate := &domain.Note{
		Title:   "",
		Content: "New Content",
		Tags:    []string{"tag1"},
	}

	dbPort := ports.NewMockDBPort(t)
	svc := services.NewNotesService(dbPort)

	id, err := svc.Create(t.Context(), "some-user-id", itemToCreate)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if id != "" {
		t.Fatalf("expected empty ID, got %s", id)
	}

	dbPort.AssertNotCalled(t, "CreateNote", mock.Anything, "some-user-id", itemToCreate)
}

// Must fail on note.Validate and not call CreateNote in the DB port.
func TestCreateNoteErrorNoTooLongTitle(t *testing.T) {
	itemToCreate := &domain.Note{
		Title:   "a" + strings.Repeat("x", 128), // Exceeds 255 characters
		Content: "New Content",
		Tags:    []string{"tag1"},
	}

	dbPort := ports.NewMockDBPort(t)
	svc := services.NewNotesService(dbPort)

	id, err := svc.Create(t.Context(), "some-user-id", itemToCreate)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if id != "" {
		t.Fatalf("expected empty ID, got %s", id)
	}

	dbPort.AssertNotCalled(t, "CreateNote", mock.Anything, "some-user-id", itemToCreate)
}

func TestDeleteNote(t *testing.T) {
	dbPort := ports.NewMockDBPort(t)
	dbPort.On("DeleteNote", mock.Anything, "some-user-id", "1").Return(nil)

	svc := services.NewNotesService(dbPort)

	err := svc.Delete(t.Context(), "some-user-id", "1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	dbPort.AssertExpectations(t)
}

func TestDeleteNoteErrorInDB(t *testing.T) {
	dbPort := ports.NewMockDBPort(t)
	dbPort.On("DeleteNote", mock.Anything, "some-user-id", "1").Return(errors.New("some error"))

	svc := services.NewNotesService(dbPort)

	err := svc.Delete(t.Context(), "some-user-id", "1")
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	dbPort.AssertExpectations(t)
}

func TestDeleteNoteErrorEmptyID(t *testing.T) {
	dbPort := ports.NewMockDBPort(t)
	svc := services.NewNotesService(dbPort)

	err := svc.Delete(t.Context(), "some-user-id", "")
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	dbPort.AssertNotCalled(t, "DeleteNote", mock.Anything, "some-user-id", "")
}

func TestUpdateNote(t *testing.T) {
	itemToUpdate := &domain.Note{
		Title:   "Updated Note",
		Content: "Updated Content",
		Tags:    []string{"tag1"},
	}

	dbPort := ports.NewMockDBPort(t)
	dbPort.On("UpdateNote", mock.Anything, "some-user-id", "1", itemToUpdate).Return(int64(1), nil)

	svc := services.NewNotesService(dbPort)

	rowsAffected, err := svc.Update(t.Context(), "some-user-id", "1", itemToUpdate)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if rowsAffected != 1 {
		t.Fatalf("expected 1 row affected, got %d", rowsAffected)
	}

	dbPort.AssertExpectations(t)
}

func TestUpdateNoteErrorInDB(t *testing.T) {
	itemToUpdate := &domain.Note{
		Title:   "Updated Note",
		Content: "Updated Content",
		Tags:    []string{"tag1"},
	}

	dbPort := ports.NewMockDBPort(t)
	dbPort.On("UpdateNote", mock.Anything, "some-user-id", "1", itemToUpdate).Return(int64(0), errors.New("some error"))

	svc := services.NewNotesService(dbPort)

	rowsAffected, err := svc.Update(t.Context(), "some-user-id", "1", itemToUpdate)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}

	dbPort.AssertExpectations(t)
}

func TestUpdateNoteErrorEmptyID(t *testing.T) {
	itemToUpdate := &domain.Note{
		Title:   "Updated Note",
		Content: "Updated Content",
		Tags:    []string{"tag1"},
	}

	dbPort := ports.NewMockDBPort(t)
	svc := services.NewNotesService(dbPort)

	rowsAffected, err := svc.Update(t.Context(), "some-user-id", "", itemToUpdate)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}

	dbPort.AssertNotCalled(t, "UpdateNote", mock.Anything, "some-user-id", "", itemToUpdate)
}

func TestUpdateNoteErrorNoTags(t *testing.T) {
	itemToUpdate := &domain.Note{
		Title:   "Updated Note",
		Content: "Updated Content",
		Tags:    []string{},
	}

	dbPort := ports.NewMockDBPort(t)
	svc := services.NewNotesService(dbPort)

	rowsAffected, err := svc.Update(t.Context(), "some-user-id", "1", itemToUpdate)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}

	dbPort.AssertNotCalled(t, "UpdateNote", mock.Anything, "some-user-id", "1", itemToUpdate)
}

func TestUpdateNoteErrorNoTitle(t *testing.T) {
	itemToUpdate := &domain.Note{
		Title:   "",
		Content: "Updated Content",
		Tags:    []string{"tag1"},
	}

	dbPort := ports.NewMockDBPort(t)
	svc := services.NewNotesService(dbPort)

	rowsAffected, err := svc.Update(t.Context(), "some-user-id", "1", itemToUpdate)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}

	dbPort.AssertNotCalled(t, "UpdateNote", mock.Anything, "some-user-id", "1", itemToUpdate)
}

func TestUpdateNoteErrorNoTooLongTitle(t *testing.T) {
	itemToUpdate := &domain.Note{
		Title:   "a" + strings.Repeat("x", 128), // Exceeds 255 characters
		Content: "Updated Content",
		Tags:    []string{"tag1"},
	}

	dbPort := ports.NewMockDBPort(t)
	svc := services.NewNotesService(dbPort)

	rowsAffected, err := svc.Update(t.Context(), "some-user-id", "1", itemToUpdate)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}

	dbPort.AssertNotCalled(t, "UpdateNote", mock.Anything, "some-user-id", "1", itemToUpdate)
}

func TestNotesGetItemsMap(t *testing.T) {
	items := []domain.Note{
		{ID: "1", Title: "Note 1", Content: "", Tags: []string{"tag1", "tag2"}},
		{ID: "2", Title: "Note 2", Content: "Content 2", Tags: []string{"tag2"}},
	}

	dbPort := ports.NewMockDBPort(t)
	dbPort.On("GetNotesMap", mock.Anything, "some-user-id", mock.Anything).Return(items, nil)

	svc := services.NewNotesService(dbPort)

	itemsMap, err := svc.GetItemsMap(t.Context(), "some-user-id", &domain.NoteSearchRequest{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(itemsMap) != len(items) {
		t.Fatalf("expected %d items, got %d", len(items), len(itemsMap))
	}

	dbPort.AssertExpectations(t)
}

func TestNotesGetItemsMapError(t *testing.T) {
	dbPort := ports.NewMockDBPort(t)
	dbPort.On("GetNotesMap", mock.Anything, "some-user-id", mock.Anything).
		Return(nil, errors.New("some error"))

	svc := services.NewNotesService(dbPort)

	itemsMap, err := svc.GetItemsMap(t.Context(), "some-user-id", &domain.NoteSearchRequest{})
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if itemsMap != nil {
		t.Fatalf("expected nil items map, got %v", itemsMap)
	}

	dbPort.AssertExpectations(t)
}

func TestNotesSearchItemsByTerm(t *testing.T) {
	req := domain.NoteRequest{
		Title: "some title",
	}

	items := []domain.Note{
		{ID: "1", Title: "Note 1", Content: "Content 1", Tags: []string{"tag1"}},
		{ID: "2", Title: "Note 2", Content: "Content 2", Tags: []string{"tag2"}},
	}

	dbPort := ports.NewMockDBPort(t)
	dbPort.On("SearchNotesByTerm", mock.Anything, "some-user-id", &req).Return(items, nil)

	svc := services.NewNotesService(dbPort)

	foundItems, err := svc.SearchItemsByTerm(t.Context(), "some-user-id", &req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(foundItems) != len(items) {
		t.Fatalf("expected %d items, got %d", len(items), len(foundItems))
	}

	dbPort.AssertExpectations(t)
}

func TestNotesSearchItemsByTermError(t *testing.T) {
	req := domain.NoteRequest{
		Title: "some title",
	}

	dbPort := ports.NewMockDBPort(t)
	dbPort.On("SearchNotesByTerm", mock.Anything, "some-user-id", &req).
		Return(nil, errors.New("some error"))

	svc := services.NewNotesService(dbPort)

	foundItems, err := svc.SearchItemsByTerm(t.Context(), "some-user-id", &req)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if foundItems != nil {
		t.Fatalf("expected nil items, got %v", foundItems)
	}

	dbPort.AssertExpectations(t)
}
