package auth

import (
	"context"

	keeperv1 "github.com/ArtShib/gophkeeper/gen/go/keeper/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	api keeperv1.AuthServiceClient
}

func New(conn *grpc.ClientConn) *Client {
	return &Client{api: keeperv1.NewAuthServiceClient(conn)}
}

func (c *Client) Login(ctx context.Context, login string, passHash []byte) (string, error) {

	resp, err := c.api.Login(ctx, keeperv1.LoginRequest_builder{
		Login:        proto.String(login),
		PasswordHash: passHash,
	}.Build())

	if err != nil {
		return "", err
	}

	return resp.GetAccessToken(), nil
}

func (c *Client) Register(ctx context.Context, login string, passHash []byte) (int64, error) {

	resp, err := c.api.Register(ctx, keeperv1.RegisterRequest_builder{
		Login:        proto.String(login),
		PasswordHash: passHash,
	}.Build())

	if err != nil {
		return 0, err
	}

	return resp.GetUserId(), nil
}
