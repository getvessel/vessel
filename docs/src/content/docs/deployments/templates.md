---
title: Templates & Examples
description: Deploy starter kits and official examples to Vessl in one click.
---

Vessl provides a rich ecosystem of templates and examples to help you get started quickly. Whether you need a simple static site, a full-stack Next.js application, or a Go microservice, you can deploy it instantly.

## One-Click Templates

Vessl's **One-Click Deploy** feature allows you to launch common open-source tools and databases instantly.

When creating a new service, select from the available templates:
- **NocoDB**
- **Directus**
- **Plausible Analytics**
- **Umami**
- **Ghost**
- **Redis**
- **PostgreSQL / MySQL / MariaDB**

These templates are pre-configured with the correct environment variables, Docker images, and persistent storage volumes out of the box.

## The Examples Repository

For code-based starters, check out the official [vesslhq/vessl-examples](https://github.com/vesslhq/vessl-examples) repository. 

This repository contains ready-to-deploy boilerplates for popular frameworks:

- **Next.js** (App Router, Pages Router)
- **Astro**
- **Vite** (React, Vue, Svelte)
- **Node.js / Express**
- **Go**
- **Python / FastAPI**

### How to use an example

1. Go to the [vesslhq/vessl-examples](https://github.com/vesslhq/vessl-examples) repository on GitHub.
2. Find the example you want to use.
3. Fork the repository or copy the code into your own Git repository.
4. In the Vessl dashboard, create a new **Source Service**.
5. Connect your GitHub account and select your new repository.
6. Vessl will automatically detect the framework and build it seamlessly using Railpack or Nixpacks.
