-- Add missing notification_alerts column to server_settings
ALTER TABLE server_settings ADD COLUMN notification_alerts TEXT DEFAULT '';

-- Fix project_webhooks: add missing columns used by the repository
ALTER TABLE project_webhooks ADD COLUMN url TEXT DEFAULT '';
ALTER TABLE project_webhooks ADD COLUMN event_types TEXT DEFAULT '';
ALTER TABLE project_webhooks ADD COLUMN include_pr_environments BOOLEAN DEFAULT FALSE;

-- Fix project_members: add missing columns used by the repository
ALTER TABLE project_members ADD COLUMN email TEXT DEFAULT '';
ALTER TABLE project_members ADD COLUMN permission TEXT DEFAULT 'Can Edit';
ALTER TABLE project_members ADD COLUMN status TEXT DEFAULT 'pending';
ALTER TABLE project_members ADD COLUMN invited_at TEXT;
ALTER TABLE project_members ADD COLUMN accepted_at TEXT;
