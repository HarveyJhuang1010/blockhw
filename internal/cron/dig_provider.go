package cron

import (
	"sync"

	"github.com/HarveyJhuang1010/blockhw/internal/models/bo"
	"github.com/HarveyJhuang1010/blockhw/internal/wrapper/evmcli"
	"go.uber.org/dig"
)

var (
	self *packet
)

func NewCronTask(in digIn) digOut {
	self = &packet{}
	self.Do(func() {
		self.in = in
		self.digOut = digOut{
			BlockTask: newBlockSyncTask(in),
		}
	})

	return self.digOut
}

type digIn struct {
	dig.In

	EvmClient *evmcli.EVMClient
	NatClient bo.NatsClient
	BlockRepo bo.BlockRepo
}

type packet struct {
	sync.Once

	in digIn

	digOut
}

type digOut struct {
	dig.Out

	BlockTask bo.CronTask `name:"block_task"`
}
