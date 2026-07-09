package project

import "context"

// Repository defines all data operations for the project domain.
type Repository interface {
	// Projects
	ListProjects(ctx context.Context) ([]ProjectConfig, error)
	GetProject(ctx context.Context, id string) (*ProjectConfig, error)
	CreateProject(ctx context.Context, p *ProjectConfig) error
	DeleteProject(ctx context.Context, id string) error

	// Environments
	ListEnvironments(ctx context.Context, projectID string) ([]EnvironmentConfig, error)
	CreateEnvironment(ctx context.Context, env *EnvironmentConfig) error
	DeleteEnvironment(ctx context.Context, id string) error

	// Domains
	ListDomains(ctx context.Context, projectID string) ([]DomainConfig, error)
	AddDomain(ctx context.Context, d *DomainConfig) error
	DeleteDomain(ctx context.Context, id string) error

	// Env vars
	GetEnvVars(ctx context.Context, projectID string) (map[string]string, error)
	SetEnvVar(ctx context.Context, projectID, key, value string) error

	// Canvas
	ListProjectCanvasSummaries(ctx context.Context) ([]ProjectCanvasSummary, error)
	GetProjectCanvasSummary(ctx context.Context, id string) (*ProjectCanvasSummary, error)
	GetEnvironmentCanvas(ctx context.Context, id string) (*EnvironmentCanvas, error)
}
