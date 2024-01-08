package block

import (
	"sync"

	"github.com/HarveyJhuang1010/blockhw/internal/models/bo"
	"go.uber.org/dig"
)

var (
	self *packet
)

func NewBlockController(in digIn) digOut {
	self = &packet{}
	self.Do(func() {
		self.in = in
		self.digOut = digOut{
			BlockController: newBlockController(in),
		}
	})

	return self.digOut
}

type digIn struct {
	dig.In

	UseCase bo.BlockUseCase
}

type packet struct {
	sync.Once

	in digIn

	digOut
}

type digOut struct {
	dig.Out

	BlockController bo.BlockController
}
