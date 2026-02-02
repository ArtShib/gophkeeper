package secret

import (
	"context"

	"github.com/ArtShib/gophkeeper/internal/server/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func getUserID(ctx context.Context) (int64, error) {
	val := ctx.Value(models.UserIDKey)

	if val == nil {
		return 0, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	userID, ok := val.(int64)
	if !ok {
		return 0, status.Error(codes.Internal, "invalid user id type in context")
	}

	return userID, nil
}
