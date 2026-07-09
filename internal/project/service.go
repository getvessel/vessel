package project

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"vessel.dev/vessel/internal/service"
	"vessel.dev/vessel/internal/utils"
)

type AppServiceRepository interface {
	CreateAppService(ctx context.Context, app *service.AppService) error
}

type Service struct {
	repo     Repository
	apps     AppServiceRepository
	domainFn func(name string) string
}

func NewService(repo Repository, apps AppServiceRepository) *Service {
	return &Service{
		repo:     repo,
		apps:     apps,
		domainFn: func(name string) string { return utils.GenerateSslipDomain(name, "") },
	}
}

func (s *Service) List(ctx context.Context) ([]ProjectConfig, error) {
	return s.repo.List(ctx)
}

func (s *Service) Get(ctx context.Context, id string) (*ProjectConfig, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) Create(ctx context.Context, req *CreateProjectRequest) (*ProjectConfig, error) {
	if req.Name == "" {
		req.Name = fmt.Sprintf("project-%s", uuid.NewString()[:8])
	}

	p := &ProjectConfig{
		ID:          req.ID,
		TeamID:      req.TeamID,
		Name:        req.Name,
		Description: req.Description,
	}
	if err := s.repo.Create(ctx, p); err != nil {
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

	app := &service.AppService{
		ProjectID:     p.ID,
		EnvironmentID: "env-prod",
		Name:          req.Name,
		RepositoryURL: repo,
		Branch:        branch,
		InternalPort:  port,
		Domain:        domain,
	}
	_ = s.apps.CreateAppService(ctx, app)

	return p, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
