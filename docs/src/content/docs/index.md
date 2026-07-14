---
title: Vessl Docs
description: Documentation for Vessl — turn any bare-metal VPS into your own private Vercel & Railway.
template: splash
hero:
  tagline: Turn any bare-metal VPS into your own private Vercel & Railway in 60 seconds.
  image:
    file: ../../assets/hero.svg
  actions:
    - text: Get Started
      link: /getting-started/
      icon: rocket
      variant: primary
    - text: View on GitHub
      link: https://github.com/vesslhq/vessl
      icon: github
      variant: secondary
---

import { Card, CardGrid, LinkCard, Icon } from '@astrojs/starlight/components';

Vessl is a self-hosted PaaS that runs on any Linux VPS — deploy apps with zero-downtime, spin up databases with one click, manage S3 storage, and collaborate with your team — all from a single dashboard on your own hardware.

## Getting Started

<CardGrid stagger>
 <Card title="Installation" icon="rocket">
  Deploy Vessl on any Linux VPS in under 60 seconds.
    <br/>
    <LinkCard title="Read Guide" href="/getting-started/" />
 </Card>
 <Card title="Deploy an App" icon="laptop">
  Docker, Railpack, Nixpacks, or Buildpacks — pick your strategy.
    <br/>
    <LinkCard title="Read Guide" href="/deployment/" />
 </Card>
 <Card title="Create a Database" icon="seti:db">
  15 managed engines — PostgreSQL, MySQL, Redis, Kafka, and more.
    <br/>
    <LinkCard title="Read Guide" href="/databases/" />
 </Card>
</CardGrid>

## Platform Features

<CardGrid>
  <LinkCard
    title="Projects"
    description="Organize services into environments with an interactive canvas view."
    href="/projects/"
  />
  <LinkCard
    title="Storage"
    description="S3-compatible object storage with auto-injected credentials."
    href="/storage/"
  />
  <LinkCard
    title="Serverless"
    description="Write Node.js, Python, or Go functions in the browser."
    href="/serverless/"
  />
  <LinkCard
    title="Integrations"
    description="Git providers, OAuth, Vercel import, and outgoing webhooks."
    href="/integrations/"
  />
</CardGrid>

## Configuration & Management

<CardGrid>
  <LinkCard
    title="Configuration"
    description="Notifications, AI diagnostics, two-factor auth, email."
    href="/configuration/"
  />
  <LinkCard
    title="Administration"
    description="Instance settings, updates, licensing, backup and restore."
    href="/admin/"
  />
  <LinkCard
    title="API Reference"
    description="Full REST API with personal access tokens and rate limits."
    href="/api/"
  />
</CardGrid>

## Community

- **[GitHub Issues](https://github.com/vesslhq/vessl/issues)** — Report bugs and request features
- **[Discord](https://discord.gg/vessl)** — Community discussions and support
- **[GitHub Discussions](https://github.com/vesslhq/vessl/discussions)** — Ask questions and share ideas
