package storage

import (
	"context"
	"fmt"

	"github.com/ArtShib/gophkeeper/internal/client/models"
	"github.com/ArtShib/gophkeeper/internal/client/storage/sqlite"
)

type Storage interface {
	Stop() error
	AddUser(ctx context.Context, id int64, login string, passHash []byte) error
	GetUser(ctx context.Context, login string, passHash []byte) (*models.User, error)
	AddSecret(ctx context.Context, secretID string, userID int64, typeSecret string, data []byte, metadata string, createdAT int64, status string) error
	GetUserSecrets(ctx context.Context, userID int64) (models.ArraySecret, error)
	UpdateSecret(ctx context.Context, secretID string, userID int64, typeSecret string, data []byte, metadata string, updatedAT int64, status string) error
	MarkDeleteSecret(ctx context.Context, secretID string, ownerID int64, updatedAT int64, status string) error
	MarkSynced(ctx context.Context, secretID string, ownerID int64, isDeleted bool, status string) error
	GetSecretsToSync(ctx context.Context, userID int64) (models.ArraySecret, error)
}

func New(ctx context.Context, path string) (Storage, error) {
	const op = "storage.NewStorage"

	store, err := sqlite.New(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return store, nil
}
