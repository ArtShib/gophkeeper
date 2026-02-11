package grpc

import (
	"context"
	"log/slog"

	"github.com/ArtShib/gophkeeper/internal/client/grpc/interceptors"
	"github.com/ArtShib/gophkeeper/internal/client/grpc/secret"
	"github.com/ArtShib/gophkeeper/internal/lib/loghelper"
	"github.com/ArtShib/gophkeeper/internal/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type SecretClient struct {
	api    *secret.Client
	logger *slog.Logger
	Conn   *grpc.ClientConn
	Token  string
}

func NewSecretClient(ctx context.Context, addr string, logger *slog.Logger, tlsConfig *models.ConfigTLS, token string) (*SecretClient, error) {
	log := loghelper.New(logger, "NewSecretClient")
	var opts []grpc.DialOption
	if tlsConfig.Cert != "" {
		//creds, err := credentials.NewClientTLSFromFile(certPath, "")
		creds, err := LoadMTLSClient(tlsConfig)
		if err != nil {
			return nil, log.LogAndReturnError(ctx, "Failed to load TLS credentials", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	opts = append(opts,
		grpc.WithUnaryInterceptor(interceptors.LoggingInterceptor(logger)),
		grpc.WithUnaryInterceptor(interceptors.AuthInterceptor(token)))

	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		return nil, log.LogAndReturnError(ctx, "failed to connect to server", err)
	}

	api := secret.New(conn)

	return &SecretClient{
		Conn:   conn,
		api:    api,
		logger: logger}, nil
}

func (c *SecretClient) Close() error {
	return c.Conn.Close()
}

func (c *SecretClient) CreateSecret(ctx context.Context, secret *models.Secret) error {
	return c.api.CreateSecret(ctx, secret)

}

func (c *SecretClient) ListSecrets(ctx context.Context, empty *emptypb.Empty) (models.ArraySecret, error) {
	return c.api.ListSecrets(ctx, empty)
}

func (c *SecretClient) UpdateSecret(ctx context.Context, secret *models.Secret) error {
	return c.api.UpdateSecret(ctx, secret)
}

func (c *SecretClient) DeleteSecret(ctx context.Context, id string, updatedAt int64) error {
	return c.api.DeleteSecret(ctx, id, updatedAt)
}
