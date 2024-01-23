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

func NewScan(config *config.Config, client *ethclient.Client, catchEventList []common.Hash) (*Scan, chan []ethType.Log, error) {
	s := &Scan{
		config: config,
		client: client,
	}

	eventlog := make(chan []ethType.Log, 100)

	scanCollection := common.HexToAddress("0x16bc64a11ff2bb9f27c1e38c07f8111af2398dd1")

	go s.lookingScan(config.Node.StartBlock, scanCollection, catchEventList, eventlog)

	return s, eventlog, nil
}

// 0x16bc64a11ff2bb9f27c1e38c07f8111af2398dd1
// block num 45072015
// 민트 블록 45072044
func (s *Scan) lookingScan(
	startBlock int64,
	scanCollection common.Address,
	catchEventList []common.Hash,
	eventLog chan<- []ethType.Log,
) {
	startReadBlock, to := startBlock, uint64(0)

	s.FilterQuery = ethereum.FilterQuery{
		Addresses: []common.Address{scanCollection},
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
				fmt.Println("from Block", s.FilterQuery.FromBlock, "to Block", to)
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
