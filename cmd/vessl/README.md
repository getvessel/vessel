# `vessl` — Remote CLI

The `vessl` binary is the remote client for your self-hosted Vessl server. It runs on your **local machine** and communicates with `vessld` over HTTP using a saved token.

## Installation

```sh
curl -fsSL https://get.vessl.dev/cli | sh
```

Or if you have Go installed:

```sh
go install vessl.dev/vessl/cmd/vessl@latest
```

After installing, authenticate against your server:

```sh
vessl login
```

This prompts for your server URL, email, and password and saves credentials to `~/.vessl/config.json`.

## Commands

### Auth

```sh
vessl login                            # Authenticate to your server
vessl logout                           # Clear saved credentials
vessl me                               # Show current logged-in user
```

### Projects

```sh
vessl project list                     # List all projects
vessl project create <name>            # Create a project
vessl project destroy <id>             # Delete a project
```

### Environments

```sh
vessl env list --project <id>          # List environments
vessl env create <name> --project <id> # Create an environment
vessl env destroy <id>                 # Delete an environment
```

### Applications

```sh
vessl apps list --environment <id>     # List apps
vessl apps create                      # Create an app
vessl apps destroy <id>                # Delete an app
```

#### Secrets (Env Vars)

```sh
vessl apps secrets list --project <id>
vessl apps secrets set KEY=VALUE --project <id>
```

#### Custom Domains

```sh
vessl apps domains list --project <id>
vessl apps domains add --domain <host> --project <id>
vessl apps domains remove <id>
```

#### Deployments & Logs

```sh
vessl apps deployments list --service <id>
vessl apps logs <deployment-id>
```

### Databases

```sh
vessl db list --project <id>           # List databases
vessl db create                        # Provision a database
vessl db destroy <id>                  # Delete a database
```

#### Backups

```sh
vessl db backups list --project <id>
vessl db backups create
vessl db backups trigger <id>
vessl db backups history <id>
```

### Deployments

```sh
vessl deploy <service-id>              # Trigger a remote deployment
```

## Config

Credentials are stored at `~/.vessl/config.json`:

```json
{
  "serverUrl": "https://your-server.com",
  "token": "<jwt>",
  "email": "you@example.com"
}
```

Run `vessl logout` to clear this file.
