package sync

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ArtShib/gophkeeper/internal/lib/loghelper"
)

type AutoSyncService struct {
	logger      *slog.Logger
	syncService *SyncService
	cancel      context.CancelFunc
	wg          sync.WaitGroup
	mu          sync.Mutex
	ticker      *time.Ticker
	isRunning   atomic.Bool
}

func New(syncService *SyncService, logger *slog.Logger) *AutoSyncService {
	return &AutoSyncService{
		logger:      logger,
		syncService: syncService,
	}
}

func (s *AutoSyncService) Start(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx, cancel := context.WithCancel(ctx)
	s.cancel = cancel

	if s.ticker != nil {
		s.ticker.Stop()
	}
	s.ticker = time.NewTicker(30 * time.Second)

	s.isRunning.Store(true)

	s.wg.Add(1)
	go s.worker(ctx)
}

func (s *AutoSyncService) worker(ctx context.Context) {
	defer func() {
		s.isRunning.Store(false)
		s.wg.Done()

		if s.ticker != nil {
			s.ticker.Stop()
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case <-s.ticker.C:
			s.syncWithContext(ctx)
		}
	}
}

func (s *AutoSyncService) syncWithContext(ctx context.Context) {
	log := loghelper.New(s.logger, "autosync.sync")

	select {
	case <-ctx.Done():
		return
	default:
	}

	if err := s.syncService.Sync(ctx); err != nil {
		log.LogError(ctx, "syncService", err)
	}

}

func (s *AutoSyncService) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning.Load() {
		return
	}

	if s.cancel != nil {
		s.cancel()
	}

	if s.ticker != nil {
		s.ticker.Stop()
		s.ticker = nil
	}

	s.wg.Wait()
}
