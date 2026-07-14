---
title: Getting Started
description: Install Vessl on your VPS and deploy your first application.
---

import { Steps, Tabs, TabItem, Aside } from '@astrojs/starlight/components';

Vessl turns any bare-metal VPS into your own private Vercel & Railway in under 60 seconds.

## One-Line Install

<Tabs>
  <TabItem label="curl" icon="seti:shell">
  ```sh
  curl -fsSL https://get.vessl.dev | sh
  ```
  </TabItem>
  <TabItem label="wget" icon="seti:shell">
  ```sh
  wget -qO- https://get.vessl.dev | sh
  ```
  </TabItem>
</Tabs>

This installs the `vessld` daemon, pulls the required Docker images, and starts the dashboard at `http://<your-ip>:8080`.

<Aside type="tip" title="Prerequisites">
- A Linux VPS (Ubuntu 22.04+, Debian 12+, or any modern distro with kernel 5.x+)
- Docker Engine 24+ installed and running
- A domain pointing to your server (recommended for automatic SSL)
- Ports 80 and 443 open (for Traefik reverse proxy)
</Aside>

### Docker Install

If you don't have Docker yet:

<Steps>

1. Run the official Docker installer script:

   ```sh
   curl -fsSL https://get.docker.com | sh
   ```

2. Add your user to the `docker` group:

   ```sh
   sudo usermod -aG docker $USER
   ```

3. Log out and back in for the group change to take effect.

</Steps>

## Post-Install

After the install script completes, you'll see:

```text
✅ Vessl vlatest installed successfully!
  📍 Dashboard:  http://203.0.113.42:8080
  🛠️  CLI:        vesslctl --help
```

### Create Your Account

You have two ways to create the admin account:

<Tabs>
  <TabItem label="Via Dashboard" icon="laptop">
  Open `http://<your-server-ip>:8080` in your browser. The first user to register becomes the **instance admin**.
  </TabItem>
  <TabItem label="Via Terminal" icon="terminal">
  Run the setup wizard directly from your server (no browser needed):
  ```sh
  vesslctl setup
  ```
  </TabItem>
</Tabs>

### CLI Admin Tool

The install script places `vesslctl` at `/usr/local/bin/vesslctl`. Use it for day-to-day admin tasks:

```sh
vesslctl status            # Show daemon health and running containers
vesslctl logs -f           # Tail daemon logs
vesslctl reset-password    # Reset admin password
vesslctl config            # View server configuration
vesslctl config site-name=MyVessl  # Update a setting
vesslctl backup            # Create a manual database backup
vesslctl update            # Upgrade to the latest version
vesslctl downgrade v0.1.0  # Downgrade to a specific version
```

### Configure a Domain (Recommended)

Without a domain, Vessl assigns a magic DNS subdomain (like `sslip.io` or `traefik.me`) to every service. For production:

<Steps>

1. Go to **Settings → Server Settings**.
2. Set your wildcard domain (e.g. `*.vessl.example.com`).
3. Point an `A` record or `CNAME` to your server.
4. Traefik provisions Let's Encrypt SSL certificates automatically.

</Steps>

## Quick Start: Deploy Your First App

### From Git

<Steps>

1. Click **New Project** in the dashboard.
2. Connect your GitHub or GitLab account.
3. Select a repository and branch.
4. Choose a build strategy (Dockerfile, Railpack, or Nixpacks).
5. Click **Deploy**.

</Steps>

<Aside>Vessl clones the repo, builds the image, and runs a health check before routing traffic.</Aside>

### One-Click Databases

<Steps>

1. Navigate to **Databases** in the sidebar.
2. Click **New Database**.
3. Select an engine (PostgreSQL, MySQL, Redis, MongoDB, etc.).
4. Click **Create**.

</Steps>

<Aside>Vessl provisions the container with persistent volumes and injects the connection string into your apps automatically.</Aside>

## What's Next

- [Follow the tutorial → Deploy your first app in 5 minutes](/tutorial/)
- [Deployment guide](/deployment/) — build strategies, domains, env vars, CI/CD
- [Add a database](/databases/)
- [Configure environment variables](/configuration/)
- [Set up notifications](/configuration/#notifications)
- [Manage Vessl via CLI](/admin/#cli-admin-tool)
