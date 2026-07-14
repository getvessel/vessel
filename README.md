# 🛰️ Vessl

**Self-hosted PaaS. Turn any VPS into your own Vercel or Railway in 60 seconds.**

---

Vessl is a lightweight, open-source Platform-as-a-Service (PaaS) designed to simplify deployments. Whether you're deploying a static site, a full-stack monorepo, or a complex microservice architecture, Vessl provides a frictionless developer experience without the vendor lock-in.

## 🚀 Quick Start

Install Vessl on any fresh Linux server (Ubuntu/Debian recommended):

```bash
curl -fsSL https://get.vessl.dev | sh
```

Once installed, your dashboard will be available at `http://your-server-ip:8080`.

## ✨ Features

Vessl is built to be simple but powerful, giving you everything you need to run production workloads out of the box.

- **Deploy Anything:** Native support for Dockerfiles, Railpack, Nixpacks, and standard Buildpacks.
- **Managed Databases:** Provision PostgreSQL, MySQL, Redis, MongoDB, and more with a single click.
- **Smart Environment:** Database credentials (`DATABASE_URL`, `REDIS_URL`) are automatically injected into your linked applications.
- **Zero-Downtime Deploys:** Seamless container swaps with built-in health checks and instant rollbacks.
- **Custom Domains & SSL:** Automatic Let's Encrypt certificates managed via Traefik v3.
- **GitOps Ready:** Connect to GitHub/GitLab for automatic deployments on push and PR preview environments.
- **Marketplace Templates:** Instantly deploy popular frameworks (Node.js, Go, Python, Ruby, PHP) from our built-in marketplace.
- **No Lock-in:** Vessl orchestrates standard Docker containers. If you ever remove Vessl, your apps keep running.

## 💻 CLI

Vessl ships two CLI tools.

### `vessld` — Server Daemon CLI

Runs **on your VPS**. Manages the server process and all resources directly.

```bash
# Server management
vessld serve                  # Start the daemon (default)
vessld setup                  # Run the initial admin configuration wizard
vessld reset-password         # Reset the admin password
vessld config                 # View or update server configuration

# Deployments
vessld deploy https://github.com/user/repo.git   # Deploy from Git
vessld deploy --image nginx:latest --port 80     # Deploy a Docker image
vessld deploy --compose docker-compose.yml       # Deploy a Compose stack

# Resource management
vessld project:list                              # List all projects
vessld apps:list                                 # List all applications
vessld apps:show <id>                            # Show app details
vessld db:list                                   # List all databases
vessld db:create my-db postgres --project <id>  # Provision a database
vessld env:list --project <id>                  # List environment variables
vessld env:set KEY=VALUE --project <id>         # Set an environment variable
vessld deployment:list --service <id>           # List deployment history
vessld deployment:logs <id>                     # View build logs
vessld domain:list --project <id>              # List custom domains
```

### `vessl` — Remote CLI

Runs **on your local machine**. Connects to your self-hosted `vessld` server over HTTP.

**Installation:**

```bash
go install vessl.dev/vessl/cmd/vessl@latest
```

**Usage:**

```bash
# Authentication
vessl login                           # Connect to your self-hosted server
vessl logout                          # Clear saved credentials
vessl me                              # Show current logged-in user

# Projects & Environments
vessl project list                    # List all projects
vessl project create <name>           # Create a project
vessl project destroy <id>            # Delete a project
vessl env list                        # List environments
vessl env create <name>               # Create an environment

# Applications
vessl apps list --environment <id>    # List apps
vessl apps create                     # Create an app
vessl apps destroy <id>               # Delete an app
vessl apps secrets list --project <id>            # List env vars
vessl apps secrets set KEY=VALUE --project <id>   # Set env var(s)
vessl apps domains list --project <id>            # List custom domains
vessl apps domains add <host> --project <id>      # Add a custom domain
vessl apps deployments list --service <id>        # Deployment history
vessl apps logs <deployment-id>                   # View build logs

# Databases
vessl db list --project <id>          # List databases
vessl db create                       # Provision a database
vessl db destroy <id>                 # Delete a database
vessl db backups list --project <id>  # List backup configs
vessl db backups trigger <id>         # Trigger a manual backup

# Deployments
vessl deploy <service-id>             # Trigger a remote deployment
```

## 🛠️ Local Development

Want to contribute or hack on Vessl locally?

```bash
# 1. Clone the repository
git clone https://github.com/vesslhq/vessl.git
cd vessl

# 2. Setup your environment
cp .env.example .env

# 3. Run the Go daemon locally (starts on :8080)
go run ./cmd/vessld
```

**Requirements:** Go 1.22+, Node.js 22+, and Docker.

## 📚 Documentation

For complete guides, API references, and advanced configuration, please visit our documentation at **[docs.vessl.dev](https://docs.vessl.dev)**.
