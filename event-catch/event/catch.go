package event

import (
	"context"
	"crypto"
	"event-catch/config"
	"event-catch/types"

	"github.com/ethereum/go-ethereum/common"
	ethType "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Catch struct {
	config *config.Config

	client *ethclient.Client

	needToCatchEvent map[common.Hash]types.NeedToCatchEvent
}

func NewCatch(config *config.Config, client *ethclient.Client, eventChan chan []ethType.Log) (*Catch, error) {
	c := &Catch{
		config: config,
		client: client,
	}

	// Transfer(address, address, uint256)
	c.needToCatchEvent = map[common.Hash]types.NeedToCatchEvent{
		common.BytesToAddress(crypto.Keccak256([]byte("Transfer(address, address, uint256)"))): {
			NeedToCatchEventFunc: c.Transfer,
		},
	}

	go c.startToCatch(eventChan)

	return c, nil
}

func (c *Catch) Transfer(e *types.Log, tx *types.Transaction) {}

// 이벤트 캐치 시작
func (c *Catch) startToCatch(events <-chan []ethType.Log) {
	for event := range events {
		ctx := context.Background()
		txList := make(map[common.Hash]*ethType.Transaction)

		for _, e := range event {
			hash := e.TxHash

			if _, ok := txList[hash]; !ok {
				if tx, pending, err := c.client.TransactionByHash(ctx, hash); err == nil {
					if !pending {
						txList[hash] = tx
					}
				}
			}

			if e.Removed {
				continue
			} else if et, ok := c.needToCatchEvent[e.Topics[0]]; !ok {
				//TODO 로그
			} else {
				et.NeedToCatchEventFunc(&e, txList[hash])
			}
		}
	}
}
