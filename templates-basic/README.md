# templates-basic

A minimal HTML web application using `gas-ui` for templates and `gas-config`
for configuration.

## What it demonstrates

- Loading configuration from `config.json` with `gas-config`.
- Rendering HTML with layouts and partials via `gas-ui`.
- Serving static files from `./static`.
- Filesystem-backed `TemplateProvider` (`gas-template/fs`).
- A scoped request logger separate from the singleton `gas.Logger`.
- Per-request logging middleware (`gas.RequestLogger`).

## Running

```bash
go run .
# in another shell:
curl http://localhost:8080/
```

## Files

- `config.json` — server, logging, and UI config (read by `gas-config`).
- `templates/` — page templates, layouts, and partials.
- `static/` — static assets served at `/static/`.

## Prerequisites

None.
