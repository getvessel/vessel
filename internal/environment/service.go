package environment

import "context"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Get(ctx context.Context, id string) (*Config, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) ListByProject(ctx context.Context, projectID string) ([]Config, error) {
	return s.repo.ListByProject(ctx, projectID)
}

func (s *Service) Create(ctx context.Context, env *Config) error {
	return s.repo.Create(ctx, env)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
