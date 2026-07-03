// File Vault API — a complete Gas example demonstrating DI, routing, auth,
// database, migrations, caching, file storage, job queues, email, and
// structured logging working together.
package main

import (
	"log"

	"github.com/gasmod/gas"
	gasconfig "github.com/gasmod/gas-config"
	"github.com/gasmod/gas-config/providers"
	database "github.com/gasmod/gas-database"
	gaslog "github.com/gasmod/gas-log"
	migrate "github.com/gasmod/gas-migrate"
	template "github.com/gasmod/gas-template/memory"

	"github.com/gasmod/gas-auth/apikey"
	"github.com/gasmod/gas-auth/jwt"
	cache "github.com/gasmod/gas-cache/memory"
	email "github.com/gasmod/gas-email/ses"
	queue "github.com/gasmod/gas-queue/sqs"
	storage "github.com/gasmod/gas-storage/s3"

	"github.com/gasmod/gas-examples/api-server/auth"
	"github.com/gasmod/gas-examples/api-server/files"
	"github.com/gasmod/gas-examples/api-server/shares"
)

func main() {
	// Load configuration from .env file + environment variables.
	// .env provides defaults for local dev; env vars override in production.
	cfg := gasconfig.New(
		gasconfig.WithProvider(providers.NewDotEnvProvider(
			providers.WithDotEnvFileNotFoundPanic(false),
		)),
		gasconfig.WithProvider(providers.NewEnvProvider()),
	)

	if err := cfg.Load(); err != nil {
		log.Fatalf("failed to load config: %s\n", err)
	}

	// In-memory template store for email templates. Services register
	// templates during Init() via gas.TemplateProvider.
	tmplStore := template.NewStore()

	app := gas.NewApp(
		// --- Infrastructure ---

		gas.WithServiceInstance[gas.ConfigProvider](cfg),

		gas.WithSingletonService[gas.Logger](gaslog.NewZeroLogLogger()),
		gas.WithScopedService[RequestLogger](requestLogger),

		gas.WithSingletonService[gas.DatabaseProvider](database.New()),
		gas.WithSingletonService[*migrate.Service](migrate.New()),
		gas.WithSingletonService[gas.CacheProvider](cache.New()),
		gas.WithSingletonService[gas.StorageProvider](storage.New()),
		gas.WithSingletonService[gas.JobQueueProvider](queue.New()),
		gas.WithServiceInstance[gas.TemplateProvider](tmplStore),
		gas.WithSingletonService[gas.EmailProvider](email.New()),

		// --- Auth ---

		// JWT and API key services are registered as singletons. They
		// manage their own config binding and (for apikey) migrations.
		gas.WithSingletonService[*jwt.Service](jwt.New()),
		gas.WithSingletonService[*apikey.Service](apikey.New()),

		// --- Application services ---

		gas.WithSingletonService[*auth.Service](auth.New),
		gas.WithSingletonService[*files.Service](files.New),
		gas.WithSingletonService[*shares.Service](shares.New),

		// --- HTTP ---

		gas.WithErrorHandler(errorHandler),
	)

	// Global middleware: security headers + request logging.
	app.Router().Use(
		gas.MiddlewareFunc(gas.SecurityHeaders()),
		gas.MiddlewareFunc(gas.RequestLogger[RequestLogger]()),
	)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
