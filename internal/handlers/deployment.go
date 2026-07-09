package handlers

import (
	"fmt"
	"net/http"
	"time"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/services"
)

type DeploymentHandler struct {
	deploymentService *services.DeploymentService
	appService        *services.AppService
}

func NewDeploymentHandler(ds *services.DeploymentService, as *services.AppService) *DeploymentHandler {
	return &DeploymentHandler{
		deploymentService: ds,
		appService:        as,
	}
}

func (h *DeploymentHandler) ListServiceDeployments(w http.ResponseWriter, r *http.Request) {
	serviceID := r.PathValue("serviceId")
	if serviceID == "" {
		WriteError(w, http.StatusBadRequest, "missing serviceId parameter")
		return
	}
	deps, err := h.deploymentService.ListByService(r.Context(), serviceID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, deps)
}

func (h *DeploymentHandler) Trigger(w http.ResponseWriter, r *http.Request) {
	serviceID := r.PathValue("serviceId")
	if serviceID == "" {
		WriteError(w, http.StatusBadRequest, "missing serviceId parameter")
		return
	}
	svc, err := h.appService.GetAppService(r.Context(), serviceID)
	if err != nil || svc == nil {
		WriteError(w, http.StatusNotFound, "service not found")
		return
	}
	dep := &models.Deployment{
		ServiceID:     serviceID,
		EnvironmentID: svc.EnvironmentID,
		ProjectID:     svc.ProjectID,
		Status:        "BUILDING",
		Branch:        svc.Branch,
		Trigger:       "Manual Deploy",
		BuildLogs:     "Initiating build...\n",
	}
	created, err := h.deploymentService.CreateDeployment(r.Context(), dep)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusAccepted)
	WriteJSON(w, http.StatusAccepted, created)
}

func (h *DeploymentHandler) Rollback(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing id parameter")
		return
	}
	targetDep, err := h.deploymentService.GetDeployment(r.Context(), id)
	if err != nil || targetDep == nil {
		WriteError(w, http.StatusNotFound, "deployment not found")
		return
	}
	newDep := &models.Deployment{
		ServiceID:     targetDep.ServiceID,
		EnvironmentID: targetDep.EnvironmentID,
		ProjectID:     targetDep.ProjectID,
		Status:        "BUILDING",
		CommitHash:    targetDep.CommitHash,
		CommitMessage: "Rollback to " + targetDep.ID,
		Branch:        targetDep.Branch,
		Trigger:       "Rollback",
		BuildLogs:     "Rolling back to deployment " + targetDep.ID + "...\n",
	}
	created, err := h.deploymentService.CreateDeployment(r.Context(), newDep)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusAccepted)
	WriteJSON(w, http.StatusAccepted, created)
}

func (h *DeploymentHandler) GetLogs(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing id parameter")
		return
	}
	dep, err := h.deploymentService.GetDeployment(r.Context(), id)
	if err != nil || dep == nil {
		WriteError(w, http.StatusNotFound, "deployment not found")
		return
	}
	WriteJSON(w, http.StatusOK, map[string]string{
		"id":        dep.ID,
		"buildLogs": dep.BuildLogs,
		"status":    dep.Status,
	})
}

func (h *DeploymentHandler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	now := time.Now().UTC()
	metrics := []map[string]any{
		{"timestamp": now.Add(-4 * time.Minute).Format(time.RFC3339), "cpuPercent": 1.2, "memoryMB": 64.5, "networkRx": 12.4, "networkTx": 8.1},
		{"timestamp": now.Add(-3 * time.Minute).Format(time.RFC3339), "cpuPercent": 2.1, "memoryMB": 66.0, "networkRx": 15.0, "networkTx": 10.2},
		{"timestamp": now.Add(-2 * time.Minute).Format(time.RFC3339), "cpuPercent": 1.8, "memoryMB": 65.2, "networkRx": 14.1, "networkTx": 9.4},
		{"timestamp": now.Add(-1 * time.Minute).Format(time.RFC3339), "cpuPercent": 3.4, "memoryMB": 68.1, "networkRx": 45.2, "networkTx": 22.0},
		{"timestamp": now.Format(time.RFC3339), "cpuPercent": 1.5, "memoryMB": 66.8, "networkRx": 18.0, "networkTx": 11.5},
	}
	WriteJSON(w, http.StatusOK, metrics)
}

func (h *DeploymentHandler) DeployProject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing project id parameter")
		return
	}
	sourceDir := fmt.Sprintf("data/builds/%s", id)
	containerID, err := h.deploymentService.DeployProject(r.Context(), id, sourceDir, nil)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, map[string]string{
		"status":       "deployed",
		"container_id": containerID,
	})
}
