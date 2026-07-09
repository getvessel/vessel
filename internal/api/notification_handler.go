package api

import (
	"encoding/json"
	"net/http"

	"vessel.dev/vessel/internal/notifier"
	"vessel.dev/vessel/internal/store"
	"vessel.dev/vessel/internal/types"
)

type NotificationHandler struct {
	store    *store.Store
	notifier *notifier.NotifierService
}

func NewNotificationHandler(s *store.Store, n *notifier.NotifierService) *NotificationHandler {
	return &NotificationHandler{store: s, notifier: n}
}

func (h *NotificationHandler) GetIntegrations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	integ, err := h.store.GetNotificationIntegration()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(integ)
}

func (h *NotificationHandler) SaveIntegrations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var integ types.NotificationIntegration
	if err := json.NewDecoder(r.Body).Decode(&integ); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.store.SaveNotificationIntegration(&integ); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(integ)
}

func (h *NotificationHandler) TestNotification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Channel   string `json:"channel"`
		ProjectID string `json:"projectId,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	event := &types.NotificationEvent{
		Title:     "🚀 Vessel Notification Test",
		Message:   "If you are reading this, your " + req.Channel + " notification channel is configured correctly!",
		Level:     "success",
		ProjectID: req.ProjectID,
		URL:       "http://localhost:8080/settings/notifications",
	}

	if err := h.notifier.Send(event); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "message": "Test notification sent successfully over " + req.Channel})
}

func (h *NotificationHandler) GetProjectPreferences(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	projectID := r.URL.Query().Get("id")
	if projectID == "" {
		http.Error(w, "Missing project id parameter", http.StatusBadRequest)
		return
	}

	pref, err := h.store.GetProjectNotificationPref(projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pref)
}

func (h *NotificationHandler) SaveProjectPreferences(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var pref types.ProjectNotificationPref
	if err := json.NewDecoder(r.Body).Decode(&pref); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.store.SaveProjectNotificationPref(&pref); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pref)
}
