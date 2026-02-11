package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/ArtShib/gophkeeper/internal/lib/loghelper"
	"github.com/ArtShib/gophkeeper/internal/models"
	mygrpc "github.com/ArtShib/gophkeeper/internal/server/grpc"
	"github.com/ArtShib/gophkeeper/internal/server/grpc/interceptors"
	authgrpc "github.com/ArtShib/gophkeeper/internal/server/grpc/server/auth"
	"github.com/ArtShib/gophkeeper/internal/server/grpc/server/secret"
	"github.com/ArtShib/gophkeeper/internal/server/service/auth"
	"github.com/ArtShib/gophkeeper/internal/server/service/keeper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	logger     *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(logger *slog.Logger, port int, authSvc *auth.Auth, keeperSvc *keeper.Keeper, config *models.ConfigTLS) (*App, error) {
	var cred credentials.TransportCredentials
	var err error
	if config.Cert != "" && config.Key != "" {
		cred, err = mygrpc.LoadMTLSServer(config)
		return nil, err
	} else {
		cred = insecure.NewCredentials()
	}

	gRPCServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.LoggerInterceptor(logger),
			interceptors.TlsInterceptor(),
			interceptors.AuthInterceptor(authSvc),
		),
		grpc.Creds(cred),
	)
	authgrpc.Register(gRPCServer, authSvc)
	secret.Register(gRPCServer, keeperSvc)

	return &App{
		logger:     logger,
		port:       port,
		gRPCServer: gRPCServer,
	}, nil
}

func (a *App) Start(ctx context.Context) error {
	logHelper := loghelper.New(a.logger, "ServerGRPC.Start")
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return logHelper.LogAndReturnError(ctx, "Error listen port", err)
	}
	logHelper.LogDebug(ctx, "Server GRPC start")
	if err := a.gRPCServer.Serve(l); err != nil {
		return logHelper.LogAndReturnError(ctx, "Error starting server", err)
	}
	return nil
}

func (a *App) Stop(ctx context.Context) {
	logHelper := loghelper.New(a.logger, "ServerGRPC.Stop")
	a.gRPCServer.GracefulStop()
	logHelper.LogDebug(ctx, "Server GRPC stop")
}
