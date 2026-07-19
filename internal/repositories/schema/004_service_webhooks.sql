CREATE TABLE IF NOT EXISTS service_webhooks (
    id TEXT PRIMARY KEY,
    service_id TEXT NOT NULL REFERENCES app_services(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    event_types TEXT DEFAULT '',
    include_pr_environments BOOLEAN DEFAULT FALSE,
    created_at DATETIME,
    updated_at DATETIME
);
