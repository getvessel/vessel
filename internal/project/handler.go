package project

import (
	"context"
	"encoding/json"
	"net/http"
)

// ProxyReloader is the minimal proxy surface used by project handlers.
type ProxyReloader interface {
	Reload(ctx context.Context) error
}

// Handler serves HTTP requests for the project domain.
type Handler struct {
	service     *Service
	proxy       ProxyReloader
	extractUser func(r *http.Request) string
}

// NewHandler creates a new project Handler.
func NewHandler(service *Service, proxy ProxyReloader, extractUser func(r *http.Request) string) *Handler {
	return &Handler{
		service:     service,
		proxy:       proxy,
		extractUser: extractUser,
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// ListProjects handles GET /api/projects.
func (h *Handler) ListProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.service.ListProjects(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, projects)
}

// CreateProject handles POST /api/projects.
func (h *Handler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var req CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid project configuration payload")
		return
	}

	p, err := h.service.CreateProject(r.Context(), &req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if h.proxy != nil {
		_ = h.proxy.Reload(r.Context())
	}
	writeJSON(w, http.StatusCreated, p)
}

// GetProject handles GET /api/projects/{id}.
func (h *Handler) GetProject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing project id parameter")
		return
	}

	p, err := h.service.GetProject(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "project not found")
		return
	}
	writeJSON(w, http.StatusOK, p)
}

// DeleteProject handles DELETE /api/projects/{id}.
func (h *Handler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing project id parameter")
		return
	}

	if err := h.service.DeleteProject(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if h.proxy != nil {
		_ = h.proxy.Reload(r.Context())
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

// ListEnvironments handles GET /api/projects/{id}/environments.
func (h *Handler) ListEnvironments(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	envs, err := h.service.ListEnvironments(r.Context(), projectID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, envs)
}

// CreateEnvironment handles POST /api/projects/{id}/environments.
func (h *Handler) CreateEnvironment(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	var env EnvironmentConfig
	if err := json.NewDecoder(r.Body).Decode(&env); err != nil {
		writeError(w, http.StatusBadRequest, "invalid environment payload")
		return
	}
	env.ProjectID = projectID
	if env.Name == "" {
		writeError(w, http.StatusBadRequest, "environment name is required")
		return
	}

	if err := h.service.CreateEnvironment(r.Context(), &env); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, env)
}

// DeleteEnvironment handles DELETE /api/environments/{id}.
func (h *Handler) DeleteEnvironment(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.service.DeleteEnvironment(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ListDomains handles GET /api/projects/{id}/domains.
func (h *Handler) ListDomains(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	domains, err := h.service.ListDomains(r.Context(), projectID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, domains)
}

// AddDomain handles POST /api/projects/{id}/domains.
func (h *Handler) AddDomain(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	var d DomainConfig
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		writeError(w, http.StatusBadRequest, "invalid domain payload")
		return
	}
	d.ProjectID = projectID
	if d.DomainName == "" {
		writeError(w, http.StatusBadRequest, "domain_name is required")
		return
	}

	if err := h.service.AddDomain(r.Context(), &d); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if h.proxy != nil {
		_ = h.proxy.Reload(r.Context())
	}
	writeJSON(w, http.StatusCreated, d)
}

// DeleteDomain handles DELETE /api/domains/{id}.
func (h *Handler) DeleteDomain(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.service.DeleteDomain(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if h.proxy != nil {
		_ = h.proxy.Reload(r.Context())
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

// GetEnvVars handles GET /api/projects/{id}/env.
func (h *Handler) GetEnvVars(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	envs, err := h.service.GetEnvVars(r.Context(), projectID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if envs == nil {
		envs = make(map[string]string)
	}
	writeJSON(w, http.StatusOK, envs)
}

// SetEnvVars handles PUT /api/projects/{id}/env.
func (h *Handler) SetEnvVars(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	var payload SetEnvVarsRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid environment variable dictionary payload")
		return
	}

	for key, value := range payload {
		if err := h.service.SetEnvVar(r.Context(), projectID, key, value); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

// ListProjectCanvasSummaries handles GET /api/canvas/projects.
func (h *Handler) ListProjectCanvasSummaries(w http.ResponseWriter, r *http.Request) {
	summaries, err := h.service.ListProjectCanvasSummaries(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, summaries)
}

// GetProjectCanvasSummary handles GET /api/projects/{id}/summary.
func (h *Handler) GetProjectCanvasSummary(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	summary, err := h.service.GetProjectCanvasSummary(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, summary)
}

// GetEnvironmentCanvas handles GET /api/environments/{id}/canvas.
func (h *Handler) GetEnvironmentCanvas(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	canvas, err := h.service.GetEnvironmentCanvas(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, canvas)
}
