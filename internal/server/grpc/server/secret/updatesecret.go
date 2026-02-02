package secret

import (
	"context"
	"encoding/json"

	keeperv1 "github.com/ArtShib/gophkeeper/gen/go/keeper/v1"
	"github.com/ArtShib/gophkeeper/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *serverAPI) UpdateSecret(ctx context.Context, req *keeperv1.UpdateSecretRequest) (*emptypb.Empty, error) {

	userID, err := getUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	reqSecret := req.GetSecret()
	if reqSecret == nil {
		return nil, status.Error(codes.InvalidArgument, "secret is required")
	}

	metaBytes, err := json.Marshal(reqSecret.GetMeta())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid metadata")
	}
	secretType := reqSecret.GetMeta().GetType().String()

	modelSecret := &models.Secret{
		ID:        reqSecret.GetId(),
		UserID:    userID,
		Type:      models.SecretType(secretType),
		Data:      reqSecret.GetEncryptedData(),
		Metadata:  metaBytes,
		UpdatedAt: reqSecret.GetUpdatedAt(),
		IsDeleted: reqSecret.GetIsDeleted(),
	}

	if err = s.service.UpdateSecret(ctx, modelSecret); err != nil {
		return nil, status.Error(codes.Internal, "failed to update secret")
	}

	return &emptypb.Empty{}, nil
}
