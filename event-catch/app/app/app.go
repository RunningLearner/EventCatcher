package app

import (
	"event-catch/config"
	"event-catch/event"
	"event-catch/repository"
	"fmt"
)

type App struct {
	config *config.Config

	repository *repository.Repository
	scan *event.Scan
}

func NewApp(config *config.Config) {
	a := App{
		config: config,
	}

	var err error

	if a.repository, err = repository.NewRepository(config); err != nil {
		panic(err)
	}

	fmt.Println(a)

	if a.scan, err = event.NewScan(config); err != nil{
		panic(err)
	}
}
