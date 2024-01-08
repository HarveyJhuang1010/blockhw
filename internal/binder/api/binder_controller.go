package api

import (
	"github.com/HarveyJhuang1010/blockhw/internal/controller/block"
	"github.com/HarveyJhuang1010/blockhw/internal/controller/transaction"
	"go.uber.org/dig"
)

func provideController(binder *dig.Container) {
	if err := binder.Provide(block.NewBlockController); err != nil {
		panic(err)
	}
	if err := binder.Provide(transaction.NewTransactionController); err != nil {
		panic(err)
	}
}
