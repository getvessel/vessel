ALTER TABLE app_services ADD COLUMN maintenance_mode BOOLEAN NOT NULL DEFAULT 0;

CREATE TABLE log_drains (
    id TEXT PRIMARY KEY,
    service_id TEXT NOT NULL,
    drain_type TEXT NOT NULL,
    endpoint_url TEXT NOT NULL,
    auth_token TEXT,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    FOREIGN KEY(service_id) REFERENCES app_services(id) ON DELETE CASCADE
);
