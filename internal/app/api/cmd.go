package api

import (
	"context"
	"os/signal"
	"sync"
	"syscall"

	"github.com/HarveyJhuang1010/blockhw/internal/appcontext"
	"github.com/HarveyJhuang1010/blockhw/internal/models/bo"
	"github.com/HarveyJhuang1010/blockhw/internal/utils/server"
	"go.uber.org/dig"
)

var (
	apiApp     *app
	srvSetOnce sync.Once
)

type app struct {
	dig.In

	ApiService bo.Service `name:"api_service"`
}

func RunServer(app app) {
	srvSetOnce.Do(func() {
		apiApp = &app
	})

	ctx := context.Background()
	logger := appcontext.GetLogger(ctx)
	sCtx, stop := signal.NotifyContext(
		ctx,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	appContext := appcontext.New(sCtx)

	logger.Info("start services")
	server.ListenAndServe(appContext, stop, app.ApiService)
}
