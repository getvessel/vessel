package service

import "context"

//go:generate mockgen -source=repository.go -destination=mock_repository_test.go -package=service

type Repository interface {
	Create(ctx context.Context, svc *AppService) error
	GetByID(ctx context.Context, id string) (*AppService, error)
	ListByEnvironment(ctx context.Context, environmentID string) ([]*AppService, error)
	ListByProject(ctx context.Context, projectID string) ([]*AppService, error)
	Update(ctx context.Context, svc *AppService) error
	Delete(ctx context.Context, id string) error
}
