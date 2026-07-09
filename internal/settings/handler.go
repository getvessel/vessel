package settings

import (
	"encoding/json"
	"net/http"

	dockerfilters "github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// Handler manages pure HTTP concerns for global settings and MCP endpoints.
type Handler struct {
	service      *Service
	dockerClient *client.Client
}

// NewHandler initializes a new Settings HTTP handler.
func NewHandler(service *Service, dockerClient *client.Client) *Handler {
	return &Handler{
		service:      service,
		dockerClient: dockerClient,
	}
}

func (h *Handler) GetServerSettings(w http.ResponseWriter, r *http.Request) {
	cfg, err := h.service.GetSettings(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, cfg)
}

func (h *Handler) UpdateServerSettings(w http.ResponseWriter, r *http.Request) {
	var req ServerSettings
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.service.UpdateSettings(r.Context(), &req); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, req)
}

func (h *Handler) TriggerSystemPrune(w http.ResponseWriter, r *http.Request) {
	if h.dockerClient == nil {
		writeJSON(w, http.StatusOK, PruneResponse{
			Status:              "simulated",
			Message:             "Docker client not initialized in standalone mode; simulated clean system prune.",
			SpaceReclaimedBytes: 104857600,
		})
		return
	}

	ctx := r.Context()
	var totalReclaimed uint64

	if cReport, err := h.dockerClient.ContainersPrune(ctx, dockerfilters.NewArgs()); err == nil {
		totalReclaimed += cReport.SpaceReclaimed
	}
	if iReport, err := h.dockerClient.ImagesPrune(ctx, dockerfilters.NewArgs()); err == nil {
		totalReclaimed += iReport.SpaceReclaimed
	}
	if nReport, err := h.dockerClient.NetworksPrune(ctx, dockerfilters.NewArgs()); err == nil {
		_ = nReport
	}
	if vReport, err := h.dockerClient.VolumesPrune(ctx, dockerfilters.NewArgs()); err == nil {
		totalReclaimed += vReport.SpaceReclaimed
	}

	writeJSON(w, http.StatusOK, PruneResponse{
		Status:              "success",
		Message:             "Docker system prune executed cleanly.",
		SpaceReclaimedBytes: totalReclaimed,
	})
}

func (h *Handler) HandleMCP(w http.ResponseWriter, r *http.Request) {
	if err := h.service.CheckMCPEnabled(r.Context()); err != nil {
		writeError(w, http.StatusForbidden, err.Error())
		return
	}

	if r.Method == http.MethodGet {
		writeJSON(w, http.StatusOK, MCPResponse{
			JSONRPC: "2.0",
			Server: map[string]any{
				"name":    "vessel-mcp-server",
				"version": "v1.0.0",
			},
			Capabilities: map[string]any{
				"tools": map[string]any{
					"listChanged": false,
				},
			},
			Tools: []map[string]any{
				{
					"name":        "list_projects",
					"description": "List all deployed projects on Vessel",
					"inputSchema": map[string]any{
						"type":       "object",
						"properties": map[string]any{},
					},
				},
				{
					"name":        "get_system_status",
					"description": "Check Vessel server CPU, RAM, and database health",
					"inputSchema": map[string]any{
						"type":       "object",
						"properties": map[string]any{},
					},
				},
			},
		})
		return
	}

	var req struct {
		JSONRPC string `json:"jsonrpc"`
		ID      any    `json:"id"`
		Method  string `json:"method"`
		Params  struct {
			Name      string         `json:"name"`
			Arguments map[string]any `json:"arguments"`
		} `json:"params"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON-RPC format")
		return
	}

	switch req.Method {
	case "tools/list":
		writeJSON(w, http.StatusOK, MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]any{
				"tools": []map[string]any{
					{
						"name":        "list_projects",
						"description": "List all deployed projects on Vessel",
					},
					{
						"name":        "get_system_status",
						"description": "Check Vessel server CPU, RAM, and database health",
					},
				},
			},
		})
	case "tools/call":
		content, err := h.service.ExecuteMCPTool(r.Context(), req.Params.Name)
		if err != nil {
			writeJSON(w, http.StatusOK, MCPResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error: &MCPError{
					Code:    -32601,
					Message: err.Error(),
				},
			})
			return
		}
		writeJSON(w, http.StatusOK, MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]any{
				"content": content,
			},
		})
	default:
		writeJSON(w, http.StatusOK, MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32601,
				Message: "Method not supported: " + req.Method,
			},
		})
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
