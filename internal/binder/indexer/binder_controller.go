package indexer

import (
	"github.com/HarveyJhuang1010/blockhw/internal/cron"
	"github.com/HarveyJhuang1010/blockhw/internal/listener"
	"go.uber.org/dig"
)

func provideController(binder *dig.Container) {
	if err := binder.Provide(cron.NewCronTask); err != nil {
		panic(err)
	}
	if err := binder.Provide(listener.NewListener); err != nil {
		panic(err)
	}
}
