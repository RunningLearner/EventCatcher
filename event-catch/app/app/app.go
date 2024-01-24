package app

import (
	"event-catch/config"
	"event-catch/event"
	"event-catch/repository"
	"fmt"

	ethType "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type App struct {
	config *config.Config

	client *ethclient.Client

	repository *repository.Repository
	scan       *event.Scan
	catch      *event.Catch
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

		//getEventsToCatch

		if a.catch, err = event.NewCatch(config, a.client); err != nil {
			panic(err)
		}

		var eventChan chan []ethType.Log

		if a.scan, eventChan, err = event.NewScan(config, a.client, a.catch.GetEventsToCatch()); err != nil {
			panic(err)
		}

		go a.catch.StartToCatch(eventChan)

		for {
		}
	}
}
