package event

import (
	"event-catch/config"

	ethType "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Catch struct {
	config *config.Config
	client *ethclient.Client
}

func NewCatch(config *config.Config, client *ethclient.Client, eventChan chan []ethType.Log) (*Catch, error) {
	c := &Catch{
		config: config,
		client: client,
	}

	//TODO 캐치해야하는 이벤트 정리
	go c.startToCatch(eventChan)

	return c, nil
}

func (c *Catch) startToCatch(events <-chan []ethType.Log) {
	for event := range events {

	}
}
