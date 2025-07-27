//go:build mysql
// +build mysql

package mysql_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gogs.utking.net/utking/spaces/internal/adapters/db/mysql"
	"gogs.utking.net/utking/spaces/internal/adapters/db/unittests"
	"gogs.utking.net/utking/spaces/internal/application/domain"
)

func TestGetUsers(t *testing.T) {
	db, dbErr := unittests.CreateMySQLTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := mysql.NewAdapterWithDB(db)

	users, err := dbAdapter.GetUsers(t.Context(), nil)
	if assert.NoError(t, err, "GetUsers should not return an error") {
		assert.Len(t, users, 3, "Expected 3 users in the test database")
	}
}

func TestGetUsersInactive(t *testing.T) {
	db, dbErr := unittests.CreateMySQLTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := mysql.NewAdapterWithDB(db)
	statusInactive := int64(domain.UserInactive)

	users, err := dbAdapter.GetUsers(t.Context(), &domain.UserRequest{
		Status: &statusInactive,
	})

	if assert.NoError(t, err, "GetUsersInactive should not return an error") {
		if assert.Len(t, users, 1, "Expected 1 inactive user in the test database") {
			assert.Zero(t, users[0].Status, "Expected the inactive user to have status 0 (inactive)")
			assert.Equal(t, "uuid-user-67890", users[0].ID)
			assert.Equal(t, "user", users[0].RoleName, "Expected the inactive user to have role 'user'")
		}
	}
}

func TestGetUsersActive(t *testing.T) {
	db, dbErr := unittests.CreateMySQLTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := mysql.NewAdapterWithDB(db)
	statusActive := int64(domain.UserActive)

	users, err := dbAdapter.GetUsers(t.Context(), &domain.UserRequest{
		Status: &statusActive,
	})

	if assert.NoError(t, err, "GetUsersActive should not return an error") {
		if assert.Len(t, users, 2, "Expected 2 active users in the test database") {
			for _, user := range users {
				assert.EqualValues(t, domain.UserActive, user.Status,
					"Expected the active users to have status 1 (active)")
				assert.NotEmpty(t, user.ID, "Expected the active users to have a non-empty ID")
				assert.Contains(t, []string{"uuid-user-12345", "uuid-user-11223"}, user.ID)
				assert.Contains(t, []string{"admin", "user"}, user.RoleName)
			}
		}
	}
}

func TestGetUsersWithFilter(t *testing.T) {
	db, dbErr := unittests.CreateMySQLTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := mysql.NewAdapterWithDB(db)

	users, err := dbAdapter.GetUsers(t.Context(), &domain.UserRequest{
		Username: "user112",
	})

	if assert.NoError(t, err, "GetUsersWithFilter should not return an error") {
		assert.Len(t, users, 1, "Expected 1 user matching the filter in the test database")
		for _, user := range users {
			assert.EqualValues(t, domain.UserActive, user.Status)
		}
	}

	users, err = dbAdapter.GetUsers(t.Context(), &domain.UserRequest{
		Email: "user112@localhost",
	})

	if assert.NoError(t, err, "GetUsersWithFilter by email should not return an error") {
		assert.Len(t, users, 1, "Expected 1 user matching the email filter in the test database")
		for _, user := range users {
			assert.Equal(t, "user112", user.Username, "Expected the user to have username 'user112'")
		}
	}

	// get non-existing user
	users, err = dbAdapter.GetUsers(t.Context(), &domain.UserRequest{
		Username: "nonexistinguser",
	})

	if assert.NoError(t, err, "GetUsersWithFilter for non-existing user should not return an error") {
		assert.Empty(t, users, "Expected no users matching the username in the test database")
	}
}

func TestGetUsersCount(t *testing.T) {
	db, dbErr := unittests.CreateMySQLTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := mysql.NewAdapterWithDB(db)

	count, err := dbAdapter.GetUsersCount(t.Context(), &domain.UserRequest{})
	if assert.NoError(t, err, "GetUsersCount should not return an error") {
		assert.Equal(t, int64(3), count, "Expected 3 users in the test database")
	}

	count, err = dbAdapter.GetUsersCount(t.Context(), &domain.UserRequest{
		Status: &[]int64{domain.UserActive}[0],
	})

	if assert.NoError(t, err, "GetUsersCount with status should not return an error") {
		assert.Equal(t, int64(2), count, "Expected 2 active users in the test database")
	}
}

func TestGetUser(t *testing.T) {
	db, dbErr := unittests.CreateMySQLTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := mysql.NewAdapterWithDB(db)

	user, err := dbAdapter.GetUser(t.Context(), "uuid-user-12345")
	if assert.NoError(t, err, "GetUser should not return an error") {
		assert.Equal(t, "user123", user.Username, "Expected the user to have username 'user123'")
		assert.Equal(t, int64(domain.UserActive), user.Status, "Expected the user to be active")
		assert.Equal(t, "admin", user.RoleName, "Expected the user to have role 'admin'")
	}

	user, err = dbAdapter.GetUser(t.Context(), "non-existing-uuid")
	if assert.Error(t, err, "GetUser for non-existing user should return an error") {
		assert.Nil(t, user, "Expected no user to be returned for non-existing UUID")
	}
}

func TestCreateUser(t *testing.T) {
	db, dbErr := unittests.CreateMySQLTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := mysql.NewAdapterWithDB(db)

	newUser := &domain.User{
		ID:       "uuid-new-user",
		Username: "newuser",
		Email:    "newuser@localhost",
		Status:   domain.UserActive,
		RoleName: "user",
	}

	userID, token, err := dbAdapter.CreateUser(t.Context(), newUser)
	if assert.NoError(t, err, "CreateUser should not return an error when there are not conflicts") {
		assert.Len(t, userID, 36, "Expected a valid UUID for the new user")
		assert.NotEmpty(t, token, "Expected a non-empty token for the new user")

		retrievedUser, getErr := dbAdapter.GetUser(t.Context(), userID)
		if assert.NoError(t, getErr, "GetUser after CreateUser should not return an error") {
			assert.Equal(t, newUser.Username, retrievedUser.Username)
			assert.Equal(t, newUser.Email, retrievedUser.Email)
			assert.Equal(t, newUser.RoleName, retrievedUser.RoleName)
			assert.EqualValues(t, domain.UserInactive, retrievedUser.Status)
		}
	}
}

func TestDeleteUser(t *testing.T) {
	db, dbErr := unittests.CreateMySQLTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := mysql.NewAdapterWithDB(db)

	// DeleteUser only suspends the user, setting Status to UserInactive
	err := dbAdapter.DeleteUser(t.Context(), "uuid-user-11223")
	if assert.NoError(t, err, "DeleteUser should not return an error for existing user") {
		user, getErr := dbAdapter.GetUser(t.Context(), "uuid-user-11223")
		if assert.NoError(t, getErr, "GetUser after DeleteUser should return no error for deleted user") {
			assert.Equal(t, "user112", user.Username, "Expected the deleted user to still have the same username")
			assert.EqualValues(t, domain.UserInactive, user.Status)
			assert.Equal(t, "user", user.RoleName, "Expected the deleted user to have role 'user'")
		}
	}

	err = dbAdapter.DeleteUser(t.Context(), "non-existing-uuid")
	assert.NoError(t, err, "DeleteUser for non-existing user should not return an error")
}

func TestUpdateUser(t *testing.T) {
	db, dbErr := unittests.CreateMySQLTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := mysql.NewAdapterWithDB(db)

	// Update existing user
	updateUser := &domain.UserUpdate{
		Username: "updateduser123",
		Email:    "updateduser123@localhost",
		Status:   domain.UserActive,
		RoleName: "admin",
	}

	affected, err := dbAdapter.UpdateUser(t.Context(), "uuid-user-12345", updateUser)
	if assert.NoError(t, err, "UpdateUser should not return an error for existing user") {
		assert.Equal(t, int64(1), affected, "Expected 1 user to be updated")
		retrievedUser, getErr := dbAdapter.GetUser(t.Context(), "uuid-user-12345")
		if assert.NoError(t, getErr, "GetUser after UpdateUser should not return an error") {
			assert.NotEqual(t, updateUser.Username, retrievedUser.Username, "Username should not be updated")
			assert.Equal(t, updateUser.Email, retrievedUser.Email)
			assert.Equal(t, updateUser.RoleName, retrievedUser.RoleName)
			assert.EqualValues(t, domain.UserActive, retrievedUser.Status)
		}
	}

	// Attempt to update non-existing user
	affected, err = dbAdapter.UpdateUser(t.Context(), "non-existing-uuid", updateUser)
	if assert.NoError(t, err, "UpdateUser for non-existing user should not return an error") {
		assert.Equal(t, int64(0), affected, "Expected no users to be updated for non-existing UUID")
	}

	// Update with new password
	updateUserWithPassword := &domain.UserUpdate{
		Email:           "user123@localhost", // email is required for updates
		Password:        "newpassword123",
		PasswordConfirm: "newpassword123",
	}

	affected, err = dbAdapter.UpdateUser(t.Context(), "uuid-user-12345", updateUserWithPassword)
	if assert.NoError(t, err, "UpdateUser with password should not return an error") {
		assert.Equal(t, int64(1), affected, "Expected 1 user to be updated with new password")
		retrievedUser, getErr := dbAdapter.GetUser(t.Context(), "uuid-user-12345")
		if assert.NoError(t, getErr, "GetUser after UpdateUser with password should not return an error") {
			assert.Equal(t, "user123", retrievedUser.Username, "Username should remain unchanged")
			// password change validation is part of testing UserService and is not handled here
		}
	}
}

func TestUpdateUserConflict(t *testing.T) {
	db, dbErr := unittests.CreateMySQLTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := mysql.NewAdapterWithDB(db)

	updateUser := &domain.UserUpdate{
		Username: "user123",
		Email:    "user123@localhost",
		Status:   domain.UserActive,
		RoleName: "admin",
	}

	affected, err := dbAdapter.UpdateUser(t.Context(), "uuid-user-67890", updateUser)
	if assert.Error(t, err, "UpdateUser should return an error due to username conflict") {
		assert.Equal(t, int64(0), affected, "Expected no users to be updated due to conflict")
	}
}

func TestChangePassword(t *testing.T) {
	db, dbErr := unittests.CreateMySQLTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := mysql.NewAdapterWithDB(db)

	// Update password for existing user
	err := dbAdapter.ChangePassword(t.Context(), "uuid-user-12345", "newpassword123")
	if assert.NoError(t, err, "UpdatePassword should not return an error for existing user") {
		retrievedUser, getErr := dbAdapter.GetUser(t.Context(), "uuid-user-12345")
		if assert.NoError(t, getErr, "GetUser after UpdatePassword should not return an error") {
			assert.Equal(t, "user123", retrievedUser.Username, "Username should remain unchanged")
			// password change validation is part of testing UserService and is not handled here
		}
	}

	// Attempt to update password for non-existing user
	err = dbAdapter.ChangePassword(t.Context(), "non-existing-uuid", "newpassword123")
	assert.NoError(t, err, "UpdatePassword for non-existing user should not return an error")
}

func TestGetUserByUsername(t *testing.T) {
	db, dbErr := unittests.CreateMySQLTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := mysql.NewAdapterWithDB(db)

	user, err := dbAdapter.GetUserByUsername(t.Context(), "user123")
	if assert.NoError(t, err, "GetUserByUsername should not return an error for existing user") {
		assert.Equal(t, "uuid-user-12345", user.ID, "Expected to retrieve the correct user by username")
		assert.Equal(t, "user123", user.Username, "Expected the user to have username 'user123'")
		assert.Equal(t, "admin", user.RoleName, "Expected the user to have role 'admin'")
	}

	user, err = dbAdapter.GetUserByUsername(t.Context(), "non-existing-user")
	if assert.Error(t, err, "GetUserByUsername for non-existing user should return an error") {
		assert.Nil(t, user, "Expected no user to be returned for non-existing username")
	}

	// get inactive user by username should return nil user and an error
	user, err = dbAdapter.GetUserByUsername(t.Context(), "uuid-user-67890")
	if assert.Error(t, err, "GetUserByUsername for inactive user should not return an error") {
		assert.Nil(t, user, "Expected no user to be returned for inactive user")
	}
}

func TestSetUserVerified(t *testing.T) {
	db, dbErr := unittests.CreateMySQLTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := mysql.NewAdapterWithDB(db)

	// can activate only inactive user
	user, err := dbAdapter.SetUserVerified(t.Context(), "activation-token-67890")
	// must be no error and the user (active now) must be returned
	if assert.NoError(t, err, "SetUserVerified should not return an error for existing user") {
		assert.Equal(t, "uuid-user-67890", user.ID, "Expected to retrieve the correct user by activation token")
		assert.Equal(t, "user678", user.Username, "Expected the user to have username 'user678'")
		assert.EqualValues(t, domain.UserActive, user.Status,
			"Expected the user to be active after verification")
	}

	// Attempt to verify an already active user
	user, err = dbAdapter.SetUserVerified(t.Context(), "activation-token-67890")
	if assert.Error(t, err, "SetUserVerified for already active user should return an error") {
		assert.Nil(t, user, "Expected no user to be returned for already active user")
	}
}

func TestGetUserAuthKey(t *testing.T) {
	db, dbErr := unittests.CreateMySQLTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := mysql.NewAdapterWithDB(db)

	authKey, err := dbAdapter.GetUserAuthKey(t.Context(), "uuid-user-12345")
	if assert.NoError(t, err, "GetUserAuthKey should not return an error for existing user") {
		assert.NotEmpty(t, authKey, "Expected a non-empty auth key for the user")
	} else {
		t.Errorf("GetUserAuthKey returned an error: %v", err)
	}

	authKey, err = dbAdapter.GetUserAuthKey(t.Context(), "non-existing-uuid")
	if assert.Error(t, err, "GetUserAuthKey for non-existing user should return an error") {
		assert.Empty(t, authKey, "Expected no auth key to be returned for non-existing UUID")
	}
}

func TestUpdateUserAuthKey(t *testing.T) {
	db, dbErr := unittests.CreateMySQLTestEngine()
	if dbErr != nil {
		t.Fatalf("test DB error, %v", dbErr)
	}

	if err := unittests.CreateTestDatabase(db); err != nil {
		t.Fatalf("test DB error, %v", err)
	}

	dbAdapter := mysql.NewAdapterWithDB(db)
	newAuthKey := make([]byte, 32)

	err := dbAdapter.UpdateUserAuthKey(t.Context(), "uuid-user-12345", newAuthKey)
	if assert.NoError(t, err, "UpdateUserAuthKey should not return an error for existing user") {
		authKey, getErr := dbAdapter.GetUserAuthKey(t.Context(), "uuid-user-12345")
		if assert.NoError(t, getErr, "GetUserAuthKey after UpdateUserAuthKey should not return an error") {
			assert.Equal(t, newAuthKey, authKey, "Expected the auth key to be updated")
		}
	}

	err = dbAdapter.UpdateUserAuthKey(t.Context(), "non-existing-uuid", newAuthKey)
	assert.NoError(t, err, "UpdateUserAuthKey for non-existing user should return no error")
}
