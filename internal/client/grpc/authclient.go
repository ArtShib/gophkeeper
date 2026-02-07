package grpc

import (
	"context"
	"log/slog"

	"github.com/ArtShib/gophkeeper/internal/client/grpc/auth"
	"github.com/ArtShib/gophkeeper/internal/client/grpc/interceptors"
	"github.com/ArtShib/gophkeeper/internal/lib/loghelper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	api    *auth.Client
	logger *slog.Logger
	Conn   *grpc.ClientConn
	//Token  string
}

func NewAuthClient(ctx context.Context, addr string, logger *slog.Logger, certPath string) (*AuthClient, error) {
	log := loghelper.New(logger, "NewAuthClient")
	var opts []grpc.DialOption
	if certPath != "" {
		creds, err := credentials.NewClientTLSFromFile(certPath, "")
		if err != nil {
			return nil, log.LogAndReturnError(ctx, "Failed to load TLS credentials", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	opts = append(opts,
		grpc.WithUnaryInterceptor(interceptors.LoggingInterceptor(logger)))

	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		return nil, log.LogAndReturnError(ctx, "failed to connect to server", err)
	}

	api := auth.New(conn)

	return &AuthClient{
		Conn:   conn,
		api:    api,
		logger: logger}, nil
}

func (c *AuthClient) Close() error {
	return c.Conn.Close()
}

func (c *AuthClient) Login(ctx context.Context, login string, passHash []byte) (string, error) {
	return c.api.Login(ctx, login, passHash)
}

func (c *AuthClient) Register(ctx context.Context, login string, passHash []byte) (int64, error) {
	return c.api.Register(ctx, login, passHash)
}
