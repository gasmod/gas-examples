package main

import (
	"log"

	"github.com/gasmod/gas"
	config "github.com/gasmod/gas-config"
	ui "github.com/gasmod/gas-ui"

	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	app := gas.NewApp(
		gas.WithService[gas.ConfigProvider](config.New(), gas.ServiceLifetimeSingleton),
		gas.WithService[gas.UIProvider](ui.New(), gas.ServiceLifetimeSingleton),
		gas.WithService[*Service](NewService, gas.ServiceLifetimeSingleton),
	)

	app.Router().Use(gas.MiddlewareFunc(middleware.Logger))
	app.Router().Use(gas.MiddlewareFunc(middleware.Recoverer))

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
