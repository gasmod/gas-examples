package main

import (
	"log"

	"github.com/gasmod/gas"
)

func main() {
	app := gas.NewApp(
		// Register app service
		gas.WithService[*Service](NewService, gas.ServiceLifetimeSingleton),
	)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
