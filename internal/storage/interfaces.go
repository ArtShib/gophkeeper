package storage

import (
	"context"

	"github.com/ArtShib/gophkeeper/internal/models"
)

type UserStorage[T any] interface {
	AddUser(ctx context.Context, user *models.User) (int64, error)
	GetUser(ctx context.Context, login string) (*models.User, error)
}

type SecretStorage[T any] interface {
	AddSecret(ctx context.Context, secretID string, userId int64, typeSecret string, data []byte, metadata string, createdAT int64) error
	GetUserSecrets(ctx context.Context, userID int64) (models.ArraySecret, error)
	UpdateSecret(ctx context.Context, secretID string, userID int64, typeSecret string, data []byte, metadata string, updatedAT int64) error
	MarkDeleteSecret(ctx context.Context, secretID string, userID int64, typeSecret string, updatedAT int64) error
	SetStatus(ctx context.Context, secretID string, status string) error
}

type SyncStorage[T any] interface {
	AddSecret(ctx context.Context, secretID string, userId int64, typeSecret string, data []byte, metadata string, createdAT int64) error
	DeleteSecret(ctx context.Context, secretID string, userID int64, typeSecret string) error
	ListUserSecrets(ctx context.Context, userID int64) (models.ListSecrets, error)
	GetSecretsToSync(ctx context.Context, userID int64) (models.ArraySecret, error)
	MarkSynced(ctx context.Context, secretID string, ownerID int64, status string) error
	UpdateSecret(ctx context.Context, secretID string, userID int64, typeSecret string, data []byte, metadata string, updatedAT int64) error
}
