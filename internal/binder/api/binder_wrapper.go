package api

import (
	"github.com/HarveyJhuang1010/blockhw/internal/wrapper/database"
	"github.com/HarveyJhuang1010/blockhw/internal/wrapper/evmcli"
	"go.uber.org/dig"
)

func provideWrapper(binder *dig.Container) {
	if err := binder.Provide(database.NewDatabaseClient); err != nil {
		panic(err)
	}
	if err := binder.Provide(evmcli.NewEVMClient); err != nil {
		panic(err)
	}
}
