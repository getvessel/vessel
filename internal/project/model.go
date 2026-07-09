package project

import (
	"time"

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

// EnvironmentConfig represents a named deployment environment (e.g. production, staging) within a project.
type EnvironmentConfig struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"projectId"`
	Name      string    `json:"name"`
	IsDefault bool      `json:"isDefault"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// DomainConfig represents a custom domain routing rule attached to a project.
type DomainConfig struct {
	ID            string    `json:"id"`
	ProjectID     string    `json:"projectId"`
	DomainName    string    `json:"domainName"`
	RedirectTo    string    `json:"redirectTo,omitempty"`
	SSLCertStatus string    `json:"sslCertStatus"`
	PathPrefix    string    `json:"pathPrefix"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// ProjectCanvasSummary is an aggregated view of a project including resource counts and service status.
// It embeds ProjectConfig directly for the JSON shape used by the dashboard.
type ProjectCanvasSummary struct {
	ProjectConfig
	EnvironmentsCount  int                `json:"environmentsCount"`
	AppsCount          int                `json:"appsCount"`
	DatabasesCount     int                `json:"databasesCount"`
	StorageCount       int                `json:"storageCount"`
	OnlineServices     int                `json:"onlineServices"`
	TotalServices      int                `json:"totalServices"`
	ServiceIcons       []string           `json:"serviceIcons"`
	DefaultEnvironment *EnvironmentConfig `json:"defaultEnvironment,omitempty"`
}

// EnvironmentCanvas holds the complete set of services running within a single environment.
// AppServiceConfig, DatabaseConfig and StorageConfig are still defined in internal/types until they are migrated.
type EnvironmentCanvas struct {
	Environment *EnvironmentConfig        `json:"environment"`
	Apps        []*types.AppServiceConfig `json:"apps"`
	Databases   []*types.DatabaseConfig   `json:"databases"`
	Storage     []*types.StorageConfig    `json:"storage"`
}
