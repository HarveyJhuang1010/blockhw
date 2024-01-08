package api

import (
	"github.com/HarveyJhuang1010/blockhw/internal/repository/block"
	"github.com/HarveyJhuang1010/blockhw/internal/repository/transaction"
	"go.uber.org/dig"
)

func provideRepository(binder *dig.Container) {
	if err := binder.Provide(block.NewBlockRepo); err != nil {
		panic(err)
	}
	if err := binder.Provide(transaction.NewTransactionRepo); err != nil {
		panic(err)
	}
}
