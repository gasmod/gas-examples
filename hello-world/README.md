# hello-world

The smallest possible Gas application: one handler, no providers, no config.

## What it demonstrates

- Building a `gas.App` with no service registrations.
- Registering an inline DI-aware handler directly on the router.
- The standard graceful-shutdown lifecycle (`app.Run`).

## Running

```bash
go run .
# in another shell:
curl http://localhost:8080/
```

## Prerequisites

None.
