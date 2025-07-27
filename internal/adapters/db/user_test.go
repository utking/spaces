package db

import "testing"

func TestUserItemValidateOk(t *testing.T) {
	u := &User{
		Username: "testuser",
		Email:    "some@email",
	}

	if err := u.Validate(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestUserItemValidateErr(t *testing.T) {
	tests := []struct {
		name string
		u    *User
	}{
		{"EmptyUsername", &User{Username: "", Email: "some@email"}},
		{"ShortUsername", &User{Username: "usr", Email: "some@email"}},
		{"LongUsername", &User{Username: "a" + string(make([]byte, 256)), Email: "some@email"}},
		{"EmptyEmail", &User{Username: "testuser", Email: ""}},
		{"LongEmail", &User{Username: "testuser", Email: "a" + string(make([]byte, 256)) + "@example.com"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.u.Validate(); err == nil {
				t.Errorf("expected error for %s, got nil", tt.name)
			}
		})
	}
}

func TestUserItemValidateUpdateOk(t *testing.T) {
	u := &User{
		Email: "updated@email",
	}

	if err := u.ValidateUpdate(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestUserItemValidateUpdateErr(t *testing.T) {
	tests := []struct {
		name string
		u    *User
	}{
		{"EmptyEmail", &User{Email: ""}},
		{"LongEmail", &User{Email: "a" + string(make([]byte, 256)) + "@example.com"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.u.ValidateUpdate(); err == nil {
				t.Errorf("expected error for %s, got nil", tt.name)
			}
		})
	}
}

func TestAuthAssignmentValidateOk(t *testing.T) {
	a := &AuthAssignment{
		UserID:   "123e4567-e89b-12d3-a456-426614174000",
		RoleName: "admin",
	}

	if err := a.Validate(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestAuthAssignmentValidateErr(t *testing.T) {
	tests := []struct {
		name string
		a    *AuthAssignment
	}{
		{"EmptyUserID", &AuthAssignment{UserID: "", RoleName: "admin"}},
		{"EmptyRoleName", &AuthAssignment{UserID: "123e4567-e89b-12d3-a456-426614174000", RoleName: ""}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.a.Validate(); err == nil {
				t.Errorf("expected error for %s, got nil", tt.name)
			}
		})
	}
}
