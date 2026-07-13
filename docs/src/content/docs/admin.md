---
title: Administration
description: Instance management, system updates, license management, and global settings.
---

Vessl administration covers instance-wide configuration available to instance admins.

## Instance Settings

Access from **Settings** in the dashboard. Only the first registered user (instance admin) can modify these.

### Global Configuration

- **Wildcard Domain**: Set the base domain for all services across all workspaces.
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
sudo vesslctl rollback
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
- Workspace limits
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
