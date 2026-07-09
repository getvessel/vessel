package env

import "context"

// Repository defines persistence operations for project-level environment variables.
type Repository interface {
	GetVars(ctx context.Context, projectID string) (map[string]string, error)
	SetVar(ctx context.Context, projectID, key, value string) error
}
