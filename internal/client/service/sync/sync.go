package sync

import (
	"context"
	"log/slog"
	"time"

	"github.com/ArtShib/gophkeeper/internal/lib/crypto"
	"github.com/ArtShib/gophkeeper/internal/lib/loghelper"
	"github.com/ArtShib/gophkeeper/internal/models"
	"google.golang.org/protobuf/types/known/emptypb"
)

type SecretGRPC interface {
	CreateSecret(ctx context.Context, secret *models.Secret) error
	ListSecrets(ctx context.Context, empty *emptypb.Empty) (models.ArraySecret, error)
	UpdateSecret(ctx context.Context, secret *models.Secret) error
	DeleteSecret(ctx context.Context, id string, updatedAt int64) error
}

type SyncService struct {
	logger     *slog.Logger
	SecretSvc  *SecretService
	SecretGRPC SecretGRPC
	userId     int64
	cryptoSvc  *crypto.CryptoService
}

func NewSyncServic(logger *slog.Logger, SecretSvc *SecretService, SecretGRPC SecretGRPC, userId int64, cryptoSvc *crypto.CryptoService) *SyncService {
	return &SyncService{
		logger:     logger,
		SecretSvc:  SecretSvc,
		SecretGRPC: SecretGRPC,
		userId:     userId,
		cryptoSvc:  cryptoSvc,
	}
}

func (s *SyncService) Sync(ctx context.Context) error {
	log := loghelper.New(s.logger, "SyncService")
	isSync := true

	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	if err := s.pushChanges(ctx); err != nil {
		log.LogError(ctx, "pushChanges", err)
		isSync = false
	}

	if err := s.pullChanges(ctx); err != nil {
		log.LogError(ctx, "pullChanges", err)
		isSync = false
	}
	if isSync {
		return nil
	}
	return models.ErrSyncService
}

func (s *SyncService) pullChanges(ctx context.Context) error {
	log := loghelper.New(s.logger, "syncservice.pullChanges")

	var outError error

	secretsServer, err := s.SecretGRPC.ListSecrets(ctx, &emptypb.Empty{})
	if err != nil {
		return log.LogAndReturnError(ctx, "Server.SecretGRPC.ListSecrets", err)
	}

	secretsClient, err := s.SecretSvc.ListUserSecrets(ctx, s.userId)
	if err != nil {
		return log.LogAndReturnError(ctx, "Client.SecretSvc.ListUserSecrets", err)
	}

	for _, secretOut := range secretsServer {
		secretIn, ok := secretsClient[secretOut.ID]
		_, err = s.cryptoSvc.Decrypt(ctx, secretOut.Data)
		if err != nil {
			log.LogError(ctx, "secretSvc.Decrypt", err, slog.String("secretOut.ID", secretOut.ID))
			outError = models.ErrSyncPullChanges
		} else {
			if !ok {
				if err = s.SecretSvc.AddSecret(ctx, &secretOut); err != nil {
					log.LogError(ctx, "SecretSvc.AddSecret", err, slog.String("secretOut.ID", secretOut.ID))
					outError = models.ErrSyncPullChanges
				}
			} else {
				if secretOut.UpdatedAt >= secretIn.UpdatedAt {
					if secretOut.IsDeleted {
						if err = s.SecretSvc.DeleteSecret(ctx, &secretOut); err != nil {
							log.LogError(ctx, "SecretSvc.DeleteSecret", err)
							outError = models.ErrSyncPullChanges
						}
					}
					if err = s.SecretSvc.UpdateSecret(ctx, &secretOut); err != nil {
						log.LogError(ctx, "Store.MarkSynced", err, slog.String("secretOut.ID", secretOut.ID))
						outError = models.ErrSyncPullChanges
					}
				} else {
					log.LogError(ctx, "fail date", models.ErrDateSecretServer, slog.String("secretOut.ID", secretOut.ID))
					outError = models.ErrSyncPullChanges
				}
			}
		}
	}
	return outError
}

func (s *SyncService) pushChanges(ctx context.Context) error {
	log := loghelper.New(s.logger, "syncservice.pushChanges")

	var outError error

	secretsClient, err := s.SecretSvc.GetSecretsToSync(ctx, s.userId)
	if err != nil {
		return log.LogAndReturnError(ctx, "Store.GetSecretsToSync", err)
	}

	for _, secretIn := range secretsClient {
		isMark := false
		switch secretIn.Status {
		case models.StatusNew:
			if err = s.SecretGRPC.CreateSecret(ctx, &secretIn); err != nil {
				log.LogError(ctx, "SecretGRPC.CreateSecret", err)
				outError = models.ErrSyncPushChanges
			}
			isMark = true
		case models.StatusModified:
			if err = s.SecretGRPC.UpdateSecret(ctx, &secretIn); err != nil {
				log.LogError(ctx, "SecretGRPC.UpdateSecret", err)
				outError = models.ErrSyncPushChanges
			}
			isMark = true
		case models.StatusDeleted:
			if err = s.SecretGRPC.DeleteSecret(ctx, secretIn.ID, secretIn.UpdatedAt); err != nil {
				log.LogError(ctx, "SecretGRPC.MarkDeleteSecret", err)
				outError = models.ErrSyncPushChanges
			}
			isMark = true
		}
		if isMark {
			if err = s.SecretSvc.MarkSynced(ctx, &secretIn); err != nil {
				log.LogError(ctx, "SecretSvc.MarkSynced", err)
				outError = models.ErrSyncPushChanges
			}
		}
	}
	return outError
}
