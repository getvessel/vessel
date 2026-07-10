package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/services"
)

type NotificationHandler struct {
	notificationService *services.NotificationService
}

func NewNotificationHandler(ns *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService: ns}
}

func (h *NotificationHandler) ListChannels(c echo.Context) error {
	if c.Request().Method != http.MethodGet {
		return c.JSON(http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
	teamID := c.QueryParam("teamId")
	if teamID == "" {
		// fallback for now
		teamID = "default"
	}
	channels, err := h.notificationService.ListChannels(c.Request().Context(), teamID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, channels)
}

func (h *NotificationHandler) SaveChannel(c echo.Context) error {
	if c.Request().Method != http.MethodPut && c.Request().Method != http.MethodPost {
		return c.JSON(http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
	var channel models.TeamNotificationChannel
	if err := c.Bind(&channel); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	if channel.TeamID == "" {
		channel.TeamID = "default"
	}
	if err := h.notificationService.SaveChannel(c.Request().Context(), &channel); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, channel)
}

func (h *NotificationHandler) DeleteChannel(c echo.Context) error {
	if c.Request().Method != http.MethodDelete {
		return c.JSON(http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing channel id"})
	}
	if err := h.notificationService.DeleteChannel(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *NotificationHandler) TestNotification(c echo.Context) error {
	if c.Request().Method != http.MethodPost {
		return c.JSON(http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
	var req struct {
		ChannelID string `json:"channelId"`
		TeamID    string `json:"teamId"`
		Provider  string `json:"provider"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	if req.Provider != "" {
		err := h.notificationService.TestGlobalNotification(c.Request().Context(), req.Provider)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "ok",
			"message": "Global test notification queued",
		})
	}

	// For testing a specific team channel
	err := h.notificationService.TestTeamNotification(c.Request().Context(), req.TeamID, req.ChannelID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status":  "ok",
		"message": "Test notification queued",
	})
}
