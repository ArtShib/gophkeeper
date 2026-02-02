package keeper

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/ArtShib/gophkeeper/internal/lib/loghelper"
	"github.com/ArtShib/gophkeeper/internal/models"
)

type StoreKeeper interface {
	AddSecret(ctx context.Context, secretID string, ownerID int64, typeSecret string, data []byte, metadata string, createdAT int64) error
	GetUserSecrets(ctx context.Context, userID int64) (models.ArraySecret, error)
	UpdateSecret(ctx context.Context, secretID string, ownerID int64, typeSecret string, data []byte, metadata string, updatedAT int64) error
	DeleteSecret(ctx context.Context, secretID string, ownerID int64, typeSecret string, updatedAT int64) error
}

type Keeper struct {
	log   *slog.Logger
	store StoreKeeper
}

func New(log *slog.Logger, store StoreKeeper) *Keeper {
	return &Keeper{
		log:   log,
		store: store,
	}
}

func (k *Keeper) AddSecret(ctx context.Context, s *models.Secret) error {
	log := loghelper.New(k.log, "Keeper.AddSecret")
	log.LogDebug(ctx, "adding secret", slog.String("id", s.ID))
	if err := k.store.AddSecret(ctx, s.ID, s.UserID, string(s.Type), s.Data, string(s.Metadata), s.CreatedAt); err != nil {
		return log.LogAndReturnError(ctx, "store.AddSecret", err)
	}
	return nil
}

func (k *Keeper) GetUserSecrets(ctx context.Context, userID int64) (models.ArraySecret, error) {
	log := loghelper.New(k.log, "Keeper.GetUserSecrets")
	log.LogDebug(ctx, "getting secrets", slog.String("id", strconv.FormatInt(userID, 10)))
	arraySecret, err := k.store.GetUserSecrets(ctx, userID)
	if err != nil {
		return models.ArraySecret{}, log.LogAndReturnError(ctx, "store.GetUserSecrets", err)
	}
	return arraySecret, nil
}

func (k *Keeper) UpdateSecret(ctx context.Context, s *models.Secret) error {
	log := loghelper.New(k.log, "Keeper.UpdateSecret")
	log.LogDebug(ctx, "updating secret", slog.String("id", s.ID))
	if err := k.store.UpdateSecret(ctx, s.ID, s.UserID, string(s.Type), s.Data, string(s.Metadata), s.UpdatedAt); err != nil {
		return log.LogAndReturnError(ctx, "store.UpdateSecret", err)
	}
	return nil
}

func (k *Keeper) DeleteSecret(ctx context.Context, s *models.Secret) error {
	log := loghelper.New(k.log, "Keeper.DeleteSecret")
	log.LogDebug(ctx, "deleting secret", slog.String("id", s.ID))
	if err := k.store.DeleteSecret(ctx, s.ID, s.UserID, string(s.Type), s.UpdatedAt); err != nil {
		return log.LogAndReturnError(ctx, "store.DeleteSecret", err)
	}
	return nil
}
