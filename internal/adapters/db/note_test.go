package db

import "testing"

func TestNoteValidateOk(t *testing.T) {
	n := &Note{
		UserID: "12345678901234567890123456789012",
		Title:  "Test Note",
	}

	if err := n.Validate(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestNoteValidateErr(t *testing.T) {
	tests := []struct {
		name string
		n    *Note
	}{
		{"EmptyUserID", &Note{UserID: "", Title: "Test Note"}},
		{"EmptyTitle", &Note{UserID: "12345678901234567890123456789012", Title: ""}},
		{"LongTitle", &Note{UserID: "12345678901234567890123456789012", Title: "a" + string(make([]byte, 129))}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.Validate(); err == nil {
				t.Errorf("expected error for %s, got nil", tt.name)
			}
		})
	}
}
