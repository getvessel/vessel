package settings

type ServerSettings struct {
	ID                 string `json:"id"`
	CaddyWildcardIP    string `json:"caddyWildcardIp"`
	DiscordWebhookURL  string `json:"discordWebhookUrl,omitempty"`
	SlackWebhookURL    string `json:"slackWebhookUrl,omitempty"`
	TelegramBotToken   string `json:"telegramBotToken,omitempty"`
	TelegramChatID     string `json:"telegramChatId,omitempty"`
	SMTPHost           string `json:"smtpHost,omitempty"`
	SMTPPort           int    `json:"smtpPort,omitempty"`
	SMTPUser           string `json:"smtpUser,omitempty"`
	SMTPPassword       string `json:"smtpPassword,omitempty"`
	SMTPFromName       string `json:"smtpFromName,omitempty"`
	SMTPFromAddress    string `json:"smtpFromAddress,omitempty"`
	NotificationAlerts bool   `json:"notificationAlerts"`

	RegistrationEnabled         bool   `json:"registrationEnabled"`
	RegistrationDomainAllowlist string `json:"registrationDomainAllowlist,omitempty"`
	CustomDNSResolvers          string `json:"customDnsResolvers"`
	DNSValidationEnabled        bool   `json:"dnsValidationEnabled"`
	IPAllowlist                 string `json:"ipAllowlist"`
	MCPServerEnabled            bool   `json:"mcpServerEnabled"`

	UpdateCheckCron   string `json:"updateCheckCron"`
	AutoUpdateEnabled bool   `json:"autoUpdateEnabled"`
	CurrentVersion    string `json:"currentVersion"`
	LatestVersion     string `json:"latestVersion"`
	LastUpdateCheck   string `json:"lastUpdateCheck"`

	UpdatedAt string `json:"updatedAt"`
}
