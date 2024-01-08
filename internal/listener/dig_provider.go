package listener

import (
	"sync"

	"github.com/HarveyJhuang1010/blockhw/internal/models/bo"
	"go.uber.org/dig"
)

var (
	self *packet
)

func NewListener(in digIn) digOut {
	self = &packet{}
	self.Do(func() {
		self.in = in
		self.digOut = digOut{
			BlockListener: newBlockListener(in),
		}
	})

	return self.digOut
}

type digIn struct {
	dig.In

	MQ    bo.NatsClient
	Block bo.BlockUseCase
}

type packet struct {
	sync.Once

	in digIn

	digOut
}

type digOut struct {
	dig.Out

	BlockListener bo.Listener `name:"block_listener"`
}
