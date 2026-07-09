package domain

import "context"

type Repository interface {
	ListByProject(ctx context.Context, projectID string) ([]Config, error)
	ListAll(ctx context.Context) ([]Config, error)
	Create(ctx context.Context, d *Config) error
	Delete(ctx context.Context, id string) error
}
