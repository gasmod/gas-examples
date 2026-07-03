# Gas Examples

Runnable sample applications demonstrating different Gas deployment shapes and
module combinations. Each example is its own Go module and can be built/run
independently.

| Example                                 | What it shows                                                                     |
|-----------------------------------------|-----------------------------------------------------------------------------------|
| [`hello-world/`](./hello-world)         | Smallest possible Gas HTTP app — single handler, no providers.                    |
| [`templates-basic/`](./templates-basic) | HTML templates with `gas-ui` layouts and partials, `gas-config`, request logging. |
| [`api-server/`](./api-server)           | Full JSON API: auth (JWT + API keys), file uploads (S3), shares, email, queue.    |
| [`lambda-worker/`](./lambda-worker)     | Headless `gas.Worker` running on AWS Lambda, consuming SQS via `gas-queue`.       |

## Running an example

```bash
cd <example-dir>
go run .
```

Examples that depend on AWS services (`api-server`, `lambda-worker`) target
[LocalStack](https://github.com/localstack/localstack) by default. Where a
`docker-compose.yaml` is present, run `docker compose up -d` first.

## Notes

- These examples track the **latest** released versions of the gasmod modules.
  If you see a build failure after updating, please open an issue.
- Configuration files (`config.json`, `.env.example`) are starter templates;
  copy and adjust for your environment.
