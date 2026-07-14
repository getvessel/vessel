---
title: CLI Reference
description: Command-line interface reference for Vessl — covering both the server daemon (vessld) and the remote client (vessl).
---

Vessl ships two CLI tools with distinct responsibilities.

| Tool | Runs on | Connects to |
|---|---|---|
| `vessld` | Your VPS / server | SQLite + Docker directly |
| `vessl` | Your local machine | `vessld` over HTTP |

---

## `vessld` — Server Daemon CLI

The `vessld` binary is the Vessl server process. It runs on your VPS and exposes the HTTP API that the dashboard and the `vessl` remote CLI consume. It also doubles as a management CLI for direct server-side operations without needing the dashboard.

### Server Commands

#### `serve`

Starts the Vessl daemon. This is the default command when no subcommand is provided.

```sh
vessld serve
```

#### `setup`

Runs the interactive setup wizard to initialise the database and create the initial admin account.

```sh
vessld setup
```

#### `reset-password`

Resets the password for the admin account. Useful if you lose access to the dashboard.

```sh
vessld reset-password
```

#### `config`

View or update global server configuration variables.

```sh
vessld config
```

#### `restart`

Gracefully restarts the Vessl daemon via Docker Compose.

```sh
vessld restart
```

#### `mcp`

Runs the Model Context Protocol (MCP) server over standard I/O for AI assistant integrations.

```sh
vessld mcp
```

#### `version`

Prints the current daemon version.

```sh
vessld version
```

### Deployment Commands

#### `deploy`

Deploy an application directly from the server terminal.

```sh
# From a Git repository
vessld deploy https://github.com/your/repo.git

# From a Docker image
vessld deploy --image nginx:latest --port 80

# From a Docker Compose file
vessld deploy --compose ./docker-compose.yml
```

### Resource Management Commands

All resource commands use a `<resource>:<action>` syntax and interact with the database directly — no HTTP, no auth required.

#### Projects

```sh
vessld project:list                        # List all projects
vessld project:show <id>                   # Show project details
vessld project:create <name>               # Create a project
vessld project:destroy <id>                # Delete a project
```

#### Applications

```sh
vessld apps:list                           # List all apps across all projects
vessld apps:show <id>                      # Show app details and env vars
vessld apps:create <name> --project <id>   # Create an app
vessld apps:destroy <id>                   # Delete an app
```

#### Databases

```sh
vessld db:list                             # List all databases
vessld db:show <id>                        # Show database details and connection string
vessld db:create <name> <engine> --project <id>  # Create a database
vessld db:destroy <id>                     # Delete a database
```

Supported engines: `postgres`, `mysql`, `mariadb`, `redis`, `mongodb`, `clickhouse`, `kafka`, `rabbitmq`, `nats`.

#### Environment Variables

```sh
vessld env:list --project <id>             # List all env vars for a project
vessld env:set KEY=VALUE --project <id>    # Set one or more env vars
vessld env:unset KEY --project <id>        # Remove an env var
```

#### Deployments & Logs

```sh
vessld deployment:list --service <id>      # List deployment history for a service
vessld deployment:show <id>                # Show deployment details
vessld deployment:logs <id>                # Print build logs for a deployment
```

#### Custom Domains

```sh
vessld domain:list --project <id>          # List custom domains for a project
vessld domain:add <hostname> --project <id> # Add a custom domain
vessld domain:remove <id>                  # Remove a custom domain
```

---

## `vessl` — Remote CLI

The `vessl` binary runs on your **local machine** and communicates with your self-hosted `vessld` server over HTTP. This is what you install and use day-to-day from your laptop.

### Installation

```sh
go install vessl.dev/vessl/cmd/vessl@latest
```

Or download a pre-built binary from the [releases page](https://github.com/vesslhq/vessl/releases).

### Authentication

Before running any command, authenticate against your self-hosted server.

#### `login`

Prompts for your server URL, email, and password. Saves a token to `~/.vessl/config.json`.

```sh
vessl login
```

#### `logout`

Clears your saved credentials.

```sh
vessl logout
```

#### `me`

Shows the currently authenticated user.

```sh
vessl me
```

### Projects

```sh
vessl project list                         # List all projects
vessl project create <name>                # Create a project
vessl project destroy <id>                 # Delete a project
```

### Environments

```sh
vessl env list --project <id>             # List environments for a project
vessl env create <name> --project <id>    # Create an environment
vessl env destroy <id>                    # Delete an environment
```

### Applications

```sh
vessl apps list --environment <id>        # List apps in an environment
vessl apps create                         # Create an app (interactive flags)
vessl apps destroy <id>                   # Delete an app
```

#### Secrets (Environment Variables)

```sh
vessl apps secrets list --project <id>             # List env vars
vessl apps secrets set KEY=VALUE --project <id>    # Set one or more env vars
```

#### Custom Domains

```sh
vessl apps domains list --project <id>             # List custom domains
vessl apps domains add --domain <host> --project <id>  # Add a domain
vessl apps domains remove <id>                     # Remove a domain
```

#### Deployments & Logs

```sh
vessl apps deployments list --service <id>         # List deployment history
vessl apps logs <deployment-id>                    # View build logs
```

### Databases

```sh
vessl db list --project <id>              # List databases
vessl db create                           # Provision a database (interactive flags)
vessl db destroy <id>                     # Delete a database
```

#### Backups

```sh
vessl db backups list --project <id>      # List backup configurations
vessl db backups create                   # Create a backup config
vessl db backups trigger <id>             # Trigger a manual backup
vessl db backups history <id>             # View backup history
```

### Trigger a Deployment

```sh
vessl deploy <service-id>                 # Trigger a remote deployment for a service
```
