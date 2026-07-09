package project_env

import (
	"encoding/json"
	"net/http"
)

// Handler serves HTTP requests for project-level environment variables.
type Handler struct {
	service *Service
}

// NewHandler creates a new project_env Handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// GetVars handles GET /api/projects/{id}/env.
func (h *Handler) GetVars(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	envs, err := h.service.GetVars(r.Context(), projectID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if envs == nil {
		envs = make(map[string]string)
	}
	writeJSON(w, http.StatusOK, envs)
}

// SetVars handles PUT /api/projects/{id}/env.
func (h *Handler) SetVars(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	var payload VarsRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid environment variable dictionary payload")
		return
	}

	for key, value := range payload {
		if err := h.service.SetVar(r.Context(), projectID, key, value); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}
