package indexer

import (
	"sync"

	"github.com/HarveyJhuang1010/blockhw/internal/models/bo"
	"go.uber.org/dig"
)

var (
	self *packet
)

func NewService(in digIn) digOut {
	self = &packet{}
	self.Do(func() {
		self.in = in
		self.digOut = digOut{
			ListenerService: newListenerService(in),
			CronService:     newCronService(in),
		}
	})

	return self.digOut
}

type digIn struct {
	dig.In

	// tasks
	BlockTask     bo.CronTask `name:"block_task"`
	BlockListener bo.Listener `name:"block_listener"`
}

type packet struct {
	sync.Once

	in digIn

	digOut
}

type digOut struct {
	dig.Out

	ListenerService bo.Service `name:"listener_service"`
	CronService     bo.Service `name:"cron_service"`
}
