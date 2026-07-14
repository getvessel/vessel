---
title: Administration
description: Instance management, system updates, license management, and global settings.
---

Vessl administration covers instance-wide configuration available to instance admins.

## CLI Administration

Vessl provides three CLI tools for different environments and use cases:

### Server CLI (`vessld`)

After installation, `vessld` is available at `/usr/local/bin/vessld`. This is a shell wrapper that manages the Vessl daemon by executing commands **inside the Docker container**. Use it for day-to-day server administration:

```sh
vessld status           # Show daemon health + running containers
vessld setup            # Interactive admin account wizard
vessld reset-password   # Reset admin password
vessld config           # View current configuration
vessld config <key>=<value>  # Update a setting (site-name, registration, telemetry)
vessld logs -f          # Tail daemon logs
vessld update           # Upgrade to the latest version
vessld downgrade <ver>  # Downgrade to a specific version (with backup + confirmation)
vessld backup           # Create a manual database backup
vessld restart          # Restart the Vessl daemon

# App management
vessld deploy <git-url>           # Deploy an app from a Git URL
vessld deploy --template nextjs   # Deploy a template from vesslhq/vessl-examples
vessld deploy --image nginx:latest --port 80  # Deploy from a Docker image
vessld apps:list                  # List all apps
vessld apps:show <id>             # Show app details
vessld apps:create <name>         # Create an app
vessld apps:destroy <id>         # Delete an app

# Database management
vessld db:list                    # List all databases
vessld db:show <id>              # Show database details
vessld db:create <name> <engine> # Create a database (postgres, mysql, redis, etc.)
vessld db:destroy <id>           # Delete a database
```

### Daemon CLI (`vessld`)

```sh
vessld serve              # Start the daemon (default)
vessld setup              # Setup wizard
vessld reset-password     # Reset admin password
vessld config             # View/update configuration
vessld deploy <url>       # Deploy from Git URL
vessld deploy --template  # Deploy a template
vessld apps:list          # List apps
vessld db:list            # List databases
vessld mcp                # Run MCP stdio server
vessld version            # Show version
```

### Remote CLI (`vessl`)

For remote management from your local machine, install the `vessl` client:

```sh
curl -fsSL https://get.vessl.dev/cli | sh
vessl login    # Connect to your server
```

For a full list of remote commands, see the [CLI Reference](/cli/).

### Update with `vessld update`

1. Shows current and latest available version.
2. Creates a pre-upgrade database backup automatically.
3. Pulls the new Docker image and recreates the container.
4. Your apps and databases experience zero downtime.

### Downgrade with `vessld downgrade`

1. Requires you to type `downgrade` to confirm (safety gate).
2. Creates a pre-downgrade database backup automatically.
3. Pulls the specified version and recreates the container.
4. If something breaks, restore from backup: `cp /vessl/data/backups/vessl-pre-downgrade-*.db /vessl/data/vessl.db`

## Instance Settings

Access from **Settings** in the dashboard. Only the first registered user (instance admin) can modify these.

### Global Configuration

- **Wildcard Domain**: Set the base domain for all services across all projects.
- **SMTP Configuration**: Instance-wide SMTP settings for transactional emails.
- **DNS Resolvers**: Custom DNS resolvers for container networking.
- **Port Ranges**: Configure the port pool for service allocation.

### Traefik Configuration

- **SSL Provider**: Let's Encrypt configuration.
- **HTTP Redirect**: Force HTTPS across all services.
- **Dashboard**: Enable or disable the Traefik dashboard.

## System Updates

### Automatic Update Checks

Vessl periodically checks GitHub releases for new versions. The dashboard displays a notification when an update is available.

### Manual Check

```sh
curl -X POST /api/settings/updates/check \
  -H "Authorization: Bearer vpt_xxx"
```

### Deploying an Update

1. Go to **Settings → Updates**.
2. Click **Check for Updates**.
3. If an update is available, click **Deploy Update**.
4. The system downloads and applies the update.
5. The dashboard displays real-time update progress.

### Update Process

1. New binary is downloaded from GitHub releases.
2. Database migrations are applied (backward-compatible).
3. Services are restarted gracefully.
4. Old binary is kept for rollback.

### Rollback

```sh
sudo vessld rollback
```

Restores the previous binary and a database backup taken before the upgrade.

## License Management

### Activating a License

1. Go to **Settings → License**.
2. Enter your license key.
3. Click **Activate**.

The license is validated against the licensing server and applied immediately.

### License Features

Plans may include:

- Seat limits (number of users)
- Instance limits
- Premium features (audit logs, SSO, advanced RBAC)

## Telemetry

By default, Vessl collects anonymized usage data to improve the product. This can be disabled in **Settings → Privacy**.

### What's Collected

- Instance version and uptime
- Count of projects, services, and users (no names or content)
- Deployment success/failure rates
- Enabled features and integrations

## Maintenance

### Data Directory

All persistent data is stored in the configured `VESSL_DATA_DIR` (default: `data/`):

```text
data/
├── vessl.db          # SQLite database
├── .vault_key        # Encryption key (keep safe)
├── databases/        # Database volumes
├── storage/          # MinIO storage volumes
└── backups/          # Backup archives
```

### Backup

Back up your Vessl instance:

```sh
# Stop the daemon
# Copy the data directory
cp -r data/ data-backup-$(date +%Y%m%d)
# Restart the daemon
```

### Restore

1. Stop the daemon.
2. Restore the data directory from your backup.
3. Restart the daemon.
4. Database migrations run automatically on startup.
