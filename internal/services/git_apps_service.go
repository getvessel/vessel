package services

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/repositories"
)

type GitAppsService struct {
	repo repositories.GitAppRepository
}

func NewGitAppsService(repo repositories.GitAppRepository) *GitAppsService {
	return &GitAppsService{repo: repo}
}

// ---- GitHub Apps ----

func (s *GitAppsService) ListGithubApps(ctx context.Context, teamID string) ([]models.GithubApp, error) {
	if teamID == "" {
		return nil, errors.New("team ID is required")
	}
	return s.repo.ListGithubApps(ctx, teamID)
}

func (s *GitAppsService) GetGithubApp(ctx context.Context, id string) (*models.GithubApp, error) {
	if id == "" {
		return nil, errors.New("app ID is required")
	}
	return s.repo.GetGithubApp(ctx, id)
}

func (s *GitAppsService) SaveGithubApp(ctx context.Context, app *models.GithubApp) error {
	if app == nil {
		return errors.New("app config is required")
	}
	if app.TeamID == "" {
		return errors.New("team ID is required")
	}
	if app.ID == "" {
		app.ID = uuid.NewString()
	}
	return s.repo.SaveGithubApp(ctx, app)
}

func (s *GitAppsService) DeleteGithubApp(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("app ID is required")
	}
	return s.repo.DeleteGithubApp(ctx, id)
}

// ---- GitLab Apps ----

func (s *GitAppsService) ListGitlabApps(ctx context.Context, teamID string) ([]models.GitlabApp, error) {
	if teamID == "" {
		return nil, errors.New("team ID is required")
	}
	return s.repo.ListGitlabApps(ctx, teamID)
}

func (s *GitAppsService) GetGitlabApp(ctx context.Context, id string) (*models.GitlabApp, error) {
	if id == "" {
		return nil, errors.New("app ID is required")
	}
	return s.repo.GetGitlabApp(ctx, id)
}

func (s *GitAppsService) SaveGitlabApp(ctx context.Context, app *models.GitlabApp) error {
	if app == nil {
		return errors.New("app config is required")
	}
	if app.TeamID == "" {
		return errors.New("team ID is required")
	}
	if app.ID == "" {
		app.ID = uuid.NewString()
	}
	return s.repo.SaveGitlabApp(ctx, app)
}

func (s *GitAppsService) DeleteGitlabApp(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("app ID is required")
	}
	return s.repo.DeleteGitlabApp(ctx, id)
}

// ---- Bitbucket Apps ----

func (s *GitAppsService) ListBitbucketApps(ctx context.Context, teamID string) ([]models.BitbucketApp, error) {
	if teamID == "" {
		return nil, errors.New("team ID is required")
	}
	return s.repo.ListBitbucketApps(ctx, teamID)
}

func (s *GitAppsService) GetBitbucketApp(ctx context.Context, id string) (*models.BitbucketApp, error) {
	if id == "" {
		return nil, errors.New("app ID is required")
	}
	return s.repo.GetBitbucketApp(ctx, id)
}

func (s *GitAppsService) SaveBitbucketApp(ctx context.Context, app *models.BitbucketApp) error {
	if app == nil {
		return errors.New("app config is required")
	}
	if app.TeamID == "" {
		return errors.New("team ID is required")
	}
	if app.ID == "" {
		app.ID = uuid.NewString()
	}
	return s.repo.SaveBitbucketApp(ctx, app)
}

func (s *GitAppsService) DeleteBitbucketApp(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("app ID is required")
	}
	return s.repo.DeleteBitbucketApp(ctx, id)
}
