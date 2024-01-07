package database

import (
	"sync"

	"go.uber.org/dig"
)

var (
	self *set
)

func NewDatabaseClient(in digIn) digOut {
	self = &set{}
	self.Do(func() {
		self.in = in
		self.digOut = digOut{
			RDB: GetDB(),
		}
	})

	return self.digOut
}

type digIn struct {
	dig.In
}

type set struct {
	sync.Once
	in digIn

	digOut
}

type digOut struct {
	dig.Out

	RDB *DB
}
