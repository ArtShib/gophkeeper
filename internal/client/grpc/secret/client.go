package secret

import (
	"context"

	keeperv1 "github.com/ArtShib/gophkeeper/gen/go/keeper/v1"
	"github.com/ArtShib/gophkeeper/internal/client/models"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Client struct {
	api keeperv1.SecretServiceClient
}

func New(conn *grpc.ClientConn) *Client {
	return &Client{api: keeperv1.NewSecretServiceClient(conn)}
}

func (c *Client) CreateSecret(ctx context.Context, secret *models.Secret) error {

	_, err := c.api.CreateSecret(ctx, keeperv1.CreateSecretRequest_builder{
		Secret: keeperv1.Secret_builder{
			Id: proto.String(secret.ID),
			Meta: keeperv1.SecretMetadata_builder{
				Type:  getProtoSecretType(secret.Metadata.Type),
				Name:  proto.String(secret.Metadata.Name),
				Extra: secret.Metadata.Extra,
			}.Build(),
			EncryptedData: secret.Data,
			CreatedAt:     proto.Int64(secret.CreatedAt),
		}.Build(),
	}.Build())

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) ListSecrets(ctx context.Context, empty *emptypb.Empty) (models.ArraySecret, error) {
	listSecrets, err := c.api.ListSecrets(ctx, empty)
	if err != nil {
		return nil, err
	}

	secrets := make(models.ArraySecret, len(listSecrets.GetSecrets()))

	for i, secret := range listSecrets.GetSecrets() {
		secrets[i] = models.Secret{
			ID:   secret.GetId(),
			Data: secret.GetEncryptedData(),
			Metadata: models.SecretMetadata{
				Type:  models.SecretType(secret.GetMeta().GetType().String()),
				Name:  secret.GetMeta().GetName(),
				Extra: secret.GetMeta().GetExtra(),
			},
			CreatedAt: secret.GetCreatedAt(),
			UpdatedAt: secret.GetUpdatedAt(),
			IsDeleted: secret.GetIsDeleted(),
		}
	}

	return secrets, nil
}

func (c *Client) UpdateSecret(ctx context.Context, secret *models.Secret) error {

	_, err := c.api.UpdateSecret(ctx, keeperv1.UpdateSecretRequest_builder{
		Secret: keeperv1.Secret_builder{
			Id: proto.String(secret.ID),
			Meta: keeperv1.SecretMetadata_builder{
				Type:  getProtoSecretType(secret.Metadata.Type),
				Name:  proto.String(secret.Metadata.Name),
				Extra: secret.Metadata.Extra,
			}.Build(),
			EncryptedData: secret.Data,
			CreatedAt:     proto.Int64(secret.CreatedAt),
			UpdatedAt:     proto.Int64(secret.UpdatedAt),
		}.Build(),
	}.Build())

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteSecret(ctx context.Context, id string, updatedAt int64) error {
	_, err := c.api.DeleteSecret(ctx, keeperv1.DeleteSecretRequest_builder{
		Id:        proto.String(id),
		UpdatedAt: proto.Int64(updatedAt),
	}.Build())

	if err != nil {
		return err
	}

	return nil
}
