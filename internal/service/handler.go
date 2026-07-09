package service

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	envID := r.PathValue("id")

	var req CreateAppServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "app service name is required")
		return
	}

	svc := &AppService{
		ProjectID:     req.ProjectID,
		EnvironmentID: envID,
		Name:          req.Name,
		RepositoryURL: req.RepositoryURL,
		Branch:        req.Branch,
		InternalPort:  req.InternalPort,
		Domain:        req.Domain,
	}
	if svc.InternalPort == 0 {
		svc.InternalPort = 3000
	}

	if err := h.repo.Create(r.Context(), svc); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, svc)
}

func (h *Handler) ListByEnvironment(w http.ResponseWriter, r *http.Request) {
	envID := r.PathValue("id")

	apps, err := h.repo.ListByEnvironment(r.Context(), envID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, apps)
}

func (h *Handler) ListByProject(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")

	apps, err := h.repo.ListByProject(r.Context(), projectID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, apps)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	svc, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, svc)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	existing, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	var req UpdateAppServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	existing.Name = req.Name
	existing.RepositoryURL = req.RepositoryURL
	existing.Branch = req.Branch
	existing.InternalPort = req.InternalPort
	existing.Domain = req.Domain
	existing.ContainerID = req.ContainerID
	existing.Status = req.Status

	if err := h.repo.Update(r.Context(), existing); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, existing)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.repo.Delete(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
