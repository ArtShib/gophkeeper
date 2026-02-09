package app

//
//import (
//	"github.com/ArtShib/gophkeeper/internal/lib/loghelper"
//	"github.com/ArtShib/gophkeeper/internal/server/app/grpc"
//	SvcKeeper "github.com/ArtShib/gophkeeper/internal/server/service/keeper"
//	"golang.org/x/sync/errgroup"
//)
//
//type App struct {
//	Logger     *slog.Logger
//	Store      storage.Storage
//	//Config     *config.Config
//	//ServerGRPC *grpc.App
//}
//
//// NewApp конструктор App
//func NewApp(ctx context.Context, cfg *config.Config, store storage.Storage, log *slog.Logger) *App {
//
//	logHelper := loghelper.New(log, "app.NewApp")
//	logHelper.LogDebug(ctx, "NewApp")
//	app := &App{
//		Config: cfg,
//		Store:  store,
//		Logger: log,
//	}
//
//	authSvc := auth.New(app.Logger, app.Store, cfg.ConfigJWT)
//	keeperSvc := SvcKeeper.New(app.Logger, app.Store)
//	app.ServerGRPC = grpc.New(log, cfg.ConfigGRPC.Port, authSvc, keeperSvc)
//	return app
//}
//
//// Run закпуск http сервера
//func (a *App) Run(ctx context.Context) error {
//
//	gr, ctx := errgroup.WithContext(ctx)
//
//	gr.Go(func() error {
//		return a.ServerGRPC.Start(ctx)
//	})
//
//	return gr.Wait()
//}
//
//// Stop остановка сервисов для реализации graceful shutdown
//func (a *App) Stop(ctx context.Context) error {
//	logHelper := loghelper.New(a.Logger, "app.Stop")
//
//	a.ServerGRPC.Stop(ctx)
//	if err := a.Store.Close(); err != nil {
//		return logHelper.LogAndReturnError(ctx, "failed to stop app gracefully", err)
//	}
//
//	return nil
//}
//
