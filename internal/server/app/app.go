package app

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/ArtShib/gophkeeper/internal/lib/loghelper"
	"github.com/ArtShib/gophkeeper/internal/models"
	"github.com/ArtShib/gophkeeper/internal/server/app/grpc"
	"github.com/ArtShib/gophkeeper/internal/server/config"
	"github.com/ArtShib/gophkeeper/internal/server/service/auth"
	SvcKeeper "github.com/ArtShib/gophkeeper/internal/server/service/keeper"
	"github.com/ArtShib/gophkeeper/internal/storage/secret"
	"github.com/ArtShib/gophkeeper/internal/storage/user"
	"golang.org/x/sync/errgroup"
)

type App struct {
	Logger     *slog.Logger
	store      *sql.DB
	Config     *config.Config
	serverGRPC *grpc.App
}

// NewApp конструктор App
func NewApp(ctx context.Context, cfg *config.Config, store *sql.DB, log *slog.Logger) (*App, error) {
	var err error

	logHelper := loghelper.New(log, "app.NewApp")
	logHelper.LogDebug(ctx, "NewApp")
	app := &App{
		Config: cfg,
		store:  store,
		Logger: log,
	}
	userStorage := user.New(store, models.DriverPostgres)
	authSvc := auth.New(app.Logger, userStorage, cfg.ConfigJWT)
	secretStorage := secret.New(store, models.DriverPostgres)
	keeperSvc := SvcKeeper.New(app.Logger, secretStorage)

	app.serverGRPC, err = grpc.New(log, cfg.ConfigGRPC.Port, authSvc, keeperSvc, cfg.ConfigTLS)
	if err != nil {
		return nil, err
	}

	return app, nil
}

// Run закпуск http сервера
func (a *App) Run(ctx context.Context) error {

	gr, ctx := errgroup.WithContext(ctx)

	gr.Go(func() error {
		return a.serverGRPC.Start(ctx)
	})

	return gr.Wait()
}

// Stop остановка сервисов для реализации graceful shutdown
func (a *App) Stop(ctx context.Context) error {
	logHelper := loghelper.New(a.Logger, "app.Stop")

	a.serverGRPC.Stop(ctx)
	if err := a.store.Close(); err != nil {
		return logHelper.LogAndReturnError(ctx, "failed to stop app gracefully", err)
	}

	return nil
}
