package notification

import "context"

// Repository defines the persistence operations for notification configuration.
type Repository interface {
	GetIntegration(ctx context.Context) (*NotificationIntegration, error)
	SaveIntegration(ctx context.Context, n *NotificationIntegration) error
	GetProjectPref(ctx context.Context, projectID string) (*ProjectNotificationPref, error)
	SaveProjectPref(ctx context.Context, pref *ProjectNotificationPref) error
}
