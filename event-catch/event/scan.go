package event

import (
	"event-catch/config"

	"github.com/ethereum/go-ethereum/ethclient"
)

type Scan struct {
	config *config.Config

	client *ethclient.Client
}

func NewScan(config *config.Config, client *ethclient.Client) (*Scan, error) {
	s := &Scan{
		config: config,
		client: client,
	}

	return s, nil
}
