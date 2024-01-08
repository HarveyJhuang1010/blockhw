package api

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
			ApiService: newApiService(in),
		}
	})

	return self.digOut
}

type digIn struct {
	dig.In

	BlockController       bo.BlockController
	TransactionController bo.TransactionController
}

type packet struct {
	sync.Once

	in digIn

	digOut
}

type digOut struct {
	dig.Out

	ApiService bo.Service `name:"api_service"`
}
