---
title: DNS Providers
description: Connect DNS providers and let Vessl create or update A records for service domains.
---

Vessl can create or update service-domain `A` records through supported DNS providers.

## Supported Providers

- Cloudflare
- Namecheap
- Spaceship

Connect providers from System Settings, DNS.

## What Automation Does

From a service Domains tab, Vessl can apply an `A` record for a public hostname. The record points at the server public IPv4 address.

Automation does not run for local loopback domains, and it currently supports IPv4 `A` records only.

## Cloudflare

Cloudflare settings include:

- API token or auth key.
- Account email when using auth key style credentials.
- Optional zone ID.

When zone ID is blank, Vessl searches Cloudflare zones from the hostname. It creates or updates an unproxied `A` record with automatic TTL.

## Namecheap

Namecheap settings include:

- API user.
- API key.
- Client IP.

Namecheap requires API access to be enabled on the account and expects the client IP used for API requests. Vessl preserves existing host records and adds or replaces the matching `A` record.

## Spaceship

Spaceship settings include:

- API key.
- API secret.

Vessl resolves the domain, removes conflicting records for the same host when necessary, then writes the `A` record with TTL `1800`.

## Common Issues

- The provider token cannot access the domain or zone.
- The hostname is not inside a domain managed by the provider account.
- The target IP is not an IPv4 address.
- Existing CNAME or ALIAS records conflict with the `A` record.
- DNS propagation has not completed yet.

After applying a record, use the service Domains tab to refresh and verify.
