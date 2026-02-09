package sync

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/ArtShib/gophkeeper/internal/lib/crypto"
	"github.com/ArtShib/gophkeeper/internal/lib/loghelper"
	"github.com/ArtShib/gophkeeper/internal/models"
)

type SyncStorage[T any] interface {
	AddSecret(ctx context.Context, secretID string, userId int64, typeSecret string, data []byte, metadata string, createdAT int64) error
	DeleteSecret(ctx context.Context, secretID string, userID int64, typeSecret string) error
	ListUserSecrets(ctx context.Context, userID int64) (models.ListSecrets, error)
	GetSecretsToSync(ctx context.Context, userID int64) (models.ArraySecret, error)
	MarkSynced(ctx context.Context, secretID string, ownerID int64, status string) error
	UpdateSecret(ctx context.Context, secretID string, userID int64, typeSecret string, data []byte, metadata string, updatedAT int64) error
}

type SecretService struct {
	store     SyncStorage[models.Secret]
	logger    *slog.Logger
	cryptoSvc *crypto.CryptoService
}

func NewSecretSvc(store SyncStorage[models.Secret], logger *slog.Logger, cryptoSvc *crypto.CryptoService) *SecretService {
	return &SecretService{store: store, logger: logger, cryptoSvc: cryptoSvc}
}

func (s *SecretService) AddSecret(ctx context.Context, secret *models.Secret) error {
	log := loghelper.New(s.logger, "secret.AddSecret")

	metadata, err := json.Marshal(secret.Metadata)
	if err != nil {
		return log.LogAndReturnError(ctx, "json.Marshal(secret.Metadata)", err)
	}

	encryptData, err := s.cryptoSvc.Encrypt(ctx, secret.Data)
	if err != nil {
		return log.LogAndReturnError(ctx, "s.cryptoSvc.Encrypt(secret.Data)", err)
	}

	if err := s.store.AddSecret(
		ctx, secret.ID,
		secret.UserID,
		string(secret.Metadata.Type),
		encryptData, string(metadata),
		secret.CreatedAt); err != nil {
		return log.LogAndReturnError(ctx, "s.store.AddSecret(ctx)", err)
	}

	return s.MarkSynced(ctx, secret)
}

func (s *SecretService) DeleteSecret(ctx context.Context, secret *models.Secret) error {
	log := loghelper.New(s.logger, "secret.AddSecret")
	if err := s.store.DeleteSecret(ctx, secret.ID, secret.UserID, string(secret.Type)); err != nil {
		return log.LogAndReturnError(ctx, "s.store.DeleteSecret(ctx)", err)
	}
	return nil
}

func (s *SecretService) ListUserSecrets(ctx context.Context, userID int64) (models.ListSecrets, error) {
	log := loghelper.New(s.logger, "secret.ListUserSecrets")
	secrets, err := s.store.ListUserSecrets(ctx, userID)
	if err != nil {
		return models.ListSecrets{}, log.LogAndReturnError(ctx, "s.store.ListUserSecrets()", err)
	}
	return secrets, nil
}

func (s *SecretService) GetSecretsToSync(ctx context.Context, userID int64) (models.ArraySecret, error) {
	log := loghelper.New(s.logger, "secret.GetSecretsToSync")
	secrets, err := s.store.GetSecretsToSync(ctx, userID)
	if err != nil {
		return models.ArraySecret{}, log.LogAndReturnError(ctx, "s.store.GetSecretsToSync()", err)
	}
	return secrets, nil
}

func (s *SecretService) MarkSynced(ctx context.Context, secret *models.Secret) error {
	log := loghelper.New(s.logger, "secret.MarkSynced")
	if err := s.store.MarkSynced(ctx, secret.ID, secret.UserID, string(models.StatusSynced)); err != nil {
		return log.LogAndReturnError(ctx, "s.store.MarkSynced()", err)
	}
	return nil
}

func (s *SecretService) UpdateSecret(ctx context.Context, secret *models.Secret) error {
	log := loghelper.New(s.logger, "secret.UpdateSecret")

	metadata, err := json.Marshal(secret.Metadata)
	if err != nil {
		return log.LogAndReturnError(ctx, "json.Marshal(secret.Metadata)", err)
	}

	encryptData, err := s.cryptoSvc.Encrypt(ctx, secret.Data)
	if err != nil {
		return log.LogAndReturnError(ctx, "s.cryptoSvc.Encrypt(secret.Data)", err)
	}

	if err := s.store.UpdateSecret(
		ctx, secret.ID,
		secret.UserID,
		string(secret.Metadata.Type),
		encryptData, string(metadata),
		secret.UpdatedAt); err != nil {
		return log.LogAndReturnError(ctx, "s.store.AddSecret(ctx)", err)
	}

	return s.MarkSynced(ctx, secret)
}
