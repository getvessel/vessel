package settings

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"
)

type SQLiteRepository struct {
	db *sql.DB
	mu sync.Mutex
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{db: db}
}

func (r *SQLiteRepository) GetServerSettings(ctx context.Context) (*ServerSettings, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var cfg ServerSettings
	err := r.db.QueryRowContext(ctx, `SELECT id, caddy_wildcard_ip, discord_webhook_url, slack_webhook_url, telegram_bot_token, telegram_chat_id, smtp_host, smtp_port, smtp_user, smtp_password, COALESCE(smtp_from_name, ''), COALESCE(smtp_from_address, ''), notification_alerts,
	                             registration_enabled, custom_dns_resolvers, dns_validation_enabled, ip_allowlist, mcp_server_enabled, update_check_cron, auto_update_enabled, current_version, latest_version, last_update_check, updated_at
	                      FROM server_settings WHERE id = 'global'`).
		Scan(&cfg.ID, &cfg.CaddyWildcardIP, &cfg.DiscordWebhookURL, &cfg.SlackWebhookURL, &cfg.TelegramBotToken, &cfg.TelegramChatID, &cfg.SMTPHost, &cfg.SMTPPort, &cfg.SMTPUser, &cfg.SMTPPassword, &cfg.SMTPFromName, &cfg.SMTPFromAddress, &cfg.NotificationAlerts,
			&cfg.RegistrationEnabled, &cfg.CustomDNSResolvers, &cfg.DNSValidationEnabled, &cfg.IPAllowlist, &cfg.MCPServerEnabled, &cfg.UpdateCheckCron, &cfg.AutoUpdateEnabled, &cfg.CurrentVersion, &cfg.LatestVersion, &cfg.LastUpdateCheck, &cfg.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		defaultSettings := &ServerSettings{
			ID:                   "global",
			CaddyWildcardIP:      "127.0.0.1",
			NotificationAlerts:   true,
			RegistrationEnabled:  true,
			DNSValidationEnabled: true,
			MCPServerEnabled:     true,
			UpdateCheckCron:      "0 * * * *",
			AutoUpdateEnabled:    false,
			CurrentVersion:       "0.1.0",
			LatestVersion:        "0.1.0",
			UpdatedAt:            time.Now().UTC().Format(time.RFC3339),
		}
		query := `INSERT INTO server_settings (id, caddy_wildcard_ip, discord_webhook_url, slack_webhook_url, telegram_bot_token, telegram_chat_id, smtp_host, smtp_port, smtp_user, smtp_password, smtp_from_name, smtp_from_address, notification_alerts,
		                                       registration_enabled, custom_dns_resolvers, dns_validation_enabled, ip_allowlist, mcp_server_enabled, update_check_cron, auto_update_enabled, current_version, latest_version, last_update_check, updated_at)
		          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		_, _ = r.db.ExecContext(ctx, query, defaultSettings.ID, defaultSettings.CaddyWildcardIP, defaultSettings.DiscordWebhookURL, defaultSettings.SlackWebhookURL, defaultSettings.TelegramBotToken, defaultSettings.TelegramChatID, defaultSettings.SMTPHost, defaultSettings.SMTPPort, defaultSettings.SMTPUser, defaultSettings.SMTPPassword, defaultSettings.SMTPFromName, defaultSettings.SMTPFromAddress, defaultSettings.NotificationAlerts,
			defaultSettings.RegistrationEnabled, defaultSettings.CustomDNSResolvers, defaultSettings.DNSValidationEnabled, defaultSettings.IPAllowlist, defaultSettings.MCPServerEnabled, defaultSettings.UpdateCheckCron, defaultSettings.AutoUpdateEnabled, defaultSettings.CurrentVersion, defaultSettings.LatestVersion, defaultSettings.LastUpdateCheck, defaultSettings.UpdatedAt)
		return defaultSettings, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get server settings: %w", err)
	}
	return &cfg, nil
}

func (r *SQLiteRepository) UpdateServerSettings(ctx context.Context, cfg *ServerSettings) error {
	if cfg.ID == "" {
		cfg.ID = "global"
	}
	cfg.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	r.mu.Lock()
	defer r.mu.Unlock()

	query := `INSERT INTO server_settings (id, caddy_wildcard_ip, discord_webhook_url, slack_webhook_url, telegram_bot_token, telegram_chat_id, smtp_host, smtp_port, smtp_user, smtp_password, smtp_from_name, smtp_from_address, notification_alerts,
	                                       registration_enabled, custom_dns_resolvers, dns_validation_enabled, ip_allowlist, mcp_server_enabled, update_check_cron, auto_update_enabled, current_version, latest_version, last_update_check, updated_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	          ON CONFLICT(id) DO UPDATE SET
	          caddy_wildcard_ip = excluded.caddy_wildcard_ip,
	          discord_webhook_url = excluded.discord_webhook_url,
	          slack_webhook_url = excluded.slack_webhook_url,
	          telegram_bot_token = excluded.telegram_bot_token,
	          telegram_chat_id = excluded.telegram_chat_id,
	          smtp_host = excluded.smtp_host,
	          smtp_port = excluded.smtp_port,
	          smtp_user = excluded.smtp_user,
	          smtp_password = excluded.smtp_password,
	          smtp_from_name = excluded.smtp_from_name,
	          smtp_from_address = excluded.smtp_from_address,
	          notification_alerts = excluded.notification_alerts,
	          registration_enabled = excluded.registration_enabled,
	          custom_dns_resolvers = excluded.custom_dns_resolvers,
	          dns_validation_enabled = excluded.dns_validation_enabled,
	          ip_allowlist = excluded.ip_allowlist,
	          mcp_server_enabled = excluded.mcp_server_enabled,
	          update_check_cron = excluded.update_check_cron,
	          auto_update_enabled = excluded.auto_update_enabled,
	          current_version = excluded.current_version,
	          latest_version = excluded.latest_version,
	          last_update_check = excluded.last_update_check,
	          updated_at = excluded.updated_at`
	_, err := r.db.ExecContext(ctx, query, cfg.ID, cfg.CaddyWildcardIP, cfg.DiscordWebhookURL, cfg.SlackWebhookURL, cfg.TelegramBotToken, cfg.TelegramChatID, cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPassword, cfg.SMTPFromName, cfg.SMTPFromAddress, cfg.NotificationAlerts,
		cfg.RegistrationEnabled, cfg.CustomDNSResolvers, cfg.DNSValidationEnabled, cfg.IPAllowlist, cfg.MCPServerEnabled, cfg.UpdateCheckCron, cfg.AutoUpdateEnabled, cfg.CurrentVersion, cfg.LatestVersion, cfg.LastUpdateCheck, cfg.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update server settings: %w", err)
	}
	return nil
}

func (r *SQLiteRepository) ListProjects(ctx context.Context) ([]map[string]any, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	rows, err := r.db.QueryContext(ctx, `SELECT id, name, repo_url, branch, status, updated_at FROM projects ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}
	defer rows.Close()

	var projects []map[string]any
	for rows.Next() {
		var id, name, repoURL, branch, status, updatedAt string
		if err := rows.Scan(&id, &name, &repoURL, &branch, &status, &updatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, map[string]any{
			"id":        id,
			"name":      name,
			"repoUrl":   repoURL,
			"branch":    branch,
			"status":    status,
			"updatedAt": updatedAt,
		})
	}
	return projects, nil
}
