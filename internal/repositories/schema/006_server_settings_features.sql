ALTER TABLE server_settings ADD COLUMN concurrent_builds INTEGER NOT NULL DEFAULT 2;
ALTER TABLE server_settings ADD COLUMN deployment_timeout INTEGER NOT NULL DEFAULT 3600;
ALTER TABLE server_settings ADD COLUMN server_timezone TEXT NOT NULL DEFAULT 'UTC';
ALTER TABLE server_settings ADD COLUMN docker_cleanup_cron TEXT NOT NULL DEFAULT '0 0 * * *';
ALTER TABLE server_settings ADD COLUMN disk_usage_threshold INTEGER NOT NULL DEFAULT 80;
ALTER TABLE server_settings ADD COLUMN disk_usage_cron TEXT NOT NULL DEFAULT '0 23 * * *';
