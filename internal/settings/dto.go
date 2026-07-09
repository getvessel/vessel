package settings

type UpdateSettingsRequest struct {
	ServerSettings
}

type PruneResponse struct {
	Status              string `json:"status"`
	Message             string `json:"message"`
	SpaceReclaimedBytes uint64 `json:"spaceReclaimedBytes"`
}

type MCPResponse struct {
	JSONRPC      string           `json:"jsonrpc"`
	ID           any              `json:"id,omitempty"`
	Result       any              `json:"result,omitempty"`
	Error        *MCPError        `json:"error,omitempty"`
	Server       map[string]any   `json:"server,omitempty"`
	Tools        []map[string]any `json:"tools,omitempty"`
	Capabilities map[string]any   `json:"capabilities,omitempty"`
}

type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
