package settings

import "context"

type Repository interface {
	GetServerSettings(ctx context.Context) (*ServerSettings, error)
	UpdateServerSettings(ctx context.Context, cfg *ServerSettings) error
	ListProjects(ctx context.Context) ([]map[string]any, error)
}
