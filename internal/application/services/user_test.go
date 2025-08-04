package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/utking/spaces/internal/application/domain"
	"github.com/utking/spaces/internal/application/services"
	"github.com/utking/spaces/internal/ports"
)

func TestCreateUser(t *testing.T) {
	newUser := &domain.User{
		Username:        "testuser",
		Email:           "test@localhost",
		Password:        "testpassword",
		PasswordConfirm: "testpassword", // Should match for creation
	}

	retUser := &domain.User{
		ID:       "some-user-id",
		AuthKey:  "some-auth-key",
		Username: newUser.Username,
		Email:    newUser.Email,
	}

	fsPort := ports.NewMockFileSystem(t)
	dbPort := ports.NewMockDBPort(t)

	dbPort.On("CreateUser", mock.Anything, mock.Anything).Return(retUser.ID, retUser.AuthKey, nil)
	dbPort.On("GetUser", mock.Anything, retUser.ID).Return(retUser, nil)

	svc := services.NewUsersService(dbPort, fsPort)

	// Create a new user
	uid, authKey, cErr := svc.Create(t.Context(), newUser)
	if assert.NoError(t, cErr) {
		assert.NotEmpty(t, uid)
		assert.NotEmpty(t, authKey)

		// Verify the user was created in the database
		user, err := dbPort.GetUser(t.Context(), uid)
		if assert.NoError(t, err) {
			assert.Equal(t, newUser.Username, user.Username)
			assert.Equal(t, newUser.Email, user.Email)
			assert.NotEmpty(t, user.ID)
			assert.NotEmpty(t, user.AuthKey)
		}
	}
}
