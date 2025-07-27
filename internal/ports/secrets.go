package ports

import (
	"context"

	"github.com/utking/spaces/internal/application/domain"
)

// SecretService is an interface that defines the methods for secrets-related operations.
type SecretService interface {
	GetTags(ctx context.Context, uid string) ([]string, error)
	GetItems(ctx context.Context, uid string, req *domain.SecretSearchRequest) ([]domain.Secret, error)
	SearchItemsByTerm(ctx context.Context, uid string, req *domain.SecretRequest) ([]domain.Secret, error)
	GetCount(ctx context.Context, uid string, req *domain.SecretSearchRequest) (int64, error)
	GetItem(ctx context.Context, uid, id string) (*domain.Secret, error)
	Create(ctx context.Context, uid string, req *domain.Secret) (string, error)
	Update(ctx context.Context, uid, id string, req *domain.Secret) (int64, error)
	Delete(ctx context.Context, uid, id string) error
	UpdateEncryptedSecrets(ctx context.Context, uid string, items map[string]domain.EncryptSecret) error

	// For import-export
	GetItemsMap(
		ctx context.Context,
		uid string,
		req *domain.SecretSearchRequest,
	) ([]domain.SecretExportItem, error)

	// Encrypt and Decrypt
	Encrypt(
		ctx context.Context,
		req *domain.SecretEncodeRequest,
		key []byte,
	) (nonce, encoded []byte, err error)
	Decrypt(
		ctx context.Context,
		nonce []byte,
		encoded []byte,
		key []byte,
	) (decoded []byte, err error)
}

// CryptoService is an interface that defines the methods for cryptographic operations.
type CryptoService interface {
	Encrypt(
		ctx context.Context,
		req *domain.SecretEncodeRequest,
		key []byte,
	) (nonce, encoded []byte, err error)
	Decrypt(
		ctx context.Context,
		nonce []byte,
		encoded []byte,
		key []byte,
	) (decoded []byte, err error)
}
