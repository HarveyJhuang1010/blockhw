package messaging

import (
	"sync"

	"github.com/HarveyJhuang1010/blockhw/internal/models/bo"
	"go.uber.org/dig"
)

var (
	self *packet
)

func NewMessageQueue(in digIn) digOut {
	self = &packet{}
	self.Do(func() {
		self.in = in
		self.digOut = digOut{
			NatClient: GetClient(),
		}
	})

	return self.digOut
}

type digIn struct {
	dig.In
}

type packet struct {
	sync.Once

	in digIn

	digOut
}

type digOut struct {
	dig.Out

	NatClient bo.NatsClient
}
