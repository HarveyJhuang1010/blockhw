package transaction

import (
	"sync"

	"github.com/HarveyJhuang1010/blockhw/internal/models/bo"
	"go.uber.org/dig"
)

var (
	self *packet
)

func NewTransactionController(in digIn) digOut {
	self = &packet{}
	self.Do(func() {
		self.in = in
		self.digOut = digOut{
			TransactionController: newTransactionController(in),
		}
	})

	return self.digOut
}

type digIn struct {
	dig.In

	UseCase bo.TransactionUseCase
}

type packet struct {
	sync.Once

	in digIn

	digOut
}

type digOut struct {
	dig.Out

	TransactionController bo.TransactionController
}
