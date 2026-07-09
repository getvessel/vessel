package deployment

import "context"

type Repository interface {
	Create(ctx context.Context, d *Deployment) error
	GetByID(ctx context.Context, id string) (*Deployment, error)
	ListByService(ctx context.Context, serviceID string) ([]*Deployment, error)
	Update(ctx context.Context, d *Deployment) error
	UpdateStatus(ctx context.Context, id, status, buildLogs, containerID string) error
}

type ServiceRepository interface {
	GetByID(ctx context.Context, id string) (any, error)
}
