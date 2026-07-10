package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/services"
)

type SettingsHandler struct {
	settingsService *services.SettingsService
}

func NewSettingsHandler(s *services.SettingsService) *SettingsHandler {
	return &SettingsHandler{settingsService: s}
}

func (h *SettingsHandler) GetSettings(c echo.Context) error {
	s, err := h.settingsService.GetSettings(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, s)
}

func (h *SettingsHandler) UpdateSettings(c echo.Context) error {
	var payload models.ServerSettings
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	if err := h.settingsService.UpdateSettings(c.Request().Context(), &payload); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, payload)
}

func (h *SettingsHandler) ListTeamNotificationChannels(c echo.Context) error {
	teamID := c.QueryParam("teamId")
	if teamID == "" {
		teamID = "default"
	}
	channels, err := h.settingsService.ListTeamNotificationChannels(c.Request().Context(), teamID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, channels)
}

func (h *SettingsHandler) SaveTeamNotificationChannel(c echo.Context) error {
	var payload models.TeamNotificationChannel
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	if payload.TeamID == "" {
		payload.TeamID = "default"
	}
	if err := h.settingsService.SaveTeamNotificationChannel(c.Request().Context(), &payload); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, payload)
}

func (h *SettingsHandler) GetTeamNotificationChannel(c echo.Context) error {
	id := c.Param("id")
	channel, err := h.settingsService.GetTeamNotificationChannel(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, channel)
}

func (h *SettingsHandler) DeleteTeamNotificationChannel(c echo.Context) error {
	id := c.Param("id")
	if err := h.settingsService.DeleteTeamNotificationChannel(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *SettingsHandler) HandleMCPRequest(c echo.Context) error {
	if c.Request().Method != http.MethodPost {
		return c.JSON(http.StatusMethodNotAllowed, map[string]string{"error": "Only POST requests allowed for MCP JSON-RPC"})
	}
	if err := h.settingsService.CheckMCPEnabled(c.Request().Context()); err != nil {
		return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
	}
	var req struct {
		JSONRPC string `json:"jsonrpc"`
		ID      any    `json:"id"`
		Method  string `json:"method"`
		Params  struct {
			Name string `json:"name"`
		} `json:"params"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}
	if req.JSONRPC != "2.0" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Only JSON-RPC 2.0 is supported"})
	}
	if req.Method == "tools/list" {
		return c.JSON(http.StatusOK, map[string]any{
			"jsonrpc": "2.0",
			"id":      req.ID,
			"result": map[string]any{
				"tools": []map[string]any{
					{
						"name":        "list_projects",
						"description": "List all deployment projects registered in this Vessel instance.",
						"inputSchema": map[string]any{
							"type":       "object",
							"properties": map[string]any{},
						},
					},
					{
						"name":        "get_system_status",
						"description": "Check basic operational and health metrics of the Vessel platform.",
						"inputSchema": map[string]any{
							"type":       "object",
							"properties": map[string]any{},
						},
					},
				},
			},
		})
	}
	if req.Method == "tools/call" {
		content, err := h.settingsService.ExecuteMCPTool(c.Request().Context(), req.Params.Name)
		if err != nil {
			return c.JSON(http.StatusOK, map[string]any{
				"jsonrpc": "2.0",
				"id":      req.ID,
				"error": map[string]any{
					"code":    -32601,
					"message": err.Error(),
				},
			})
		}
		return c.JSON(http.StatusOK, map[string]any{
			"jsonrpc": "2.0",
			"id":      req.ID,
			"result": map[string]any{
				"content": content,
			},
		})
	}
	return c.JSON(http.StatusOK, map[string]any{
		"jsonrpc": "2.0",
		"id":      req.ID,
		"error": map[string]any{
			"code":    -32601,
			"message": "Method not found: " + req.Method,
		},
	})
}
