package block

import (
	"sync"

	"github.com/HarveyJhuang1010/blockhw/internal/models/bo"
	"github.com/HarveyJhuang1010/blockhw/internal/wrapper/database"
	"go.uber.org/dig"
)

var (
	self *packet
)

func NewBlockRepo(in digIn) digOut {
	self = &packet{}
	self.Do(func() {
		self.in = in
		self.digOut = digOut{
			BlockRepo: newBlockRepo(in),
		}
	})

	return self.digOut
}

type digIn struct {
	dig.In

	RDB *database.DB
}

type packet struct {
	sync.Once

	in digIn

	digOut
}

type digOut struct {
	dig.Out

	BlockRepo bo.BlockRepo
}
