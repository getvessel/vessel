package settings

// UpdateSettingsRequest represents the payload to update global server settings.
type UpdateSettingsRequest struct {
	ServerSettings
}

// PruneResponse represents the result of a Docker system prune operation.
type PruneResponse struct {
	Status              string `json:"status"`
	Message             string `json:"message"`
	SpaceReclaimedBytes uint64 `json:"spaceReclaimedBytes"`
}

// MCPResponse represents JSON-RPC responses for MCP server queries.
type MCPResponse struct {
	JSONRPC      string           `json:"jsonrpc"`
	ID           any              `json:"id,omitempty"`
	Result       any              `json:"result,omitempty"`
	Error        *MCPError        `json:"error,omitempty"`
	Server       map[string]any   `json:"server,omitempty"`
	Tools        []map[string]any `json:"tools,omitempty"`
	Capabilities map[string]any   `json:"capabilities,omitempty"`
}

// MCPError represents a JSON-RPC error block.
type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
