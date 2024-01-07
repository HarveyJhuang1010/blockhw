package server

import (
	"context"
	"sync"

	"github.com/HarveyJhuang1010/blockhw/internal/models/bo"
)

func ListenAndServe(ctx context.Context, stop context.CancelFunc, svc ...bo.Service) {
	if len(svc) == 0 {
		panic("nil service")
	}

	for _, s := range svc {
		go s.Run(ctx, stop)
	}

	<-ctx.Done()

	wg := &sync.WaitGroup{}
	wg.Add(len(svc))
	for _, s := range svc {
		go s.Shutdown(ctx, wg)
	}
	wg.Wait()
}
