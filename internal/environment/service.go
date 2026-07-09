package environment

import "context"

// Service provides business logic for deployment environments.
type Service struct {
	repo Repository
}

// NewService creates a new environment Service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// Get returns a single environment by ID.
func (s *Service) Get(ctx context.Context, id string) (*Config, error) {
	return s.repo.Get(ctx, id)
}

// ListByProject returns all environments for a project.
func (s *Service) ListByProject(ctx context.Context, projectID string) ([]Config, error) {
	return s.repo.ListByProject(ctx, projectID)
}

// Create creates a new environment.
func (s *Service) Create(ctx context.Context, env *Config) error {
	return s.repo.Create(ctx, env)
}

// Delete removes an environment by ID.
func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
