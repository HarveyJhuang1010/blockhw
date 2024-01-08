package indexer

import (
	"context"
	"sync"
	"time"

	"github.com/HarveyJhuang1010/blockhw/internal/appcontext"
	"github.com/HarveyJhuang1010/blockhw/internal/models/bo"
	"go.uber.org/zap"
)

type listenerService struct {
	in digIn
}

func newListenerService(in digIn) bo.Service {
	return &listenerService{
		in: in,
	}
}

func (s *listenerService) Run(ctx context.Context, stop context.CancelFunc) {
	logger := appcontext.GetLogger()

	defer func() {
		if r := recover(); r != nil {
			logger.Error("panic", zap.Any("panic", r))
		}
		stop()
	}()

	logger.Info("start listen")
	s.in.BlockListener.Listen(ctx)
}

func (s *listenerService) Shutdown(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	logger := appcontext.GetLogger()

	logger.Info("wait for shutting down listener...")
	time.Sleep(1 * time.Minute)
	logger.Info("shutdown listener")
}
