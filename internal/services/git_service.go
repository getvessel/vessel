package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/repositories"
)

type GitService struct {
	repo repositories.GitRepository
}

func NewGitService(r repositories.GitRepository) *GitService {
	return &GitService{repo: r}
}

func (s *GitService) SaveProvider(ctx context.Context, gp *models.GitProviderConfig) error {
	if gp == nil || gp.UserID == "" || gp.Provider == "" {
		return errors.New("valid git provider config with userId and provider required")
	}
	if gp.ID == "" {
		gp.ID = uuid.New().String()
	}
	gp.UpdatedAt = time.Now()
	if gp.CreatedAt.IsZero() {
		gp.CreatedAt = gp.UpdatedAt
	}
	return s.repo.SaveProvider(ctx, gp)
}

func (s *GitService) GetProvider(ctx context.Context, userID, provider string) (*models.GitProviderConfig, error) {
	if userID == "" || provider == "" {
		return nil, errors.New("userId and provider required")
	}
	return s.repo.GetProvider(ctx, userID, provider)
}

func (s *GitService) GetAnyProviderByType(ctx context.Context, provider string) (*models.GitProviderConfig, error) {
	if provider == "" {
		return nil, errors.New("provider required")
	}
	return s.repo.GetAnyProviderByType(ctx, provider)
}

func (s *GitService) ListProvidersByUser(ctx context.Context, userID string) ([]*models.GitProviderConfig, error) {
	if userID == "" {
		return nil, errors.New("userId required")
	}
	return s.repo.ListProvidersByUser(ctx, userID)
}

func (s *GitService) DeleteProvider(ctx context.Context, userID, provider string) error {
	if userID == "" || provider == "" {
		return errors.New("userId and provider required")
	}
	return s.repo.DeleteProvider(ctx, userID, provider)
}
