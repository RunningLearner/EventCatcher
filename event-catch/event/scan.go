package event

import (
	"event-catch/config"

	"github.com/ethereum/go-ethereum/ethclient"
)

type Scan struct {
	config *config.Config

	client *ethclient.Client
}

func NewScan(config *config.Config) (*Scan, error) {
	s := &Scan{
		config: config,
	}

	var err error

	if s.client, err = ethclient.Dial(config.Node.Uri); err != nil {
		return nil, err
	} else {
		return s, nil
	}
}
