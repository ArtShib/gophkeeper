package auth

import (
	"context"

	keeperv1 "github.com/ArtShib/gophkeeper/gen/go/keeper/v1"
)

type AuthService interface {
	RegisterNewUser(ctx context.Context, login string, passHash string) (int64, error)
	Login(ctx context.Context, login string, passHash string) (string, error)
}

type serverAPI struct {
	keeperv1.UnimplementedAuthServiceServer
	service AuthService
}
