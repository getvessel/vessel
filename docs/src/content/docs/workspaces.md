---
title: Workspaces & Teams
description: Multi-tenant workspaces with team management, audit logs, and SSO.
---

Workspaces enable multi-tenancy — each workspace is an isolated team with its own projects, members, and settings.

## Creating a Workspace

1. Click the workspace switcher in the top bar.
2. Click **New Workspace**.
3. Enter a name and optional slug.
4. Click **Create**.

You can belong to multiple workspaces and switch between them.

## Members & Roles

### Inviting Members

1. Go to **Workspace Settings → Members**.
2. Click **Invite Member**.
3. Enter their email address.
4. Select a role:

| Role          | Permissions                                       |
| ------------- | ------------------------------------------------- |
| **Owner**     | Full access, can delete workspace, manage billing |
| **Admin**     | Manage projects, members, and settings            |
| **Developer** | Create and manage projects and services           |
| **Viewer**    | Read-only access to all resources                 |

1. Click **Send Invite**.

The invitee receives an email with an acceptance link. Invites can be revoked before acceptance.

### Removing Members

1. Go to **Workspace Settings → Members**.
2. Click the remove button next to the member.
3. Confirm the action.

## Trusted Domains

Restrict sign-ups and invitations to specific email domains:

1. Go to **Workspace Settings → Trusted Domains**.
2. Click **Add Domain**.
3. Enter a domain (e.g. `company.com`).
4. Only users with emails matching the domain can join.

## SSH Keys

Add SSH keys for Git access to private repositories:

1. Go to **Workspace Settings → SSH Keys**.
2. Click **Add SSH Key**.
3. Paste your public key.
4. Vessl uses this key when cloning repositories from connected Git providers.

## Audit Logs

Track every action in the workspace:

1. Go to **Audit Logs** in the sidebar.
2. View a chronological list of events:

- Member invitations and role changes
- Project creation and deletion
- Deployment triggers and rollbacks
- Database and storage operations
- Settings modifications

Audit logs include timestamps, actor information, action type, and contextual metadata.

## Deleting a Workspace

1. Go to **Workspace Settings**.
2. Scroll to the bottom and click **Delete Workspace**.
3. Confirm the deletion.

All projects, data, and configurations within the workspace are removed permanently.
