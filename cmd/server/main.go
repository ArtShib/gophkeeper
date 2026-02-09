package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ArtShib/gophkeeper/internal/lib/exit"
	myLogger "github.com/ArtShib/gophkeeper/internal/lib/logger"
	"github.com/ArtShib/gophkeeper/internal/lib/loghelper"
	"github.com/ArtShib/gophkeeper/internal/server/app"
	"github.com/ArtShib/gophkeeper/internal/server/config"
	"github.com/ArtShib/gophkeeper/internal/storage/postgres"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {

	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)

	var err error
	logger := myLogger.New()
	logHelper := loghelper.New(logger, "main")

	ctx := context.Background()

	cfg, err := config.MustLoadConfig(ctx, logger)
	if err != nil {
		logHelper.LogError(ctx, "run MustLoadConfig", err)
		//exit.Code(exit.ExitConfigError).Exit()
	}

	initStorCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	store, err := postgres.New(initStorCtx, cfg.ConfigStore.DatabaseDSN)
	if err != nil {
		logHelper.LogError(ctx, "initStor", err)
		exit.Code(exit.ExitStoreError).Exit()
	}

	application := app.NewApp(ctx, cfg, store.DB, logger)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		if err = application.Run(ctx); err != nil {
			logHelper.LogError(ctx, "run application", err)
			exit.Code(exit.ExitGRPCSrvError).Exit()
		}
	}()
	<-quit
	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 10*time.Second)
	defer shutdownCancel()

	if err := application.Stop(shutdownCtx); err != nil {
		logHelper.LogError(shutdownCtx, "application shutdown error", err)
	}
}
