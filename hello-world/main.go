package main

import (
	"log"
	"net/http"

	"github.com/gasmod/gas"
)

func main() {
	app := gas.NewApp()

	app.Router().Handle("", http.MethodGet, "/", func(ctx gas.Context) error {
		return ctx.Text(http.StatusOK, "Hello, World!")
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
