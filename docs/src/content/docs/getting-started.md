---
title: Getting Started
description: Install Vessl on your VPS and deploy your first application.
---

Vessl turns any bare-metal VPS into your own private Vercel & Railway in under 60 seconds.

## One-Line Install

```sh
curl -fsSL https://get.vessl.dev | sh
```

This installs the `vessld` daemon, pulls the required Docker images, and starts the dashboard at `http://<your-ip>:8080`.

## Prerequisites

- A Linux VPS (Ubuntu 22.04+, Debian 12+, or any modern distro with kernel 5.x+)
- Docker Engine 24+ installed and running
- A domain pointing to your server (recommended for automatic SSL)
- Ports 80 and 443 open (for Traefik reverse proxy)

### Docker Install

If you don't have Docker yet:

```sh
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER
```

Log out and back in for the group change to take effect.

## Post-Install

After the install script completes, open `http://<your-server-ip>:8080` in your browser.

### Create Your Account

1. The first user to register becomes the **instance admin**.
2. Enter your name, email, and password.
3. You'll be prompted to create your first workspace.

### Configure a Domain (Recommended)

Without a domain, Vessl assigns an `sslip.io` subdomain to every service. For production:

1. Go to **Settings → Server Settings**.
2. Set your wildcard domain (e.g. `*.vessl.example.com`).
3. Point an `A` record or `CNAME` to your server.
4. Traefik provisions Let's Encrypt SSL certificates automatically.

## Quick Start: Deploy Your First App

### From Git

1. Click **New Project** in the dashboard.
2. Connect your GitHub or GitLab account.
3. Select a repository and branch.
4. Choose a build strategy (Dockerfile, Railpack, or Nixpacks).
5. Click **Deploy**.

Vessl clones the repo, builds the image, and runs a health check before routing traffic.

### From a Public Git URL

1. Click **New Project → Deploy from Git URL**.
2. Paste a public repository URL (e.g. `https://github.com/user/repo.git`).
3. Configure the build command if needed.
4. Click **Deploy**.

### One-Click Databases

1. Navigate to **Databases** in the sidebar.
2. Click **New Database**.
3. Select an engine (PostgreSQL, MySQL, Redis, MongoDB, etc.).
4. Click **Create**.

Vessl provisions the container with persistent volumes and injects the connection string into your apps automatically.

## What's Next

- [Deploy your first app](/deployment/)
- [Add a database](/databases/)
- [Configure environment variables](/configuration/)
- [Set up notifications](/configuration/#notifications)
