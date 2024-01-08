package transaction

import (
	"sync"

	"github.com/HarveyJhuang1010/blockhw/internal/models/bo"
	"go.uber.org/dig"
)

var (
	self *packet
)

func NewTransactionUseCase(in digIn) digOut {
	self = &packet{}
	self.Do(func() {
		self.in = in
		self.digOut = digOut{
			TransactionUseCase: newTransactionUseCase(in),
		}
	})

	return self.digOut
}

type digIn struct {
	dig.In

	Repo bo.TransactionRepo
}

type packet struct {
	sync.Once

	in digIn

	digOut
}

type digOut struct {
	dig.Out

	TransactionUseCase bo.TransactionUseCase
}
