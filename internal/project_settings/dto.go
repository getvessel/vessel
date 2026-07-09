package project_settings

type CreateWebhookRequest struct {
	URL                   string   `json:"url"`
	EventTypes            []string `json:"eventTypes"`
	IncludePREnvironments bool     `json:"includePrEnvironments"`
}

type CreateTokenRequest struct {
	Name          string `json:"name"`
	EnvironmentID string `json:"environmentId"`
}

type AddMemberRequest struct {
	Email      string `json:"email"`
	Permission string `json:"permission"`
}
