# cmd/vesseld

Entrypoint for the Vessel self-hosted daemon.

## What it does

- Loads environment variables from `.env`
- Initializes the embedded SQLite database and vault
- Sets up the Docker engine client and Caddy reverse proxy
- Starts an Echo HTTP server on `PORT` (default `8080`)
- Serves the self-hosted dashboard static assets from `dashboard/dist`
- Wires OSS handlers for auth, projects, databases, deployments, backups, terminal, and settings

## Run locally

```bash
cp .env.example .env
go run ./cmd/vesseld
```

## Build

```bash
go build -o bin/vesseld ./cmd/vesseld
```
