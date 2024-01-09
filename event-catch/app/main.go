package main

import (
	"event-catch/app/app"
	"event-catch/config"
	"flag"
)

var configFlag = flag.String("config", "./config.toml", "toml env file not found")

func main() {
	flag.Parse()

	app.NewApp(config.NewConfig(*configFlag))
}
