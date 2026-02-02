package auth

import (
	"context"
	"errors"

	keeperv1 "github.com/ArtShib/gophkeeper/gen/go/keeper/v1"
	"github.com/ArtShib/gophkeeper/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func (s *serverAPI) Login(ctx context.Context, req *keeperv1.LoginRequest) (*keeperv1.LoginResponse, error) {
	login := req.GetLogin()
	passHash := req.GetPasswordHash()

	if login == "" || passHash == "" {
		return nil, status.Error(codes.InvalidArgument, "login and password_hash required")
	}

	token, err := s.service.Login(ctx, login, passHash)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}

	response := keeperv1.LoginResponse_builder{
		AccessToken: proto.String(token),
	}.Build()

	return response, nil
}
