package notification

import "context"

// Dispatcher sends notification events to configured channels.
type Dispatcher interface {
	Send(event *NotificationEvent) error
}

// Service provides business logic for notification configuration and dispatch.
type Service struct {
	repo       Repository
	dispatcher Dispatcher
}

// NewService creates a new notification Service.
func NewService(repo Repository, dispatcher Dispatcher) *Service {
	return &Service{repo: repo, dispatcher: dispatcher}
}

// GetIntegration returns the global notification integration settings.
func (s *Service) GetIntegration(ctx context.Context) (*NotificationIntegration, error) {
	return s.repo.GetIntegration(ctx)
}

// SaveIntegration persists the global notification integration settings.
func (s *Service) SaveIntegration(ctx context.Context, n *NotificationIntegration) error {
	return s.repo.SaveIntegration(ctx, n)
}

// GetProjectPref returns notification preferences for a project.
func (s *Service) GetProjectPref(ctx context.Context, projectID string) (*ProjectNotificationPref, error) {
	return s.repo.GetProjectPref(ctx, projectID)
}

// SaveProjectPref persists notification preferences for a project.
func (s *Service) SaveProjectPref(ctx context.Context, pref *ProjectNotificationPref) error {
	return s.repo.SaveProjectPref(ctx, pref)
}

// SendTest sends a test notification event through the given channel.
func (s *Service) SendTest(channel, projectID string) error {
	event := &NotificationEvent{
		Title:     "Vessel Notification Test",
		Message:   "If you are reading this, your " + channel + " notification channel is configured correctly!",
		Level:     "success",
		ProjectID: projectID,
		URL:       "http://localhost:8080/settings/notifications",
	}
	return s.dispatcher.Send(event)
}
