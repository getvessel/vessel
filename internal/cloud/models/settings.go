package models

import "time"

type TeamAISettings struct {
	ID        string    `json:"id"`
	TeamID    string    `json:"teamId"`
	Provider  string    `json:"provider"`
	APIKey    string    `json:"apiKey,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type TeamEmailSettings struct {
	ID              string    `json:"id"`
	TeamID          string    `json:"teamId"`
	SMTPHost        string    `json:"smtpHost,omitempty"`
	SMTPPort        int       `json:"smtpPort,omitempty"`
	SMTPUser        string    `json:"smtpUser,omitempty"`
	SMTPPassword    string    `json:"smtpPassword,omitempty"`
	SMTPFromName    string    `json:"smtpFromName,omitempty"`
	SMTPFromAddress string    `json:"smtpFromAddress,omitempty"`
	ResendAPIKey    string    `json:"resendApiKey,omitempty"`
	UseResend       bool      `json:"useResend"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}
