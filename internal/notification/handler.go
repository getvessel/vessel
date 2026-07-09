package notification

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
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	http.Error(w, msg, status)
}

func (h *Handler) GetIntegrations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	integ, err := h.service.GetIntegration(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, integ)
}

func (h *Handler) SaveIntegrations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var integ NotificationIntegration
	if err := json.NewDecoder(r.Body).Decode(&integ); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.service.SaveIntegration(r.Context(), &integ); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, integ)
}

func (h *Handler) TestNotification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		Channel   string `json:"channel"`
		ProjectID string `json:"projectId,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.service.SendTest(req.Channel, req.ProjectID); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"message": "Test notification sent successfully over " + req.Channel,
	})
}

func (h *Handler) GetProjectPreferences(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	projectID := r.PathValue("id")
	if projectID == "" {
		writeError(w, http.StatusBadRequest, "Missing project id parameter")
		return
	}

	pref, err := h.service.GetProjectPref(r.Context(), projectID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, pref)
}

func (h *Handler) SaveProjectPreferences(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var pref ProjectNotificationPref
	if err := json.NewDecoder(r.Body).Decode(&pref); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.service.SaveProjectPref(r.Context(), &pref); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, pref)
}
