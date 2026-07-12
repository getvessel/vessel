---
title: Deployment
description: Deploy applications with zero-downtime.
---

Vessl supports multiple build strategies to deploy your applications.

## Deploy from Git

Connect your GitHub or GitLab account, select a repository, and deploy with one click.

## Build Strategies

- **Dockerfile** — uses your existing Dockerfile
- **Railpack / Nixpacks** — auto-detects your language and framework
- **Buildpacks** — Cloud Native Buildpacks support

## Zero-Downtime Deployments

Vessl performs health-checked container swaps. The new container must pass `/health` before traffic is routed to it.
