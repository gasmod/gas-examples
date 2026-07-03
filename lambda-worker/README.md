# lambda-worker

A `gas.Worker` (no HTTP server, no router) running on AWS Lambda and consuming
SQS messages.

## What it demonstrates

- Building a headless Gas application with `gas.NewWorker` instead of
  `gas.NewApp` — same DI container, lifecycle, and migrations, no router.
- Initializing the worker in `init()` so the container survives across
  warm Lambda invocations.
- Using `lambda.WithEnableSIGTERM` to call `worker.Shutdown` cleanly when
  Lambda freezes the execution environment.
- Wiring `gas-database`, `gas-queue/sqs`, `gas-log`, and `gas-config` via DI.

## Building for Lambda

```bash
GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o bootstrap .
zip lambda.zip bootstrap
```

Upload `lambda.zip` to a Lambda function configured with:

- Runtime: `provided.al2` (custom runtime / `bootstrap` binary).
- Architecture: `arm64`.
- Trigger: an SQS queue.
- Env vars: `APP__DATABASE__DSN`, `APP__QUEUE__REGION`, and any others your
  `gas-config` env provider reads.

## Local testing

You can run this against [LocalStack](https://github.com/localstack/localstack)
with a SQS event payload piped to the handler. See the AWS Lambda Go runtime
docs for local-emulation options.

## Prerequisites

- AWS account or LocalStack.
- A Postgres-compatible database reachable from the function.
