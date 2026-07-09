package env

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	service *Service
}

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
