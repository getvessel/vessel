package project

import (
	"time"

	"vessel.dev/vessel/internal/domain"
)

type ProjectConfig struct {
	ID          string    `json:"id"`
	WorkspaceID string    `json:"workspaceId,omitempty"`
	TeamID      string    `json:"teamId,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type DomainConfig = domain.Config
