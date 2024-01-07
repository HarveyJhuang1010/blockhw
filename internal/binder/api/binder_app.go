package api

import (
	"github.com/HarveyJhuang1010/blockhw/internal/app/api"
	"go.uber.org/dig"
)

func provideApp(binder *dig.Container) {
	if err := binder.Provide(api.NewService); err != nil {
		panic(err)
	}
}
