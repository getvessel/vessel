package job

import (
	"context"
	"time"
)

type Repository interface {
	Create(ctx context.Context, j *Job) error
	GetByID(ctx context.Context, id string) (*Job, error)
	ListByProject(ctx context.Context, projectID string) ([]Job, error)
	Update(ctx context.Context, j *Job) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id, status string, lastRunAt *time.Time, output string) error
}
