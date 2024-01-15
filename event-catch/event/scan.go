package event

import (
	"context"
	"event-catch/config"
	"fmt"
	"math/big"
	"time"

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

	go s.lookingScan(config.Node.StartBlock, eventlog)

	return s, eventlog, nil
}

func (s *Scan) lookingScan(
	startBlock int64,
	// TODO scan해야하는 collection
	// TODO 캐치애햐하는 이벤트
	eventLog chan<- []ethType.Log,
) {
	startReadBlock, to := startBlock, uint64(0)

	s.FilterQuery = ethereum.FilterQuery{
		Addresses: []common.Address{},
		Topics:    [][]common.Hash{},
		FromBlock: big.NewInt(startReadBlock),
	}

	for {
		time.Sleep(1e8)

		ctx := context.Background()

		if maxBlock, err := s.client.BlockNumber(ctx); err != nil {
			fmt.Println("Get BlockNumber", "err", err.Error())
			continue
		} else {
			to = maxBlock

			if to > uint64(startReadBlock) {
				// s.FilterQuery.FromBlock(startReadBlock)
				s.FilterQuery.FromBlock = big.NewInt(startReadBlock)
				s.FilterQuery.ToBlock = big.NewInt(int64(to))

				tryCount := 1

			Retry:
				if logs, err := s.client.FilterLogs(ctx, s.FilterQuery); err != nil {

					if tryCount == 3 {
						fmt.Println("failed to get Filter", "err", err.Error())
						break
					} else {
						newTo := big.NewInt(int64(to) - 1)
						newFrom := big.NewInt(startBlock - 1)

						s.FilterQuery.ToBlock = newTo
						s.FilterQuery.FromBlock = newFrom

						tryCount++

						goto Retry
					}

				} else if len(logs) > 0 {
					eventLog <- logs

					startReadBlock = int64(to)
				}
			}
		}
	}
}
