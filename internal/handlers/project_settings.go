package handlers

import (
	"encoding/json"
	"net/http"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/services"
)

type ProjectSettingsHandler struct {
	settingsService *services.ProjectSettingsService
}

func NewProjectSettingsHandler(s *services.ProjectSettingsService) *ProjectSettingsHandler {
	return &ProjectSettingsHandler{settingsService: s}
}

func (h *ProjectSettingsHandler) ListWebhooks(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	if projectID == "" {
		WriteError(w, http.StatusBadRequest, "missing projectId")
		return
	}
	list, err := h.settingsService.ListWebhooks(r.Context(), projectID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, list)
}

func (h *ProjectSettingsHandler) CreateWebhook(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	if projectID == "" {
		WriteError(w, http.StatusBadRequest, "missing projectId")
		return
	}
	var req models.Webhook
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.ProjectID = projectID
	created, err := h.settingsService.CreateWebhook(r.Context(), &req)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, created)
}

func (h *ProjectSettingsHandler) DeleteWebhook(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	id := r.PathValue("id")
	if projectID == "" || id == "" {
		WriteError(w, http.StatusBadRequest, "missing projectId or id")
		return
	}
	if err := h.settingsService.DeleteWebhook(r.Context(), id, projectID); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ProjectSettingsHandler) ListTokens(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	if projectID == "" {
		WriteError(w, http.StatusBadRequest, "missing projectId")
		return
	}
	list, err := h.settingsService.ListTokens(r.Context(), projectID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, list)
}

func (h *ProjectSettingsHandler) CreateToken(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	if projectID == "" {
		WriteError(w, http.StatusBadRequest, "missing projectId")
		return
	}
	var req models.ProjectToken
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.ProjectID = projectID
	token, raw, err := h.settingsService.CreateToken(r.Context(), &req)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, map[string]any{
		"id":        token.ID,
		"name":      token.Name,
		"token":     raw,
		"createdAt": token.CreatedAt,
	})
}

func (h *ProjectSettingsHandler) DeleteToken(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	id := r.PathValue("id")
	if projectID == "" || id == "" {
		WriteError(w, http.StatusBadRequest, "missing projectId or id")
		return
	}
	if err := h.settingsService.DeleteToken(r.Context(), id, projectID); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ProjectSettingsHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	if projectID == "" {
		WriteError(w, http.StatusBadRequest, "missing projectId")
		return
	}
	list, err := h.settingsService.ListMembers(r.Context(), projectID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, list)
}

func (h *ProjectSettingsHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	if projectID == "" {
		WriteError(w, http.StatusBadRequest, "missing projectId")
		return
	}
	var req models.ProjectMember
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.ProjectID = projectID
	added, err := h.settingsService.AddMember(r.Context(), &req)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, added)
}

func (h *ProjectSettingsHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	id := r.PathValue("id")
	if projectID == "" || id == "" {
		WriteError(w, http.StatusBadRequest, "missing projectId or id")
		return
	}
	if err := h.settingsService.RemoveMember(r.Context(), id, projectID); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
