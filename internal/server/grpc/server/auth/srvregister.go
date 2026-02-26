package auth

import (
	keeperv1 "github.com/ArtShib/gophkeeper/gen/go/keeper/v1"
	"google.golang.org/grpc"
)

func Register(gRPC *grpc.Server, svc AuthService) {
	keeperv1.RegisterAuthServiceServer(gRPC, &serverAPI{
		service: svc,
	})
}
