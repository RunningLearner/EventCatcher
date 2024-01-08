package main

import (
	"event-catch/app/app"
	"event-catch/config"
	"fmt"
)

func main() {
	config.NewConfig("./config.toml")
	a := app.NewApp()

	fmt.Println(a)
}
