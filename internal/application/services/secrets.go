package services

import (
	"context"
	"errors"

	"gogs.utking.net/utking/spaces/internal/application/domain"
	"gogs.utking.net/utking/spaces/internal/ports"
)

// SecretService is a struct that implements the SecretService interface.
type SecretService struct {
	db      ports.DBPort
	cryptor ports.CryptoService
}

// NewSecretService creates a new instance of the SecretService struct.
func NewSecretService(db ports.DBPort, cryptor ports.CryptoService) *SecretService {
	return &SecretService{
		db:      db,
		cryptor: cryptor,
	}
}

// GetTags retrieves a list of secret tags from the database service.
func (a *SecretService) GetTags(
	ctx context.Context,
	uid string,
) ([]string, error) {
	return a.db.GetSecretTags(ctx, uid)
}

// GetItems retrieves a list of secrets from the database service.
func (a *SecretService) GetItems(
	ctx context.Context,
	uid string,
	req *domain.SecretSearchRequest,
) ([]domain.Secret, error) {
	return a.db.GetSecrets(ctx, uid, req)
}

func (a *SecretService) GetCount(
	ctx context.Context,
	uid string,
	req *domain.SecretSearchRequest,
) (int64, error) {
	return a.db.GetSecretsCount(ctx, uid, req)
}

func (a *SecretService) GetItem(
	ctx context.Context,
	uid, id string,
) (*domain.Secret, error) {
	return a.db.GetSecret(ctx, uid, id)
}

func (a *SecretService) Create(
	ctx context.Context,
	uid string,
	req *domain.Secret,
) (string, error) {
	return a.db.CreateSecret(ctx, uid, req)
}

func (a *SecretService) Update(
	ctx context.Context,
	uid, id string,
	req *domain.Secret,
) (int64, error) {
	// id must be given
	if id == "" {
		return 0, errors.New("secret ID must be provided")
	}

	return a.db.UpdateSecret(ctx, uid, id, req)
}

func (a *SecretService) Delete(ctx context.Context, uid, id string) error {
	// id must be given
	if id == "" {
		return errors.New("secret ID must be provided")
	}

	return a.db.DeleteSecret(ctx, uid, id)
}

func (a *SecretService) Encrypt(
	ctx context.Context,
	req *domain.SecretEncodeRequest,
	key []byte,
) (nonce, encoded []byte, err error) {
	return a.cryptor.Encrypt(ctx, req, key)
}

func (a *SecretService) Decrypt(
	ctx context.Context,
	nonce []byte,
	encoded []byte,
	key []byte,
) (decoded []byte, err error) {
	return a.cryptor.Decrypt(ctx, nonce, encoded, key)
}

func (a *SecretService) GetItemsMap(
	ctx context.Context,
	uid string,
	req *domain.SecretSearchRequest,
) ([]domain.SecretExportItem, error) {
	return a.db.GetSecretsMap(ctx, uid, req)
}

// SearchItemsByTerm searches for secrets by a term in the database service.
func (a *SecretService) SearchItemsByTerm(
	ctx context.Context,
	uid string,
	req *domain.SecretRequest,
) ([]domain.Secret, error) {
	return a.db.SearchSecretsByTerm(ctx, uid, req)
}

func (a *SecretService) UpdateEncryptedSecrets(
	ctx context.Context,
	uid string,
	items map[string]domain.EncryptSecret,
) error {
	return a.db.UpdateEncryptedSecrets(ctx, uid, items)
}
