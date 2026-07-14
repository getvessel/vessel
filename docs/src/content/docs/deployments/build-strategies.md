---
title: Build Strategies
description: Vessl auto-detects the best build strategy based on your project.
---

Vessl supports multiple build strategies and deployment workflows to get your applications online. It auto-detects the best build strategy based on your project, but you can override it per deployment.

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
