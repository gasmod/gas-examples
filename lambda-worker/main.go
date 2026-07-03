// Lambda worker example showing Gas DI, config, database, and SQS queue
// without the Gas App or HTTP server. gas.NewWorker provides the same DI
// container, service lifecycle, and migrations as gas.NewApp but without
// a router or HTTP server — ideal for AWS Lambda, CLI tools, or background
// workers.
package main

import (
	"log"

	"github.com/gasmod/gas"
	gasconfig "github.com/gasmod/gas-config"
	"github.com/gasmod/gas-config/providers"
	database "github.com/gasmod/gas-database"
	gaslog "github.com/gasmod/gas-log"
	queue "github.com/gasmod/gas-queue/sqs"

	"github.com/aws/aws-lambda-go/lambda"
)

// worker holds the Gas worker at package level so it survives across
// Lambda invocations (Lambda reuses the execution environment).
var worker *gas.Worker

func init() {
	// Load config from environment variables. In Lambda, configuration
	// comes from env vars set in the function's deployment config.
	cfg := gasconfig.New(
		gasconfig.WithProvider(providers.NewEnvProvider(
			providers.WithEnvPrefix("APP"),
		)),
	)

	if err := cfg.Load(); err != nil {
		log.Fatalf("failed to load config: %s\n", err)
	}

	// NewWorker creates an EventBus and DI container internally, just
	// like NewApp but without a Router or HTTP server. Only WorkerOption
	// values are accepted — passing an AppOption panics.
	worker = gas.NewWorker(
		// Pre-built config instance. Other services (database, queue)
		// bind their settings from this provider automatically.
		gas.WithServiceInstance[gas.ConfigProvider](cfg),

		// Singleton logger — no scoped logger needed since there's no
		// per-request lifecycle in Lambda.
		gas.WithSingletonService[gas.Logger](gaslog.NewZeroLogLogger()),

		// Database connection registered as gas.DatabaseProvider. The
		// singleton is created once and reused across invocations —
		// Lambda best practice for connection pooling.
		gas.WithSingletonService[gas.DatabaseProvider](database.New()),

		// SQS queue client registered as gas.JobQueueProvider. Services
		// consume the provider interface, never the concrete backend.
		gas.WithSingletonService[gas.JobQueueProvider](queue.New()),

		// Handler is a singleton that receives all deps via DI. Resolved
		// once after Start, then passed to lambda.Start.
		gas.WithSingletonService[*Handler](NewHandler),
	)

	// Start runs: InitServices (BuildAll + Init) → migrations → ready
	// hooks. Non-blocking — does not start an HTTP server or wait for
	// signals. If the database is unreachable, config is invalid, or
	// NotificationQueueURL is missing, this fails fast.
	if err := worker.Start(); err != nil {
		log.Fatalf("failed to start worker: %s\n", err)
	}
}

func main() {
	// Resolve the handler once — all deps were injected by the container
	// during Start. This is the only manual resolve needed.
	h := gas.MustResolve[*Handler](worker.ServiceContainer())

	// WithEnableSIGTERM registers a callback that Lambda invokes before
	// freezing the execution environment. Worker.Shutdown emits
	// SystemShuttingDown and closes services in reverse init order.
	lambda.StartWithOptions(h.Handle, lambda.WithEnableSIGTERM(func() {
		if err := worker.Shutdown(); err != nil {
			log.Printf("shutdown error: %s\n", err)
		}
	}))
}
