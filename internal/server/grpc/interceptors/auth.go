package interceptors

import (
	"context"

	"github.com/ArtShib/gophkeeper/internal/server/models"
	"github.com/ArtShib/gophkeeper/internal/server/service/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthInterceptor(auth *auth.Auth) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {

		switch info.FullMethod {
		case "/secret.v1.AuthService/Register",
			"/secret.v1.AuthService/Login":
			return handler(ctx, req)
		}

		mb, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Internal, "failed to fetch metadata")
		}

		values := mb.Get("access_token")
		if len(values) == 0 {
			return nil, status.Error(codes.PermissionDenied, "permission denied")
		}

		userID, err := auth.ParseToken(ctx, values[0])
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token")
		}

		ctx = context.WithValue(ctx, models.UserIDKey, userID)
		return handler(ctx, req)
	}
}
