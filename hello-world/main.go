package main

import (
	"log"

	"github.com/gasmod/gas"
	config "github.com/gasmod/gas-config"
)

func main() {
	cfg := gas.DefaultConfig()

	cfgMod := config.New()
	if err := cfgMod.Init(); err != nil {
		log.Fatal(err)
	}

	reg := gas.NewMiddlewareRegistry()
	router := gas.NewRouter(reg)
	bus := gas.NewEventBus()

	app := gas.NewApp(
		gas.WithConfig(cfg),
		gas.WithRouter(router),
		gas.WithMiddlewareRegistry(reg),
		gas.WithEventBus(bus),
		gas.WithModule(cfgMod),
		gas.WithModule(NewModule(WithRouter(router))),
	)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
