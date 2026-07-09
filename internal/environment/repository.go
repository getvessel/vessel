package environment

import "context"

type Repository interface {
	Get(ctx context.Context, id string) (*Config, error)
	ListByProject(ctx context.Context, projectID string) ([]Config, error)
	Create(ctx context.Context, env *Config) error
	Delete(ctx context.Context, id string) error
}
