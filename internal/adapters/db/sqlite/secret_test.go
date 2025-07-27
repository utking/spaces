package sqlite_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/utking/spaces/internal/adapters/db/sqlite"
	"github.com/utking/spaces/internal/adapters/db/unittests"
	"github.com/utking/spaces/internal/application/domain"
)

func TestGetSecretTags(t *testing.T) {
	db, dbErr := unittests.CreateTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := sqlite.NewAdapterWithDB(db)
	userID := "uuid-user-12345"

	tags, err := dbAdapter.GetSecretTags(t.Context(), userID)
	if assert.NoError(t, err) {
		assert.Len(t, tags, 5)
		assert.Equal(t, []string{"main", "personal", "test3", "test4", "work"}, tags)
	}

	// Test for empty user ID (all tags)
	tags, err = dbAdapter.GetSecretTags(t.Context(), "")
	if assert.NoError(t, err) {
		assert.Len(t, tags, 8)
		assert.Equal(
			t,
			[]string{"example", "example4", "main", "personal", "test", "test3", "test4", "work"},
			tags,
		)
	}

	// Test for non-existing user ID
	tags, err = dbAdapter.GetSecretTags(t.Context(), "non-existing-user")
	if assert.NoError(t, err) {
		assert.Empty(t, tags)
	}
}

func TestGetSecrets(t *testing.T) {
	db, dbErr := unittests.CreateTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := sqlite.NewAdapterWithDB(db)
	userID := "uuid-user-12345"

	secrets, err := dbAdapter.GetSecrets(t.Context(), userID, nil)
	if assert.NoError(t, err) {
		assert.Len(t, secrets, 3)
		for _, secret := range secrets {
			assert.NotEmpty(t, secret.ID)
			assert.NotEmpty(t, secret.Name)
		}
	}

	// Test for empty user ID (all secrets)
	secrets, err = dbAdapter.GetSecrets(t.Context(), "", nil)
	if assert.NoError(t, err) {
		assert.Empty(t, secrets, "Expected no secrets for empty user ID")
	}

	// Test for non-existing user ID
	secrets, err = dbAdapter.GetSecrets(t.Context(), "non-existing-user", nil)
	if assert.NoError(t, err) {
		assert.Empty(t, secrets, "Expected no secrets for non-existing user ID")
	}
}

func TestGetSecretsCount(t *testing.T) {
	db, dbErr := unittests.CreateTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := sqlite.NewAdapterWithDB(db)
	userID := "uuid-user-12345"

	count, err := dbAdapter.GetSecretsCount(t.Context(), userID, nil)
	if assert.NoError(t, err) {
		assert.EqualValues(t, 3, count)
	}

	// Test for empty user ID (all secrets)
	count, err = dbAdapter.GetSecretsCount(t.Context(), "", nil)
	if assert.NoError(t, err) {
		assert.EqualValues(t, 6, count)
	}

	// Test for non-existing user ID
	count, err = dbAdapter.GetSecretsCount(t.Context(), "non-existing-user", nil)
	if assert.NoError(t, err) {
		assert.EqualValues(t, 0, count)
	}
}

func TestGetSecret(t *testing.T) {
	db, dbErr := unittests.CreateTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := sqlite.NewAdapterWithDB(db)
	userID := "uuid-user-12345"
	secretID := "uuid-password-12345"

	secret, err := dbAdapter.GetSecret(t.Context(), userID, secretID)
	if assert.NoError(t, err) {
		assert.NotNil(t, secret)
		assert.Equal(t, secretID, secret.ID)
		assert.NotEmpty(t, secret.Name)
		assert.NotEmpty(t, secret.EncodedUsername)
		assert.NotEmpty(t, secret.URL)
		assert.NotEmpty(t, secret.EncodedSecret, "Expected encoded secret to be set")
		assert.NotEmpty(t, secret.Description)
		assert.NotEmpty(t, secret.Tags)
	}

	// Test for non-existing secret ID. must be not found error
	secret, err = dbAdapter.GetSecret(t.Context(), userID, "non-existing-secret")
	if assert.Error(t, err) {
		assert.Nil(t, secret)
		assert.EqualError(t, err, "secret not found")
	}
}

func TestDeleteSecret(t *testing.T) {
	db, dbErr := unittests.CreateTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := sqlite.NewAdapterWithDB(db)
	userID := "uuid-user-12345"
	secretID := "uuid-password-12345"

	err := dbAdapter.DeleteSecret(t.Context(), userID, secretID)
	if assert.NoError(t, err) {
		// Verify that the secret was deleted
		secret, getErr := dbAdapter.GetSecret(t.Context(), userID, secretID)
		assert.Error(t, getErr)
		assert.Nil(t, secret)
	}

	// Test for non-existing secret ID. must be not found error
	assert.NoError(
		t, dbAdapter.DeleteSecret(t.Context(), userID, "non-existing-secret-y"),
		"Expected no error when deleting a non-existing secret")

	// delete a secret with an empty user ID
	err = dbAdapter.DeleteSecret(t.Context(), "", secretID)
	assert.NoError(t, err, "Expected no error when deleting a secret with an empty user ID")
}

func TestCreateSecret(t *testing.T) {
	db, dbErr := unittests.CreateTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := sqlite.NewAdapterWithDB(db)
	userID := "uuid-user-12345"

	// Create a new secret
	newSecret := &domain.Secret{
		Name:            "New Secret",
		Description:     "This is a new secret",
		EncodedUsername: []byte("newuser"),
		URL:             "https://example.com",
		Tags:            []string{"newtag1", "newtag2"},
		EncodedSecret:   []byte("encoded-data"),
	}

	secID, err := dbAdapter.CreateSecret(t.Context(), userID, newSecret)
	if assert.NoError(t, err) {
		assert.NotEmpty(t, secID, "Expected a valid secret ID")
		savedSecret, getErr := dbAdapter.GetSecret(t.Context(), userID, secID)
		if assert.NoError(t, getErr) {
			assert.NotNil(t, savedSecret)
			assert.Equal(t, newSecret.Name, savedSecret.Name)
			assert.Equal(t, newSecret.Description, savedSecret.Description)
			assert.Equal(t, newSecret.EncodedUsername, savedSecret.EncodedUsername)
			assert.Equal(t, newSecret.URL, savedSecret.URL)
			assert.Equal(t, newSecret.Tags, savedSecret.Tags)
			// Check that encoded and nonce are set
			assert.NotEmpty(t, savedSecret.EncodedSecret, "Expected encoded secret to be set")
		}
	}

	// test create the same name secret
	secID, err = dbAdapter.CreateSecret(t.Context(), userID, newSecret)
	if assert.Error(t, err) {
		assert.EqualError(t, err, "secret with this name already exists")
		assert.Empty(t, secID, "Expected no secret ID for duplicate creation")
	}

	// test create a secret without nonce and encoded data
	newSecret.Name = "Secret Without Encoded"
	newSecret.EncodedSecret = []byte{}

	secID, err = dbAdapter.CreateSecret(t.Context(), userID, newSecret)
	if assert.NoError(t, err) {
		assert.NotEmpty(t, secID, "Expected a valid secret ID")
		savedSecret, getErr := dbAdapter.GetSecret(t.Context(), userID, secID)
		if assert.NoError(t, getErr) {
			assert.NotNil(t, savedSecret)
			assert.Equal(t, newSecret.Name, savedSecret.Name)
			assert.Equal(t, newSecret.Description, savedSecret.Description)
			assert.Equal(t, newSecret.EncodedUsername, savedSecret.EncodedUsername)
			assert.Equal(t, newSecret.URL, savedSecret.URL)
			assert.Equal(t, newSecret.Tags, savedSecret.Tags)
			// Check that encoded and nonce are nil
			assert.Nil(t, savedSecret.EncodedSecret, "Expected encoded secret to be nil")
		}
	}

	// test create with encoded nil - will create a secret with empty encoded
	newSecret.Name = "Secret With Encoded Nil"
	newSecret.EncodedSecret = nil

	secID, err = dbAdapter.CreateSecret(t.Context(), userID, newSecret)
	if assert.NoError(t, err) {
		assert.NotEmpty(t, secID, "Expected a valid secret ID")
		savedSecret, getErr := dbAdapter.GetSecret(t.Context(), userID, secID)
		if assert.NoError(t, getErr) {
			assert.NotNil(t, savedSecret)
			assert.Equal(t, newSecret.Name, savedSecret.Name)
			assert.Equal(t, newSecret.Description, savedSecret.Description)
			assert.Equal(t, newSecret.EncodedUsername, savedSecret.EncodedUsername)
			assert.Equal(t, newSecret.URL, savedSecret.URL)
			assert.Equal(t, newSecret.Tags, savedSecret.Tags)
			// Check that encoded and nonce are nil
			assert.Empty(t, savedSecret.EncodedSecret, "Expected encoded secret to be empty")
		}
	}

	// create a secret with only the name
	newSecret = &domain.Secret{}

	newSecret.Name = "just a name secret"
	newSecret.Description = ""
	newSecret.EncodedUsername = nil
	newSecret.EncodedSecret = nil
	newSecret.URL = ""
	newSecret.Tags = []string{}

	// must fail on empty tags list
	secID, err = dbAdapter.CreateSecret(t.Context(), userID, newSecret)
	if assert.Error(t, err) {
		assert.Empty(t, secID, "Expected an empty secret ID for invalid creation")
	}
}

func TestUpdateSecret(t *testing.T) {
	db, dbErr := unittests.CreateTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := sqlite.NewAdapterWithDB(db)
	userID := "uuid-user-12345"
	secretID := "uuid-password-12345"

	// Get the existing secret
	secret, err := dbAdapter.GetSecret(t.Context(), userID, secretID)
	if assert.NoError(t, err) {
		assert.NotNil(t, secret)
		assert.NotEmpty(t, secret.ID, "Expected a valid secret ID")
		assert.NotEmpty(t, secret.Name, "Expected a valid secret name")
		assert.NotEmpty(t, secret.EncodedUsername, "Expected a valid secret username")
		assert.NotEmpty(t, secret.URL, "Expected a valid secret URL")
		assert.NotEmpty(t, secret.EncodedSecret, "Expected a valid encoded secret")
	}

	// Update the secret
	secret.Name = "Updated Secret Name"
	secret.Description = "Updated description"
	secret.EncodedUsername = []byte("updateduser")
	secret.URL = "https://updated-example.com"
	secret.Tags = []string{"updatedtag1", "updatedtag2"}

	affected, updateErr := dbAdapter.UpdateSecret(t.Context(), userID, secretID, secret)
	if assert.NoError(t, updateErr) {
		// Verify that the secret was updated
		assert.EqualValues(t, 1, affected, "Expected one row to be updated")
		updatedSecret, getErr := dbAdapter.GetSecret(t.Context(), userID, secretID)
		if assert.NoError(t, getErr) {
			assert.Equal(t, secret.Name, updatedSecret.Name)
			assert.Equal(t, secret.Description, updatedSecret.Description)
			assert.Equal(t, secret.EncodedUsername, updatedSecret.EncodedUsername)
			assert.Equal(t, secret.URL, updatedSecret.URL)
			assert.Equal(t, secret.Tags, updatedSecret.Tags)
		}

		// Test updating a secret without changes (affected rows should be 1). Should not return an error
		affected, updateErr = dbAdapter.UpdateSecret(t.Context(), userID, secretID, secret)
		if assert.NoError(t, updateErr) {
			assert.EqualValues(t, 1, affected,
				"Expected one row to be affected when updating with the same data")
		}
	}

	// Test updating a non-existing secret ID. it will create the secret and no error
	affected, updateErr = dbAdapter.UpdateSecret(t.Context(), userID, "non-existing-secret-x", secret)
	if assert.Error(t, updateErr, "Expected an error when updating a non-existing secret") {
		assert.EqualError(t, updateErr, "the secret does not exist/belong to current user")
		assert.EqualValues(t, 0, affected, "Expected no rows to be affected for non-existing secret")
	}

	// Test updating a secret with an empty user ID
	_, updateErr = dbAdapter.UpdateSecret(t.Context(), "", secret.ID, secret)
	if assert.Error(t, updateErr) {
		assert.EqualError(t, updateErr, "the secret does not exist/belong to current user")
	}

	// Test update secret with nil encoded data
	secret.EncodedSecret = nil
	affected, updateErr = dbAdapter.UpdateSecret(t.Context(), userID, secretID, secret)
	if assert.NoError(t, updateErr) {
		assert.EqualValues(t, 1, affected, "Expected one row to be updated with empty encoded data")
		updatedSecret, getErr := dbAdapter.GetSecret(t.Context(), userID, secretID)
		if assert.NoError(t, getErr) {
			assert.NotNil(t, updatedSecret)
			assert.Equal(t, secret.Name, updatedSecret.Name)
			assert.Equal(t, secret.Description, updatedSecret.Description)
			assert.Equal(t, secret.EncodedUsername, updatedSecret.EncodedUsername)
			assert.Equal(t, secret.URL, updatedSecret.URL)
			assert.Equal(t, secret.Tags, updatedSecret.Tags)
			assert.Empty(t, updatedSecret.EncodedSecret, "Expected encoded secret to be empty")
		}
	}

	// Test update secret with empty encoded data
	secret.EncodedSecret = []byte{}
	affected, updateErr = dbAdapter.UpdateSecret(t.Context(), userID, secretID, secret)
	if assert.NoError(t, updateErr) {
		assert.EqualValues(t, 1, affected, "Expected one row to be updated with empty encoded data")
		updatedSecret, getErr := dbAdapter.GetSecret(t.Context(), userID, secretID)
		if assert.NoError(t, getErr) {
			assert.NotNil(t, updatedSecret)
			assert.Equal(t, secret.Name, updatedSecret.Name)
			assert.Equal(t, secret.Description, updatedSecret.Description)
			assert.Equal(t, secret.EncodedUsername, updatedSecret.EncodedUsername)
			assert.Equal(t, secret.URL, updatedSecret.URL)
			assert.Equal(t, secret.Tags, updatedSecret.Tags)
			assert.Empty(t, updatedSecret.EncodedSecret, "Expected encoded secret to be empty")
		}
	}
}

func TestGetSecretsMap(t *testing.T) {
	db, dbErr := unittests.CreateTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := sqlite.NewAdapterWithDB(db)
	userID := "uuid-user-12345"

	secretsMap, err := dbAdapter.GetSecretsMap(t.Context(), userID, nil)
	if assert.NoError(t, err) {
		assert.Len(t, secretsMap, 3)
		for _, secret := range secretsMap {
			assert.NotEmpty(t, secret.Name)
			assert.NotEmpty(t, secret.EncodedUsername)
			assert.NotEmpty(t, secret.URL)
			assert.NotEmpty(t, secret.Description)
		}
	}

	// Test for empty user ID (all secrets)
	secretsMap, err = dbAdapter.GetSecretsMap(t.Context(), "", nil)
	if assert.NoError(t, err) {
		assert.Empty(t, secretsMap, "Expected no secrets for empty user ID")
	}

	// Test for non-existing user ID
	secretsMap, err = dbAdapter.GetSecretsMap(t.Context(), "non-existing-user-z", nil)
	if assert.NoError(t, err) {
		assert.Empty(t, secretsMap, "Expected no secrets for non-existing user ID")
	}
}

func TestSearchSecretsByTerm(t *testing.T) {
	db, dbErr := unittests.CreateTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := sqlite.NewAdapterWithDB(db)
	userID := "uuid-user-12345"
	term := "user"
	req := &domain.SecretRequest{
		Name:        term,
		Username:    term,
		URL:         term,
		Description: term,
	}

	// Search for secrets by term
	secrets, err := dbAdapter.SearchSecretsByTerm(t.Context(), userID, req)
	if assert.NoError(t, err) {
		assert.Len(t, secrets, 3) // Assuming there are 2 secrets with "test" in their name or description
		for _, secret := range secrets {
			assert.NotEmpty(t, secret.Name)
		}
	}

	// Test for empty user ID (all secrets)
	secrets, err = dbAdapter.SearchSecretsByTerm(t.Context(), "", req)
	if assert.NoError(t, err) {
		assert.Empty(t, secrets, "Expected no secrets for empty user ID")
	}

	// Test for non-existing user ID
	secrets, err = dbAdapter.SearchSecretsByTerm(t.Context(), "non-existing-user-x", req)
	if assert.NoError(t, err) {
		assert.Empty(t, secrets, "Expected no secrets for non-existing user ID")
	}
}
