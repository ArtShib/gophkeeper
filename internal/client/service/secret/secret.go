package secret

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/ArtShib/gophkeeper/internal/lib/crypto"
	"github.com/ArtShib/gophkeeper/internal/lib/loghelper"
	"github.com/ArtShib/gophkeeper/internal/models"
)

//type StoreSecret interface {
//	AddSecret(ctx context.Context, secretID string, userID int64, typeSecret string, data []byte, metadata string, createdAT int64, status string) error
//	GetUserSecrets(ctx context.Context, userID int64) (models.ArraySecret, error)
//	UpdateSecret(ctx context.Context, secretID string, userID int64, typeSecret string, data []byte, metadata string, updatedAT int64, status string) error
//	MarkDeleteSecret(ctx context.Context, secretID string, ownerID int64, updatedAT int64, status string) error
//	ListUserSecrets(ctx context.Context, userID int64) (models.ListSecrets, error)
//	GetSecretsToSync(ctx context.Context, userID int64) (models.ArraySecret, error)
//	MarkSynced(ctx context.Context, secretID string, ownerID int64, isDeleted bool, status string) error
//}

type SecretStorage[T any] interface {
	AddSecret(ctx context.Context, secretID string, userId int64, typeSecret string, data []byte, metadata string, createdAT int64) error
	GetUserSecrets(ctx context.Context, userID int64) (models.ArraySecret, error)
	UpdateSecret(ctx context.Context, secretID string, userID int64, typeSecret string, data []byte, metadata string, updatedAT int64) error
	MarkDeleteSecret(ctx context.Context, secretID string, userID int64, typeSecret string, updatedAT int64) error
	SetStatus(ctx context.Context, secretID string, status string) error
}

type SecretSvc struct {
	store     SecretStorage[models.Secret]
	logger    *slog.Logger
	cryptoSvc *crypto.CryptoService
}

func New(store SecretStorage[models.Secret], logger *slog.Logger, cryptoSvc *crypto.CryptoService) *SecretSvc {
	return &SecretSvc{store: store, logger: logger, cryptoSvc: cryptoSvc}
}

func (s *SecretSvc) CreateSecret(
	ctx context.Context,
	userID int64,
	typeSecret models.SecretType,
	data []byte,
	metadata models.SecretMetadata,
) (*models.Secret, error) {

	log := loghelper.New(s.logger, "CreateSecret")

	id, err := crypto.GenerateUUID()
	if err != nil {
		return nil, log.LogAndReturnError(ctx, "secret create failed", err)
	}

	encryptData, err := s.cryptoSvc.Encrypt(ctx, data)
	if err != nil {
		return nil, log.LogAndReturnError(ctx, "secret encrypt failed", err)
	}

	return &models.Secret{
		ID:        id,
		UserID:    userID,
		Type:      typeSecret,
		Data:      encryptData,
		Metadata:  metadata,
		CreatedAt: time.Now().Unix(),
		Status:    models.StatusNew,
	}, nil
}

func (s *SecretSvc) AddSecret(ctx context.Context, secret *models.Secret) error {
	log := loghelper.New(s.logger, "secret.AddSecret")

	metadata, err := json.Marshal(secret.Metadata)
	if err != nil {
		return log.LogAndReturnError(ctx, "json.Marshal(secret.Metadata)", err)
	}

	encryptData, err := s.cryptoSvc.Encrypt(ctx, secret.Data)
	if err != nil {
		return log.LogAndReturnError(ctx, "s.cryptoSvc.Encrypt(secret.Data)", err)
	}

	createdAt := time.Now().Unix()

	if err := s.store.AddSecret(
		ctx, secret.ID,
		secret.UserID,
		string(secret.Metadata.Type),
		encryptData, string(metadata),
		createdAt); err != nil {
		return log.LogAndReturnError(ctx, "s.store.AddSecret(ctx)", err)
	}

	if err := s.store.SetStatus(ctx, secret.ID, string(secret.Status)); err != nil {
		return log.LogAndReturnError(ctx, "s.store.AddSecret(ctx)", err)
	}

	return nil
}

func (s *SecretSvc) GetUserSecrets(ctx context.Context, userID int64) (models.ArraySecret, error) {
	log := loghelper.New(s.logger, "secret.GetUserSecrets")
	secrets, err := s.store.GetUserSecrets(ctx, userID)
	if err != nil {
		return models.ArraySecret{}, log.LogAndReturnError(ctx, "s.store.GetUserSecrets(ctx)", err)
	}
	return secrets, nil
}

func (s *SecretSvc) UpdateSecret(ctx context.Context, secret *models.Secret) error {
	log := loghelper.New(s.logger, "secret.UpdateSecret")
	metadata, err := json.Marshal(secret.Metadata)
	if err != nil {
		return log.LogAndReturnError(ctx, "json.Marshal(secret.Metadata)", err)
	}
	updatedAT := time.Now().Unix()

	if err := s.store.UpdateSecret(
		ctx, secret.ID,
		secret.UserID,
		string(secret.Metadata.Type),
		secret.Data, string(metadata),
		updatedAT); err != nil {
		return log.LogAndReturnError(ctx, "s.store.UpdateSecret(ctx)", err)
	}

	status := string(models.StatusModified)
	if err := s.store.SetStatus(ctx, secret.ID, status); err != nil {
		return log.LogAndReturnError(ctx, "s.store.AddSecret(ctx)", err)
	}
	return nil
}

func (s *SecretSvc) MarkDeleteSecret(ctx context.Context, secretID string, userId int64, secretType models.SecretType) error {
	log := loghelper.New(s.logger, "secret.MarkDeleteSecret")

	updatedAT := time.Now().Unix()

	if err := s.store.MarkDeleteSecret(ctx, secretID, userId, string(secretType), updatedAT); err != nil {
		return log.LogAndReturnError(ctx, "s.store.MarkDeleteSecret(ctx)", err)
	}

	status := string(models.StatusDeleted)
	if err := s.store.SetStatus(ctx, secretID, status); err != nil {
		return log.LogAndReturnError(ctx, "s.store.AddSecret(ctx)", err)
	}

	return nil
}
