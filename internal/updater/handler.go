package updater

import (
	"encoding/json"
	"net/http"
)

// Handler manages pure HTTP endpoints for version checks and auto-update deployments.
type Handler struct {
	service *UpdaterService
}

// NewHandler initializes a new Updater HTTP handler.
func NewHandler(service *UpdaterService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetUpdateStatus(w http.ResponseWriter, r *http.Request) {
	info, err := h.service.CheckForUpdate(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, info)
}

func (h *Handler) CheckUpdate(w http.ResponseWriter, r *http.Request) {
	info, err := h.service.CheckForUpdate(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, info)
}

func (h *Handler) DeployUpdate(w http.ResponseWriter, r *http.Request) {
	if err := h.service.DeployUpdate(r.Context()); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deploying"})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
