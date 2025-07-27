// Package cryptor provides an implementation of the Cryptor interface
package cryptor

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"gogs.utking.net/utking/spaces/internal/application/domain"
)

type CryptoKey []byte

func (k CryptoKey) String() string {
	return "(redacted)"
}

type Cryptor struct{}

func New() *Cryptor {
	return &Cryptor{}
}

// Encrypt encrypts the given plaintext using AES-GCM with a random nonce.
// INFO: Don't use more than 2^32 random nonces with a given key
// because of the risk of a repeat.
func (a *Cryptor) Encrypt(
	_ context.Context,
	req *domain.SecretEncodeRequest,
	key []byte,
) (nonce, encoded []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	nonce = make([]byte, 12)
	if _, rErr := io.ReadFull(rand.Reader, nonce); rErr != nil {
		return nil, nil, rErr
	}

	aesgcm, gcmErr := cipher.NewGCM(block)
	if gcmErr != nil {
		panic(gcmErr.Error())
	}

	encoded = aesgcm.Seal(nil, nonce, req.PlainText, nil)

	return nonce, encoded, nil
}

func (a *Cryptor) Decrypt(_ context.Context, nonce, encoded, key []byte) (decoded []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	decoded, err = aesgcm.Open(nil, nonce, encoded, nil)
	if err != nil {
		return nil, err
	}

	return decoded, nil
}
