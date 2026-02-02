package secret

import (
	keeperv1 "github.com/ArtShib/gophkeeper/gen/go/keeper/v1"
	"google.golang.org/grpc"
)

func Register(gRPC *grpc.Server, svc KeeperService) {
	keeperv1.RegisterSecretServiceServer(gRPC, &serverAPI{
		service: svc,
	})
}
