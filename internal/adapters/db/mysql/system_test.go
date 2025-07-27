//go:build mysql
// +build mysql

package mysql_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/utking/spaces/internal/adapters/db/mysql"
	"github.com/utking/spaces/internal/adapters/db/unittests"
)

func TestGetSystemStats(t *testing.T) {
	db, dbErr := unittests.CreateMySQLTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := mysql.NewAdapterWithDB(db)
	// get for a specific user ID
	userID := "uuid-user-12345"

	stats, err := dbAdapter.GetSystemStats(t.Context(), userID)
	if assert.NoError(t, err) {
		assert.NotNil(t, stats)
		assert.EqualValues(t, 4, stats.NoteTags)
		assert.EqualValues(t, 3, stats.Notes)
		assert.EqualValues(t, 5, stats.SecretTags)
		assert.EqualValues(t, 3, stats.Secrets)
		// user does not have bookmarks or bookmark tags in this test setup
		assert.EqualValues(t, 0, stats.Bookmarks)
		assert.EqualValues(t, 0, stats.BookmarkTags)
	}

	// test for empty user ID (all items)
	stats, err = dbAdapter.GetSystemStats(t.Context(), "")
	if assert.NoError(t, err) {
		assert.NotNil(t, stats)
		assert.EqualValues(t, 2, stats.ActiveUsers)
		assert.EqualValues(t, 1, stats.InactiveUsers)
		assert.EqualValues(t, 7, stats.NoteTags)
		assert.EqualValues(t, 6, stats.Notes)
		assert.EqualValues(t, 8, stats.SecretTags)
		assert.EqualValues(t, 6, stats.Secrets)
		// for all users, we have bookmarks and bookmark tags
		assert.EqualValues(t, 4, stats.Bookmarks)
		assert.EqualValues(t, 5, stats.BookmarkTags)
	}

	// test for non-existing user ID
	stats, err = dbAdapter.GetSystemStats(t.Context(), "non-existing-user")
	if assert.NoError(t, err) {
		assert.NotNil(t, stats)
		assert.EqualValues(t, 0, stats.NoteTags)
		assert.EqualValues(t, 0, stats.Notes)
		assert.EqualValues(t, 0, stats.SecretTags)
		assert.EqualValues(t, 0, stats.Secrets)
		assert.EqualValues(t, 0, stats.Bookmarks)
		assert.EqualValues(t, 0, stats.BookmarkTags)
	}
}
