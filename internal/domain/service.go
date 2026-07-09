package domain

import "context"

// Service provides business logic for custom domains.
type Service struct {
	repo Repository
}

// NewService creates a new domain Service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// ListByProject returns all domains for a project.
func (s *Service) ListByProject(ctx context.Context, projectID string) ([]Config, error) {
	return s.repo.ListByProject(ctx, projectID)
}

// ListAll returns every custom domain across all projects.
func (s *Service) ListAll(ctx context.Context) ([]Config, error) {
	return s.repo.ListAll(ctx)
}

// Create adds a custom domain to a project.
func (s *Service) Create(ctx context.Context, d *Config) error {
	return s.repo.Create(ctx, d)
}

// Delete removes a custom domain by ID.
func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
