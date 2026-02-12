package grpc

import (
	"context"
	"log/slog"

	"github.com/ArtShib/gophkeeper/internal/client/grpc/auth"
	"github.com/ArtShib/gophkeeper/internal/client/grpc/interceptors"
	"github.com/ArtShib/gophkeeper/internal/lib/loghelper"
	"github.com/ArtShib/gophkeeper/internal/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	api    *auth.Client
	logger *slog.Logger
	Conn   *grpc.ClientConn
	//Token  string
}

func NewAuthClient(ctx context.Context, addr string, logger *slog.Logger, tlsConfig *models.ConfigTLS) (*AuthClient, error) {
	log := loghelper.New(logger, "NewAuthClient")
	var opts []grpc.DialOption
	if tlsConfig.Cert != "" {
		//creds, err := credentials.NewClientTLSFromFile(certPath, "localhost")
		creds, err := LoadMTLSClient(tlsConfig)
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

//[87 65 9 232 251 15 212 127 110 172 155 192 88 141 187 72 48 15 188 38 216 109 102 164 170 212 51 93 101 91 251 187]
//[87 65 9 232 251 15 212 127 110 172 155 192 88 141 187 72 48 15 188 38 216 109 102 164 170 212 51 93 101 91 251 187]
