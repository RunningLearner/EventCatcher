package event

import (
	"context"
	"event-catch/config"
	"event-catch/repository"
	"event-catch/types"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	ethType "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Catch struct {
	config *config.Config

	client     *ethclient.Client
	repository *repository.Repository

	needToCatchEvent map[common.Hash]types.NeedToCatchEvent
}

func NewCatch(config *config.Config, client *ethclient.Client, repository *repository.Repository) (*Catch, error) {
	c := &Catch{
		config:     config,
		client:     client,
		repository: repository,
	}

	c.needToCatchEvent = map[common.Hash]types.NeedToCatchEvent{
		common.BytesToHash(crypto.Keccak256([]byte("Transfer(address, address, uint256)"))): {
			NeedToCatchEventFunc: c.Transfer,
		},
	}

	return c, nil
}

func (c *Catch) Transfer(e *ethType.Log, tx *ethType.Transaction) {
	fmt.Println("들어왔습니다.")

	// e.Topics[1][:] 인덱싱이 되어있는경우 이벤트 가져오는 방법
	// e.Topics[2][:]

	//인덱스가 안되있는 경우
	// e.Data[:0x20]
	// e.Data[0x20:0x40]

	from := common.BytesToAddress(e.Topics[1][:])
	to := common.BytesToAddress(e.Topics[2][:])
	tokenID := new(big.Int).SetBytes(e.Topics[3][:])

	chainId, _ := c.client.ChainID(context.Background())
	sender, _ := ethType.Sender(ethType.NewLondonSigner(big.NewInt(80001)), tx)

	// 1. tx에 대한 이벤트 넣어주기
	c.repository.UpsertTxEvent(from, to, sender, tokenID, e.TxHash.Hex())

	fmt.Println(from, to, tokenID)

}

// 이벤트 캐치 시작
func (c *Catch) StartToCatch(events <-chan []ethType.Log) {
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
