package handlers

import (
	"github.com/labstack/echo/v4"

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

func (h *ProjectSettingsHandler) ListWebhooks(c echo.Context) error {
	projectID := c.Param("projectId")
	if projectID == "" {
		WriteError(w, http.StatusBadRequest, "missing projectId")
		return nil
	}
	list, err := h.settingsService.ListWebhooks(r.Context(), projectID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusOK, list)
}

func (h *ProjectSettingsHandler) CreateWebhook(c echo.Context) error {
	projectID := c.Param("projectId")
	if projectID == "" {
		WriteError(w, http.StatusBadRequest, "missing projectId")
		return nil
	}
	var req models.Webhook
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	req.ProjectID = projectID
	created, err := h.settingsService.CreateWebhook(r.Context(), &req)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusCreated, created)
}

func (h *ProjectSettingsHandler) DeleteWebhook(c echo.Context) error {
	projectID := c.Param("projectId")
	id := c.Param("id")
	if projectID == "" || id == "" {
		WriteError(w, http.StatusBadRequest, "missing projectId or id")
		return nil
	}
	if err := h.settingsService.DeleteWebhook(r.Context(), id, projectID); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ProjectSettingsHandler) ListTokens(c echo.Context) error {
	projectID := c.Param("projectId")
	if projectID == "" {
		WriteError(w, http.StatusBadRequest, "missing projectId")
		return nil
	}
	list, err := h.settingsService.ListTokens(r.Context(), projectID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusOK, list)
}

func (h *ProjectSettingsHandler) CreateToken(c echo.Context) error {
	projectID := c.Param("projectId")
	if projectID == "" {
		WriteError(w, http.StatusBadRequest, "missing projectId")
		return nil
	}
	var req models.ProjectToken
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	req.ProjectID = projectID
	token, raw, err := h.settingsService.CreateToken(r.Context(), &req)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusCreated, map[string]any{
		"id":        token.ID,
		"name":      token.Name,
		"token":     raw,
		"createdAt": token.CreatedAt,
	})
}

func (h *ProjectSettingsHandler) DeleteToken(c echo.Context) error {
	projectID := c.Param("projectId")
	id := c.Param("id")
	if projectID == "" || id == "" {
		WriteError(w, http.StatusBadRequest, "missing projectId or id")
		return nil
	}
	if err := h.settingsService.DeleteToken(r.Context(), id, projectID); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ProjectSettingsHandler) ListMembers(c echo.Context) error {
	projectID := c.Param("projectId")
	if projectID == "" {
		WriteError(w, http.StatusBadRequest, "missing projectId")
		return nil
	}
	list, err := h.settingsService.ListMembers(r.Context(), projectID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusOK, list)
}

func (h *ProjectSettingsHandler) AddMember(c echo.Context) error {
	projectID := c.Param("projectId")
	if projectID == "" {
		WriteError(w, http.StatusBadRequest, "missing projectId")
		return nil
	}
	var req models.ProjectMember
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	req.ProjectID = projectID
	added, err := h.settingsService.AddMember(r.Context(), &req)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	WriteJSON(w, http.StatusCreated, added)
}

func (h *ProjectSettingsHandler) RemoveMember(c echo.Context) error {
	projectID := c.Param("projectId")
	id := c.Param("id")
	if projectID == "" || id == "" {
		WriteError(w, http.StatusBadRequest, "missing projectId or id")
		return nil
	}
	if err := h.settingsService.RemoveMember(r.Context(), id, projectID); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}
	w.WriteHeader(http.StatusNoContent)
}
