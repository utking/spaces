package cryptor

import (
	"testing"

	"github.com/utking/spaces/internal/application/domain"
)

func TestCryptor(t *testing.T) {
	// Initialize the cryptor
	cryptor := New()

	// Define a test key and plaintext
	key := []byte("thisis32bitlongpassphraseimusing")
	plaintext := []byte("Hello, World!")

	// Encrypt the plaintext
	nonce, encoded, err := cryptor.Encrypt(t.Context(), &domain.SecretEncodeRequest{PlainText: plaintext}, key)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	// Decrypt the encoded data
	decoded, err := cryptor.Decrypt(t.Context(), nonce, encoded, key)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	// Check if the decrypted data matches the original plaintext
	if string(decoded) != string(plaintext) {
		t.Errorf("Decrypted text does not match original: got %s, want %s", decoded, plaintext)
	}
}
