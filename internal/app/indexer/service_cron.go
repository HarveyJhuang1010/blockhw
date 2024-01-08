package indexer

import (
	"context"
	"fmt"
	"runtime"
	"sync"

	"github.com/HarveyJhuang1010/blockhw/internal/appcontext"
	"github.com/HarveyJhuang1010/blockhw/internal/models/bo"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type cronService struct {
	in digIn
}

func newCronService(in digIn) bo.Service {
	return &cronService{
		in: in,
	}
}

func (s *cronService) Run(ctx context.Context, stop context.CancelFunc) {
	logger := appcontext.GetLogger()

	defer func() {
		if r := recover(); r != nil {
			logger.Error("panic", zap.Any("panic", r))
		}
		stop()
	}()

	// initiate cron server and configure with secondOptional for compatibility from previous version
	server := cron.New(
		cron.WithParser(
			cron.NewParser(
				cron.SecondOptional|
					cron.Minute|
					cron.Hour|
					cron.Dom|
					cron.Month|
					cron.Dow|
					cron.Descriptor,
			),
		),
		cron.WithChain(
			recoverJob(logger),
			skipIfStillRunning(logger),
		),
	)

	//tasks
	tasks := []bo.CronTask{
		self.in.BlockTask,
	}

	for _, task := range tasks {
		server.AddJob(task.Schedule(), task)
	}

	server.Start()

	<-ctx.Done()
}

func (s *cronService) Shutdown(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	logger := appcontext.GetLogger()
	logger.Info("shutdown cron")
}

// recoverJob panics in wrapped jobs and log them with the provided logger.
func recoverJob(logger *zap.Logger) cron.JobWrapper {
	return func(j cron.Job) cron.Job {
		return cron.FuncJob(func() {
			defer func() {
				if r := recover(); r != nil {
					const size = 64 << 10
					buf := make([]byte, size)
					buf = buf[:runtime.Stack(buf, false)]
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					logger.Error("panic", zap.ByteString("stack", buf), zap.Error(err))
				}
			}()
			j.Run()
		})
	}
}

// skipIfStillRunning skips an invocation of the Job if a previous invocation is
// still running. It logs skips to the given logger at Info level.
func skipIfStillRunning(logger *zap.Logger) cron.JobWrapper {
	return func(j cron.Job) cron.Job {
		var ch = make(chan struct{}, 1)
		ch <- struct{}{}
		return cron.FuncJob(func() {
			select {
			case v := <-ch:
				j.Run()
				ch <- v
			default:
				logger.Info("skip")
			}
		})
	}
}
