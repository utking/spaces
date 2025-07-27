package domain_test

import (
	"testing"

	"github.com/utking/spaces/internal/application/domain"
)

func TestSecretExportRequestValidateOk(t *testing.T) {
	s := &domain.SecretExportRequest{
		// any 32 char long password
		Password: "12345678901234567890123456789012",
	}

	if err := s.Validate(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestSecretExportRequestValidateErr(t *testing.T) {
	s := &domain.SecretExportRequest{
		Password: "",
	}

	if err := s.Validate(); err == nil {
		t.Error("expected error for empty password, got nil")
	}

	s.Password = "short"
	if err := s.Validate(); err == nil {
		t.Error("expected error for short password, got nil")
	}
}
