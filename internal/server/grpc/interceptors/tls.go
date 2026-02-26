package interceptors

import (
	"context"

	mygrpc "github.com/ArtShib/gophkeeper/internal/server/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TlsInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {

		switch info.FullMethod {
		case "/keeper.v1.AuthService/Register",
			"/keeper.v1.AuthService/Login":
			return handler(ctx, req)
		}

		clientCN := mygrpc.GetClientCN(ctx)
		if clientCN != "" {
			ctx = context.WithValue(ctx, "client_id", clientCN)
			return handler(ctx, req)
		}
		return nil, status.Errorf(codes.Unauthenticated, "failed mTLS auth")
	}
}
