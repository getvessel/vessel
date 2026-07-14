---
title: Environments
description: Create isolated environments for development and production.
---

Environments provide isolation between development stages. Each environment has its own services, databases, and variables.

### Default Environments

Every new project gets a **production** environment. You can add more:

- **staging** — pre-production testing
- **dev** — development work
- **preview** — ephemeral PR previews

### Creating Environments

1. Open your project.
2. Go to **Environments**.
3. Click **New Environment**.
4. Enter a name (e.g. `staging`, `dev`).
5. Click **Create**.

### Environment Variables

Set environment variables per environment:

1. Open your project.
2. Go to **Environment Variables**.
3. Add key-value pairs.
4. These are injected into every service and job in the environment.
