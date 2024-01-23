package event

import (
	"context"
	"event-catch/config"
	"event-catch/types"
	"fmt"

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

	c.needToCatchEvent = map[common.Hash]types.NeedToCatchEvent{
		common.BytesToHash(crypto.Keccak256([]byte("Transfer(address, address, uint256)"))): {
			NeedToCatchEventFunc: c.Transfer,
		},
	}

	go c.startToCatch(eventChan)

	return c, nil
}

func (c *Catch) Transfer(e *ethType.Log, tx *ethType.Transaction) {
	fmt.Println("들어왔습니다.")
}

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

func (c *Catch) GetEventsToCatch() []common.Hash {
	eventsToCatch := make([]common.Hash, 0)

	for e := range c.needToCatchEvent {
		eventsToCatch = append(eventsToCatch, e)
	}

	return eventsToCatch
}
