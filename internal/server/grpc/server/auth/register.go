package auth

import (
	"context"

	keeperv1 "github.com/ArtShib/gophkeeper/gen/go/keeper/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func (s *serverAPI) Register(ctx context.Context, req *keeperv1.RegisterRequest) (*keeperv1.RegisterResponse, error) {
	login := req.GetLogin()
	passHash := req.GetPasswordHash()

	if login == "" || passHash == "" {
		return nil, status.Error(codes.InvalidArgument, "login and password_hash required")
	}

	userID, err := s.service.RegisterNewUser(ctx, login, passHash)

	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	response := keeperv1.RegisterResponse_builder{
		UserId: proto.Int64(userID),
	}.Build()

	return response, nil
}
