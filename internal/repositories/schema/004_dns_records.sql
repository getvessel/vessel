CREATE TABLE IF NOT EXISTS dns_records (
    id TEXT PRIMARY KEY,
    domain_name TEXT NOT NULL,
    record_type TEXT NOT NULL,
    record_name TEXT NOT NULL,
    record_value TEXT NOT NULL,
    ttl INTEGER DEFAULT 3600,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
