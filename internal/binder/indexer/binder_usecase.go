package indexer

import (
	"github.com/HarveyJhuang1010/blockhw/internal/usecase/block"
	"github.com/HarveyJhuang1010/blockhw/internal/usecase/transaction"
	"go.uber.org/dig"
)

func provideUseCase(binder *dig.Container) {
	if err := binder.Provide(block.NewBlockUseCase); err != nil {
		panic(err)
	}
	if err := binder.Provide(transaction.NewTransactionUseCase); err != nil {
		panic(err)
	}
}
