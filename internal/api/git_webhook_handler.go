package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"vessel.dev/vessel/internal/git"
	"vessel.dev/vessel/internal/types"
)

func (s *Server) handleGitWebhook(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	if projectID == "" {
		writeError(w, http.StatusBadRequest, "missing projectId parameter")
		return
	}

	project, err := s.store.GetProject(projectID)
	if err != nil || project == nil {
		writeError(w, http.StatusNotFound, "project not found")
		return
	}

	writeJSON(w, http.StatusAccepted, map[string]string{
		"status":  "accepted",
		"message": fmt.Sprintf("triggering background build & deployment for %s", project.Name),
	})

	go func() {
		ctx := context.Background()
		sourceDir := filepath.Join("data", "builds", project.ID)
		if s.gitService != nil {
			if err := s.gitService.CloneOrPullRepository(ctx, project.ID, sourceDir, nil); err != nil {
				log.Printf("[GitWebhook] Git clone/pull failed for project %s (%s): %v", project.Name, project.ID, err)
				return
			}
		}
		if s.deployer != nil {
			containerID, err := s.deployer.Deploy(ctx, project, sourceDir, nil)
			if err != nil {
				log.Printf("[GitWebhook] Deployment failed for project %s (%s): %v", project.Name, project.ID, err)
				return
			}
			log.Printf("[GitWebhook] Successfully rolled out container %s for project %s (%s)", containerID, project.Name, project.ID)
		}
		if s.proxyManager != nil {
			_ = s.proxyManager.Reload(ctx)
		}
	}()
}

func (s *Server) handleServiceGitWebhook(w http.ResponseWriter, r *http.Request) {
	serviceID := r.PathValue("serviceId")
	if serviceID == "" {
		writeError(w, http.StatusBadRequest, "missing serviceId parameter")
		return
	}

	if s.store == nil {
		writeError(w, http.StatusInternalServerError, "store unavailable")
		return
	}

	appService, err := s.store.GetAppService(serviceID)
	if err != nil || appService == nil {
		writeError(w, http.StatusNotFound, "service not found")
		return
	}

	writeJSON(w, http.StatusAccepted, map[string]string{
		"status":  "accepted",
		"message": fmt.Sprintf("triggering background build & rollout for service %s", appService.Name),
	})

	go func() {
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
		_ = s.store.CreateDeployment(dep)

		sourceDir := filepath.Join("data", "builds", "services", appService.ID)
		if s.gitService != nil && appService.RepositoryURL != "" {
			app := &git.AppService{
				ID:            appService.ID,
				ProjectID:     appService.ProjectID,
				EnvironmentID: appService.EnvironmentID,
				Name:          appService.Name,
				RepositoryURL: appService.RepositoryURL,
				Branch:        appService.Branch,
				ContainerID:   appService.ContainerID,
			}
			if err := s.gitService.CloneOrPullAppRepository(ctx, app, sourceDir, nil); err != nil {
				log.Printf("[ServiceGitWebhook] Git clone/pull failed for service %s (%s): %v", appService.Name, appService.ID, err)
				_ = s.store.UpdateDeploymentStatus(dep.ID, "FAILED", dep.BuildLogs+fmt.Sprintf("Error cloning repository: %v\n", err), "")
				return
			}
		}

		_ = s.store.UpdateDeploymentStatus(dep.ID, "ACTIVE", dep.BuildLogs+"Deployment rollout triggered via Webhook.\n", appService.ContainerID)
		if s.proxyManager != nil {
			_ = s.proxyManager.Reload(ctx)
		}
	}()
}
