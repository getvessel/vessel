---
title: System Settings
description: Reference for Vessl system-wide settings and what they affect.
---

System Settings control server-wide behavior across your entire Vessl instance. You can access these by clicking on **Settings** in the bottom left of the Vessl dashboard.

## General Settings

- **Site Name**: The global name of your Vessl instance (e.g., used in emails).
- **Public IPs**: Explicitly define the IPv4 and IPv6 addresses for your server.
- **Server Timezone**: Sets the timezone for cron jobs, metrics, and logs.
- **Deployment Timeout**: Maximum time (in seconds) a deployment can run before being forcibly timed out (default: 3600).
- **Concurrent Builds**: Controls how many builds can happen at the exact same time (default: 2).

## Registrations

- **Allow Registrations**: Toggles whether new users can create an account on your Vessl instance.
- **Domain Allowlist**: A comma-separated list of allowed email domains (e.g., `@example.com`). If set, only users with matching emails can register.

## External Notifications

Vessl supports deep integration with external platforms to notify you of deployment statuses, database backups, and system events.

- **Discord**: Provide a webhook URL. You can optionally enable @ping notifications.
- **Slack**: Provide a Slack webhook URL.
- **Telegram**: Provide a Bot Token and Chat ID.
- **Pushover**: Provide your User Key and API Token for mobile push notifications.
- **Generic Webhooks**: Send POST requests to any custom URL with a JSON payload of the event.

## Email Providers

Vessl can send transactional emails (e.g., invites, password resets) using one of two methods:

- **SMTP**: Configure standard SMTP host, port, user, password, and from-address.
- **Resend**: Provide a Resend API key for zero-configuration email delivery.

## Maintenance & Cleanup

Vessl has built-in garbage collection to prevent your server from filling up with old images and logs.

- **Docker Cleanup Cron**: A cron expression (default: `0 0 * * *`) that dictates when to run `docker system prune` to clear unused images, stopped containers, and dangling build caches.
- **Disk Usage Alerts**: A cron expression to check disk space.
- **Disk Usage Threshold**: A percentage (e.g., 80%). If disk usage exceeds this, Vessl will alert you via your configured notification channels.

## AI Integrations

You can provide global API keys for Vessl's AI features (if enabled):

- **OpenAI API Key**: Used as the default model provider.
- **Anthropic API Key**: Alternative model provider.

## Updates & Telemetry

- **Update Check Cron**: How frequently Vessl checks for a new version of itself.
- **Auto Update Enabled**: If true, Vessl will automatically pull and deploy the latest version when detected.
- **Telemetry**: Opt-in/opt-out of anonymous usage data collection.

## Advanced Settings

- **Custom DNS Resolvers**: Override the DNS resolvers used by Vessl's internal checks.
- **MCP Server Enabled**: Toggles whether the Model Context Protocol (MCP) server is enabled on your instance.
- **IP Allowlist**: A global allowlist of IPs that are allowed to access the Vessl dashboard.
