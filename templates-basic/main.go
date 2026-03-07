package main

import (
	"log"

	"github.com/gasmod/gas"
	gasconfig "github.com/gasmod/gas-config"
	gasenv "github.com/gasmod/gas-config/extensions/gas-env"
	gaslog "github.com/gasmod/gas-log"
	ui "github.com/gasmod/gas-ui"
)

func main() {
	cfgProvider := gasconfig.New()

	// Environment defaults
	cfgProvider.SetDefault(gasenv.DefaultConfigKey, gasenv.Development)

	// UI defaults
	cfgProvider.SetDefault("UI.StaticPath", "/static/*")
	cfgProvider.SetDefault("UI.StaticStripPrefix", "/static/")

	if err := cfgProvider.Load(); err != nil {
		log.Fatalf("failed to load config: %s\n", err)
	}

	app := gas.NewApp(
		gas.WithServiceInstance[gas.ConfigProvider](cfgProvider),
		gas.WithSingletonService[gas.UIProvider](ui.New()),

		gas.WithSingletonService[gas.Logger](gaslog.NewSlogLogger()),
		gas.WithScopedService[RequestLogger](RequestLoggerCtor()),

		gas.WithAppModule[*Module](NewModule),
	)

	app.Router().Use(gas.MiddlewareFunc(gas.RequestLogger[RequestLogger]()))

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
