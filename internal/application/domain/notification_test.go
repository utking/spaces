package domain_test

import (
	"testing"

	"gogs.utking.net/utking/spaces/internal/application/domain"
)

func TestNotificationValidationOk(t *testing.T) {
	n := &domain.Notification{
		Title:   "Test Notification",
		Message: "This is a test notification message.",
		To:      "to@email.local",
	}

	if err := n.Validate(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestNotificationValidationErr(t *testing.T) {
	tests := []struct {
		name string
		n    *domain.Notification
	}{
		{"EmptyTitle", &domain.Notification{Title: "", Message: "some", To: "some@email"}},
		{"EmptyMessage", &domain.Notification{Title: "some", Message: "", To: "some@email"}},
		{"EmptyTo", &domain.Notification{Title: "some", Message: "some", To: ""}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.Validate(); err == nil {
				t.Errorf("expected error for %s, got nil", tt.name)
			}
		})
	}
}
