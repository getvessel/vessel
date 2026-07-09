package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/services"
)

type SettingsHandler struct {
	settingsService *services.SettingsService
}

func NewSettingsHandler(s *services.SettingsService) *SettingsHandler {
	return &SettingsHandler{settingsService: s}
}

func (h *SettingsHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	s, err := h.settingsService.GetSettings(r.Context())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, s)
}

func (h *SettingsHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	var payload models.ServerSettings
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	if err := h.settingsService.UpdateSettings(r.Context(), &payload); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, payload)
}

func (h *SettingsHandler) GetNotificationIntegration(w http.ResponseWriter, r *http.Request) {
	n, err := h.settingsService.GetNotificationIntegration(r.Context())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, n)
}

func (h *SettingsHandler) SaveNotificationIntegration(w http.ResponseWriter, r *http.Request) {
	var payload models.NotificationIntegration
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	if err := h.settingsService.SaveNotificationIntegration(r.Context(), &payload); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, payload)
}

func (h *SettingsHandler) GetProjectNotificationPref(w http.ResponseWriter, r *http.Request) {
	projectID := strings.TrimPrefix(r.URL.Path, "/api/settings/notifications/project/")
	if projectID == "" || projectID == r.URL.Path {
		WriteError(w, http.StatusBadRequest, "missing project id parameter")
		return
	}
	pref, err := h.settingsService.GetProjectNotificationPref(r.Context(), projectID)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, pref)
}

func (h *SettingsHandler) SaveProjectNotificationPref(w http.ResponseWriter, r *http.Request) {
	var payload models.ProjectNotificationPref
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		WriteError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	if err := h.settingsService.SaveProjectNotificationPref(r.Context(), &payload); err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, payload)
}

func (h *SettingsHandler) HandleMCPRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "Only POST requests allowed for MCP JSON-RPC")
		return
	}
	if err := h.settingsService.CheckMCPEnabled(r.Context()); err != nil {
		WriteError(w, http.StatusForbidden, err.Error())
		return
	}
	var req struct {
		JSONRPC string `json:"jsonrpc"`
		ID      any    `json:"id"`
		Method  string `json:"method"`
		Params  struct {
			Name string `json:"name"`
		} `json:"params"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid JSON-RPC 2.0 payload")
		return
	}
	if req.JSONRPC != "2.0" {
		WriteError(w, http.StatusBadRequest, "Only JSON-RPC 2.0 is supported")
		return
	}
	if req.Method == "tools/list" {
		WriteJSON(w, http.StatusOK, map[string]any{
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
		return
	}
	if req.Method == "tools/call" {
		content, err := h.settingsService.ExecuteMCPTool(r.Context(), req.Params.Name)
		if err != nil {
			WriteJSON(w, http.StatusOK, map[string]any{
				"jsonrpc": "2.0",
				"id":      req.ID,
				"error": map[string]any{
					"code":    -32601,
					"message": err.Error(),
				},
			})
			return
		}
		WriteJSON(w, http.StatusOK, map[string]any{
			"jsonrpc": "2.0",
			"id":      req.ID,
			"result": map[string]any{
				"content": content,
			},
		})
		return
	}
	WriteJSON(w, http.StatusOK, map[string]any{
		"jsonrpc": "2.0",
		"id":      req.ID,
		"error": map[string]any{
			"code":    -32601,
			"message": "Method not found: " + req.Method,
		},
	})
}
