---
title: No Lock-In
description: Vessl apps are standard Docker containers — they survive removal of the Vessl daemon.
---

Vessl is designed so you never lose access to your applications. Every app and database runs as a standard Docker container with persistent volumes. Removing Vessl leaves your containers running.

## How It Works

- **App containers** — deployed with `--restart unless-stopped` via standard Docker.
- **Database containers** — run on named volumes that persist independently.
- **Traefik reverse proxy** — is managed by Vessl, but your containers keep running if Vessl stops.

## Uninstall Vessl Without Losing Apps

```sh
vesslctl uninstall
```

This command:
1. Stops the Vessl daemon container.
2. Removes the systemd service.
3. Removes `vesslctl` from PATH.
4. **Leaves all your app and database containers running.**

After uninstall, your apps continue serving traffic if you set up your own reverse proxy. The Traefik routing will stop, but your containers are still running on the Vessl Docker network with their assigned ports.

## Adopt Your Containers (After Uninstall)

To take manual control of your containers:

```sh
# List all running containers
docker ps --filter network=vessl-network

# Inspect an app container
docker inspect <container-name>

# View logs
docker logs <container-name>

# Set up your own reverse proxy (nginx example):
# docker run -d --name my-proxy -p 80:80 -p 443:443 ...
# Point it to your app containers on the vessl-network
```

### Databases

Database containers have persistent volumes:

```sh
# List volumes
docker volume ls | grep vessl-db

# Backup a database volume
docker run --rm -v vessl-db-data-<id>:/data -v $(pwd):/backup alpine tar czf /backup/db-backup.tar.gz -C /data .
```

## Migration to Another Platform

Since everything is standard Docker, migrating is straightforward:

1. List all running containers: `docker ps --filter network=vessl-network`
2. For each container, note the image, env vars, and volume mounts.
3. Recreate them on your new platform with the same configuration.

```sh
# Example: recreate a database container manually
docker run -d \
  --name my-postgres \
  --network vessl-network \
  -e POSTGRES_USER=vessl \
  -e POSTGRES_PASSWORD=<password> \
  -e POSTGRES_DB=vessl \
  -v vessl-db-data-<id>:/var/lib/postgresql/data \
  postgres:16-alpine
```

## Backup Before Changes

Before any major operation:

```sh
vesslctl backup
```

This creates a timestamped copy of the Vessl database at `/vessl/data/backups/`.
