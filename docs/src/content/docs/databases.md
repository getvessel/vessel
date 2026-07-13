---
title: Databases
description: Managed database engines with one-click provisioning, automated backups, and connection string injection.
---

Spin up managed databases directly from the Vessl dashboard. Each database runs in its own Docker container with persistent volumes, automatic health checks, and daily backups.

## Supported Engines

### Relational

| Engine | Version | Default Port |
|---|---|---|
| PostgreSQL | 16-alpine | 5432 |
| MySQL | 8.0 | 3306 |
| MariaDB | 11 | 3306 |
| ClickHouse | latest | 9000 |

### NoSQL

| Engine | Version | Default Port |
|---|---|---|
| MongoDB | 7.0 | 27017 |
| Redis | 7-alpine | 6379 |
| Dragonfly | latest | 6379 |
| KeyDB | latest | 6379 |

### Message Brokers

| Engine | Version | Default Port |
|---|---|---|
| Kafka | latest | 9092 |
| RabbitMQ | 4-alpine | 5672 |
| NATS | 2-alpine | 4222 |

### One-Click Deployers

| Service | Purpose | Port |
|---|---|---|
| NocoDB | Open-source Airtable alternative | 8080 |
| Plausible | Web analytics | 8000 |
| WordPress | CMS | 80 |
| Gitea | Self-hosted Git service | 3000 |

## Creating a Database

1. Navigate to **Databases** in the sidebar.
2. Click **New Database**.
3. Select an engine from the list.
4. Optionally set a custom name and port.
5. Click **Create**.

Vessl provisions the container, creates a default database and user, and mounts a persistent volume at `/var/lib/data`.

## Connection Strings

Once created, the connection string is automatically injected into every service in the same project:

```
DATABASE_URL=postgresql://vessl:<password>@<service-name>:5432/vessl
REDIS_URL=redis://<service-name>:6379
MONGO_URL=mongodb://vessl:<password>@<service-name>:27017/vessl
```

You can also find the connection details on the database's detail page in the dashboard.

## Managing Databases

### Start / Stop

Databases can be started and stopped from the dashboard. Stopping a database frees resources while preserving the volume data.

### SQL Studio

Vessl includes an in-browser SQL query editor for PostgreSQL, MySQL, and MariaDB databases:

1. Open the database detail page.
2. Click **SQL Studio**.
3. Write and execute queries directly in the browser.

### Configuration

Each database has sensible defaults:
- **Port**: Assigned from the engine's default port range
- **Username**: `vessl`
- **Database name**: `vessl`
- **Data volume**: Persisted at `<data-dir>/databases/<id>/`

## Backups

### Automated Backups

Backups run daily by default with configurable retention:

1. Navigate to **Backups** in the sidebar.
2. Click **New Backup Configuration**.
3. Select the database and set the schedule.
4. Configure retention (number of backups to keep).

### Manual Backups

Trigger a backup at any time from the database detail page or the Backups section.

### S3 Destinations

Backups can be uploaded to S3-compatible storage for offsite redundancy:

1. Go to **Backups → S3 Destinations**.
2. Add your S3-compatible endpoint (AWS S3, MinIO, Backblaze B2, etc.).
3. Provide access key, secret key, and bucket name.
4. Select the S3 destination when configuring a backup.

### Supported Backup Commands

| Engine | Backup Command |
|---|---|
| PostgreSQL | `pg_dump` |
| MySQL / MariaDB | `mysqldump` |
| MongoDB | `mongodump` |
| Redis | `redis-cli SAVE` |

### Restore

To restore from a backup:

1. Create a new database instance.
2. Download the backup record from the dashboard.
3. Restore using the engine's native restore command:

```sh
# PostgreSQL
pg_restore -U vessl -d vessl < backup.sql

# MySQL
mysql -u vessl -p vessl < backup.sql
```
