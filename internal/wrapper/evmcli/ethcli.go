package evmcli

import (
	"sync"

	"github.com/HarveyJhuang1010/blockhw/internal/config"
	"github.com/ethereum/go-ethereum/ethclient"
)

var defaultEVMClient EVMClient

type EVMClient struct {
	*ethclient.Client

	config *config.Ethereum

	connLock *sync.Mutex
}

func Initialize(cfg *config.Ethereum) {
	if defaultEVMClient.Client == nil {
		defaultEVMClient.initialize(cfg)
	}
}

func Finalize() {
	defaultEVMClient.Close()
}

func (c *EVMClient) initialize(cfg *config.Ethereum) {
	c.config = cfg

	if c.connLock == nil {
		c.connLock = &sync.Mutex{}
	}
	c.connLock.Lock()
	defer c.connLock.Unlock()

	if c.Client == nil {
		c.Client = c.dial(cfg.Endpoint)
	}
}

func (c *EVMClient) dial(endPoint string) *ethclient.Client {
	client, err := ethclient.Dial(endPoint)
	if err != nil {
		panic(err)
	}
	return client
}

func GetClient() *EVMClient {
	return defaultEVMClient.getClient()
}

// GetDB returns the database handle.
func (c *EVMClient) getClient() *EVMClient {
	c.connLock.Lock()
	conn := c.Client
	c.connLock.Unlock()

	if conn == nil {
		panic("uninitialized client")
	}

	return &defaultEVMClient
}
