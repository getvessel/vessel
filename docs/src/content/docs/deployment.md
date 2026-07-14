---
title: Deployment
description: Deploy applications with zero-downtime, multiple build strategies, and custom domains.
---

Vessl supports multiple build strategies and deployment workflows to get your applications online.

## Build Strategies

Vessl auto-detects the best build strategy based on your project. You can override it per deployment.

### Dockerfile

If your repository contains a `Dockerfile` at the root, Vessl uses it by default.

```dockerfile
FROM node:22-alpine
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
CMD ["node", "index.js"]
```

No additional configuration needed — just push and deploy.

### Railpack

[Railpack](https://railpack.com) auto-detects your language and framework. Supported stacks:

- Node.js
- Go
- Python
- Rust
- PHP
- Ruby
- Static sites (HTML/CSS/JS)

Railpack generates an optimal Dockerfile for your project without you writing one.

### Nixpacks

[Nixpacks](https://nixpacks.com) uses Nix expressions to build reproducible environments. It supports the same languages as Railpack plus additional ecosystem tools.

### Buildpacks

Cloud Native Buildpacks support is available for OCI-compliant builds. Select **Buildpacks** in the deployment settings.

### Serverless

For lightweight functions, Vessl offers a serverless deployment mode. Write code in Node.js, Python, or Go directly in the dashboard editor. See the [Serverless guide](/serverless/) for details.

## Custom Domains

Attach custom domains with automatic Let's Encrypt SSL via Traefik v3.

1. Go to your **Project → Domains**.
2. Click **Add Domain**.
3. Enter the domain (e.g. `app.example.com`).
4. Add the displayed `A` or `CNAME` record at your DNS provider.
5. SSL is provisioned automatically within seconds.

Wildcard domains are supported — configure the base domain in **Settings → Server Settings**.

## Environment Variables

Set environment variables at three levels:

| Level | Scope | Use Case |
|---|---|---|
| **Project** | All services in the project | Shared config, API keys |
| **Service** | A single service | Service-specific secrets |
| **Auto-linked** | Injected by Vessl | Database URLs, storage endpoints |

### Auto-Linked Connection Strings

When you create a database or storage instance in the same project, Vessl automatically injects connection strings as environment variables:

- `DATABASE_URL` — PostgreSQL, MySQL, MariaDB
- `REDIS_URL` — Redis, Dragonfly, KeyDB
- `MONGO_URL` — MongoDB
- `CLICKHOUSE_URL` — ClickHouse
- `KAFKA_BROKERS` — Kafka
- `RABBITMQ_URL` — RabbitMQ
- `NATS_URL` — NATS
- `S3_ENDPOINT`, `S3_ACCESS_KEY`, `S3_SECRET_KEY` — MinIO storage

### Encrypted Vault

Environment variables are encrypted at rest using AES-256-GCM. Access them from the **Variables** tab of any service or project.

## CI/CD & Git Integration

### Git Providers

Connect GitHub or GitLab to enable automatic deployments:

1. Go to **Settings → Git Apps**.
2. Install the Vessl GitHub App or configure a GitLab App.
3. Grant repository access to the repos you want to deploy.

### Automatic Deployments

Once connected, every push to the configured branch triggers a new deployment:

1. Vessl receives the webhook.
2. Clones the latest commit.
3. Builds a new container image.
4. Runs a health check on the new container.
5. Swaps traffic to the new container (zero-downtime).
6. Cleans up the old container.

### Manual Deploy

Trigger a deployment from the dashboard or CLI:

```sh
curl -X POST /api/projects/:id/deploy \
  -H "Authorization: Bearer vpt_xxx"
```

### Webhooks

Configure outgoing webhooks for deployment events:

- `deployment.started`
- `deployment.completed`
- `deployment.failed`

Webhooks can notify external services, chat platforms, or your own automation.

## PR Previews

When a pull request is opened against your connected repository, Vessl can spin up an ephemeral preview environment:

1. The webhook triggers a new deployment on the PR branch.
2. A preview URL is generated (e.g. `pr-42.myproject.vessl.example.com`).
3. The preview is automatically destroyed when the PR is merged or closed.

Enable PR previews in **Project Settings → Git Integration**.

## Rollbacks

Every deployment is a versioned release. To roll back:

1. Go to **Project → Deployments**.
2. Find the deployment you want to revert to.
3. Click **Rollback**.

Vessl redeploys the previous image with the original environment variables.

## Health Checks

Vessl performs health checks before routing traffic to a new container:

- **Endpoint**: `GET /health` (configurable)
- **Timeout**: 30 seconds
- **Retries**: 3 attempts

If the health check fails, the deployment is marked as failed and the old container continues serving traffic.

## Deployment Logs

Stream live logs during and after deployment:

1. Open the **Deployments** tab of your service.
2. Click on a deployment to view its logs.
3. Logs are available in real-time via SSE streaming.

## Build Diagnostics

If a deployment fails, use **AI Diagnose** to analyze the build logs:

1. Open the failed deployment.
2. Click **AI Diagnose**.
3. Vessl sends the logs to your configured AI provider and suggests a fix.

Configure your AI provider in **Server Settings → AI**.
