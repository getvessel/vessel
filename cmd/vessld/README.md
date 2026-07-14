# `vessld` — Server Daemon

`vessld` is the Vessl server process. It runs on your VPS, manages the SQLite database, orchestrates Docker containers, and exposes the HTTP API consumed by the dashboard and the `vessl` remote CLI.

## Running the Server

```sh
vessld serve          # Start the daemon (default when no subcommand is given)
```

By default it listens on `:8080`. Configure with environment variables:

| Variable          | Default | Description                               |
| ----------------- | ------- | ----------------------------------------- |
| `PORT`            | `8080`  | HTTP port to listen on                    |
| `HOST`            | ``      | Bind address                              |
| `VESSL_DATA_DIR`  | `data/` | Directory for SQLite DB and secrets vault |
| `VESSL_TLS_EMAIL` | ``      | Email for Let's Encrypt (Traefik)         |

## Setup & Maintenance

```sh
vessld setup            # Interactive first-time setup wizard
vessld reset-password   # Reset the admin account password
vessld config           # View or update server configuration
vessld restart          # Gracefully restart the daemon
vessld version          # Print the current version
```

## Deployment Commands

Deploy applications directly from the server terminal without using the dashboard.

```sh
# From a Git repository
vessld deploy https://github.com/your/repo.git

# From a template (e.g. go-fiber, nextjs)
vessld deploy --template go-fiber

# From a Docker image
vessld deploy --image nginx:latest --port 80

# From a Docker Compose file
vessld deploy --compose ./docker-compose.yml
```

## Resource Management

All commands below operate directly on the database — no HTTP, no auth token required. Useful for admin recovery or scripting.

### Projects

```sh
vessld project:list
vessld project:show <id>
vessld project:create <name>
vessld project:destroy <id>
```

### Applications

```sh
vessld apps:list
vessld apps:show <id>
vessld apps:create <name> --project <id>
vessld apps:destroy <id>
```

### Databases

```sh
vessld db:list
vessld db:show <id>
vessld db:create <name> <engine> --project <id>
vessld db:destroy <id>
```

Supported engines: `postgres`, `mysql`, `mariadb`, `redis`, `mongodb`, `clickhouse`, `kafka`, `rabbitmq`, `nats`.

### Environment Variables

```sh
vessld env:list --project <id>
vessld env:set KEY=VALUE --project <id>
vessld env:unset KEY --project <id>
```

### Deployments & Logs

```sh
vessld deployment:list --service <id>
vessld deployment:show <id>
vessld deployment:logs <id>
```

### Custom Domains

```sh
vessld domain:list --project <id>
vessld domain:add <hostname> --project <id>
vessld domain:remove <id>
```

## Advanced

```sh
vessld mcp              # Run the MCP stdio server for AI integrations
```
