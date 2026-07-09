package env

import "context"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetVars(ctx context.Context, projectID string) (map[string]string, error) {
	return s.repo.GetVars(ctx, projectID)
}

func (s *Service) SetVar(ctx context.Context, projectID, key, value string) error {
	return s.repo.SetVar(ctx, projectID, key, value)
}
