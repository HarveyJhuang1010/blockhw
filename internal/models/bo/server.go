package bo

import (
	"context"
	"sync"
)

type Service interface {
	Run(ctx context.Context, stop context.CancelFunc)
	Shutdown(ctx context.Context, wg *sync.WaitGroup)
}
