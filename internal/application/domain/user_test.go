package domain_test

import (
	"testing"

	"gogs.utking.net/utking/spaces/internal/application/domain"
)

func TestRandInt63(t *testing.T) {
	// Test if RandInt63 returns a non-negative value
	value := domain.Int63()
	if value < 0 {
		t.Errorf("expected non-negative value, got %d", value)
	}
}

func TestGenerateRandomString(t *testing.T) {
	length := 10
	randomString := domain.GenerateRandomString(length)

	// Test if the generated string has the correct length
	if len(randomString) != length {
		t.Errorf("expected length %d, got %d", length, len(randomString))
	}

	// Test if the generated string is not empty
	if randomString == "" {
		t.Error("expected non-empty string, got empty string")
	}
}

func TestGetPasswordHash(t *testing.T) {
	pass := "testpassword"

	hash, err := domain.GetPasswordHash(pass)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if hash == "" {
		t.Error("expected non-empty hash, got empty string")
	}

	// Check if the hash is not the same as the original password
	if hash == pass {
		t.Error("expected hash to be different from password")
	}
}

// TestPasswordVerify tests the PasswordVerify function
// It checks if the function correctly verifies a password against its calculated hash.
func TestPasswordVerify(t *testing.T) {
	pass := "testpassword"
	hash, err := domain.GetPasswordHash(pass)
	if err != nil {
		t.Fatalf("failed to get password hash: %v", err)
	}

	tests := []struct {
		name     string
		password string
		expected bool
	}{
		{"Correct Password", pass, true},
		{"Incorrect Password", "wrongpassword", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if ok := domain.PasswordVerify(tt.password, hash); ok != tt.expected {
				t.Errorf("expected %v for password '%s', got %v", tt.expected, tt.password, ok)
			}
		})
	}
}

func TestUserValidateOk(t *testing.T) {
	u := &domain.User{
		Username:        "testuser",
		Password:        "testpassword",
		PasswordConfirm: "testpassword",
		Email:           "some@email",
	}

	if err := u.Validate(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestUserValidateErr(t *testing.T) {
	tests := []struct {
		name string
		u    *domain.User
	}{
		{
			"EmptyUsername",
			&domain.User{Username: "", Password: "testpassword", PasswordConfirm: "testpassword", Email: "some@email"},
		},
		{"EmptyPassword", &domain.User{Username: "testuser", Password: "", PasswordConfirm: "", Email: "some@email"}},
		{
			"WrongConfirmPassword",
			&domain.User{
				Username:        "testuser",
				Password:        "testpassword",
				PasswordConfirm: "wrongpassword",
				Email:           "some@email",
			},
		},
		{
			"EmptyEmail",
			&domain.User{Username: "testuser", Password: "testpassword", PasswordConfirm: "testpassword", Email: ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.u.Validate(); err == nil {
				t.Errorf("expected error for %s, got nil", tt.name)
			}
		})
	}
}

func TestValidateUserUpdateOk(t *testing.T) {
	u := &domain.UserUpdate{
		Email:    "updated@email",
		RoleName: "admin",
	}

	// Test with no password update
	if err := u.Validate(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Test with password update
	u.Password = "newpassword"
	u.PasswordConfirm = "newpassword"

	if err := u.Validate(); err != nil {
		t.Errorf("expected no error with password, got %v", err)
	}
}

func TestValidateUserUpdateErr(t *testing.T) {
	tests := []struct {
		name string
		u    *domain.UserUpdate
	}{
		{"EmptyRoleName", &domain.UserUpdate{Username: "updateduser", Email: "updated@email", RoleName: ""}},
		{"EmptyEmail", &domain.UserUpdate{Username: "updateduser", Email: ""}},
		{
			"WrongConfirmPassword",
			&domain.UserUpdate{
				Username:        "updateduser",
				Email:           "updated@email",
				Password:        "newpassword",
				PasswordConfirm: "wrongpassword",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.u.Validate(); err == nil {
				t.Errorf("expected error for %s, got nil", tt.name)
			}
		})
	}
}
