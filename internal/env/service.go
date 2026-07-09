package env

import "context"

// Service provides business logic for project-level environment variables.
type Service struct {
	repo Repository
}

// NewService creates a new env Service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// GetVars returns all project-level environment variables.
func (s *Service) GetVars(ctx context.Context, projectID string) (map[string]string, error) {
	return s.repo.GetVars(ctx, projectID)
}

// SetVar upserts a single project-level environment variable.
func (s *Service) SetVar(ctx context.Context, projectID, key, value string) error {
	return s.repo.SetVar(ctx, projectID, key, value)
}
