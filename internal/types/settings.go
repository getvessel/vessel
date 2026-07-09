package types

import "time"

type ServerSettings struct {
	ID                 string `json:"id"` // singleton "global"
	CaddyWildcardIP    string `json:"caddyWildcardIp"`
	DiscordWebhookURL  string `json:"discordWebhookUrl,omitempty"`
	SlackWebhookURL    string `json:"slackWebhookUrl,omitempty"`
	TelegramBotToken   string `json:"telegramBotToken,omitempty"`
	TelegramChatID     string `json:"telegramChatId,omitempty"`
	SMTPHost           string `json:"smtpHost,omitempty"`
	SMTPPort           int    `json:"smtpPort,omitempty"`
	SMTPUser           string `json:"smtpUser,omitempty"`
	SMTPPassword       string `json:"smtpPassword,omitempty"`
	NotificationAlerts bool   `json:"notificationAlerts"` // enable/disable global push alerts
	UpdatedAt          string `json:"updatedAt"`
}

type PersonalAccessToken struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Name      string    `json:"name"`
	TokenHash string    `json:"-"`
	Prefix    string    `json:"prefix"` // e.g. vsl_user_
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
}
