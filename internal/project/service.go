package project

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"vessel.dev/vessel/internal/types"
	"vessel.dev/vessel/internal/utils"
)

// AppServiceRepository is the minimal surface project.Service needs from the app-service domain.
type AppServiceRepository interface {
	CreateAppService(ctx context.Context, app *types.AppServiceConfig) error
}

// Service implements the project domain business logic.
type Service struct {
	repo     Repository
	apps     AppServiceRepository
	domainFn func(name string) string
}

// NewService creates a new project Service.
func NewService(repo Repository, apps AppServiceRepository) *Service {
	return &Service{
		repo:     repo,
		apps:     apps,
		domainFn: func(name string) string { return utils.GenerateSslipDomain(name, "") },
	}
}

// ListProjects returns all projects.
func (s *Service) ListProjects(ctx context.Context) ([]ProjectConfig, error) {
	return s.repo.ListProjects(ctx)
}

// GetProject returns a single project by ID.
func (s *Service) GetProject(ctx context.Context, id string) (*ProjectConfig, error) {
	return s.repo.GetProject(ctx, id)
}

// CreateProject creates a project and its default application service.
func (s *Service) CreateProject(ctx context.Context, req *CreateProjectRequest) (*ProjectConfig, error) {
	if req.Name == "" {
		req.Name = fmt.Sprintf("project-%s", uuid.NewString()[:8])
	}

	p := &ProjectConfig{
		ID:          req.ID,
		TeamID:      req.TeamID,
		Name:        req.Name,
		Description: req.Description,
	}
	if err := s.repo.CreateProject(ctx, p); err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	port := req.InternalPort
	if port <= 0 {
		port = req.InternalPortSnake
	}
	if port <= 0 {
		port = 3000
	}

	repo := req.RepositoryURL
	if repo == "" {
		repo = req.RepositoryURLSnake
	}

	domain := req.Domain
	if domain == "" {
		domain = s.domainFn(req.Name)
	}

	branch := req.Branch
	if branch == "" {
		branch = "main"
	}

	envs, _ := s.repo.ListEnvironments(ctx, p.ID)
	envID := "env-prod"
	if len(envs) > 0 {
		envID = envs[0].ID
	}

	app := &types.AppServiceConfig{
		ProjectID:     p.ID,
		EnvironmentID: envID,
		Name:          req.Name,
		RepositoryURL: repo,
		Branch:        branch,
		InternalPort:  port,
		Domain:        domain,
	}
	_ = s.apps.CreateAppService(ctx, app)

	return p, nil
}

// DeleteProject removes a project by ID.
func (s *Service) DeleteProject(ctx context.Context, id string) error {
	return s.repo.DeleteProject(ctx, id)
}

// ListEnvironments returns all environments for a project.
func (s *Service) ListEnvironments(ctx context.Context, projectID string) ([]EnvironmentConfig, error) {
	return s.repo.ListEnvironments(ctx, projectID)
}

// CreateEnvironment creates a new environment.
func (s *Service) CreateEnvironment(ctx context.Context, env *EnvironmentConfig) error {
	return s.repo.CreateEnvironment(ctx, env)
}

// DeleteEnvironment removes an environment by ID.
func (s *Service) DeleteEnvironment(ctx context.Context, id string) error {
	return s.repo.DeleteEnvironment(ctx, id)
}

// ListDomains returns all custom domains for a project.
func (s *Service) ListDomains(ctx context.Context, projectID string) ([]DomainConfig, error) {
	return s.repo.ListDomains(ctx, projectID)
}

// AddDomain adds a custom domain to a project.
func (s *Service) AddDomain(ctx context.Context, d *DomainConfig) error {
	return s.repo.AddDomain(ctx, d)
}

// DeleteDomain removes a custom domain by ID.
func (s *Service) DeleteDomain(ctx context.Context, id string) error {
	return s.repo.DeleteDomain(ctx, id)
}

// GetEnvVars returns all project-level environment variables.
func (s *Service) GetEnvVars(ctx context.Context, projectID string) (map[string]string, error) {
	return s.repo.GetEnvVars(ctx, projectID)
}

// SetEnvVar upserts a single project environment variable.
func (s *Service) SetEnvVar(ctx context.Context, projectID, key, value string) error {
	return s.repo.SetEnvVar(ctx, projectID, key, value)
}

// ListProjectCanvasSummaries returns aggregated dashboard summaries for all projects.
func (s *Service) ListProjectCanvasSummaries(ctx context.Context) ([]ProjectCanvasSummary, error) {
	return s.repo.ListProjectCanvasSummaries(ctx)
}

// GetProjectCanvasSummary returns an aggregated dashboard summary for a single project.
func (s *Service) GetProjectCanvasSummary(ctx context.Context, id string) (*ProjectCanvasSummary, error) {
	return s.repo.GetProjectCanvasSummary(ctx, id)
}

// GetEnvironmentCanvas returns all resources running in an environment.
func (s *Service) GetEnvironmentCanvas(ctx context.Context, id string) (*EnvironmentCanvas, error) {
	return s.repo.GetEnvironmentCanvas(ctx, id)
}
