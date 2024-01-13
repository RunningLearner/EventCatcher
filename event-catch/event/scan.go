package event

import (
	"event-catch/config"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	ethType "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Scan struct {
	config *config.Config

	FilterQuery ethereum.FilterQuery
	client      *ethclient.Client
}

func NewScan(config *config.Config, client *ethclient.Client) (*Scan, chan []ethType.Log, error) {
	s := &Scan{
		config: config,
		client: client,
	}

	eventlog := make(chan []ethType.Log, 100)

	s.lookingScan(config.Node.StartBlock, eventlog)

	return s, eventlog, nil
}

func (s *Scan) lookingScan(
	startBlock int64,
	// TODO scan해야하는 collection
	// TODO 캐치애햐하는 이벤트
	eventLog chan<- []ethType.Log,
) {
	startReadBlock, to := startBlock, int64(0)

	s.FilterQuery = ethereum.FilterQuery{
		Addresses: []common.Address{},
		Topics:    [][]common.Hash{},
	}

}
