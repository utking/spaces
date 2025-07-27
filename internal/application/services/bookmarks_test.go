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

func TestBookmarksGetTags(t *testing.T) {
	tagsInDB := []string{"tag1", "tag2"}

	dbPort := ports.NewMockDBPort(t)
	dbPort.On("GetBookmarkTags", mock.Anything, "some-user-id").Return(tagsInDB, nil)

	svc := services.NewBookmarkService(dbPort)

	tags, err := svc.GetTags(t.Context(), "some-user-id")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(tags) != len(tagsInDB) {
		t.Fatalf("expected some tags, got none")
	}

	dbPort.AssertExpectations(t)
}

func TestBookmarksGetTagsError(t *testing.T) {
	dbPort := ports.NewMockDBPort(t)
	dbPort.On("GetBookmarkTags", mock.Anything, "some-user-id").Return([]string{}, errors.New("some error"))

	svc := services.NewBookmarkService(dbPort)

	_items, err := svc.GetTags(t.Context(), "some-user-id")
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if len(_items) != 0 {
		t.Fatalf("expected no tags, got %d", len(_items))
	}

	dbPort.AssertExpectations(t)
}

func TestBookmarksSearchItemsByTerm(t *testing.T) {
	req := domain.BookmarkSearchRequest{
		Title: "term",
		URL:   "term",
	}

	items := []domain.Bookmark{
		{ID: "1", Title: "Note 1", URL: "https://url-1", Tags: []string{"tag1"}},
		{ID: "2", Title: "Note 2", URL: "https://url-1/term", Tags: []string{"tag1"}},
		{ID: "3", Title: "Note 3", URL: "https://url-2/", Tags: []string{"term"}},
	}

	dbPort := ports.NewMockDBPort(t)
	dbPort.On("SearchBookmarksByTerm", mock.Anything, "some-user-id", &req).Return(items, nil)

	svc := services.NewBookmarkService(dbPort)

	foundItems, err := svc.SearchItemsByTerm(t.Context(), "some-user-id", &req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(foundItems) != len(items) {
		t.Fatalf("expected %d items, got %d", len(items), len(foundItems))
	}

	dbPort.AssertExpectations(t)
}

func TestBookmarkssSearchItemsByTermError(t *testing.T) {
	req := domain.BookmarkSearchRequest{
		Title: "term",
		URL:   "term",
	}

	dbPort := ports.NewMockDBPort(t)
	dbPort.On("SearchBookmarksByTerm", mock.Anything, "some-user-id", &req).
		Return(nil, errors.New("some error"))

	svc := services.NewBookmarkService(dbPort)

	foundItems, err := svc.SearchItemsByTerm(t.Context(), "some-user-id", &req)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if foundItems != nil {
		t.Fatalf("expected nil items, got %v", foundItems)
	}

	dbPort.AssertExpectations(t)
}

func TestBookmarksGetItems(t *testing.T) {
	req := &domain.BookmarkSearchRequest{}

	items := []domain.Bookmark{
		{ID: "1", Title: "Note 1", URL: "https://url-1", Tags: []string{"tag1"}},
		{ID: "2", Title: "Note 2", URL: "https://url-2", Tags: []string{"tag2"}},
	}

	dbPort := ports.NewMockDBPort(t)
	dbPort.On("GetBookmarks", mock.Anything, "some-user-id", req).Return(items, nil)

	svc := services.NewBookmarkService(dbPort)

	foundItems, err := svc.GetItems(t.Context(), "some-user-id", req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(foundItems) != len(items) {
		t.Fatalf("expected %d items, got %d", len(items), len(foundItems))
	}

	dbPort.AssertExpectations(t)
}

func TestBookmarksGetItemsError(t *testing.T) {
	req := &domain.BookmarkSearchRequest{}
	items := []domain.Bookmark{}

	dbPort := ports.NewMockDBPort(t)
	dbPort.On("GetBookmarks", mock.Anything, "some-user-id", req).
		Return(items, errors.New("some error"))

	svc := services.NewBookmarkService(dbPort)

	foundItems, err := svc.GetItems(t.Context(), "some-user-id", req)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if len(foundItems) != 0 {
		t.Fatalf("expected no items, got %d", len(foundItems))
	}

	dbPort.AssertExpectations(t)
}

func TestBookmarksGetCount(t *testing.T) {
	req := &domain.BookmarkSearchRequest{}

	expectedCount := int64(42)

	dbPort := ports.NewMockDBPort(t)
	dbPort.On("GetBookmarksCount", mock.Anything, "some-user-id", req).Return(expectedCount, nil)

	svc := services.NewBookmarkService(dbPort)

	count, err := svc.GetCount(t.Context(), "some-user-id", req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if count != expectedCount {
		t.Fatalf("expected count %d, got %d", expectedCount, count)
	}

	dbPort.AssertExpectations(t)
}

func TestBookmarksGetCountError(t *testing.T) {
	req := &domain.BookmarkSearchRequest{}

	dbPort := ports.NewMockDBPort(t)
	dbPort.On("GetBookmarksCount", mock.Anything, "some-user-id", req).
		Return(int64(0), errors.New("some error"))

	svc := services.NewBookmarkService(dbPort)

	count, err := svc.GetCount(t.Context(), "some-user-id", req)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if count != 0 {
		t.Fatalf("expected count 0, got %d", count)
	}

	dbPort.AssertExpectations(t)
}

func TestBookmarksGetItem(t *testing.T) {
	item := &domain.Bookmark{
		ID:    "1",
		Title: "Note 1",
		URL:   "https://url-1",
		Tags:  []string{"tag1"},
	}

	dbPort := ports.NewMockDBPort(t)
	dbPort.On("GetBookmark", mock.Anything, "some-user-id", "1").Return(item, nil)

	svc := services.NewBookmarkService(dbPort)

	foundItem, err := svc.GetItem(t.Context(), "some-user-id", "1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if foundItem == nil {
		t.Fatalf("expected item, got nil")
	}

	if foundItem.ID != item.ID || foundItem.Title != item.Title || foundItem.URL != item.URL {
		t.Fatalf("expected item %v, got %v", item, foundItem)
	}

	dbPort.AssertExpectations(t)
}

func TestBookmarksGetItemError(t *testing.T) {
	dbPort := ports.NewMockDBPort(t)
	dbPort.On("GetBookmark", mock.Anything, "some-user-id", "1").
		Return(nil, errors.New("some error"))

	svc := services.NewBookmarkService(dbPort)

	foundItem, err := svc.GetItem(t.Context(), "some-user-id", "1")
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if foundItem != nil {
		t.Fatalf("expected nil item, got %v", foundItem)
	}

	dbPort.AssertExpectations(t)
}

func TestBookmarksGetItemEMptyIDError(t *testing.T) {
	svc := services.NewBookmarkService(ports.NewMockDBPort(t))

	foundItem, err := svc.GetItem(t.Context(), "some-user-id", "")
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if foundItem != nil {
		t.Fatalf("expected nil item, got %v", foundItem)
	}
}

func TestBookmarksCreate(t *testing.T) {
	req := &domain.Bookmark{
		ID:    "new-id",
		Title: "Note 1",
		URL:   "https://url-1",
		Tags:  []string{"tag1"},
	}

	expectedID := "new-id"

	dbPort := ports.NewMockDBPort(t)
	dbPort.On("CreateBookmark", mock.Anything, "some-user-id", req).Return(expectedID, nil)

	svc := services.NewBookmarkService(dbPort)

	id, err := svc.Create(t.Context(), "some-user-id", req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if id != expectedID {
		t.Fatalf("expected ID %s, got %s", expectedID, id)
	}

	dbPort.AssertExpectations(t)
}

func TestBookmarksCreateError(t *testing.T) {
	req := &domain.Bookmark{
		ID:    "new-id",
		Title: "Note 1",
		URL:   "https://url-1",
		Tags:  []string{"tag1"},
	}

	dbPort := ports.NewMockDBPort(t)
	dbPort.On("CreateBookmark", mock.Anything, "some-user-id", req).
		Return("", errors.New("some error"))

	svc := services.NewBookmarkService(dbPort)

	id, err := svc.Create(t.Context(), "some-user-id", req)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if id != "" {
		t.Fatalf("expected empty ID, got %s", id)
	}

	dbPort.AssertExpectations(t)
}

func TestBookmarksCreateEmptyTagsError(t *testing.T) {
	req := &domain.Bookmark{
		ID:    "new-id",
		Title: "Note 1",
		URL:   "https://url-1",
		Tags:  []string{},
	}

	svc := services.NewBookmarkService(ports.NewMockDBPort(t))

	id, err := svc.Create(t.Context(), "some-user-id", req)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if id != "" {
		t.Fatalf("expected empty ID, got %s", id)
	}
}

func TestBookmarksCreateEmptyTitleError(t *testing.T) {
	req := &domain.Bookmark{
		ID:    "new-id",
		Title: "",
		URL:   "https://url-1",
		Tags:  []string{"tag1"},
	}

	svc := services.NewBookmarkService(ports.NewMockDBPort(t))

	id, err := svc.Create(t.Context(), "some-user-id", req)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if id != "" {
		t.Fatalf("expected empty ID, got %s", id)
	}
}

func TestBookmarksCreateEmptyURL(t *testing.T) {
	req := &domain.Bookmark{
		ID:    "new-id",
		Title: "Note 1",
		URL:   "",
		Tags:  []string{"tag1"},
	}

	svc := services.NewBookmarkService(ports.NewMockDBPort(t))

	id, err := svc.Create(t.Context(), "some-user-id", req)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if id != "" {
		t.Fatalf("expected empty ID, got %s", id)
	}
}

func TestBookmarksCreateTooLongURLError(t *testing.T) {
	req := &domain.Bookmark{
		ID:    "new-id",
		Title: "Note 1",
		URL:   "a" + strings.Repeat("a", 4096), // URL too long
		Tags:  []string{"tag1"},
	}

	svc := services.NewBookmarkService(ports.NewMockDBPort(t))

	id, err := svc.Create(t.Context(), "some-user-id", req)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if id != "" {
		t.Fatalf("expected empty ID, got %s", id)
	}
}

func TestBookmarksCreateTooLongTitleError(t *testing.T) {
	req := &domain.Bookmark{
		ID:    "new-id",
		Title: strings.Repeat("a", 256), // Title too long
		URL:   "https://url-1",
		Tags:  []string{"tag1"},
	}

	svc := services.NewBookmarkService(ports.NewMockDBPort(t))

	id, err := svc.Create(t.Context(), "some-user-id", req)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if id != "" {
		t.Fatalf("expected empty ID, got %s", id)
	}
}

func TestBookmarksUpdate(t *testing.T) {
	req := &domain.Bookmark{
		ID:    "1",
		Title: "Updated Note",
		URL:   "https://updated-url",
		Tags:  []string{"updated-tag"},
	}

	expectedRowsAffected := int64(1)

	dbPort := ports.NewMockDBPort(t)
	dbPort.On("UpdateBookmark", mock.Anything, "some-user-id", "1", req).
		Return(expectedRowsAffected, nil)

	svc := services.NewBookmarkService(dbPort)

	rowsAffected, err := svc.Update(t.Context(), "some-user-id", "1", req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if rowsAffected != expectedRowsAffected {
		t.Fatalf("expected %d rows affected, got %d", expectedRowsAffected, rowsAffected)
	}

	dbPort.AssertExpectations(t)
}

func TestBookmarksUpdateError(t *testing.T) {
	req := &domain.Bookmark{
		ID:    "1",
		Title: "Updated Note",
		URL:   "https://updated-url",
		Tags:  []string{"updated-tag"},
	}

	dbPort := ports.NewMockDBPort(t)
	dbPort.On("UpdateBookmark", mock.Anything, "some-user-id", "1", req).
		Return(int64(0), errors.New("some error"))

	svc := services.NewBookmarkService(dbPort)

	rowsAffected, err := svc.Update(t.Context(), "some-user-id", "1", req)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}

	dbPort.AssertExpectations(t)
}

func TestBookmarksUpdateEmptyIDError(t *testing.T) {
	req := &domain.Bookmark{
		ID:    "",
		Title: "Updated Note",
		URL:   "https://updated-url",
		Tags:  []string{"updated-tag"},
	}

	svc := services.NewBookmarkService(ports.NewMockDBPort(t))

	rowsAffected, err := svc.Update(t.Context(), "some-user-id", "", req)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}
}

func TestBookmarksUpdateEmptyTagsError(t *testing.T) {
	req := &domain.Bookmark{
		ID:    "1",
		Title: "Updated Note",
		URL:   "https://updated-url",
		Tags:  []string{},
	}

	svc := services.NewBookmarkService(ports.NewMockDBPort(t))

	rowsAffected, err := svc.Update(t.Context(), "some-user-id", "1", req)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}
}

func TestBookmarksUpdateEmptyTitleError(t *testing.T) {
	req := &domain.Bookmark{
		ID:    "1",
		Title: "",
		URL:   "https://updated-url",
		Tags:  []string{"updated-tag"},
	}

	svc := services.NewBookmarkService(ports.NewMockDBPort(t))

	rowsAffected, err := svc.Update(t.Context(), "some-user-id", "1", req)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}
}

func TestBookmarksUpdateEmptyURL(t *testing.T) {
	req := &domain.Bookmark{
		ID:    "1",
		Title: "Updated Note",
		URL:   "",
		Tags:  []string{"updated-tag"},
	}

	svc := services.NewBookmarkService(ports.NewMockDBPort(t))

	rowsAffected, err := svc.Update(t.Context(), "some-user-id", "1", req)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}
}

func TestBookmarksUpdateTooLongURLError(t *testing.T) {
	req := &domain.Bookmark{
		ID:    "1",
		Title: "Updated Note",
		URL:   "a" + strings.Repeat("a", 4096), // URL too long
		Tags:  []string{"updated-tag"},
	}

	svc := services.NewBookmarkService(ports.NewMockDBPort(t))

	rowsAffected, err := svc.Update(t.Context(), "some-user-id", "1", req)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}
}

func TestBookmarksUpdateTooLongTitleError(t *testing.T) {
	req := &domain.Bookmark{
		ID:    "1",
		Title: strings.Repeat("a", 256), // Title too long
		URL:   "https://updated-url",
		Tags:  []string{"updated-tag"},
	}

	svc := services.NewBookmarkService(ports.NewMockDBPort(t))

	rowsAffected, err := svc.Update(t.Context(), "some-user-id", "1", req)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}
}

func TestBookmarksUpdateEmptyID(t *testing.T) {
	svc := services.NewBookmarkService(ports.NewMockDBPort(t))

	rowsAffected, err := svc.Update(t.Context(), "some-user-id", "", nil)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if rowsAffected != 0 {
		t.Fatalf("expected 0 rows affected, got %d", rowsAffected)
	}
}

func TestBookmarksDelete(t *testing.T) {
	dbPort := ports.NewMockDBPort(t)
	dbPort.On("DeleteBookmark", mock.Anything, "some-user-id", "1").Return(nil)

	svc := services.NewBookmarkService(dbPort)

	err := svc.Delete(t.Context(), "some-user-id", "1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	dbPort.AssertExpectations(t)
}

func TestBookmarksDeleteError(t *testing.T) {
	dbPort := ports.NewMockDBPort(t)
	dbPort.On("DeleteBookmark", mock.Anything, "some-user-id", "1").
		Return(errors.New("some error"))

	svc := services.NewBookmarkService(dbPort)

	err := svc.Delete(t.Context(), "some-user-id", "1")
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	dbPort.AssertExpectations(t)
}

func TestBookmarksDeleteEmptyIDError(t *testing.T) {
	svc := services.NewBookmarkService(ports.NewMockDBPort(t))

	err := svc.Delete(t.Context(), "some-user-id", "")
	if err == nil {
		t.Fatalf("expected error, got none")
	}
}

func TestBookmarksGetItemsMap(t *testing.T) {
	req := &domain.BookmarkSearchRequest{}

	items := []domain.Bookmark{
		{ID: "1", Title: "Note 1", URL: "https://url-1", Tags: []string{"tag1"}},
		{ID: "2", Title: "Note 2", URL: "https://url-2", Tags: []string{"tag2"}},
	}

	dbPort := ports.NewMockDBPort(t)
	dbPort.On("GetBookmarksMap", mock.Anything, "some-user-id", req).Return(items, nil)

	svc := services.NewBookmarkService(dbPort)

	foundItems, err := svc.GetItemsMap(t.Context(), "some-user-id", req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(foundItems) != len(items) {
		t.Fatalf("expected %d items, got %d", len(items), len(foundItems))
	}

	dbPort.AssertExpectations(t)
}

func TestBookmarksGetItemsMapError(t *testing.T) {
	req := &domain.BookmarkSearchRequest{}

	dbPort := ports.NewMockDBPort(t)
	dbPort.On("GetBookmarksMap", mock.Anything, "some-user-id", req).
		Return(nil, errors.New("some error"))

	svc := services.NewBookmarkService(dbPort)

	foundItems, err := svc.GetItemsMap(t.Context(), "some-user-id", req)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if foundItems != nil {
		t.Fatalf("expected nil items, got %v", foundItems)
	}

	dbPort.AssertExpectations(t)
}
