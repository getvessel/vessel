package git

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"vessel.dev/vessel/internal/types"
)

type Store interface {
	GetProject(id string) (*types.ProjectConfig, error)
	GetAppService(id string) (*types.AppServiceConfig, error)
	CreateDeployment(dep *types.DeploymentRecord) error
	UpdateDeploymentStatus(id, status, buildLogs, containerID string) error
}

type GitService interface {
	CloneOrPullRepository(ctx context.Context, projectID, targetDir string, logWriter io.Writer) error
	CloneOrPullAppRepository(ctx context.Context, app *AppService, targetDir string, logWriter io.Writer) error
}

type Deployer interface {
	Deploy(ctx context.Context, project *types.ProjectConfig, sourceDir string, logWriter io.Writer) (string, error)
}

type ProxyManager interface {
	Reload(ctx context.Context) error
}

type WebhookHandler struct {
	store        Store
	gitService   GitService
	deployer     Deployer
	proxyManager ProxyManager
}

func NewWebhookHandler(store Store, gitService GitService, deployer Deployer, proxyManager ProxyManager) *WebhookHandler {
	return &WebhookHandler{
		store:        store,
		gitService:   gitService,
		deployer:     deployer,
		proxyManager: proxyManager,
	}
}

func (h *WebhookHandler) HandleProjectWebhook(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	if projectID == "" {
		writeError(w, http.StatusBadRequest, "missing projectId parameter")
		return
	}

	project, err := h.store.GetProject(projectID)
	if err != nil || project == nil {
		writeError(w, http.StatusNotFound, "project not found")
		return
	}

	writeJSON(w, http.StatusAccepted, map[string]string{
		"status":  "accepted",
		"message": fmt.Sprintf("triggering background build & deployment for %s", project.Name),
	})

	ctx := context.Background()
	sourceDir := filepath.Join("data", "builds", project.ID)
	if h.gitService != nil {
		if err := h.gitService.CloneOrPullRepository(ctx, project.ID, sourceDir, nil); err != nil {
			log.Printf("[GitWebhook] Git clone/pull failed for project %s (%s): %v", project.Name, project.ID, err)
			return
		}
	}
	if h.deployer != nil {
		containerID, err := h.deployer.Deploy(ctx, project, sourceDir, nil)
		if err != nil {
			log.Printf("[GitWebhook] Deployment failed for project %s (%s): %v", project.Name, project.ID, err)
			return
		}
		log.Printf("[GitWebhook] Successfully rolled out container %s for project %s (%s)", containerID, project.Name, project.ID)
	}
	if h.proxyManager != nil {
		_ = h.proxyManager.Reload(ctx)
	}
}

func (h *WebhookHandler) HandleServiceWebhook(w http.ResponseWriter, r *http.Request) {
	serviceID := r.PathValue("serviceId")
	if serviceID == "" {
		writeError(w, http.StatusBadRequest, "missing serviceId parameter")
		return
	}

	if h.store == nil {
		writeError(w, http.StatusInternalServerError, "store unavailable")
		return
	}

	appService, err := h.store.GetAppService(serviceID)
	if err != nil || appService == nil {
		writeError(w, http.StatusNotFound, "service not found")
		return
	}

	writeJSON(w, http.StatusAccepted, map[string]string{
		"status":  "accepted",
		"message": fmt.Sprintf("triggering background build & rollout for service %s", appService.Name),
	})

	ctx := context.Background()
	dep := &types.DeploymentRecord{
		ID:            uuid.NewString(),
		ServiceID:     appService.ID,
		EnvironmentID: appService.EnvironmentID,
		ProjectID:     appService.ProjectID,
		Status:        "BUILDING",
		Branch:        appService.Branch,
		Trigger:       "Git Webhook Push",
		BuildLogs:     fmt.Sprintf("Initiating automated build from %s branch %s...\n", appService.RepositoryURL, appService.Branch),
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}
	_ = h.store.CreateDeployment(dep)

	sourceDir := filepath.Join("data", "builds", "services", appService.ID)
	if h.gitService != nil && appService.RepositoryURL != "" {
		app := &AppService{
			ID:            appService.ID,
			ProjectID:     appService.ProjectID,
			EnvironmentID: appService.EnvironmentID,
			Name:          appService.Name,
			RepositoryURL: appService.RepositoryURL,
			Branch:        appService.Branch,
			ContainerID:   appService.ContainerID,
		}
		if err := h.gitService.CloneOrPullAppRepository(ctx, app, sourceDir, nil); err != nil {
			log.Printf("[ServiceGitWebhook] Git clone/pull failed for service %s (%s): %v", appService.Name, appService.ID, err)
			_ = h.store.UpdateDeploymentStatus(dep.ID, "FAILED", dep.BuildLogs+fmt.Sprintf("Error cloning repository: %v\n", err), "")
			return
		}
	}

	_ = h.store.UpdateDeploymentStatus(dep.ID, "ACTIVE", dep.BuildLogs+"Deployment rollout triggered via Webhook.\n", appService.ContainerID)
	if h.proxyManager != nil {
		_ = h.proxyManager.Reload(ctx)
	}
}
