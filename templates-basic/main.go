package main

import (
	"log"

	"github.com/gasmod/gas"
	config "github.com/gasmod/gas-config"
	ui "github.com/gasmod/gas-ui"

	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := gas.DefaultConfig()

	cfgMod := config.New()
	if err := cfgMod.Init(); err != nil {
		log.Fatal(err)
	}

	if err := cfgMod.Bind(cfg); err != nil {
		log.Fatal(err)
	}

	router := gas.NewRouter()
	bus := gas.NewEventBus()

	uiMod := ui.New(
		ui.WithConfig(&ui.Config{
			TemplateDir: "templates",
			StaticDir:   "static",
			StaticPath:  "/static/",
			LayoutName:  "base",
		}),
		ui.WithRouter(router),
	)

	_ = router.Use(gas.MiddlewareFunc(middleware.Logger))
	_ = router.Use(gas.MiddlewareFunc(middleware.Recoverer))
	_ = router.Use(gas.MiddlewareFunc(middleware.Compress(5)))

	app := gas.NewApp(
		gas.WithConfig(cfg),
		gas.WithRouter(router),
		gas.WithEventBus(bus),
		gas.WithModule(cfgMod),
		gas.WithModule(uiMod),
		gas.WithModule(NewModule(
			WithRouter(router),
			WithUIProvider(uiMod),
		)),
	)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
