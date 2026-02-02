package secret

import (
	"context"

	keeperv1 "github.com/ArtShib/gophkeeper/gen/go/keeper/v1"
	"github.com/ArtShib/gophkeeper/internal/models"
)

type KeeperService interface {
	AddSecret(ctx context.Context, s *models.Secret) error
	GetUserSecrets(ctx context.Context, userID int64) (models.ArraySecret, error)
	UpdateSecret(ctx context.Context, s *models.Secret) error
	DeleteSecret(ctx context.Context, s *models.Secret) error
}

type serverAPI struct {
	keeperv1.UnimplementedSecretServiceServer
	service KeeperService
}
