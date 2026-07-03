# api-server

A JSON file-vault API showcasing most of the Gas ecosystem in a single app.

## What it demonstrates

- **gas-auth**: JWT for user login, API keys for programmatic access, a `Chain`
  authenticator that accepts either.
- **gas-database**: PostgreSQL via `pgx`, with `sqlc`-generated queries.
- **gas-migrate**: embedded migrations registered per service.
- **gas-storage**: file uploads to S3 (LocalStack), with presigned download URLs.
- **gas-queue**: async share-notification jobs via SQS.
- **gas-email**: templated emails through SES.
- **gas-cache**: presigned URLs cached in-memory.
- **gas-template** + **gas-ui**: email templates rendered from an embedded FS.

## Running

LocalStack and Postgres are provided via Docker Compose:

```bash
docker compose up -d
cp .env.example .env       # adjust if needed
go run .
```

The server listens on `:8080` by default.

## Endpoints

| Method | Path                          | Auth     | Description                       |
|--------|-------------------------------|----------|-----------------------------------|
| POST   | `/api/auth/register`          | -        | Create a user.                    |
| POST   | `/api/auth/login`             | -        | Get a JWT.                        |
| POST   | `/api/auth/keys`              | JWT      | Mint an API key.                  |
| GET    | `/api/auth/keys`              | JWT      | List your API keys.               |
| DELETE | `/api/auth/keys/{id}`         | JWT      | Revoke an API key.                |
| POST   | `/api/files`                  | JWT/Key  | Upload a file.                    |
| GET    | `/api/files/{id}`             | JWT/Key  | Get a presigned download URL.     |
| POST   | `/api/files/{id}/share`       | JWT/Key  | Share a file by email.            |
| GET    | `/api/shares/{token}`         | -        | Redeem a share token.             |

## Prerequisites

- Docker (LocalStack + Postgres).
- Go 1.25+ (matches `go.mod`).
