package project

import (
	"time"

	"vessel.dev/vessel/internal/domain"
	"vessel.dev/vessel/internal/environment"
	"vessel.dev/vessel/internal/types"
)

// ProjectConfig holds the core metadata for a Vessel project.
type ProjectConfig struct {
	ID          string    `json:"id"`
	WorkspaceID string    `json:"workspaceId,omitempty"`
	TeamID      string    `json:"teamId,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// CanvasSummary is an aggregated view of a project including resource counts and service status.
type CanvasSummary struct {
	ProjectConfig
	EnvironmentsCount  int                 `json:"environmentsCount"`
	AppsCount          int                 `json:"appsCount"`
	DatabasesCount     int                 `json:"databasesCount"`
	StorageCount       int                 `json:"storageCount"`
	OnlineServices     int                 `json:"onlineServices"`
	TotalServices      int                 `json:"totalServices"`
	ServiceIcons       []string            `json:"serviceIcons"`
	DefaultEnvironment *environment.Config `json:"defaultEnvironment,omitempty"`
}

// EnvironmentCanvas holds the complete set of services running within a single environment.
type EnvironmentCanvas struct {
	Environment *environment.Config       `json:"environment"`
	Apps        []*types.AppServiceConfig `json:"apps"`
	Databases   []*types.DatabaseConfig   `json:"databases"`
	Storage     []*types.StorageConfig    `json:"storage"`
}

// DomainConfig is an alias for domain.Config kept for internal convenience.
type DomainConfig = domain.Config
