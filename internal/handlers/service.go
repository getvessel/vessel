package handlers

import (
	"encoding/json"
	"net/http"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/services"
)

type AppHandler struct {
	appService *services.AppService
}

func NewAppHandler(s *services.AppService) *AppHandler {
	return &AppHandler{appService: s}
}

func (h *AppHandler) Create(w http.ResponseWriter, r *http.Request) {
	envID := r.PathValue("id")
	var req models.AppService
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	if req.Name == "" {
		WriteError(w, http.StatusBadRequest, "app service name is required")
		return
	}
	req.EnvironmentID = envID
	if req.InternalPort == 0 {
		req.InternalPort = 3000
	}
	created, err := h.appService.CreateAppService(r.Context(), &req)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, created)
}

func (h *AppHandler) ListByEnvironment(w http.ResponseWriter, r *http.Request) {
	envID := r.PathValue("id")
	apps, err := h.appService.ListByEnvironment(r.Context(), envID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, apps)
}

func (h *AppHandler) ListByProject(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	apps, err := h.appService.ListByProject(r.Context(), projectID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, apps)
}

func (h *AppHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	svc, err := h.appService.GetAppService(r.Context(), id)
	if err != nil || svc == nil {
		WriteError(w, http.StatusNotFound, "app service not found")
		return
	}
	WriteJSON(w, http.StatusOK, svc)
}

func (h *AppHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	existing, err := h.appService.GetAppService(r.Context(), id)
	if err != nil || existing == nil {
		WriteError(w, http.StatusNotFound, "app service not found")
		return
	}
	var req models.AppService
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	existing.Name = req.Name
	existing.RepositoryURL = req.RepositoryURL
	existing.Branch = req.Branch
	existing.InternalPort = req.InternalPort
	existing.Domain = req.Domain
	existing.ContainerID = req.ContainerID
	existing.Status = req.Status
	if err := h.appService.UpdateAppService(r.Context(), existing); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, existing)
}

func (h *AppHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.appService.DeleteAppService(r.Context(), id); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type ServiceVarHandler struct {
	appService *services.AppService
}

func NewServiceVarHandler(s *services.AppService) *ServiceVarHandler {
	return &ServiceVarHandler{appService: s}
}

func (h *ServiceVarHandler) List(w http.ResponseWriter, r *http.Request) {
	serviceID := r.PathValue("serviceId")
	list, err := h.appService.ListVariablesByService(r.Context(), serviceID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, list)
}

func (h *ServiceVarHandler) Create(w http.ResponseWriter, r *http.Request) {
	serviceID := r.PathValue("serviceId")
	var req models.Variable
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	svc, err := h.appService.GetAppService(r.Context(), serviceID)
	if err != nil || svc == nil {
		WriteError(w, http.StatusNotFound, "service not found")
		return
	}
	req.ServiceID = serviceID
	req.ProjectID = svc.ProjectID
	req.EnvironmentID = svc.EnvironmentID
	created, err := h.appService.CreateVariable(r.Context(), &req)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, created)
}

func (h *ServiceVarHandler) Update(w http.ResponseWriter, r *http.Request) {
	serviceID := r.PathValue("serviceId")
	id := r.PathValue("id")
	var req models.Variable
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.ID = id
	req.ServiceID = serviceID
	if err := h.appService.UpdateVariable(r.Context(), &req); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, req)
}

func (h *ServiceVarHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.appService.DeleteVariable(r.Context(), id); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
