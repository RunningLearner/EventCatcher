package app

import (
	"event-catch/config"
	"event-catch/event"
	"event-catch/repository"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
)

type App struct {
	config *config.Config

	client *ethclient.Client

	repository *repository.Repository
	scan       *event.Scan
}

func NewApp(config *config.Config) {
	a := App{
		config: config,
	}

	var err error

	fmt.Println(a)

	if a.client, err = ethclient.Dial(config.Node.Uri); err != nil {
		panic(err)
	} else {
		if a.repository, err = repository.NewRepository(config); err != nil {
			panic(err)
		}

		if a.scan, err = event.NewScan(config, a.client); err != nil {
			panic(err)
		}
	}
}
