package service_var

import "context"

type ServiceRepository interface {
	GetByID(ctx context.Context, id string) (*ServiceDTO, error)
}

type ServiceDTO struct {
	ID            string
	ProjectID     string
	EnvironmentID string
}

type Repository interface {
	ListByService(ctx context.Context, serviceID string) ([]*Variable, error)
	Create(ctx context.Context, v *Variable) error
	Update(ctx context.Context, v *Variable) error
	Delete(ctx context.Context, id, serviceID string) error
	BulkSet(ctx context.Context, serviceID, environmentID string, vars []*Variable) error
}
