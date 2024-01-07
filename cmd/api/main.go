package main

import (
	"runtime/debug"

	"go.uber.org/zap"

	"github.com/HarveyJhuang1010/blockhw/internal/app/api"
	"github.com/HarveyJhuang1010/blockhw/internal/appcontext"
	binderApi "github.com/HarveyJhuang1010/blockhw/internal/binder/api"
	"github.com/HarveyJhuang1010/blockhw/internal/config"
	"github.com/HarveyJhuang1010/blockhw/internal/logging"
	"github.com/HarveyJhuang1010/blockhw/internal/wrapper/database"
	"github.com/HarveyJhuang1010/blockhw/internal/wrapper/evmcli"
)

func main() {
	defer func() {
		database.Finalize()
		evmcli.Finalize()
	}()

	cfg := config.NewConfig()
	l := logging.NewZapLogger("api")
	appcontext.SetLogger(l)

	if bi, ok := debug.ReadBuildInfo(); ok {
		l.Info("build info", zap.String("path", bi.Path))
		l.Info("build info", zap.Any("setting", bi.Settings))
	}

	database.Initialize(cfg.Database)
	evmcli.Initialize(cfg.Ethereum)

	binder := binderApi.New()
	if err := binder.Invoke(api.RunServer); err != nil {
		panic(err)
	}
}
