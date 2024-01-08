package indexer

import (
	"github.com/HarveyJhuang1010/blockhw/internal/app/indexer"
	"go.uber.org/dig"
)

func provideApp(binder *dig.Container) {
	if err := binder.Provide(indexer.NewService); err != nil {
		panic(err)
	}
}
