package main

import (
	"github.com/HarveyJhuang1010/blockhw/internal/app/api"
	binderApi "github.com/HarveyJhuang1010/blockhw/internal/binder/api"
	"github.com/HarveyJhuang1010/blockhw/internal/config"
	"github.com/HarveyJhuang1010/blockhw/internal/wrapper/database"
	"github.com/HarveyJhuang1010/blockhw/internal/wrapper/evmcli"
)

func main() {
	defer func() {
		database.Finalize()
		evmcli.Finalize()
	}()

	cfg := config.NewConfig()

	database.Initialize(cfg.Database)
	evmcli.Initialize(cfg.Ethereum)

	if err := database.AutoMigrate(); err != nil {
		panic(err)
	}

	binder := binderApi.New()
	if err := binder.Invoke(api.RunServer); err != nil {
		panic(err)
	}
}
