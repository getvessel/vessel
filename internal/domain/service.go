package domain

import "context"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListByProject(ctx context.Context, projectID string) ([]Config, error) {
	return s.repo.ListByProject(ctx, projectID)
}

func (s *Service) ListAll(ctx context.Context) ([]Config, error) {
	return s.repo.ListAll(ctx)
}

func (s *Service) Create(ctx context.Context, d *Config) error {
	return s.repo.Create(ctx, d)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
