package project_settings

import "time"

type Webhook struct {
	ID                    string    `json:"id"`
	ProjectID             string    `json:"projectId"`
	URL                   string    `json:"url"`
	EventTypes            []string  `json:"eventTypes"`
	IncludePREnvironments bool      `json:"includePrEnvironments"`
	CreatedAt             time.Time `json:"createdAt"`
	UpdatedAt             time.Time `json:"updatedAt"`
}

type Token struct {
	ID            string    `json:"id"`
	ProjectID     string    `json:"projectId"`
	EnvironmentID string    `json:"environmentId"`
	Name          string    `json:"name"`
	TokenPrefix   string    `json:"tokenPrefix"`
	CreatedAt     time.Time `json:"createdAt"`
}

type ProjectMember struct {
	ID         string    `json:"id"`
	ProjectID  string    `json:"projectId"`
	UserID     string    `json:"userId,omitempty"`
	Email      string    `json:"email"`
	Permission string    `json:"permission"`
	Status     string    `json:"status"`
	InvitedAt  time.Time `json:"invitedAt"`
	AcceptedAt time.Time `json:"acceptedAt,omitempty"`
}
