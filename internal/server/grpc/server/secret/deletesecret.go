package secret

import (
	"context"

	keeperv1 "github.com/ArtShib/gophkeeper/gen/go/keeper/v1"
	"github.com/ArtShib/gophkeeper/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *serverAPI) DeleteSecret(ctx context.Context, req *keeperv1.DeleteSecretRequest) (*emptypb.Empty, error) {
	userID, err := getUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	secretID := req.GetId()
	updatedAt := req.GetUpdatedAt()

	modelSecret := &models.Secret{
		ID:        secretID,
		UserID:    userID,
		UpdatedAt: updatedAt,
	}

	if err := s.service.DeleteSecret(ctx, modelSecret); err != nil {
		return nil, status.Error(codes.Internal, "failed to delete secret")
	}

	return &emptypb.Empty{}, nil
}
