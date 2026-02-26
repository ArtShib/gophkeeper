package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ArtShib/gophkeeper/internal/client/config"
	mygrpc "github.com/ArtShib/gophkeeper/internal/client/grpc"
	"github.com/ArtShib/gophkeeper/internal/client/service/auth"
	"github.com/ArtShib/gophkeeper/internal/client/ui"
	"github.com/ArtShib/gophkeeper/internal/lib/exit"
	mylogger "github.com/ArtShib/gophkeeper/internal/lib/logger"
	"github.com/ArtShib/gophkeeper/internal/lib/loghelper"
	"github.com/ArtShib/gophkeeper/internal/models"
	"github.com/ArtShib/gophkeeper/internal/storage/secret"
	"github.com/ArtShib/gophkeeper/internal/storage/sqlite"
	"github.com/ArtShib/gophkeeper/internal/storage/sync"
	"github.com/ArtShib/gophkeeper/internal/storage/user"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	var err error
	logger := mylogger.New()
	logHelper := loghelper.New(logger, "main")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conf := config.MustLoadConfig()
	store, err := sqlite.New(ctx, conf.ConfigStore.DatabaseDSN)

	if err != nil {
		logHelper.LogError(ctx, "initStor", err)
		exit.Code(exit.ExitStoreError).Exit()
	}

	userStorage := user.New(store.DB, models.DriverSQLite)
	secretStore := secret.New(store.DB, models.DriverSQLite)
	syncStore := sync.New(store.DB, models.DriverSQLite)

	grpcAuthClient, err := mygrpc.NewAuthClient(ctx, conf.ConfigGRPC.Server, logger, &models.ConfigTLS{})
	authService := auth.New(logger, userStorage, grpcAuthClient, conf.ConfigCrypto)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		cancel()
	}()

	p := tea.NewProgram(
		ui.InitialChoiceModel(ctx, authService, conf, secretStore, syncStore, logger),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	go func() {
		if _, err := p.Run(); err != nil {
			fmt.Printf("Ошибка запуска TUI: %v\n", err)
			cancel()
		}
	}()

	<-ctx.Done()

	time.Sleep(100 * time.Millisecond)
	p.Quit()
	store.DB.Close()
}
