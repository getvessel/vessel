package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"vessel.dev/vessel/internal/models"
)

type SettingsRepository interface {
	GetServerSettings(ctx context.Context) (*models.ServerSettings, error)
	UpdateServerSettings(ctx context.Context, cfg *models.ServerSettings) error
	ListProjects(ctx context.Context) ([]map[string]any, error)
}

type NotificationRepository interface {
	GetIntegration(ctx context.Context) (*models.NotificationIntegration, error)
	SaveIntegration(ctx context.Context, n *models.NotificationIntegration) error
	GetProjectPref(ctx context.Context, projectID string) (*models.ProjectNotificationPref, error)
	SaveProjectPref(ctx context.Context, pref *models.ProjectNotificationPref) error
}

type SettingsSQLiteRepository struct {
	db *sql.DB
	mu sync.Mutex
}

func NewSettingsSQLiteRepository(db *sql.DB) *SettingsSQLiteRepository {
	return &SettingsSQLiteRepository{db: db}
}

func (r *SettingsSQLiteRepository) GetServerSettings(ctx context.Context) (*models.ServerSettings, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var cfg models.ServerSettings
	err := r.db.QueryRowContext(ctx, `SELECT id, caddy_wildcard_ip, discord_webhook_url, slack_webhook_url, telegram_bot_token, telegram_chat_id, 
		 smtp_host, smtp_port, smtp_user, smtp_password, smtp_from_name, smtp_from_address, notification_alerts, 
		 registration_enabled, registration_domain_allowlist, custom_dns_resolvers, dns_validation_enabled, ip_allowlist, mcp_server_enabled, default_wildcard_domain, 
		 update_check_cron, auto_update_enabled, current_version, latest_version, last_update_check, updated_at 
		 FROM server_settings LIMIT 1`).
		Scan(&cfg.ID, &cfg.CaddyWildcardIP, &cfg.DiscordWebhookURL, &cfg.SlackWebhookURL, &cfg.TelegramBotToken, &cfg.TelegramChatID,
			&cfg.SMTPHost, &cfg.SMTPPort, &cfg.SMTPUser, &cfg.SMTPPassword, &cfg.SMTPFromName, &cfg.SMTPFromAddress, &cfg.NotificationAlerts,
			&cfg.RegistrationEnabled, &cfg.RegistrationDomainAllowlist, &cfg.CustomDNSResolvers, &cfg.DNSValidationEnabled, &cfg.IPAllowlist, &cfg.MCPServerEnabled, &cfg.DefaultWildcardDomain,
			&cfg.UpdateCheckCron, &cfg.AutoUpdateEnabled, &cfg.CurrentVersion, &cfg.LatestVersion, &cfg.LastUpdateCheck, &cfg.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		defaultSettings := &models.ServerSettings{
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
		                                       registration_enabled, registration_domain_allowlist, custom_dns_resolvers, dns_validation_enabled, ip_allowlist, mcp_server_enabled, default_wildcard_domain, update_check_cron, auto_update_enabled, current_version, latest_version, last_update_check, updated_at)
		          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		_, _ = r.db.ExecContext(ctx, query, defaultSettings.ID, defaultSettings.CaddyWildcardIP, defaultSettings.DiscordWebhookURL, defaultSettings.SlackWebhookURL, defaultSettings.TelegramBotToken, defaultSettings.TelegramChatID, defaultSettings.SMTPHost, defaultSettings.SMTPPort, defaultSettings.SMTPUser, defaultSettings.SMTPPassword, defaultSettings.SMTPFromName, defaultSettings.SMTPFromAddress, defaultSettings.NotificationAlerts,
			defaultSettings.RegistrationEnabled, defaultSettings.RegistrationDomainAllowlist, defaultSettings.CustomDNSResolvers, defaultSettings.DNSValidationEnabled, defaultSettings.IPAllowlist, defaultSettings.MCPServerEnabled, defaultSettings.DefaultWildcardDomain, defaultSettings.UpdateCheckCron, defaultSettings.AutoUpdateEnabled, defaultSettings.CurrentVersion, defaultSettings.LatestVersion, defaultSettings.LastUpdateCheck, defaultSettings.UpdatedAt)
		return defaultSettings, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get server settings: %w", err)
	}
	return &cfg, nil
}

func (r *SettingsSQLiteRepository) UpdateServerSettings(ctx context.Context, cfg *models.ServerSettings) error {
	if cfg.ID == "" {
		cfg.ID = "global"
	}
	cfg.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	r.mu.Lock()
	defer r.mu.Unlock()

	query := `INSERT INTO server_settings (id, caddy_wildcard_ip, discord_webhook_url, slack_webhook_url, telegram_bot_token, telegram_chat_id, smtp_host, smtp_port, smtp_user, smtp_password, smtp_from_name, smtp_from_address, notification_alerts,
	                                       registration_enabled, registration_domain_allowlist, custom_dns_resolvers, dns_validation_enabled, ip_allowlist, mcp_server_enabled, default_wildcard_domain, update_check_cron, auto_update_enabled, current_version, latest_version, last_update_check, updated_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
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
	          registration_domain_allowlist = excluded.registration_domain_allowlist,
	          custom_dns_resolvers = excluded.custom_dns_resolvers,
	          dns_validation_enabled = excluded.dns_validation_enabled,
	          ip_allowlist = excluded.ip_allowlist,
	          mcp_server_enabled = excluded.mcp_server_enabled,
	          default_wildcard_domain = excluded.default_wildcard_domain,
	          update_check_cron = excluded.update_check_cron,
	          auto_update_enabled = excluded.auto_update_enabled,
	          current_version = excluded.current_version,
	          latest_version = excluded.latest_version,
	          last_update_check = excluded.last_update_check,
	          updated_at = excluded.updated_at`
	_, err := r.db.ExecContext(ctx, query, cfg.ID, cfg.CaddyWildcardIP, cfg.DiscordWebhookURL, cfg.SlackWebhookURL, cfg.TelegramBotToken, cfg.TelegramChatID, cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPassword, cfg.SMTPFromName, cfg.SMTPFromAddress, cfg.NotificationAlerts,
		cfg.RegistrationEnabled, cfg.RegistrationDomainAllowlist, cfg.CustomDNSResolvers, cfg.DNSValidationEnabled, cfg.IPAllowlist, cfg.MCPServerEnabled, cfg.DefaultWildcardDomain, cfg.UpdateCheckCron, cfg.AutoUpdateEnabled, cfg.CurrentVersion, cfg.LatestVersion, cfg.LastUpdateCheck, cfg.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update server settings: %w", err)
	}
	return nil
}

func (r *SettingsSQLiteRepository) ListProjects(ctx context.Context) ([]map[string]any, error) {
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

type NotificationSQLiteRepository struct {
	db *sql.DB
}

func NewNotificationSQLiteRepository(db *sql.DB) *NotificationSQLiteRepository {
	return &NotificationSQLiteRepository{db: db}
}

func (r *NotificationSQLiteRepository) GetIntegration(ctx context.Context) (*models.NotificationIntegration, error) {
	query := `SELECT id, smtp_enabled, COALESCE(smtp_host, ''), COALESCE(smtp_port, 587), COALESCE(smtp_user, ''), COALESCE(smtp_password, ''), COALESCE(smtp_from_name, ''), COALESCE(smtp_from_address, ''), resend_enabled, COALESCE(resend_api_key, ''), slack_enabled, COALESCE(slack_webhook_url, ''), discord_enabled, COALESCE(discord_webhook_url, ''), discord_ping_enabled, telegram_enabled, COALESCE(telegram_bot_token, ''), COALESCE(telegram_chat_id, ''), pushover_enabled, COALESCE(pushover_user_key, ''), COALESCE(pushover_api_token, ''), webhook_enabled, COALESCE(webhook_url, ''), COALESCE(updated_at, '') FROM notification_integrations WHERE id = 'global'`

	row := r.db.QueryRowContext(ctx, query)
	var n models.NotificationIntegration
	var smtpHost, smtpUser, smtpPassword, smtpFromName, smtpFromAddress, resendKey, slackURL, discordURL, telegramBot, telegramChat, pushoverUser, pushoverToken, webhookURL, updatedAt string
	var smtpPort int

	err := row.Scan(
		&n.ID, &n.SMTPEnabled, &smtpHost, &smtpPort, &smtpUser, &smtpPassword, &smtpFromName, &smtpFromAddress,
		&n.ResendEnabled, &resendKey, &n.SlackEnabled, &slackURL,
		&n.DiscordEnabled, &discordURL, &n.DiscordPingEnabled,
		&n.TelegramEnabled, &telegramBot, &telegramChat,
		&n.PushoverEnabled, &pushoverUser, &pushoverToken,
		&n.WebhookEnabled, &webhookURL, &updatedAt,
	)
	if err == sql.ErrNoRows {
		return &models.NotificationIntegration{ID: "global"}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan notification integration: %w", err)
	}

	n.SMTPHost = smtpHost
	n.SMTPPort = smtpPort
	n.SMTPUser = smtpUser
	n.SMTPPassword = smtpPassword
	n.SMTPFromName = smtpFromName
	n.SMTPFromAddress = smtpFromAddress
	n.ResendAPIKey = resendKey
	n.SlackWebhookURL = slackURL
	n.DiscordWebhookURL = discordURL
	n.TelegramBotToken = telegramBot
	n.TelegramChatID = telegramChat
	n.PushoverUserKey = pushoverUser
	n.PushoverAPIToken = pushoverToken
	n.WebhookURL = webhookURL
	n.UpdatedAt = updatedAt

	return &n, nil
}

func (r *NotificationSQLiteRepository) SaveIntegration(ctx context.Context, n *models.NotificationIntegration) error {
	n.ID = "global"
	n.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	query := `INSERT INTO notification_integrations (
		id, smtp_enabled, smtp_host, smtp_port, smtp_user, smtp_password, smtp_from_name, smtp_from_address,
		resend_enabled, resend_api_key, slack_enabled, slack_webhook_url,
		discord_enabled, discord_webhook_url, discord_ping_enabled,
		telegram_enabled, telegram_bot_token, telegram_chat_id,
		pushover_enabled, pushover_user_key, pushover_api_token,
		webhook_enabled, webhook_url, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		smtp_enabled = excluded.smtp_enabled,
		smtp_host = excluded.smtp_host,
		smtp_port = excluded.smtp_port,
		smtp_user = excluded.smtp_user,
		smtp_password = excluded.smtp_password,
		smtp_from_name = excluded.smtp_from_name,
		smtp_from_address = excluded.smtp_from_address,
		resend_enabled = excluded.resend_enabled,
		resend_api_key = excluded.resend_api_key,
		slack_enabled = excluded.slack_enabled,
		slack_webhook_url = excluded.slack_webhook_url,
		discord_enabled = excluded.discord_enabled,
		discord_webhook_url = excluded.discord_webhook_url,
		discord_ping_enabled = excluded.discord_ping_enabled,
		telegram_enabled = excluded.telegram_enabled,
		telegram_bot_token = excluded.telegram_bot_token,
		telegram_chat_id = excluded.telegram_chat_id,
		pushover_enabled = excluded.pushover_enabled,
		pushover_user_key = excluded.pushover_user_key,
		pushover_api_token = excluded.pushover_api_token,
		webhook_enabled = excluded.webhook_enabled,
		webhook_url = excluded.webhook_url,
		updated_at = excluded.updated_at`

	_, err := r.db.ExecContext(ctx, query,
		n.ID, n.SMTPEnabled, n.SMTPHost, n.SMTPPort, n.SMTPUser, n.SMTPPassword, n.SMTPFromName, n.SMTPFromAddress,
		n.ResendEnabled, n.ResendAPIKey, n.SlackEnabled, n.SlackWebhookURL,
		n.DiscordEnabled, n.DiscordWebhookURL, n.DiscordPingEnabled,
		n.TelegramEnabled, n.TelegramBotToken, n.TelegramChatID,
		n.PushoverEnabled, n.PushoverUserKey, n.PushoverAPIToken,
		n.WebhookEnabled, n.WebhookURL, n.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to save notification integration: %w", err)
	}

	return nil
}

func (r *NotificationSQLiteRepository) GetProjectPref(ctx context.Context, projectID string) (*models.ProjectNotificationPref, error) {
	query := `SELECT project_id, email_enabled, slack_enabled, discord_enabled, telegram_enabled, pushover_enabled, webhook_enabled, COALESCE(events, 'deploy.success,deploy.failure,invite'), updated_at FROM project_notification_prefs WHERE project_id = ?`

	row := r.db.QueryRowContext(ctx, query, projectID)
	var pref models.ProjectNotificationPref
	err := row.Scan(&pref.ProjectID, &pref.EmailEnabled, &pref.SlackEnabled, &pref.DiscordEnabled, &pref.TelegramEnabled, &pref.PushoverEnabled, &pref.WebhookEnabled, &pref.Events, &pref.UpdatedAt)
	if err == sql.ErrNoRows {
		return &models.ProjectNotificationPref{
			ProjectID:       projectID,
			EmailEnabled:    true,
			SlackEnabled:    true,
			DiscordEnabled:  true,
			TelegramEnabled: true,
			PushoverEnabled: true,
			WebhookEnabled:  true,
			Events:          "deploy.success,deploy.failure,invite",
			UpdatedAt:       time.Now().UTC(),
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan project notification preferences: %w", err)
	}

	return &pref, nil
}

func (r *NotificationSQLiteRepository) SaveProjectPref(ctx context.Context, pref *models.ProjectNotificationPref) error {
	pref.UpdatedAt = time.Now().UTC()

	query := `INSERT INTO project_notification_prefs (
		project_id, email_enabled, slack_enabled, discord_enabled, telegram_enabled, pushover_enabled, webhook_enabled, events, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(project_id) DO UPDATE SET
		email_enabled = excluded.email_enabled,
		slack_enabled = excluded.slack_enabled,
		discord_enabled = excluded.discord_enabled,
		telegram_enabled = excluded.telegram_enabled,
		pushover_enabled = excluded.pushover_enabled,
		webhook_enabled = excluded.webhook_enabled,
		events = excluded.events,
		updated_at = excluded.updated_at`

	_, err := r.db.ExecContext(ctx, query,
		pref.ProjectID, pref.EmailEnabled, pref.SlackEnabled, pref.DiscordEnabled,
		pref.TelegramEnabled, pref.PushoverEnabled, pref.WebhookEnabled, pref.Events, pref.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to save project notification preferences: %w", err)
	}

	return nil
}
