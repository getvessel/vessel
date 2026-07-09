package handlers

import (
	"encoding/json"
	"net/http"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/services"
)

type EnvironmentHandler struct {
	envService *services.EnvironmentService
}

func NewEnvironmentHandler(s *services.EnvironmentService) *EnvironmentHandler {
	return &EnvironmentHandler{envService: s}
}

func (h *EnvironmentHandler) ListByProject(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	if projectID == "" {
		WriteError(w, http.StatusBadRequest, "missing project id parameter")
		return
	}
	envs, err := h.envService.ListByProject(r.Context(), projectID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, envs)
}

func (h *EnvironmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	if projectID == "" {
		WriteError(w, http.StatusBadRequest, "missing project id parameter")
		return
	}
	var env models.EnvironmentConfig
	if err := json.NewDecoder(r.Body).Decode(&env); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid environment payload")
		return
	}
	env.ProjectID = projectID
	if env.Name == "" {
		WriteError(w, http.StatusBadRequest, "environment name is required")
		return
	}
	created, err := h.envService.CreateEnvironment(r.Context(), &env)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, created)
}

func (h *EnvironmentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing id parameter")
		return
	}
	if err := h.envService.DeleteEnvironment(r.Context(), id); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type DomainHandler struct {
	envService *services.EnvironmentService
}

func NewDomainHandler(s *services.EnvironmentService) *DomainHandler {
	return &DomainHandler{envService: s}
}

func (h *DomainHandler) ListByProject(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	if projectID == "" {
		WriteError(w, http.StatusBadRequest, "missing project id parameter")
		return
	}
	domains, err := h.envService.ListDomainsByProject(r.Context(), projectID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, domains)
}

func (h *DomainHandler) Create(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	if projectID == "" {
		WriteError(w, http.StatusBadRequest, "missing project id parameter")
		return
	}
	var d models.DomainConfig
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid domain payload")
		return
	}
	d.ProjectID = projectID
	if d.DomainName == "" {
		WriteError(w, http.StatusBadRequest, "domainName is required")
		return
	}
	created, err := h.envService.CreateDomain(r.Context(), &d)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, created)
}

func (h *DomainHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "missing id parameter")
		return
	}
	if err := h.envService.DeleteDomain(r.Context(), id); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type ProjectEnvHandler struct {
	envService *services.EnvironmentService
}

func NewProjectEnvHandler(s *services.EnvironmentService) *ProjectEnvHandler {
	return &ProjectEnvHandler{envService: s}
}

func (h *ProjectEnvHandler) GetVars(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	if projectID == "" {
		WriteError(w, http.StatusBadRequest, "missing project id parameter")
		return
	}
	vars, err := h.envService.GetVars(r.Context(), projectID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if vars == nil {
		vars = map[string]string{}
	}
	WriteJSON(w, http.StatusOK, vars)
}

func (h *ProjectEnvHandler) SetVars(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	if projectID == "" {
		WriteError(w, http.StatusBadRequest, "missing project id parameter")
		return
	}
	var vars map[string]string
	if err := json.NewDecoder(r.Body).Decode(&vars); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid env variables payload")
		return
	}
	for k, v := range vars {
		if err := h.envService.SetVar(r.Context(), projectID, k, v); err != nil {
			WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	WriteJSON(w, http.StatusOK, vars)
}
