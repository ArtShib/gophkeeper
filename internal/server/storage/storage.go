package storage

import (
	"context"
	"fmt"

	"github.com/ArtShib/gophkeeper/internal/models"
	"github.com/ArtShib/gophkeeper/internal/server/storage/postgres"
)

type Storage interface {
	Close() error
	AddUser(ctx context.Context, login string, passHash []byte, createdAT int64) (*models.User, error)
	GetUser(ctx context.Context, login string) (*models.User, error)
	AddSecret(ctx context.Context, secretID string, ownerID int64, typeSecret string, data []byte, metadata string, createdAT int64) error
	GetUserSecrets(ctx context.Context, userID int64) (models.ArraySecret, error)
	UpdateSecret(ctx context.Context, secretID string, ownerID int64, typeSecret string, data []byte, metadata string, updatedAT int64) error
	DeleteSecret(ctx context.Context, secretID string, ownerID int64, typeSecret string, updatedAT int64) error
}

func New(ctx context.Context, dsn string) (Storage, error) {
	const op = "storage.NewStorage"

	store, err := postgres.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return store, nil
}
