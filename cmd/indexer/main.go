package main

import (
	"github.com/HarveyJhuang1010/blockhw/internal/app/indexer"
	binderIndexer "github.com/HarveyJhuang1010/blockhw/internal/binder/indexer"
	"github.com/HarveyJhuang1010/blockhw/internal/config"
	"github.com/HarveyJhuang1010/blockhw/internal/wrapper/database"
	"github.com/HarveyJhuang1010/blockhw/internal/wrapper/evmcli"
	"github.com/HarveyJhuang1010/blockhw/internal/wrapper/messaging"
)

func main() {
	defer func() {
		database.Finalize()
		evmcli.Finalize()
		messaging.Finalize()
	}()

	cfg := config.NewConfig()

	database.Initialize(cfg.Database)
	evmcli.Initialize(cfg.Ethereum)
	if err := messaging.Initialize(cfg.Nats); err != nil {
		panic(err)
	}

	if err := database.AutoMigrate(); err != nil {
		panic(err)
	}

	binder := binderIndexer.New()
	if err := binder.Invoke(indexer.RunServer); err != nil {
		panic(err)
	}
}
