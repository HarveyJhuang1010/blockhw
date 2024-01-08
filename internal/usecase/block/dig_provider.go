package block

import (
	"sync"

	"github.com/HarveyJhuang1010/blockhw/internal/models/bo"
	"github.com/HarveyJhuang1010/blockhw/internal/wrapper/evmcli"
	"go.uber.org/dig"
)

var (
	self *packet
)

func NewBlockUseCase(in digIn) digOut {
	self = &packet{}
	self.Do(func() {
		self.in = in
		self.digOut = digOut{
			BlockUseCase: newBlockUseCase(in),
		}
	})

	return self.digOut
}

type digIn struct {
	dig.In

	Repo      bo.BlockRepo
	EvmClient *evmcli.EVMClient
}

type packet struct {
	sync.Once

	in digIn

	digOut
}

type digOut struct {
	dig.Out

	BlockUseCase bo.BlockUseCase
}
