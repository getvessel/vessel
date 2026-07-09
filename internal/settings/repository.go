package settings

import "context"

// Repository defines ONLY the data operations required by the Settings Service (`Accept interfaces, return structs`).
type Repository interface {
	GetServerSettings(ctx context.Context) (*ServerSettings, error)
	UpdateServerSettings(ctx context.Context, cfg *ServerSettings) error
	ListProjects(ctx context.Context) ([]map[string]any, error)
}
