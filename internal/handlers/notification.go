package handlers

import (
	"github.com/labstack/echo/v4"

	"encoding/json"
	"net/http"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/services"
)

type NotificationHandler struct {
	notificationService *services.NotificationService
}

func NewNotificationHandler(ns *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService: ns}
}

func (h *NotificationHandler) GetIntegrations(c echo.Context) error {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return nil
	}

	integ, err := h.notificationService.GetIntegration(r.Context())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}

	WriteJSON(w, http.StatusOK, integ)
}

func (h *NotificationHandler) SaveIntegrations(c echo.Context) error {
	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return nil
	}

	var integ models.NotificationIntegration
	if err := c.Bind(&integ); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	if err := h.notificationService.SaveIntegration(r.Context(), &integ); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}

	WriteJSON(w, http.StatusOK, integ)
}

func (h *NotificationHandler) TestNotification(c echo.Context) error {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return nil
	}

	var req struct {
		Channel   string `json:"channel"`
		ProjectID string `json:"projectId,omitempty"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	if err := h.notificationService.SendTest(req.Channel, req.ProjectID); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return nil
	}

	WriteJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"message": "Test notification sent successfully over " + req.Channel,
	})
}

func (h *NotificationHandler) GetProjectPreferences(c echo.Context) error {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return nil
	}

	projectID := c.Param("id")
	if projectID == "" {
		WriteError(w, http.StatusBadRequest, "Missing project id parameter")
		return nil
	}

	pref, err := h.notificationService.GetProjectPref(r.Context(), projectID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}

	WriteJSON(w, http.StatusOK, pref)
}

func (h *NotificationHandler) SaveProjectPreferences(c echo.Context) error {
	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return nil
	}

	projectID := c.Param("id")
	if projectID == "" {
		WriteError(w, http.StatusBadRequest, "Missing project id parameter")
		return nil
	}

	var pref models.ProjectNotificationPref
	if err := c.Bind(&pref); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	pref.ProjectID = projectID

	if err := h.notificationService.SaveProjectPref(r.Context(), &pref); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return nil
	}

	WriteJSON(w, http.StatusOK, pref)
}
