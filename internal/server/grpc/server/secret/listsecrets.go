package secret

import (
	"context"

	"github.com/ArtShib/gophkeeper/internal/lib/customprototype"

	keeperv1 "github.com/ArtShib/gophkeeper/gen/go/keeper/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *serverAPI) GetUserSecrets(ctx context.Context, empty *emptypb.Empty) (*keeperv1.ListSecretsResponse, error) {
	userID, err := getUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	arraySecret, err := s.service.GetUserSecrets(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to fetch secrets")
	}

	secrets := make([]*keeperv1.Secret, len(arraySecret))

	for i, secret := range arraySecret {
		//var meta keeperv1.SecretMetadata
		//if err := json.Unmarshal(secret.Metadata, &meta); err != nil {
		//
		//	return nil, status.Errorf(codes.Internal, "failed to unmarshal meta for secret %s", secret.ID)
		//}
		secrets[i] = keeperv1.Secret_builder{
			Id: proto.String(secret.ID),
			Meta: keeperv1.SecretMetadata_builder{
				Type:  customprototype.GetProtoSecretType(secret.Metadata.Type),
				Name:  proto.String(secret.Metadata.Name),
				Extra: secret.Metadata.Extra,
			}.Build(),
			EncryptedData: secret.Data,
			CreatedAt:     proto.Int64(secret.CreatedAt),
			UpdatedAt:     proto.Int64(secret.UpdatedAt),
		}.Build()

	}

	return keeperv1.ListSecretsResponse_builder{
		Secrets: secrets,
	}.Build(), nil
}
