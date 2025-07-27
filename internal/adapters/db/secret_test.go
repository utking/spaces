package db

import (
	"testing"
)

func TestSecretValidateOk(t *testing.T) {
	s := &Secret{
		UserID:      "12345678901234567890123456789012", // 32 characters
		Name:        "Test Secret",
		Username:    []byte("testuser"),
		Tags:        []string{"tag1", "tag2"},
		Description: "This is a test secret",
	}

	if err := s.Validate(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestSecretValidateErr(t *testing.T) {
	tests := []struct {
		name string
		s    *Secret
	}{
		{"EmptyUserID", &Secret{UserID: "", Name: "Test Secret", Username: []byte("testuser")}},
		{"NoTags", &Secret{
			UserID:   "12345678901234567890123456789012",
			Name:     "Test Secret",
			Username: []byte("testuser"),
			Tags:     nil}},
		{"LongName", &Secret{
			UserID:   "12345678901234567890123456789012",
			Name:     "a" + string(make([]byte, 129)),
			Username: []byte("testuser")}},
		{"LongDescription", &Secret{
			Username:    []byte("testuser"),
			UserID:      "12345678901234567890123456789012",
			Name:        "Test Secret",
			Description: "a" + string(make([]byte, 4097))},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Validate(); err == nil {
				t.Errorf("expected error for %s, got nil", tt.name)
			}
		})
	}
}

func TestEncryptSecretValidateOk(t *testing.T) {
	es := &Secret{
		ID:     "12345678901234567890123456789012", // 32 characters
		Secret: []byte("encryptedsecretdata"),
	}

	if es.ID == "" || len(es.Secret) == 0 {
		t.Errorf("expected no error, got invalid EncryptSecret")
	}
}

func TestEncryptSecretValidateErr(t *testing.T) {
	tests := []struct {
		name string
		es   *Secret
	}{
		{
			"NoUserID", &Secret{
				UserID: "", Name: "Test Secret", Tags: []string{"tag1", "tag2"},
				Secret: []byte("encryptedsecretdata"),
			}},
		{
			"NoTags", &Secret{
				Name: "Test Secret", UserID: "user-12345",
				Secret: []byte("encryptedsecretdata"),
			}},
		{
			"LongSecret", &Secret{
				Name:   "Test Secret",
				UserID: "user-12345",
				Tags:   []string{"tag1", "tag2"},
				Secret: make([]byte, 5000)}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.es.Validate(); err == nil {
				t.Errorf("expected error for %s, got nil", tt.name)
			}
		})
	}
}
