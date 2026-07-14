---
title: CLI Reference
description: Command-line interface reference for managing Vessl from the terminal.
---

Vessl provides a powerful command-line interface (`vesslctl` or `vessld`) to manage your instance directly from the terminal without using the web dashboard.

## Global Commands

### `serve`

Starts the Vessl daemon (this is the default command).

```sh
vesslctl serve
```

### `setup`

Runs the interactive setup wizard to initialize the database, create the first workspace, and provision the master admin account.

```sh
vesslctl setup
```

### `reset-password`

Resets the password for the master admin account. Useful if you lose access to the dashboard.

```sh
vesslctl reset-password
```

### `config`

View or update global system configuration variables interactively.

```sh
vesslctl config
```

### `restart`

Gracefully restarts the Vessl daemon by triggering a Docker Compose restart.

```sh
vesslctl restart
```

## Application Management

### `deploy`

Deploy an application directly from the CLI. This bypasses the dashboard and streams logs to your terminal.

**Deploy from a Git repository:**

```sh
vesslctl deploy https://github.com/your/repo
```

**Deploy from a Docker image:**

```sh
vesslctl deploy --image nginx:latest --port 80
```

**Deploy from a local Docker Compose file:**

```sh
vesslctl deploy --compose ./docker-compose.yml
```

### `apps:list`

List all deployed applications and their current status.

```sh
vesslctl apps:list
```

### `apps:show <id>`

Show detailed information, environment variables, and recent deployment history for a specific application.

```sh
vesslctl apps:show <app_id>
```

### `apps:create`

Create a new application interactively.

```sh
vesslctl apps:create
```

### `apps:destroy`

Permanently delete an application and its associated containers.

```sh
vesslctl apps:destroy
```

## Database Management

### `db:list`

List all provisioned databases across all projects.

```sh
vesslctl db:list
```

### `db:show <id>`

Display details and connection strings for a specific database.

```sh
vesslctl db:show <db_id>
```

### `db:create`

Create a new managed database interactively (PostgreSQL, MySQL, Redis, etc.).

```sh
vesslctl db:create
```

### `db:destroy`

Permanently delete a database. Warning: This action destroys all data irreversibly!

```sh
vesslctl db:destroy
```

## Advanced

### `mcp`

Run the Model Context Protocol (MCP) server over standard I/O. This allows AI assistants to securely interface with the Vessl backend for diagnostics and auto-healing.

```sh
vesslctl mcp
```

### `version`

Show the current Vessl daemon version.

```sh
vesslctl version
```
