package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"vessel.dev/vessel/internal/engine"
	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/repositories"
)

type PRPreviewService struct {
	repo       repositories.PRPreviewRepository
	appService *AppService
	gitService *GitService
	deployer   *engine.Deployer
}

func NewPRPreviewService(
	repo repositories.PRPreviewRepository,
	appService *AppService,
	gitService *GitService,
	deployer *engine.Deployer,
) *PRPreviewService {
	return &PRPreviewService{
		repo:       repo,
		appService: appService,
		gitService: gitService,
		deployer:   deployer,
	}
}

func (s *PRPreviewService) DeployPRPreview(ctx context.Context, appID string, prNumber int, commitHash, branch string) (*models.PRPreview, error) {
	app, err := s.appService.GetAppService(ctx, appID)
	if err != nil || app == nil {
		return nil, errors.New("app service not found")
	}

	previewDomain := fmt.Sprintf("pr-%d.%s", prNumber, app.Domain)
	if app.Domain == "" {
		previewDomain = fmt.Sprintf("pr-%d.%s.sslip.io", prNumber, app.Name) // Fallback
	}

	preview := &models.PRPreview{
		ID:            uuid.NewString(),
		ServiceID:     app.ID,
		ProjectID:     app.ProjectID,
		PRNumber:      prNumber,
		Branch:        branch,
		CommitHash:    commitHash,
		Status:        "BUILDING",
		PreviewDomain: previewDomain,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, preview); err != nil {
		return nil, err
	}

	go func() {
		bgCtx := context.Background()
		sourceDir := filepath.Join("data", "builds", "pr-previews", preview.ID)

		// 1. Clone the specific PR branch
		clonedApp := *app
		clonedApp.Branch = branch
		if err := s.gitService.CloneOrPullAppRepository(bgCtx, &clonedApp, sourceDir, nil); err != nil {
			log.Printf("[PRPreview] failed to clone PR branch %s: %v", branch, err)
			preview.Status = "FAILED"
			_ = s.repo.Update(bgCtx, preview)
			return
		}

		// 2. Deploy using engine deployer, mapping to the preview domain
		// For ephemeral previews, we pass the previewDomain.
		// Note: The Deploy method should ideally accept domain overrides.
		// For now we'll update the clonedApp domain to the previewDomain
		clonedApp.Domain = previewDomain
		clonedApp.Name = fmt.Sprintf("%s-pr-%d", app.Name, prNumber)

		containerID, deployErr := s.deployer.DeployAppService(bgCtx, &clonedApp, sourceDir, nil)
		if deployErr != nil {
			log.Printf("[PRPreview] failed to deploy: %v", deployErr)
			preview.Status = "FAILED"
			_ = s.repo.Update(bgCtx, preview)
			return
		}

		preview.ContainerID = containerID
		preview.Status = "READY"
		_ = s.repo.Update(bgCtx, preview)

		// 3. Post Commit Status Check Update to GitHub/GitLab
		// (Integration goes here - currently mocked)
		log.Printf("[PRPreview] Posted status 'success' for commit %s. Preview available at %s", commitHash, previewDomain)
	}()

	return preview, nil
}

func (s *PRPreviewService) DestroyPRPreview(ctx context.Context, appID string, prNumber int) error {
	previews, err := s.repo.GetByAppAndPR(ctx, appID, prNumber)
	if err != nil {
		return err
	}

	for _, p := range previews {
		if p.ContainerID != "" {
			_ = s.deployer.Stop(ctx, p.ContainerID)
			_ = s.deployer.Remove(ctx, p.ContainerID)
		}
		_ = s.repo.Delete(ctx, p.ID)
	}

	return nil
}
