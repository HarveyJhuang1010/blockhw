package indexer

import (
	"context"
	"os/signal"
	"sync"
	"syscall"

	"github.com/HarveyJhuang1010/blockhw/internal/appcontext"
	"github.com/HarveyJhuang1010/blockhw/internal/logging"
	"github.com/HarveyJhuang1010/blockhw/internal/models/bo"
	"github.com/HarveyJhuang1010/blockhw/internal/utils/server"
	"go.uber.org/dig"
)

var (
	cronApp    *app
	srvSetOnce sync.Once
)

type app struct {
	dig.In

	ListenerService bo.Service `name:"listener_service"`
	CronService     bo.Service `name:"cron_service"`
}

func RunServer(app app) {
	srvSetOnce.Do(func() {
		cronApp = &app
	})

	ctx := context.Background()
	sCtx, stop := signal.NotifyContext(
		ctx,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	appContext := appcontext.New(sCtx)
	l := logging.NewZapLogger("indexer")
	appcontext.SetLogger(l)

	l.Info("start services")
	server.ListenAndServe(appContext, stop, app.ListenerService, app.CronService)
}
